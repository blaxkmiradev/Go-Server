# 🚀 Go-Server

<p align="center">
  <img src="https://img.shields.io/go/version/go1.21+-blue?style=flat-square" alt="Go Version">
  <img src="https://img.shields.io/github/license/anomalyco/go-server?style=flat-square" alt="License">
  <img src="https://img.shields.io/github/size/anomalyco/go-server/go-server.exe?style=flat-square" alt="Binary Size">
</p>

**All-in-One Local & Public Server** - Host HTML files, proxy to Flask/Python servers, or make your server public with Cloudflare Tunnel.

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 📂 **HTML Hosting** | Serve static HTML files from any folder |
| 🔄 **Proxy Server** | Forward requests to Flask, Django, or any server |
| 🌍 **Cloudflare Tunnel** | Make your local server publicly accessible |
| 🎨 **Cool Design** | Beautiful color terminal UI |
| ⚡ **Fast** | Built with Go for maximum performance |

---

## 📥 Download

**Pre-built exe:** `go-server.exe` (just run it!)

Or build from source:
```bash
go build -o go-server.exe main.go
```

---

## 🚀 Quick Start

### Interactive Mode (Recommended)

Just run without arguments:
```bash
go-server.exe
```

Then follow the menu:
```
┌─ SELECT MODE ───────────────────────────────┐
│ 1. Serve HTML Files                         │
│ 2. Proxy to Server (Flask/Python)           │
│ 3. Make Public (Cloudflare)                │
└────────────────────────────────────────────┘

> Choose [1-3]: 1
> Enter port [8080]: 8080
> Enter HTML folder path: my-website
```

### Command Line Mode

```bash
# Serve HTML files
go-server.exe -path "C:\my-html" -port 8080

# Proxy to Python Flask
go-server.exe -proxy "localhost:5000"

# Make public with Cloudflare
go-server.exe -path "C:\html" -port 8080 -token "your-tunnel-token"
```

---

## 📖 Usage Examples

### 1️⃣ Host HTML Website

```bash
go-server.exe -path "C:\Users\RIkixz\Documents\my-website" -port 3000
```

Access at: http://localhost:3000

### 2️⃣ Proxy to Python Flask

```bash
# Start your Flask app first
python app.py

# In another terminal, start Go-Server proxy
go-server.exe -proxy "localhost:5000"
```

Your Flask app is now accessible via Go-Server!

### 3️⃣ Make Public (Cloudflare Tunnel)

```bash
go-server.exe -path "C:\html" -port 8080 -token "eyJhIjoi..."
```

Your server is now publicly accessible on the internet!

---

## 🔧 Command Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `-path` | `-p` | HTML folder path | (none) |
| `-port` | `-o` | Server port | `8080` |
| `-proxy` | `-x` | Proxy to server | (none) |
| `-token` | `-t` | Cloudflare tunnel token | (none) |
| `-help` | `-h` | Show help | |

---

## ☁️ Cloudflare Setup

### Get Tunnel Token

1. Go to [Cloudflare Zero Dashboard](https://dash.cloudflare.com/)
2. Access → Tunnel
3. Create a new tunnel
4. Copy the token

### Auto Install

When you use `-token`, Go-Server automatically:
- ✅ Detects your OS (Windows/Linux/Mac)
- ✅ Downloads cloudflared
- ✅ Starts the tunnel

No manual installation needed!

---

## 🛠️ Build from Source

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)

### Build

```bash
# Clone or navigate to folder
cd go-server

# Download dependencies
go mod tidy

# Build exe
go build -o go-server.exe main.go
```

### Run Tests

```bash
go run main.go -path test-html -port 8080
```

---

## 📁 Project Structure

```
go-server/
├── main.go           # Main source code
├── go.mod           # Go modules
├── go-server.exe    # Built executable
├── test-html/      # Test HTML files
│   └── index.html
└── README.md       # This file
```

---

## 🎯 Use Cases

| Use Case | Command |
|----------|---------|
| Personal website | `-path "C:\website" -port 80` |
| Flask development | `-proxy "localhost:5000"` |
| API testing | `-proxy "localhost:8000"` |
| Share with friends | `-path "C:\html" -token "..."` |

---

## 📝 License

MIT License - Feel free to use!

---

## 🤝 Support

- ⭐ Star this repo
- 🐛 Report bugs
- 💡 Feature requests

---

**Made with ❤️ using Go**