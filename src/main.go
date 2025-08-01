package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Použití: gopunch https://example.com")
		os.Exit(1)
	}

	url := os.Args[1]
	start := time.Now()

	resp, err := http.Get(url)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Chyba: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Printf("✅ Status: %s (%d ms)\n", resp.Status, elapsed.Milliseconds())
}
