# 📤 File Uploader

Web application for uploading files to configured folders via a simple web interface.

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Available-2496ED?style=flat&logo=docker)](https://hub.docker.com/r/phd59fr/fileuploader)

## 📝 Description

Lightweight web server that provides a clean drag-and-drop interface to upload files directly into pre-configured directories. Each directory is defined by a label and a path in a YAML configuration file.

**Key Features:**
- 📁 Upload files to configurable folders via web UI
- 🎯 Drag & drop or file selection
- 🚫 File type restriction (`.torrent` only)
- 📏 Configurable file size limit (Go-style byte format)
- 🧠 In-memory upload history (since app start)
- 🐳 Docker image ~11 MB

## ⚙️ Configuration

Edit the `config.yaml` file at the root of the project:

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  max_file_size: "1MiB"

folders:
  - label: "PH"
    path: "/storage/ph"
  - label: "Jean"
    path: "/storage/jean"
  - label: "HTTP"
    path: "/storage/http"
```

**Size format:** Supports `B`, `KB`, `MB`, `GB`, `TB` (decimal) and `KiB`, `MiB`, `GiB`, `TiB` (binary), as well as short forms `K`, `M`, `G`, `T` (binary).

## 🚀 Usage

### 🐳 Docker

```bash
# Pull image
docker pull phd59fr/fileuploader:latest

# Run with mounted storage directories
docker run -d \
  -p 8080:8080 \
  -v /storage/ph:/storage/ph \
  -v /storage/jean:/storage/jean \
  -v /storage/http:/storage/http \
  -v $(pwd)/config.yaml:/app/config.yaml \
  phd59fr/fileuploader:latest

# Build locally
docker build -t fileuploader .
```

### Local Development

```bash
# Install dependencies
go mod tidy

# Run
go run .

# Build
go build -o fileuploader .
```

### Environment Variables

| Variable     | Description                |
|--------------|----------------------------|
| CONFIG_PATH  | Path to config.yaml (optional, defaults to `config.yaml`) |

## 📦 Project Structure

```
.
├── main.go                  # Application entry point & HTTP routes
├── config.go                # YAML config loader & ByteSize parser
├── handlers.go              # HTTP handlers (upload, config, history)
├── uploads.go               # In-memory upload history (thread-safe)
├── templates/
│   └── index.html           # Web UI (drag & drop, folder selection)
├── config.yaml              # Folder and server configuration
├── Dockerfile               # Multi-stage Docker build
├── go.mod / go.sum          # Dependencies
└── LICENSE
```

## 🔧 How It Works

1. **Startup**: Reads `config.yaml` and registers configured folders
2. **Web UI**: Lists folders as clickable cards with colored icons
3. **Upload**: User drags or selects a file → validated against config
4. **Validation**: File extension check (`.torrent` only) + size limit
5. **Storage**: File written to the matching folder's path on disk
6. **History**: File metadata (name, folder, size, date) kept in memory
7. **Security**: Path traversal prevention, `path.Clean`, absolute path verification

## 📊 API Endpoints

| Route                 | Method | Description                          |
|------------------------|--------|--------------------------------------|
| `/`                    | GET    | Web UI                               |
| `/api/config`          | GET    | List configured folders (label only) |
| `/api/uploads`         | GET    | Upload history (since app start)     |
| `/api/upload/{label}`  | POST   | Upload a file to the specified folder|

## 📦 Dependencies

- **[YAML.v3](https://gopkg.in/yaml.v3)** - YAML parsing library
- Standard library only (net/http, encoding/json, html/template, etc.)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🍰 Contributing
Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

## ❤️ Support
A simple star to this project repo is enough to keep me motivated on this project for days. If you find your self very much excited with this project let me know with a tweet.

If you have any questions, feel free to reach out to me on [X](https://twitter.com/xxPHDxx).