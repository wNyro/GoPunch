package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/wnyro/gopunch/internal/checker"
)

func main() {
	interval := flag.Int("interval", 0, "Interval between requests in seconds (0 = run once)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: gopunch --interval 3 https://example.com")
		os.Exit(1)
	}

	url := args[len(args)-1]

	runCheck := func() {
		status, elapsed, err := checker.CheckURL(url)
		if err != nil {
			fmt.Printf("❌ Error: %s\n", err)
			return
		}
		fmt.Printf("✅ Status: %s (%d ms)\n", status, elapsed)
	}

	if *interval == 0 {
		runCheck()
	} else {
		for {
			runCheck()
			time.Sleep(time.Duration(*interval) * time.Second)
		}
	}
}
