package cachedPageDownloader

import "time"

type Options struct {
	CacheDir               string
	ShouldKeepCacheOnClose bool
	MaxCacheDuration       time.Duration
}
