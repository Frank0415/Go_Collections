package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type safeMap struct {
	mu sync.Mutex
	v  map[string]int
}

func (Map *safeMap) val(str string) (int, bool) {
	Map.mu.Lock()
	defer Map.mu.Unlock()
	val, ok := (Map.v)[str]
	return val, ok
}

func (Map *safeMap) inc(str string, val int) {
	Map.mu.Lock()
	defer Map.mu.Unlock()
	(Map.v)[str] += val
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(UrlMap *safeMap, url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		return
	}
	if _, ok := UrlMap.val(url); ok {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		UrlMap.inc(url, -1)
		fmt.Println(err)
		return
	}
	UrlMap.inc(url, 1)
	fmt.Printf("found: %s %q\n", url, body)

	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Go(func() { Crawl(UrlMap, u, depth-1, fetcher) })
	}
	wg.Wait()

}

func Crawl_run() {
	var UrlMap safeMap
	UrlMap.v = make(map[string]int)
	Crawl(&UrlMap, "https://golang.org/", 4, fetcher)
}

func AccessMap(UrlMap *map[string]int, str []string) {

}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
