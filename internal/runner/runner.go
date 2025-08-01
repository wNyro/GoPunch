package runner

import (
	"context"
	"fmt"
	"github.com/wnyro/gopunch/internal/checker"
	"github.com/wnyro/gopunch/internal/stats"
	"golang.org/x/sync/semaphore"
	"log"
	"os"
	"sync"
	"time"
)

var (
	concurrency int64
	sem         *semaphore.Weighted
	once        sync.Once
)

func SetConcurrency(c int) {
	concurrency = int64(c)
	once.Do(func() {
		sem = semaphore.NewWeighted(concurrency)
	})
}

func RunCheck(urls []string, logger *log.Logger, stats *stats.Stats, verbose bool) {
	var wg sync.WaitGroup
	for _, url := range urls {
		if err := sem.Acquire(context.Background(), 1); err != nil {
			fmt.Fprintf(os.Stderr, "Semaphore acquire error: %v\n", err)
			continue
		}
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			defer sem.Release(1)
			status, elapsed, size, headers, err := checker.CheckURL(u)
			stats.AddCheck(u, elapsed, size, status, err != nil)
			if err != nil {
				fmt.Printf("❌ Error for %s: %s\n", u, err)
				if logger != nil {
					logger.Printf("ERROR %s %v", u, err)
				}
				return
			}
			if verbose {
				msg := fmt.Sprintf("✅ %s - Status: %s, Time: %d ms, Size: %d bytes\nHeaders: %v\n", u, status, elapsed, size, headers)
				fmt.Print(msg)
				if logger != nil {
					logger.Printf("%s %s %dms %dbytes headers:%v", u, status, elapsed, size, headers)
				}
			} else {
				msg := fmt.Sprintf("✅ %s - Status: %s, Time: %d ms, Size: %d bytes\n", u, status, elapsed, size)
				fmt.Print(msg)
				if logger != nil {
					logger.Printf("%s %s %dms %dbytes", u, status, elapsed, size)
				}
			}
		}(url)
	}
	wg.Wait()
}

func RunWithInterval(urls []string, interval time.Duration, logger *log.Logger, statistics *stats.Stats, verbose bool, sigChan chan os.Signal) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-sigChan:
			fmt.Println("\nExiting...")
			stats.PrintStats(statistics)
			os.Exit(0)
		default:
			RunCheck(urls, logger, statistics, verbose)
			select {
			case <-sigChan:
				fmt.Println("\nExiting...")
				stats.PrintStats(statistics)
				os.Exit(0)
			case <-ticker.C:
			}
		}
	}
}
