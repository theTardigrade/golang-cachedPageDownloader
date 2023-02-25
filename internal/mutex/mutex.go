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

// GetLocked returns a mutex from the collection,
// based on the hash value of the given key,
// after locking it.
func GetLocked(key string) (mutex *sync.Mutex) {
	mutex = get(key)

	mutex.Lock()

	return mutex
}

// GetUniqueLocked attempts to returns a mutex
// from the collection,
// based on the hash value of the given primary key,
// after locking it.
// However, if the mutex found using the primary key
// is identical to any of the mutexes found using any
// of the secondary keys, then no mutex is returned
// or locked.
func GetUniqueLocked(primaryKey string, secondaryKeys ...string) (mutex *sync.Mutex, found bool) {
	mutex = get(primaryKey)

	for i := 0; i < len(secondaryKeys); i++ {
		mutex2 := get(secondaryKeys[i])

		if mutex == mutex2 {
			return
		}
	}

	found = true

	mutex.Lock()

	return
}
