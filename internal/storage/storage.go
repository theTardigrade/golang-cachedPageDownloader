package storage

import "time"

type Datum struct {
	SetTime time.Time `json:"t"`
	Content []byte    `json:"c"`
}

func NewDatum(content []byte) *Datum {
	return &Datum{
		SetTime: time.Now().UTC(),
		Content: content,
	}
}
