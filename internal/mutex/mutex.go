package mutex

import (
	"sync"

	hash "github.com/theTardigrade/golang-hash"
)

const (
	count = 1 << 16
)

var (
	collection = make([]*sync.Mutex, count)
)

func init() {
	for i := 0; i < count; i++ {
		collection[i] = &sync.Mutex{}
	}
}

func Get(key string) *sync.Mutex {
	index := hash.Uint64String(key) % count

	return collection[index]
}
