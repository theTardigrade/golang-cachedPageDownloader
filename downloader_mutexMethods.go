package cachedPageDownloader

import (
	"fmt"
	"strings"
	"sync"

	"github.com/theTardigrade/golang-cachedPageDownloader/internal/mutex"
)

const (
	mutexKeySeparator byte = '|'
)

func (downloader *Downloader) mutexKeyDefaultParts() []string {
	return []string{
		downloader.options.CacheDir,
	}
}

func (downloader *Downloader) mutexKey(keyParts []string) string {
	var builder strings.Builder

	if keyPartsLen := len(keyParts); keyPartsLen > 0 {
		builder.WriteString(keyParts[0])

		for i := 1; i < keyPartsLen; i++ {
			builder.WriteByte(mutexKeySeparator)
			builder.WriteString(keyParts[i])
		}
	}

	for _, part := range downloader.mutexKeyDefaultParts() {
		builder.WriteByte(mutexKeySeparator)
		builder.WriteString(part)
	}

	return builder.String()
}

func (downloader *Downloader) mutexLocked(primaryKeyParts ...string) (currentMutex *sync.Mutex) {
	primaryKey := downloader.mutexKey(primaryKeyParts)
	fmt.Println("KEY", primaryKey)
	currentMutex = mutex.GetLocked(primaryKey)

	return

}

func (downloader *Downloader) mutexUniqueLocked(
	primaryKeyParts []string,
	secondaryKeyParts ...[]string,
) (currentMutex *sync.Mutex, found bool) {
	primaryKey := downloader.mutexKey(primaryKeyParts)

	if len(secondaryKeyParts) == 0 {
		currentMutex = mutex.GetLocked(primaryKey)
		found = true

		return
	}

	secondaryKeys := make([]string, len(secondaryKeyParts))

	for i, parts := range secondaryKeyParts {
		secondaryKeys[i] = downloader.mutexKey(parts)
	}

	currentMutex, found = mutex.GetUniqueLocked(primaryKey, secondaryKeys...)

	return
}
