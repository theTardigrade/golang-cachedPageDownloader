package cachedPageDownloader

import (
	"os"
	"path/filepath"
)

const (
	downloaderCacheFileExt = ".tmp"
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
