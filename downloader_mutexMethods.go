package cachedPageDownloader

import (
	"strings"
	"sync"

	"github.com/theTardigrade/golang-cachedPageDownloader/internal/mutex"
)

const (
	mutexKeySeparator byte = ':'
)

func (downloader *Downloader) mutexKeyDefaultParts() []string {
	return []string{
		downloader.options.CacheDir,
	}
}

func (downloader *Downloader) mutexKey(keyParts []string) string {
	var builder strings.Builder
	var i int

	for _, part := range keyParts {
		if i > 0 {
			builder.WriteByte(mutexKeySeparator)
		}

		builder.WriteString(part)

		i++
	}

	for _, part := range downloader.mutexKeyDefaultParts() {
		if i > 0 {
			builder.WriteByte(mutexKeySeparator)
		}

		builder.WriteString(part)

		i++
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
