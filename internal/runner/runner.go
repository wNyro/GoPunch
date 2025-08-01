package runner

import (
	"context"
	"fmt"
	"github.com/wnyro/gopunch/internal/checker"
	"github.com/wnyro/gopunch/internal/config"
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

func RunCheck(urls []string, urlConfigs map[string]config.URLConfig, logger *log.Logger, stats *stats.Stats, verbose bool) {
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
			uc, exists := urlConfigs[u]
			if !exists {
				uc = config.URLConfig{URL: u, Method: "GET", Data: ""}
			}
			checker.SetMethod(uc.Method, uc.Data)
			status, elapsed, size, headers, err := checker.CheckURL(u)
			stats.AddCheck(u, elapsed, size, status, err != nil)
			if err != nil {
				fmt.Printf("❌ Error for %s: %s\n", u, err)
				if logger != nil {
					logger.Printf("ERROR %s %v (Method: %s)", u, err, uc.Method)
				}
				return
			}
			if verbose {
				msg := fmt.Sprintf("✅ %s - Status: %s, Time: %d ms, Size: %d bytes\nHeaders: %v\n", u, status, elapsed, size, headers)
				fmt.Print(msg)
				if logger != nil {
					logger.Printf("%s %s %dms %dbytes headers:%v (Method: %s)", u, status, elapsed, size, headers, uc.Method)
				}
			} else {
				msg := fmt.Sprintf("✅ %s - Status: %s, Time: %d ms, Size: %d bytes\n", u, status, elapsed, size)
				fmt.Print(msg)
				if logger != nil {
					logger.Printf("%s %s %dms %dbytes (Method: %s)", u, status, elapsed, size, uc.Method)
				}
			}
		}(url)
	}
	wg.Wait()
}

func RunWithInterval(urls []string, urlConfigs map[string]config.URLConfig, interval time.Duration, logger *log.Logger, statistics *stats.Stats, verbose bool, sigChan chan os.Signal) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-sigChan:
			return
		default:
			RunCheck(urls, urlConfigs, logger, statistics, verbose)
			select {
			case <-sigChan:
				return
			case <-ticker.C:
			}
		}
	}
}
