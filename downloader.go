package cachedPageDownloader

import (
	"os"
	"path/filepath"
)

const (
	downloaderCacheFileExt        = ".cache.tmp"
	downloaderCacheDirTempPattern = "golang-cachedPageDownloader-*"
)

// Downloader provides methods that do the main
// work of this package.
type Downloader struct {
	options                          Options
	isCacheDirTemp                   bool
	mutexNamespaceDefaultPartsCached []string
}

// NewDownloader returns a pointer to a newly
// allocated Downloader struct.
// An error will also be returned if the cache
// directory cannot be found or created.
func NewDownloader(options *Options) (downloader *Downloader, err error) {
	downloader = &Downloader{}

	if options != nil {
		downloader.options = *options
		options = &downloader.options
	}

	if options.CacheDir == "" {
		options.CacheDir, err = os.MkdirTemp("", downloaderCacheDirTempPattern)
		if err != nil {
			return
		}
		downloader.isCacheDirTemp = true
	} else {
		if err = os.MkdirAll(options.CacheDir, os.ModeDir); err != nil {
			return
		}
	}

	if !filepath.IsAbs(options.CacheDir) {
		if options.CacheDir, err = filepath.Abs(options.CacheDir); err != nil {
			return
		}
	}

	downloader.mutexNamespaceDefaultPartsCached = downloader.mutexNamespaceDefaultParts()

	return
}
