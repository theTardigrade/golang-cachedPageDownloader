package storage

import "time"

// Datum is used to structure the information contained
// in each of the files used in the cache.
type Datum struct {
	SetTime time.Time `json:"t"`
	Content []byte    `json:"c"`
}

// NewDatum returns a pointer to a newly allocated
// Datum struct, populating it with the given content
// and the current time.
func NewDatum(content []byte) *Datum {
	return &Datum{
		SetTime: time.Now().UTC(),
		Content: content,
	}
}
