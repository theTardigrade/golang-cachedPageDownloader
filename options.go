package cachedPageDownloader

import "time"

type Options struct {
	CacheDir               string
	MaxCacheDuration       time.Duration
	ShouldKeepCacheOnClose bool
}
