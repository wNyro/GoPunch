package checker

import (
	"net/http"
	"time"
)

func CheckURL(url string) (string, int64, error) {
	start := time.Now()

	resp, err := http.Get(url)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		return "", elapsed, err
	}
	defer resp.Body.Close()

	return resp.Status, elapsed, nil
}
