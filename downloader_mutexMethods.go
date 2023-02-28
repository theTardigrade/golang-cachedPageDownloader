package cachedPageDownloader

import (
	"strings"

	namespacedMutex "github.com/theTardigrade/golang-namespacedMutex"
)

const (
	mutexNamespaceSeparator byte = '|'
)

var (
	mutexCollection = namespacedMutex.New(&namespacedMutex.Options{
		BucketCount:           4_194_319, // first prime number after 1 << 22
		MaxUniqueAttemptCount: 1 << 16,
	})
)

func mutexWriteNamespacePart(namespacePart string, builder *strings.Builder) {
	namespacePart = strings.TrimSpace(namespacePart)

	for _, r := range namespacePart {
		switch r {
		case '\\', rune(mutexNamespaceSeparator):
			builder.WriteByte('\\')
		}

		builder.WriteRune(r)
	}
}

func (downloader *Downloader) mutexNamespaceDefaultParts() (namespaceParts []string) {
	namespaceParts = []string{
		downloader.options.CacheDir,
	}

	for i, part := range namespaceParts {
		var builder strings.Builder

		mutexWriteNamespacePart(part, &builder)

		namespaceParts[i] = builder.String()
	}

	return
}

func (downloader *Downloader) mutexNamespace(namespaceParts []string) string {
	var builder strings.Builder

	if namespacePartsLen := len(namespaceParts); namespacePartsLen > 0 {
		mutexWriteNamespacePart(namespaceParts[0], &builder)

		for i := 1; i < namespacePartsLen; i++ {
			builder.WriteByte(mutexNamespaceSeparator)
			mutexWriteNamespacePart(namespaceParts[i], &builder)
		}
	}

	for _, part := range downloader.mutexNamespaceDefaultPartsCached {
		builder.WriteByte(mutexNamespaceSeparator)
		builder.WriteString(part)
	}

	return builder.String()
}

func (downloader *Downloader) mutexGetLocked(namespaceParts ...string) (currentMutex *namespacedMutex.MutexWrapper) {
	namespace := downloader.mutexNamespace(namespaceParts)
	currentMutex = mutexCollection.GetLocked(false, namespace)

	return
}

func (downloader *Downloader) mutexGetLockedIfUnique(
	namespaceParts []string,
	comparisonNamespaceParts ...[]string,
) (currentMutex *namespacedMutex.MutexWrapper, found bool) {
	namespace := downloader.mutexNamespace(namespaceParts)

	if len(comparisonNamespaceParts) == 0 {
		currentMutex = mutexCollection.GetLocked(false, namespace)
		found = true

		return
	}

	comparisonNamespaces := make([]string, len(comparisonNamespaceParts))

	for i, parts := range comparisonNamespaceParts {
		comparisonNamespaces[i] = downloader.mutexNamespace(parts)
	}

	currentMutex, found = mutexCollection.GetLockedIfUnique(false, namespace, comparisonNamespaces...)

	return
}
