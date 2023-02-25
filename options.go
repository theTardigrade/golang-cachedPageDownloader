package cachedPageDownloader

import "time"

// Options is used by the NewDownloader function
// to determine how it should be initialized.
type Options struct {
	CacheDir               string
	MaxCacheDuration       time.Duration
	ShouldKeepCacheOnClose bool
}
