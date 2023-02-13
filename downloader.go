package cachedPageDownloader

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/theTardigrade/golang-cachedPageDownloader/internal/mutex"
	hash "github.com/theTardigrade/golang-hash"
)

const (
	fileExt = ".tmp"
)

type Downloader struct {
	options   *Options
	isTempDir bool
}

func NewDownloader(options *Options) (downloader *Downloader, err error) {
	if options == nil {
		options = &Options{}
	}

	downloader = &Downloader{
		options: options,
	}

	if options.CacheDir == "" {
		options.CacheDir, err = os.MkdirTemp("", "golang-cachedPageDownloader-*")
		if err != nil {
			return
		}
		downloader.isTempDir = true
	} else {
		if err = os.MkdirAll(options.CacheDir, os.ModeDir); err != nil {
			if os.IsExist(err) {
				err = nil
			} else {
				return
			}
		}
	}

	if !filepath.IsAbs(options.CacheDir) {
		if options.CacheDir, err = filepath.Abs(options.CacheDir); err != nil {
			return
		}
	}

	return
}

func (downloader *Downloader) Close() (err error) {
	options := downloader.options

	if !options.ShouldKeepCacheOnClose {
		currentMutex := mutex.Get("C:" + options.CacheDir)

		defer currentMutex.Unlock()
		currentMutex.Lock()

		if downloader.isTempDir {
			if err = os.RemoveAll(options.CacheDir); err != nil {
				return
			}
		} else {
			var cacheDirContents []string

			cacheDirContents, err = filepath.Glob(filepath.Join(options.CacheDir, "*"+fileExt))
			if err != nil {
				return
			}

			for _, item := range cacheDirContents {
				if err = os.RemoveAll(item); err != nil {
					return
				}
			}
		}
	}

	return
}

func (downloader *Downloader) CacheDir() string {
	return downloader.options.CacheDir
}

func (downloader *Downloader) MaxCacheDuration() time.Duration {
	return downloader.options.MaxCacheDuration
}

func (downloader *Downloader) ShouldKeepCacheOnClose() bool {
	return downloader.options.ShouldKeepCacheOnClose
}

func (downloader *Downloader) Clean() (err error) {
	options := downloader.options

	if options.MaxCacheDuration == 0 {
		return
	}

	currentMutex := mutex.Get("C:" + options.CacheDir)

	defer currentMutex.Unlock()
	currentMutex.Lock()

	var cacheDirContents []string

	cacheDirContents, err = filepath.Glob(filepath.Join(options.CacheDir, "*"+fileExt))
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
		}
	}

	return
}

func (downloader *Downloader) Download(rawURL string) (content []byte, isFromCache bool, err error) {
	options := downloader.options

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	rawURL = parsedURL.String()

	fileHash := hash.Uint256String(rawURL).Text(62)
	fileName := fileHash + fileExt
	filePath := filepath.Join(options.CacheDir, fileName)

	currentMutexKey := fmt.Sprintf("D:%s:%s", downloader.options.CacheDir, rawURL)
	currentMutex := mutex.Get(currentMutexKey)

	defer currentMutex.Unlock()
	currentMutex.Lock()

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

	var fileInfo fs.FileInfo

	if fileInfo, err = os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
	} else if options.MaxCacheDuration == 0 || time.Since(fileInfo.ModTime()) <= options.MaxCacheDuration {
		content, err = os.ReadFile(filePath)
		if err != nil {
			return
		}

		contentReader := bytes.NewReader(content)

		var contentGzipReader *gzip.Reader

		if contentGzipReader, err = gzip.NewReader(contentReader); err != nil {
			return
		}

		if content, err = io.ReadAll(contentGzipReader); err != nil {
			return
		}

		found = true
	}

	return
}

func (downloader *Downloader) downloadFromInternet(rawURL string) (content []byte, err error) {
	var resp *http.Response

	if resp, err = http.Get(rawURL); err != nil {
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

	if _, err = fileWriter.Write(content); err != nil {
		return
	}

	return
}
