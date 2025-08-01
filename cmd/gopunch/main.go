package main

import (
	"flag"
	"fmt"
	"github.com/wnyro/gopunch/internal/checker"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	interval := flag.Int("interval", 0, "Interval between requests in seconds (0 = run once)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: gopunch --interval 3 https://example.com [https://another.com ...]")
		os.Exit(1)
	}

	urls := args

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	runCheck := func(urls []string) {
		var wg sync.WaitGroup
		wg.Add(len(urls))

		for _, url := range urls {
			go func(u string) {
				defer wg.Done()
				status, elapsed, err := checker.CheckURL(u)
				if err != nil {
					fmt.Printf("❌ Error for %s: %s\n", u, err)
					return
				}
				fmt.Printf("✅ %s - Status: %s (%d ms)\n", u, status, elapsed)
			}(url)
		}
		wg.Wait()
	}

	if *interval == 0 {
		runCheck(urls)
	} else {
		ticker := time.NewTicker(time.Duration(*interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-sigChan:
				fmt.Println("\nExiting...")
				return
			default:
				runCheck(urls)
				select {
				case <-sigChan:
					fmt.Println("\nExiting...")
					return
				case <-ticker.C:
				}
			}
		}
	}
}
