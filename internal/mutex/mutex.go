package mutex

import (
	"math/big"
	"strconv"
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

func index(key string) uint64 {
	hashedKey := hash.Uint256String(key)

	return hashedKey.Mod(hashedKey, countBig).Uint64()
}

func get(key string) (mutex *sync.Mutex) {
	mutex = collection[index(key)]

	return mutex
}

// GetLocked returns a mutex from the collection,
// based on a hashed value for the given key,
// after locking it.
func GetLocked(key string) (mutex *sync.Mutex) {
	mutex = get(key)

	mutex.Lock()

	return mutex
}

const (
	getUniqueLockedMaxAttempts = 1 << 11
)

func getUniqueLockedAttempt(primaryKey string, secondaryKeys ...string) (mutex *sync.Mutex, found bool) {
	mutex = get(primaryKey)

	for i := 0; i < len(secondaryKeys); i++ {
		secondaryMutex := get(secondaryKeys[i])

		if mutex == secondaryMutex {
			mutex = nil
			return
		}
	}

	found = true

	mutex.Lock()

	return
}

// GetUniqueLocked attempts to return a mutex
// from the collection,
// based on a hashed value for the given primary key,
// after locking it.
// However, if the mutex found using the primary key
// is identical to any of the mutexes found using any
// of the secondary keys, then no mutex is returned
// or locked.
func GetUniqueLocked(primaryKey string, secondaryKeys ...string) (mutex *sync.Mutex, found bool) {
	mutex, found = getUniqueLockedAttempt(primaryKey, secondaryKeys...)
	if found {
		return
	}

	for i := 2; i <= getUniqueLockedMaxAttempts; i++ {
		mutex, found = getUniqueLockedAttempt(primaryKey+strconv.Itoa(i), secondaryKeys...)
		if found {
			return
		}
	}

	return
}
