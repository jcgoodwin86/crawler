package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

// addPageVisit increments the visit count for a URL and returns true if this is the first time we see it.
// We insert a placeholder PageData to mark it as visited; it'll be replaced later with full data.
func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if pageData, visited := cfg.pages[normalizedURL]; visited {
		// Increment visit count for existing page
		pageData.VisitCount++
		cfg.pages[normalizedURL] = pageData
		return false
	}

	// First visit - create new entry with count of 1
	cfg.pages[normalizedURL] = PageData{URL: normalizedURL, VisitCount: 1}
	return true
}

// setPageData safely stores the final PageData for a URL.
func (cfg *config) setPageData(normalizedURL string, data PageData) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.pages[normalizedURL] = data
}

func configure(rawBaseURL string, maxConcurrency int, maxPages int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	return &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}, nil
}
