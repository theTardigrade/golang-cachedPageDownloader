package mutex

import (
	"math/big"
	"sync"

	hash "github.com/theTardigrade/golang-hash"
)

const (
	count = 4_194_319 // first prime number after 1<<22
)

var (
	collection = make([]*sync.Mutex, count)
	countBig   = big.NewInt(count)
)

func init() {
	for i := 0; i < count; i++ {
		collection[i] = &sync.Mutex{}
	}
}

func GetLocked(key string) *sync.Mutex {
	hashedKey := hash.Uint256String(key)
	index := hashedKey.Mod(hashedKey, countBig).Uint64()
	mutex := collection[index]

	mutex.Lock()

	return mutex
}
