package stats

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Stats struct {
	mu            sync.Mutex
	totalChecks   int64
	errorCount    int64
	totalElapsed  map[string]int64
	countPerURL   map[string]int64
	totalSize     map[string]int64
	statusCodes   map[string]map[string]int
	responseTimes map[string][]int64
}

func NewStats() *Stats {
	return &Stats{
		totalElapsed:  make(map[string]int64),
		countPerURL:   make(map[string]int64),
		totalSize:     make(map[string]int64),
		statusCodes:   make(map[string]map[string]int),
		responseTimes: make(map[string][]int64),
	}
}

func (s *Stats) AddCheck(url string, elapsed int64, size int64, status string, isError bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.totalChecks++
	if isError {
		s.errorCount++
		return
	}
	s.totalElapsed[url] += elapsed
	s.countPerURL[url]++
	s.totalSize[url] += size
	s.responseTimes[url] = append(s.responseTimes[url], elapsed)
	if s.statusCodes[url] == nil {
		s.statusCodes[url] = make(map[string]int)
	}
	s.statusCodes[url][status]++
}

func PrintStats(s *Stats) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Println("\n--- Statistics ---")
	fmt.Printf("Total checks: %d\n", s.totalChecks)
	fmt.Printf("Total errors: %d\n", s.errorCount)
	for url, count := range s.countPerURL {
		avgTime := float64(s.totalElapsed[url]) / float64(count)
		avgSize := float64(s.totalSize[url]) / float64(count)
		fmt.Printf("\n%s - Checks: %d, Avg response time: %.2f ms, Avg size: %.2f bytes\n", url, count, avgTime, avgSize)
		fmt.Println("Status codes:")
		for status, scount := range s.statusCodes[url] {
			fmt.Printf("  %s: %d (%.2f%%)\n", status, scount, float64(scount)/float64(count)*100)
		}
		fmt.Println("Response time histogram:")
		times := s.responseTimes[url]
		if len(times) > 0 {
			var buckets [5]int
			for _, t := range times {
				if t < 100 {
					buckets[0]++
				} else if t < 500 {
					buckets[1]++
				} else if t < 1000 {
					buckets[2]++
				} else if t < 5000 {
					buckets[3]++
				} else {
					buckets[4]++
				}
			}
			fmt.Printf("  <100ms: %d | 100-500ms: %d | 500-1000ms: %d | 1000-5000ms: %d | >5000ms: %d\n",
				buckets[0], buckets[1], buckets[2], buckets[3], buckets[4])
		}
	}
	fmt.Println("------------------")
}

func (s *Stats) Export(filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().Format("20060102_150405")
	ext := filepath.Ext(filename)
	base := filename[:len(filename)-len(ext)]
	timedFilename := fmt.Sprintf("%s_%s%s", base, timestamp, ext)

	switch ext {
	case ".csv":
		return s.exportCSV(timedFilename)
	case ".json":
		return s.exportJSON(timedFilename)
	default:
		return fmt.Errorf("unsupported export format: %s", ext)
	}
}

func (s *Stats) exportCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"URL", "Checks", "AvgTime_ms", "AvgSize_bytes", "StatusCode", "Count"})
	for url, count := range s.countPerURL {
		avgTime := float64(s.totalElapsed[url]) / float64(count)
		avgSize := float64(s.totalSize[url]) / float64(count)
		for status, scount := range s.statusCodes[url] {
			writer.Write([]string{url, fmt.Sprintf("%d", count), fmt.Sprintf("%.2f", avgTime), fmt.Sprintf("%.2f", avgSize), status, fmt.Sprintf("%d", scount)})
		}
	}
	return nil
}

func (s *Stats) exportJSON(filename string) error {
	data := map[string]interface{}{
		"total_checks": s.totalChecks,
		"total_errors": s.errorCount,
		"urls":         make(map[string]interface{}),
	}
	for url, count := range s.countPerURL {
		data["urls"].(map[string]interface{})[url] = map[string]interface{}{
			"checks":         count,
			"avg_time_ms":    float64(s.totalElapsed[url]) / float64(count),
			"avg_size_bytes": float64(s.totalSize[url]) / float64(count),
			"status_codes":   s.statusCodes[url],
			"response_times": s.responseTimes[url],
		}
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
