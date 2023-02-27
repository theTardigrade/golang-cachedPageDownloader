package cachedPageDownloader

import "time"

// storageDatum is used to structure the information contained
// in each of the files used in the cache.
type storageDatum struct {
	SetTime time.Time `json:"t"`
	Content []byte    `json:"c"`
}

// newStorageDatum returns a pointer to a newly allocated
// storageDatum struct, populating it with the given content
// and the current time.
func newStorageDatum(content []byte) *storageDatum {
	return &storageDatum{
		SetTime: time.Now().UTC(),
		Content: content,
	}
}
