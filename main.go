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

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, history *fetchHistory, wg *sync.WaitGroup) {
	defer wg.Done()
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		if err := history.Set(u); err != nil {
			// println(err.Error())
			continue
		}
		wg.Add(1)
		go Crawl(u, depth-1, fetcher, history, wg)
	}
	return
}

func main() {
	wg := &sync.WaitGroup{}
	var history fetchHistory
	wg.Add(1)
	go Crawl("http://golang.org/", 4, fetcher, &history, wg)
	wg.Wait()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

// fetchHistory store for visited urls
type fetchHistory struct {
	sync.Mutex
	history_map map[string]struct{}
}

// Set saves url to history map
// if already exists returns error
func (h *fetchHistory) Set(url string) error {
	if h.history_map == nil {
		h.history_map = make(map[string]struct{}, 0)
	}
	if _, ok := h.history_map[url]; ok {
		return fmt.Errorf("already visited: %s", url)
	}
	// lock map before modify
	h.Lock()
	h.history_map[url] = struct{}{}
	// unlock map after
	h.Unlock()
	return nil
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	// time.Sleep(2 * time.Second)
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
