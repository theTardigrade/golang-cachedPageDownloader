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

func get(key string) (mutex *sync.Mutex) {
	hashedKey := hash.Uint256String(key)
	index := hashedKey.Mod(hashedKey, countBig).Uint64()
	mutex = collection[index]

	return mutex
}

func GetLocked(key string) (mutex *sync.Mutex) {
	mutex = get(key)

	mutex.Lock()

	return mutex
}

func GetUniqueLocked(primaryKey string, otherKeys ...string) (mutex *sync.Mutex, found bool) {
	mutex = get(primaryKey)

	for i := 0; i < len(otherKeys); i++ {
		mutex2 := get(otherKeys[i])

		if mutex == mutex2 {
			return
		}
	}

	found = true

	mutex.Lock()

	return
}
