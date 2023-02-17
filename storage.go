package cachedPageDownloader

import "time"

type storage struct {
	SetTime time.Time `json:"t"`
	Content []byte    `json:"c"`
}
