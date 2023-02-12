# golang-cachedPageDownloader

This package makes it easy to download &mdash; or access a previously cached version of &mdash; webpages.

## Example

```golang
package main

import (
	"fmt"
	"time"

	cachedPageDownloader "github.com/theTardigrade/golang-cachedPageDownloader"
)

const (
	exampleURL = "https://google.com/"
)

func main() {
	downloader, err := cachedPageDownloader.NewDownloader(&cachedPageDownloader.Options{
		CacheDir:               "./cache",
		ShouldKeepCacheOnClose: true,
		MaxCacheDuration:       time.Minute * 5,
	})
	if err != nil {
		panic(err)
	}
	defer downloader.Close()

	// calling the function below will retrieve the content of the webpage from the internet
	content, isFromCache, err := downloader.Download(exampleURL)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(content))
	fmt.Println(isFromCache) // false

	fmt.Println("*****")

	// calling the function again will retrieve the content of the webpage from our cache
	content, isFromCache, err = downloader.Download(exampleURL)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(content))
	fmt.Println(isFromCache) // true
}
```

## Support

If you use this package, or find any value in it, please consider donating at [Ko-fi](https://ko-fi.com/thetardigrade).