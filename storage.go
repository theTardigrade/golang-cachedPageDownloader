package cachedPageDownloader

import "time"

type storage struct {
	SetTime time.Time `json:"t"`
	Content []byte    `json:"c"`
}

func newStorage(content []byte) *storage {
	return &storage{
		SetTime: time.Now().UTC(),
		Content: content,
	}
}
