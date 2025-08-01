package main

import (
	"flag"
	"fmt"
	"github.com/wnyro/gopunch/internal/logger"
	"github.com/wnyro/gopunch/internal/runner"
	"github.com/wnyro/gopunch/internal/stats"
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

	logger := logger.InitLogger(*logfile)
	statistics := stats.NewStats()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	if *interval == 0 {
		runner.RunCheck(urls, logger, statistics, *verbose)
		stats.PrintStats(statistics)
	} else {
		runner.RunWithInterval(urls, time.Duration(*interval)*time.Second, logger, statistics, *verbose, sigChan)
	}
}
