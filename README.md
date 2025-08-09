# GoPunch

**Dead Simple Uptime Monitoring Tool**  
GoPunch is a lightweight CLI utility written in Go for checking the uptime and response times of websites. It supports one-off checks or recurring interval-based monitoring with detailed logging, statistics, and export to CSV/JSON.

---

## ✨ Features

- ✅ Supports checking multiple URLs with customizable HTTP methods (GET, POST, HEAD, etc.)
- ⏱️ Configurable request interval, timeout, and concurrency
- 📊 Detailed statistics (average response times, error counts, status codes, response time histograms)
- 📝 Optional logging to a file or stdout
- 📦 Export statistics to CSV or JSON
- 🧠 Thread-safe parallel requests with configurable concurrency
- ⚙️ JSON configuration file for streamlined setup
- 🚀 Dependency-light and easy to use

---

## 🔧 Installation

1. Clone the repository:

```bash
git clone https://github.com/wnyro/GoPunch.git
cd GoPunch
```

2. Install dependency:

```bash
go get golang.org/x/sync
```

3. Run it directly:

```bash
go run . https://postman-echo.com/get
```

4. (Optional) Build the binary:

```bash
go build -o gopunch ./cmd/gopunch
```

---

## 🚀 Usage

```bash
gopunch [flags] <url1> [url2 ...]
```

Run with a JSON config file:

```bash
gopunch --config=config.json
```

### Example Config (`config.json`)

```json
{
  "urls": ["https://postman-echo.com/get", "https://postman-echo.com/post"],
  "interval": 3,
  "timeout": 15,
  "logfile": "results.log",
  "verbose": true,
  "export": "stats.csv",
  "concurrency": 2,
  "url_configs": [
    {
      "url": "https://postman-echo.com/get",
      "method": "GET",
      "data": ""
    },
    {
      "url": "https://postman-echo.com/post",
      "method": "POST",
      "data": "{\"test\":\"value\"}"
    }
  ]
}
```

---

## 📘 Flags

| Flag         | Description                                          | Default   |
|--------------|-----------------------------------------------------|-----------|
| `--config`   | Path to JSON config file                            | `""`      |
| `--interval` | Interval between checks in seconds (0 = once)       | `0`       |
| `--timeout`  | Timeout for each request in seconds                 | `5`       |
| `--verbose`  | Enable verbose output (shows HTTP headers)          | `false`   |
| `--logfile`  | Write output to a file (stdout if empty)            | `""`      |
| `--method`   | HTTP method (GET, POST, HEAD, etc.)                 | `"GET"`   |
| `--data`     | Request body for POST/PUT methods                   | `""`      |
| `--concurrency` | Maximum concurrent requests                      | `10`      |
| `--export`   | Export stats to file (e.g., `stats.csv`, `stats.json`) | `""`    |

**Note**: Flags override config values if specified.

---

## 📤 Output Example

### Console Output

```
✅ https://postman-echo.com/post - Status: 200 OK, Time: 500 ms, Size: 446 bytes
Headers: map[Content-Type:[application/json] ...]
✅ https://postman-echo.com/get - Status: 200 OK, Time: 300 ms, Size: 200 bytes
Headers: map[Content-Type:[application/json] ...]

--- Statistics ---
Total checks: 2
Total errors: 0
https://postman-echo.com/post - Checks: 1, Avg response time: 500.00 ms, Avg size: 446.00 bytes
Status codes:
  200 OK: 1 (100.00%)
Response time histogram:
  <100ms: 0 | 100-500ms: 1 | 500-1000ms: 0 | 1000-5000ms: 0 | >5000ms: 0
https://postman-echo.com/get - Checks: 1, Avg response time: 300.00 ms, Avg size: 200.00 bytes
Status codes:
  200 OK: 1 (100.00%)
Response time histogram:
  <100ms: 0 | 100-500ms: 1 | 500-1000ms: 0 | 1000-5000ms: 0 | >5000ms: 0
------------------
```

### Log File (`results.log`)

```
2025/08/01 14:51:16 https://postman-echo.com/post 200 OK 500ms 446bytes headers:map[...] (Method: POST)
2025/08/01 14:51:16 https://postman-echo.com/get 200 OK 300ms 200bytes headers:map[...] (Method: GET)
```

### Export File (`stats.csv`)

```
URL,Checks,AvgTime_ms,AvgSize_bytes,StatusCode,Count
https://postman-echo.com/get,1,300.00,200.00,200 OK,1
https://postman-echo.com/post,1,500.00,446.00,200 OK,1
```

### Export File (`stats.json`)

```json
{
  "total_checks": 2,
  "total_errors": 0,
  "urls": {
    "https://postman-echo.com/get": {
      "checks": 1,
      "avg_time_ms": 300,
      "avg_size_bytes": 200,
      "status_codes": { "200 OK": 1 },
      "response_times": [300]
    },
    "https://postman-echo.com/post": {
      "checks": 1,
      "avg_time_ms": 500,
      "avg_size_bytes": 446,
      "status_codes": { "200 OK": 1 },
      "response_times": [500]
    }
  }
}
```

---

## 🛑 Shutdown

For interval-based checks (`--interval`), press `Ctrl+C` to stop. GoPunch will print final statistics and export to the specified file (if set).

---

## 🧪 Example Commands

| Purpose                        | Command Example                                                                 |
|--------------------------------|---------------------------------------------------------------------------------|
| Single check with config       | `gopunch --config=config.json`                                                  |
| Config with JSON export        | `gopunch --config=config.json --export=stats.json`                              |
| Config with custom interval     | `gopunch --config=config.json --interval=5`                                     |
| Single GET check               | `gopunch --timeout=15 --method=GET --concurrency=5 --export=stats.csv https://postman-echo.com/get` |
| Single POST check with verbose | `gopunch --timeout=15 --method=POST --data='{"test":"value"}' --concurrency=2 --verbose --logfile=results.log --export=stats.json https://postman-echo.com/post` |
| Periodic GET check             | `gopunch --interval=3 --timeout=15 --method=GET --concurrency=2 --verbose --logfile=results.log --export=stats.csv https://postman-echo.com/get https://postman-echo.com/post` |
| HEAD check                    | `gopunch --timeout=15 --method=HEAD --concurrency=5 --logfile=results.log --export=stats.csv https://postman-echo.com/get` |
| Minimal check                 | `gopunch https://postman-echo.com/get`                                          |
| Low concurrency with verbose   | `gopunch --timeout=15 --method=GET --concurrency=1 --verbose --logfile=results.log --export=stats.json https://postman-echo.com/get` |
| Periodic check without log     | `gopunch --interval=3 --timeout=15 --method=GET --concurrency=2 --verbose --export=stats.csv https://postman-echo.com/get` |

---

## 📝 Notes

- Use `https://postman-echo.com/get` and `https://postman-echo.com/post` for reliable testing, as `httpbin.org` may return `503 Service Unavailable`.
- Ensure `config.json` is in the working directory when using `--config`.
- Increase `--timeout` (e.g., to 20s) if you encounter `context deadline exceeded` errors.
- Logs default to stdout if `--logfile` is not specified.

---

## 📄 License

wNyro Restricted Software License © 2025 | Made by wNyro 🐼


---
