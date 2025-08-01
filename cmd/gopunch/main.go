package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wnyro/gopunch/internal/checker"
)

func main() {
	interval := flag.Int("interval", 0, "Interval between requests in seconds (0 = run once)")
	timeout := flag.Int("timeout", 5, "Request timeout in seconds")
	verbose := flag.Bool("verbose", false, "Verbose output")
	logfile := flag.String("logfile", "", "File to log results (optional)")
	flag.Parse()

	urls := flag.Args()
	if len(urls) < 1 {
		fmt.Println("Usage: gopunch --interval 3 --timeout 5 --verbose --logfile out.log https://example.com [https://another.com ...]")
		os.Exit(1)
	}
	
	checker.SetTimeout(time.Duration(*timeout) * time.Second)

	logger := InitLogger(*logfile)
	stats := NewStats()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	if *interval == 0 {
		RunCheck(urls, logger, stats, *verbose)
		PrintStats(stats)
	} else {
		RunWithInterval(urls, time.Duration(*interval)*time.Second, logger, stats, *verbose, sigChan)
	}
}
