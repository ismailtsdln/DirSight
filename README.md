# DirSight 🎯

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**DirSight** is a high-performance, next-generation web directory discovery and 403 Forbidden bypass tool. Engineered in Go, it leverages concurrent goroutine pools and sophisticated bypass techniques to uncover hidden resources that traditional scanners miss.

---

## 🚀 Key Features

- **⚡ High-Speed Scanning**: Optimized goroutine pool for extremely fast concurrent requests.
- **🛡️ 403 Bypass Engine**: Automatic application of header-based and path-based bypass techniques.
- **🔍 WAF Awareness**: Built-in fingerprinting to detect Cloudflare, Akamai, AWS WAF, and more.
- **📈 Real-time Progress**: Visual feedback with a persistent progress bar and live status updates.
- **📑 Tabular Results**: Clean, aligned output using `tabwriter` for easy readability.
- **📦 Smart Filtering**: Automatic 404 suppression and response length deduplication.
- **💾 Professional Exporting**: Detailed results exported in structured JSON for integration.

---

## 🛠️ Installation

Ensure you have **Go** installed on your system.

```bash
go install github.com/ismailtsdln/DirSight/cmd/dirsight@latest
```

Alternatively, clone and build locally:

```bash
git clone https://github.com/ismailtsdln/DirSight.git
cd DirSight
go build -o dirsight cmd/dirsight/main.go
```

---

## 📖 Usage

### Basic Scan
```bash
dirsight -u https://example.com -w wordlist.txt
```

### Advanced Scan with Bypass and Export
```bash
dirsight -u https://example.com -w common.txt -expand -t 50 -json results.json
```

### CLI Options

| Flag | Description | Default |
|------|-------------|---------|
| `-u` | Target URL (required) | `""` |
| `-w` | Path to wordlist file (required) | `""` |
| `-t` | Number of concurrent threads | `10` |
| `-timeout` | Request timeout duration | `10s` |
| `-proxy` | Proxy URL (e.g., http://127.0.0.1:8080) | `""` |
| `-k` | Allow insecure SSL connections | `false` |
| `-expand` | Enable path-based bypass variations | `false` |
| `-json` | Path to export results in JSON format | `""` |

---

## 🏗️ Architecture

DirSight is designed with modularity in mind:

- **Engine**: Handles HTTP logic, retry mechanisms, and goroutine pooling.
- **Bypass Logic**: Modular definitions for header and path manipulation.
- **Wordlist Manager**: Dynamic loading and expansion of scan targets.
- **UX/CLI**: Clean separation between logic and terminal presentation.

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📜 License

Distributed under the MIT License. See `LICENSE` for more information.

## ⚠️ Disclaimer

This tool is for educational and authorized penetration testing purposes only. The developer is not responsible for any misuse or damage caused by this utility.
