package main

import (
	"fmt"
	"sync"
)

type Stats struct {
	mu           sync.Mutex
	totalChecks  int64
	errorCount   int64
	totalElapsed map[string]int64
	countPerURL  map[string]int64
}

func NewStats() *Stats {
	return &Stats{
		totalElapsed: make(map[string]int64),
		countPerURL:  make(map[string]int64),
	}
}

func (s *Stats) AddCheck(url string, elapsed int64, isError bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.totalChecks++
	if isError {
		s.errorCount++
		return
	}

	s.totalElapsed[url] += elapsed
	s.countPerURL[url]++
}

func PrintStats(s *Stats) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("\n--- Statistics ---")
	fmt.Printf("Total checks: %d\n", s.totalChecks)
	fmt.Printf("Total errors: %d\n", s.errorCount)
	for url, count := range s.countPerURL {
		avg := float64(s.totalElapsed[url]) / float64(count)
		fmt.Printf("%s - Checks: %d, Average response time: %.2f ms\n", url, count, avg)
	}
	fmt.Println("------------------")
}
