# GoPunch

**Dead Simple Uptime Monitoring Tool**  
GoPunch is a lightweight CLI utility written in Go for checking the uptime and response times of websites. It's great for quick one-off checks or recurring interval-based monitoring with logging and statistics.

---

## ✨ Features

- ✅ Supports checking multiple URLs
- ⏱️ Customizable request interval and timeout
- 📊 Collects statistics (average response times, error count)
- 📝 Optional logging to a file
- 🧠 Thread-safe concurrency (each URL is checked in parallel)
- 📦 Simple and dependency-free

---

## 🔧 Installation

1. Clone the repo:

```bash
git clone https://github.com/wnyro/GoPunch.git
cd gopunch
````

2. Run it:

```bash
go run . https://example.com
```

3. (Optional) Build it:

```bash
go build -o gopunch ./cmd/gopunch
```

---

## 🚀 Usage

```bash
gopunch [flags] <url1> [url2 ...]
```

### Example:

```bash
gopunch --interval 5 --timeout 3 --verbose --logfile out.log https://google.com https://github.com
```

---

## 📘 Flags

| Flag         | Description                                   | Default |
| ------------ | --------------------------------------------- | ------- |
| `--interval` | Interval between checks in seconds (0 = once) | `0`     |
| `--timeout`  | Timeout for each request in seconds           | `5`     |
| `--verbose`  | Enable verbose output (shows HTTP status)     | `false` |
| `--logfile`  | Write output to a file                        | `""`    |

---

## 📤 Output Example

```
✅ https://google.com - Status: 200 OK (94 ms)
✅ https://github.com - Status: 200 OK (122 ms)

--- Statistics ---
Total checks: 2
Total errors: 0
https://google.com - Checks: 1, Average response time: 94.00 ms
https://github.com - Checks: 1, Average response time: 122.00 ms
------------------
```

---

## 🛑 Graceful Shutdown

If running with `--interval`, simply press `Ctrl+C` to stop. GoPunch will print a final summary with statistics.

---

## 🧪 Example Commands

| Purpose             | Command Example                                                              |
| ------------------- | ---------------------------------------------------------------------------- |
| Single check        | `gopunch https://example.com`                                                |
| Interval check (5s) | `gopunch --interval 5 https://example.com https://github.com`                |
| Custom timeout (2s) | `gopunch --timeout 2 https://httpstat.us/200?sleep=3000`                     |
| Save to file        | `gopunch --logfile output.log https://example.com`                           |
| Verbose + log       | `gopunch --verbose --logfile out.log https://example.com https://google.com` |

---

## 📄 License

MIT License © 2025 \ Made by wNyro 🐼