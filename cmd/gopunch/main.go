package main

import (
	"flag"
	"fmt"
	"github.com/wnyro/gopunch/internal/checker"
	"github.com/wnyro/gopunch/internal/config"
	"github.com/wnyro/gopunch/internal/logger"
	"github.com/wnyro/gopunch/internal/runner"
	"github.com/wnyro/gopunch/internal/stats"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"time"
)

func main() {
	configFile := flag.String("config", "", "Path to JSON config file")
	interval := flag.Int("interval", 0, "Interval between requests in seconds (0 = run once)")
	timeout := flag.Int("timeout", 5, "Request timeout in seconds")
	verbose := flag.Bool("verbose", false, "Verbose output")
	logfile := flag.String("logfile", "", "File to log results (optional)")
	method := flag.String("method", "GET", "HTTP method (GET, POST, HEAD, etc.)")
	data := flag.String("data", "", "Request body for POST/PUT methods")
	concurrency := flag.Int("concurrency", 10, "Maximum concurrent requests")
	export := flag.String("export", "", "Export stats to file (e.g., stats.csv, stats.json)")
	flag.Parse()

	var cfg *config.Config
	urls := []string{}
	urlConfigs := make(map[string]config.URLConfig)
	if *configFile != "" {
		var err error
		cfg, err = config.LoadConfig(*configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
			os.Exit(1)
		}
		urls = cfg.URLs
		for _, uc := range cfg.URLConfigs {
			urlConfigs[uc.URL] = uc
		}
		if *interval == 0 {
			*interval = cfg.Interval
		}
		if *timeout == 5 {
			*timeout = cfg.Timeout
		}
		if *logfile == "" {
			*logfile = cfg.Logfile
		}
		if *verbose == false {
			*verbose = cfg.Verbose
		}
		if *export == "" {
			*export = cfg.Export
		}
		if *concurrency == 10 {
			*concurrency = cfg.Concurrency
		}
	} else {
		for _, arg := range flag.Args() {
			if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
				urls = append(urls, arg)
				urlConfigs[arg] = config.URLConfig{URL: arg, Method: *method, Data: *data}
			}
		}
	}

	if len(urls) < 1 {
		fmt.Println("Usage: gopunch --config config.json [--interval 3] [--timeout 5] [--verbose] [--logfile out.log] [--method GET] [--data '{}'] [--concurrency 10] [--export stats.csv] https://example.com [https://another.com ...]")
		os.Exit(1)
	}

	checker.SetTimeout(time.Duration(*timeout) * time.Second)
	logger := logger.InitLogger(*logfile)
	statistics := stats.NewStats()
	runner.SetConcurrency(*concurrency)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	if *interval == 0 {
		runner.RunCheck(urls, urlConfigs, logger, statistics, *verbose)
		stats.PrintStats(statistics)
		if *export != "" {
			if err := statistics.Export(*export); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to export stats: %v\n", err)
			}
		}
	} else {
		runner.RunWithInterval(urls, urlConfigs, time.Duration(*interval)*time.Second, logger, statistics, *verbose, sigChan)
		go func() {
			<-sigChan
			if *export != "" {
				if err := statistics.Export(*export); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to export stats: %v\n", err)
				}
			}
			fmt.Println("\nExiting...")
			stats.PrintStats(statistics)
			os.Exit(0)
		}()
	}
}
