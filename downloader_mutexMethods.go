package cachedPageDownloader

import (
	"strings"
	"sync"

	"github.com/theTardigrade/golang-cachedPageDownloader/internal/mutex"
)

const (
	mutexKeySeparator = ":"
)

func (downloader *Downloader) mutexKey(keyParts []string) string {
	var builder strings.Builder

	keyParts = append(keyParts, downloader.options.CacheDir)

	for i, part := range keyParts {
		if i > 0 {
			builder.WriteString(mutexKeySeparator)
		}

		builder.WriteString(part)
	}

	return builder.String()
}

func (downloader *Downloader) mutexLocked(primaryKeyParts ...string) (currentMutex *sync.Mutex) {
	primaryKey := downloader.mutexKey(primaryKeyParts)
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
