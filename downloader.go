package cachedPageDownloader

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	hash "github.com/theTardigrade/golang-hash"
)

const (
	fileExt = ".tmp"
)

type Downloader struct {
	options *Options
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
		options.isTempDir = true
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
		if options.isTempDir {
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

	var fileInfo fs.FileInfo

	if fileInfo, err = os.Stat(filePath); err != nil {
		if !os.IsNotExist(err) {
			return
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

		isFromCache = true
	}

	if !isFromCache {
		var resp *http.Response

		if resp, err = http.Get(rawURL); err != nil {
			return
		}
		defer resp.Body.Close()

		if content, err = io.ReadAll(resp.Body); err != nil {
			return
		}

		var file *os.File

		if file, err = os.Create(filePath); err != nil {
			return
		}

		var fileWriter *gzip.Writer

		if fileWriter, err = gzip.NewWriterLevel(file, gzip.BestCompression); err != nil {
			return
		}

		if _, err = fileWriter.Write(content); err != nil {
			return
		}

		if err = fileWriter.Close(); err != nil {
			return
		}
	}

	return
}
