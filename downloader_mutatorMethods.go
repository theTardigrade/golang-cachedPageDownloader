package cachedPageDownloader

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/theTardigrade/golang-cachedPageDownloader/internal/storage"

	hash "github.com/theTardigrade/golang-hash"
)

// Close should be called whenever the Downloader is
// no longer needed.
func (downloader *Downloader) Close() (err error) {
	options := downloader.options

	currentMutex, currentMutexFound := downloader.mutexUniqueLocked(
		[]string{"CL"},
		[]string{"C"},
	)

	if currentMutexFound {
		defer currentMutex.Unlock()
	}

	if !options.ShouldKeepCacheOnClose {
		if err = downloader.Clear(); err != nil {
			return
		}
	}

	return
}

// Clear will remove all files from the cache directory.
func (downloader *Downloader) Clear() (err error) {
	options := downloader.options

	currentMutex := downloader.mutexLocked("C")

	defer currentMutex.Unlock()

	if downloader.isCacheDirTemp {
		err = os.RemoveAll(options.CacheDir)

		return
	}

	var cacheDirContents []string

	cacheDirContents, err = filepath.Glob(filepath.Join(options.CacheDir, "*"+downloaderCacheFileExt))
	if err != nil {
		return
	}

	for _, item := range cacheDirContents {
		if err = os.RemoveAll(item); err != nil {
			return
		}
	}

	dirEntries, err := os.ReadDir(options.CacheDir)
	if err != nil {
		return
	}

	if len(dirEntries) == 0 {
		if err = os.Remove(options.CacheDir); err != nil {
			return
		}
	}

	return
}

// Clean will remove all files that have existed for longer than
// the maximum duration from the cache directory.
func (downloader *Downloader) Clean() (err error) {
	options := downloader.options

	if options.MaxCacheDuration == 0 {
		return
	}

	currentMutex := downloader.mutexLocked("C")

	defer currentMutex.Unlock()

	var cacheDirContents []string

	cacheDirContents, err = filepath.Glob(filepath.Join(options.CacheDir, "*"+downloaderCacheFileExt))
	if err != nil {
		return
	}

	for _, filePath := range cacheDirContents {
		var fileInfo fs.FileInfo

		if fileInfo, err = os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				err = nil
			} else {
				return
			}
		} else if time.Since(fileInfo.ModTime()) > options.MaxCacheDuration {
			if err = os.RemoveAll(filePath); err != nil {
				return
			}
		} else {
			var fileDatum *storage.Datum

			fileDatum, err = downloader.readCachedJSON(filePath)
			if err != nil {
				return
			}

			if time.Since(fileDatum.SetTime) > options.MaxCacheDuration {
				if err = os.RemoveAll(filePath); err != nil {
					return
				}
			}
		}
	}

	return
}

// Download will retrieve an HTML page as a byte slice from the internet.
// If the webpage has been downloaded previously, the function will attempt
// to retrieve it from the cache, reducing the number of network requests.
func (downloader *Downloader) Download(rawURL string) (content []byte, isFromCache bool, err error) {
	options := downloader.options

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	rawURL = parsedURL.String()

	fileHash := hash.Uint256String(rawURL).Text(62)
	fileName := fileHash + downloaderCacheFileExt
	filePath := filepath.Join(options.CacheDir, fileName)

	currentMutex := downloader.mutexLocked("D", rawURL)

	defer currentMutex.Unlock()

	content, isFromCache, err = downloader.readFromCache(filePath)
	if err != nil || isFromCache {
		return
	}

	content, err = downloader.downloadFromInternet(rawURL)
	if err != nil {
		return
	}

	err = downloader.writeToCache(filePath, content)
	if err != nil {
		return
	}

	return
}

func (downloader *Downloader) readFromCache(filePath string) (content []byte, found bool, err error) {
	options := downloader.options

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}

		return
	}

	defer func() {
		if err == nil && !found {
			os.Remove(filePath)
		}
	}()

	if options.MaxCacheDuration != 0 && time.Since(fileInfo.ModTime()) > options.MaxCacheDuration {
		return
	}

	fileDatum, err := downloader.readCachedJSON(filePath)
	if err != nil {
		return
	}

	if options.MaxCacheDuration != 0 && time.Since(fileDatum.SetTime) > options.MaxCacheDuration {
		return
	}

	content = fileDatum.Content
	found = true

	return
}

func (downloader *Downloader) readCachedJSON(filePath string) (fileDatum *storage.Datum, err error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	contentReader := bytes.NewReader(content)

	contentGzipReader, err := gzip.NewReader(contentReader)
	if err != nil {
		return
	}
	defer contentGzipReader.Close()

	if content, err = io.ReadAll(contentGzipReader); err != nil {
		return
	}

	fileDatum = new(storage.Datum)

	if err = json.Unmarshal(content, fileDatum); err != nil {
		return
	}

	return
}

func (downloader *Downloader) downloadFromInternet(rawURL string) (content []byte, err error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = ErrStatusCodeNotOK
		return
	}

	if content, err = io.ReadAll(resp.Body); err != nil {
		return
	}

	return
}

func (downloader *Downloader) writeToCache(filePath string, content []byte) (err error) {
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	fileWriter, err := gzip.NewWriterLevel(file, gzip.BestCompression)
	if err != nil {
		return
	}
	defer fileWriter.Close()

	fileDatum := storage.NewDatum(content)

	content, err = json.Marshal(fileDatum)
	if err != nil {
		return
	}

	if _, err = fileWriter.Write(content); err != nil {
		return
	}

	return
}
