# Gofi

**Gofi** is a lightweight, high-performance RSS/Atom-to-JSON proxy service written in Go. It fetches any public feed URL, normalizes it into a simple JSON schema, and serves it via an HTTP API—with built-in Redis caching, per-token access control, and rate-limiting.

---

## ✨ Key Features

- **🔗 RSS/Atom Support**  
  Parses RSS 2.0, Atom, and RDF feeds using [gofeed](https://github.com/mmcdole/gofeed).

- **⚡ Fast JSON Output**  
  Returns a structured JSON payload with feed metadata and items:
  ```jsonc
  {
    "status": "ok",
    "feed": {
      "title": "Example Blog",
      "link": "https://example.com",
      "description": "An example feed",
      "image": "https://example.com/logo.png",
      "updated": "2025-05-13T10:00:00Z"
    },
    "items": [
      {
        "title": "Post title",
        "link": "https://example.com/post",
        "author": "Author Name",
        "published": "2025-05-12T08:30:00Z",
        "content": "...",
        "description": "...",
        "categories": ["tag1","tag2"],
        "enclosure": {
          "url": "...", "type": "audio/mpeg", "length": 12345
        }
      },
      …
    ]
  }
  ```

- **🔒 Token-Based Auth**  
  Admins can issue tokens for clients. Every request to `/v1/feed` must present a valid `Bearer` token.

- **⏱️ Rate-Limiting**  
  Per-token rate limits (e.g. “1000 requests per hour”) enforced via [ulule/limiter](https://github.com/ulule/limiter) + Redis.

- **📦 Redis Caching**  
  Feed responses are cached in Redis (15 min TTL by default) to reduce upstream calls.

- **🐳 Docker-Ready**  
  Single Dockerfile + `docker-compose.yml` spins up Gofi + Redis in seconds.

---

## 📂 Project Structure

```
feedjson/
├── cmd/feedjson/
│   └── main.go            — application entrypoint
├── internal/
│   ├── api/
│   │   └── router.go      — HTTP routes & handlers
│   ├── auth/
│   │   ├── middleware.go  — bearer-token check
│   │   └── store.go       — Redis token store
│   ├── cache/
│   │   └── cache.go       — Redis client & helpers
│   ├── config/
│   │   └── config.go      — env-var loader
│   ├── parser/
│   │   └── parser.go      — feed fetching & parsing
│   ├── ratelimit/
│   │   └── middleware.go  — rate-limit setup
│   └── model/
│       └── feed.go        — JSON response schemas
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── README.md
└── .gitignore
```

---

## 🚀 Getting Started

### Prerequisites

- [Go 1.24+](https://golang.org/dl/)  
- [Docker & Docker Compose](https://docs.docker.com/get-started/)  
- (Optional) Local Redis instance if not using Docker

### 1. Clone & Build

```bash
git clone https://github.com/zeroCoder1/feedjson.git
cd feedjson
go mod tidy
```

### 2. Configure Environment

Create a `.env` (or set in your shell / compose) with:

```dotenv
# Redis connection
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

# Rate limit (e.g. 100 requests per hour)
RATE_LIMIT=100-H

# Admin secret for token issuance
ADMIN_SECRET=your-admin-secret

# HTTP bind port
PORT=8080
```

### 3. Launch with Docker Compose

```bash
docker-compose up --build -d
```

- **Redis** on port 6379  
- **Gofi** on port 8080

---

## 🛠 API Usage

### 1. Issue a Client Token (Admin only)

```bash
curl -i -X POST http://localhost:8080/v1/tokens      -H "X-Admin-Token: your-admin-secret"
```

**Response** (HTTP 201):

```json
{"token":"3fa85f64-5717-4562-b3fc-2c963f66afa6"}
```

### 2. Fetch a Feed

```bash
curl -i   -H "Authorization: Bearer 3fa85f64-5717-4562-b3fc-2c963f66afa6"   "http://localhost:8080/v1/feed?rss_url=https://blog.golang.org/feed.atom&count=5"
```

- **200 OK** with JSON payload.  
- **401 Unauthorized** if token is missing/invalid.  
- **429 Too Many Requests** if rate limit exceeded.

---

## ⚙️ Configuration

| Env Var         | Default   | Description                             |
|-----------------|-----------|-----------------------------------------|
| `REDIS_ADDR`    | `localhost:6379` | Redis server address           |
| `REDIS_PASSWORD`| ―         | Redis password (if any)                 |
| `REDIS_DB`      | `0`       | Redis database index                    |
| `RATE_LIMIT`    | `1000-H`  | Rate (e.g. `100-H`, `5000-D`, etc.)     |
| `ADMIN_SECRET`  | ―         | Secret to protect token-issuance API    |
| `PORT`          | `8080`    | HTTP server port                        |

---

## 🧪 Testing

- **Unit tests**:  
  ```bash
  go test ./internal/parser
  go test ./internal/cache
  go test ./internal/auth
  ```

- **Integration**:  
  Run against local Redis and verify endpoints with `curl`.

---

## 🤝 Contributing

1. Fork the repo & create a feature branch  
2. `git checkout -b feature/your-feature`  
3. Write code + tests  
4. `go fmt` & `go vet`  
5. Open a PR & reference related issue  

Please follow the [Go project layout](https://github.com/golang-standards/project-layout) and write clear, concise commit messages.

---

## 📜 License

This project is licensed under the **MIT License**.  
See [LICENSE](LICENSE) for details.
