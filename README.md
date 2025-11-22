# tinysvc 🚀

A lightweight, self-hosted utility service providing IP detection and pastebin functionality. Built with Go and Clean Architecture principles, designed to run efficiently on minimal hardware (like a Samsung N150 with 2GB RAM).

## ✨ Features

- **IP Detection**: Get your public IP address with Cloudflare header support
- **Pastebin**: Create, retrieve, and manage text pastes
  - Markdown rendering support
  - Configurable expiration (30 days default, or never)
  - 1MB size limit per paste
  - Automatic cleanup of expired pastes
- **Rate Limiting**: Built-in protection against abuse
- **RESTful API**: Clean, documented API with OpenAPI 3.0 spec
- **Minimal Resources**: Optimized for low-memory environments
- **Cloudflare Tunnel Ready**: Perfect for exposing services securely

## 🏗️ Architecture

tinysvc follows **Clean Architecture** principles:

```
┌─────────────────────────────────────────┐
│          Delivery (HTTP)                │
├─────────────────────────────────────────┤
│          Use Cases (Business Logic)     │
├─────────────────────────────────────────┤
│          Domain (Entities)              │
├─────────────────────────────────────────┤
│          Repository (Interfaces)        │
├─────────────────────────────────────────┤
│          Infrastructure (SQLite, etc)   │
└─────────────────────────────────────────┘
```

See [ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed architecture documentation.

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- SQLite3

### Installation

```bash
# Clone the repository
git clone https://github.com/aprksy/tinysvc.git
cd tinysvc

# Download dependencies
go mod download

# Build the application
make build

# Run the server
make run
```

The server will start on `http://localhost:8080` by default.

### Using Docker

```bash
# Build Docker image
docker build -t tinysvc:latest .

# Run container
docker run -p 8080:8080 -v $(pwd)/data:/app/data tinysvc:latest
```

## 📖 API Usage

### Get Your IP

```bash
curl http://localhost:8080/api/v1/ip
```

**Response:**
```json
{
  "ip": "203.0.113.42"
}
```

### Create a Paste

```bash
# Simple text paste
curl -X POST http://localhost:8080/api/v1/paste \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello, world!",
    "is_markdown": false,
    "expiry_days": 30
  }'
```

**Response:**
```json
{
  "id": "abc12345",
  "content": "Hello, world!",
  "is_markdown": false,
  "expires_at": "2024-02-15T10:30:00Z",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Create a Markdown Paste

```bash
curl -X POST http://localhost:8080/api/v1/paste \
  -H "Content-Type: application/json" \
  -d '{
    "content": "# Hello\n\nThis is **markdown**",
    "is_markdown": true,
    "expiry_days": 7
  }'
```

### Get a Paste

```bash
curl http://localhost:8080/api/v1/paste/abc12345
```

### Delete a Paste

```bash
curl -X DELETE http://localhost:8080/api/v1/paste/abc12345
```

See [API.md](docs/API.md) for complete API documentation.

## ⚙️ Configuration

Configuration is done via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | Server port |
| `SERVER_HOST` | `0.0.0.0` | Server host |
| `DB_PATH` | `./data/tinysvc.db` | SQLite database path |

Example:

```bash
export SERVER_PORT=3000
export DB_PATH=/var/lib/tinysvc/db.sqlite
./bin/tinysvc
```

## 🛠️ Development

### Prerequisites for Development

```bash
# Install Air for hot-reload
go install github.com/cosmtrek/air@latest

# Install golangci-lint for linting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Development Commands

```bash
# Run with hot-reload
make dev

# Run tests
make test

# Run tests with coverage
make coverage

# Lint code
make lint

# Clean build artifacts
make clean
```

### Project Structure

```
tinysvc/
├── cmd/server/          # Application entry point
├── internal/
│   ├── domain/          # Business entities & logic
│   ├── usecase/         # Application services
│   ├── repository/      # Storage interfaces
│   ├── delivery/http/   # HTTP handlers & routing
│   └── infrastructure/  # External dependencies (DB, config)
├── api/                 # OpenAPI specification
├── docs/                # Documentation
└── migrations/          # Database migrations
```

## 🔒 Security

- **Rate Limiting**: 10 requests/second per IP with burst of 20
- **Request Timeouts**: 60 second timeout for all requests
- **Content Size Limits**: 10MB maximum per paste
- **Automatic Cleanup**: Expired pastes are cleaned every hour

## 🚢 Deployment with Cloudflare Tunnel

1. Install `cloudflared`:
```bash
curl -L https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-386.deb -o cloudflared.deb
sudo dpkg -i cloudflared.deb
```

2. Authenticate:
```bash
cloudflared tunnel login
```

3. Create tunnel:
```bash
cloudflared tunnel create tinysvc
```

4. Configure tunnel (`~/.cloudflared/config.yml`):
```yaml
tunnel: <tunnel-id>
credentials-file: /home/user/.cloudflared/<tunnel-id>.json

ingress:
  - hostname: tinysvc.yourdomain.com
    service: http://localhost:8080
  - service: http_status:404
```

5. Run tunnel:
```bash
cloudflared tunnel run tinysvc
```

## 🎨 Frontend Features

### Landing Page (`/`)
- Service registry with live status
- Clean, responsive design
- Quick navigation to all services

### What's My IP (`/ip.html`)
- Instant IP detection
- Support for Cloudflare headers
- Copy to clipboard
- IP version and type detection

### Pastebin (`/paste.html`)
- Create pastes with markdown support
- Syntax highlighting for code blocks
- Flexible expiration (1 day to never)
- Copy, download, and share functionality
- Raw and rendered view toggle
- Direct paste URLs (e.g., `/paste.html?id=abc123`)

### Markdown Support
Full GitHub Flavored Markdown support including:
- Headers, lists, links, images
- Code blocks with syntax highlighting
- Tables, blockquotes
- Bold, italic, strikethrough
- And more!

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📚 Resources

- [OpenAPI Specification](api/openapi.yaml)
- [Architecture Documentation](docs/ARCHITECTURE.md)
- [API Documentation](docs/API.md)

---

Built with ❤️ for the Model Playground (Cloud Sonnet 4.5) contributors on Outlier