package checker

import (
	"bytes"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

var method string
var body []byte

func SetTimeout(t time.Duration) {
	client.Timeout = t
}

func SetMethod(m string, data string) {
	method = strings.ToUpper(m)
	if data != "" {
		body = []byte(data)
	} else {
		body = nil
	}
}

func CheckURL(url string) (status string, elapsed int64, size int64, headers http.Header, err error) {
	start := time.Now()
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return "", time.Since(start).Milliseconds(), 0, nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	elapsed = time.Since(start).Milliseconds()
	if err != nil {
		return "", elapsed, 0, nil, err
	}
	defer resp.Body.Close()
	return resp.Status, elapsed, resp.ContentLength, resp.Header, nil
}
