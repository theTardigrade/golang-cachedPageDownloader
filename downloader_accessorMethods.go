package cachedPageDownloader

import (
	"time"
)

func (downloader *Downloader) CacheDir() string {
	return downloader.options.CacheDir
}

func (downloader *Downloader) MaxCacheDuration() time.Duration {
	return downloader.options.MaxCacheDuration
}

func (downloader *Downloader) ShouldKeepCacheOnClose() bool {
	return downloader.options.ShouldKeepCacheOnClose
}
