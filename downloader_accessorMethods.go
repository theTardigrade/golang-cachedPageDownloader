package cachedPageDownloader

import (
	"time"
)

// CacheDir returns the absolute path to the directory
// where the cache is located.
func (downloader *Downloader) CacheDir() string {
	return downloader.options.CacheDir
}

// MaxCacheDuration returns the maximum duration during
// which items in the cache will be considered valid.
func (downloader *Downloader) MaxCacheDuration() time.Duration {
	return downloader.options.MaxCacheDuration
}

// ShouldKeepCacheOnClose returns false if the cache
// directory will be deleted when the Close function
// is called.
func (downloader *Downloader) ShouldKeepCacheOnClose() bool {
	return downloader.options.ShouldKeepCacheOnClose
}
