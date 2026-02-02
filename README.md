# 🧙 ProxyWiz


A high-performance, concurrent proxy checker and scraper written in Go. Automatically scrape proxies from 40+ public sources or use your own list, then validate them with detailed geolocation data.

![Animation](https://github.com/user-attachments/assets/b25b8917-b1ee-435d-8736-a7ff1ce27a24)

## ✨ Features

- 🔥 **Blazing Fast** - Concurrent worker pool for maximum throughput
- 🌐 **Auto-Scraping** - Scrapes from 40+ public proxy sources automatically
- 📁 **Custom Lists** - Use your own proxy files
- 🗺️ **Geolocation Data** - Get detailed IP info (country, city, ISP, timezone)
- 📊 **JSON Export** - Full detailed JSON output for working proxies
- 🎨 **Beautiful CLI** - Colored output with progress bar
- 💾 **Auto-Save** - Saves progress on Ctrl+C interrupt
- 🔄 **Deduplication** - Automatically removes duplicate proxies
- ⚙️ **Configurable** - Adjust workers, timeout, and verbosity

## 📦 Installation

### Prerequisites

- Go 1.21 or higher

### Build from Source

```bash
git clone https://github.com/Syn2Much/Proxy-Wiz.git
cd Proxy-Wiz
go build -o proxywiz .
```

### Run Directly

```bash
go run .
```

## 🚀 Usage

```bash
./proxywiz
```

You'll be prompted for:

| Option | Description | Default |
|--------|-------------|---------|
| **Proxy file** | Path to your proxy list (leave empty to scrape) | Auto-scrape |
| **Workers** | Number of concurrent workers | 10 |
| **Timeout** | Connection timeout in seconds | 8 |
| **Verbose** | Show detailed logging | No |

## 📄 Output Files

### `checked.txt`

Simple list of working proxies:

```
192.168.1.1:8080
proxy.example.com:3128
```

### `working_proxies.json`

Detailed JSON with geolocation data:

```json
[
  {
    "proxy": "192.168.1.1:8080",
    "tested_at": "2026-02-02T14:32:16Z",
    "ip": "104.12.255.182",
    "type": "IPv4",
    "country": "United States",
    "country_code": "US",
    "region": "North Carolina",
    "city": "Charlotte",
    "latitude": 35.2270869,
    "longitude": -80.8431267,
    "is_eu": false,
    "connection": {
      "asn": 7018,
      "org": "AT&T Enterprises, LLC",
      "isp": "AT&T Enterprises, LLC",
      "domain": "att.com"
    },
    "timezone": {
      "id": "America/New_York",
      "abbr": "EST",
      "utc": "-05:00"
    }
  }
]
```

## 🔧 Proxy File Format

Your proxy file should contain one proxy per line:

```
192.168.1.1:8080
http://proxy.example.com:3128
http://user:password@proxy.example.com:8080
https://secure-proxy.example.com:443
```

Supported formats:

- `ip:port`
- `http://ip:port`
- `https://ip:port`
- `http://user:pass@ip:port`

## ⚡ Performance Tips

| Workers | Use Case |
|---------|----------|
| 10-50 | Conservative, low bandwidth |
| 100-200 | Standard usage |
| 500+ | High-speed connections |

**Note:** More workers = faster checking but more memory/bandwidth usage.


## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This tool is for educational and legitimate purposes only. Always ensure you have permission to use proxies and comply with applicable laws and terms of service.

## 👤 Author

**@Syn2Much**

---

⭐ Star this repo if you find it useful!
