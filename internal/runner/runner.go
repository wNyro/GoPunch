package runner

import (
	"fmt"
	"github.com/wnyro/gopunch/internal/checker"
	"github.com/wnyro/gopunch/internal/stats"
	"log"
	"os"
	"sync"
	"time"
)

func RunCheck(urls []string, logger *log.Logger, stats *stats.Stats, verbose bool) {
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go func(u string) {
			defer wg.Done()
			status, elapsed, err := checker.CheckURL(u)

			stats.AddCheck(u, elapsed, err != nil)

			if err != nil {
				fmt.Printf("❌ Error for %s: %s\n", u, err)
				if logger != nil {
					logger.Printf("ERROR %s %v", u, err)
				}
				return
			}

			msg := fmt.Sprintf("✅ %s - Status: %s (%d ms)\n", u, status, elapsed)
			if verbose {
				fmt.Print(msg)
			} else {
				fmt.Printf("✅ %s (%d ms)\n", u, elapsed)
			}

			if logger != nil {
				logger.Printf("%s %s %dms", u, status, elapsed)
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
