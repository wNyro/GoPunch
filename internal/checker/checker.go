package checker

import (
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func SetTimeout(t time.Duration) {
	client.Timeout = t
}

func CheckURL(url string) (string, int64, error) {
	start := time.Now()
	resp, err := client.Get(url)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		return "", elapsed, err
	}
	defer resp.Body.Close()

	return resp.Status, elapsed, nil
}
