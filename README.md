# Pack Calculator — Backend

A Go REST API that determines the optimal combination of pre-defined pack sizes to fulfil any order quantity. Given an order of N items and a set of available pack sizes, the API returns the minimum number of packs (and minimum total items) needed to ship without sending fewer items than ordered.

## Business Logic

Orders cannot be split into individual items — only complete packs can be shipped. The algorithm must:

1. **Never ship fewer items than ordered** — always round up to the nearest achievable quantity.
2. **Minimise total items shipped** (primary goal) — send as close to the order quantity as possible.
3. **Minimise the number of packs** (secondary tie-breaker) — use fewer, larger packs when totals are equal.

A greedy approach cannot guarantee the global optimum, so the API uses a **dynamic programming** (unbounded knapsack variant) solution running in **O(order × K)** time, where K is the number of configured pack sizes.

**Example** — packs: `[250, 500, 1000, 2000, 5000]`, order: `12001`
- Result: `2×5000 + 1×2000 + 1×250 = 12250 items, 4 packs`

> **Note on storage:** For simplicity, no persistent storage solution has been used. Pack configuration is held entirely in-memory and protected by a `sync.RWMutex`, making it safe for concurrent reads and writes. State resets to the default pack sizes `[250, 500, 1000, 2000, 5000]` on every server restart.

## Tech Stack

| Technology | Version |
|---|---|
| Go | 1.26 |
| Fiber | 3 |
| Testify | 1 |
| Docker (runtime) | scratch |
| CI/CD | GitHub Actions |
| Registry | AWS ECR |
| Compute | AWS ECS Fargate |

## Requirements

| Tool | Minimum version |
|---|---|
| Go | 1.26 |
| Docker | any recent version |
| golangci-lint | v2.x (for linting) |
| make | any version |

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `3000` | Port the HTTP server listens on |
| `ALLOW_ORIGINS` | `*` | Comma-separated list of allowed CORS origins |

## Running Without Docker

```bash
# Install / tidy dependencies
go mod tidy

# Run the development server (port 3000 by default)
make run

# Or with custom environment variables
PORT=8080 ALLOW_ORIGINS=http://localhost:5173 go run ./cmd/server
```

Build and run the binary directly:

```bash
make build
./bin/pack-calculator
```

## Running With Docker

```bash
# Build the image using make (defaults to image name 'pack-calculator-backend')
make docker-build

# Build with a custom image name
make docker-build IMAGE=my-image-name

# Or build directly with Docker
docker build -t pack-calculator-backend .

# Run
docker run -p 3000:3000 pack-calculator-backend

# Run with custom environment
docker run -p 3000:3000 \
  -e PORT=3000 \
  -e ALLOW_ORIGINS=http://localhost:5173 \
  pack-calculator-backend
```

The image uses a minimal `scratch` base — no shell, no OS utilities, smallest possible attack surface.

## API Reference

All endpoints are prefixed with `/api/v1`.

### Health check

```
GET /health
```

```json
{ "status": "ok" }
```

### Get pack configuration

```
GET /api/v1/config/packs
```

```json
{ "packs": [250, 500, 1000, 2000, 5000] }
```

### Update pack configuration

```
POST /api/v1/config/packs
Content-Type: application/json

{ "packs": [250, 500, 1000, 2000, 5000] }
```

Rules: at least one pack size, all values must be positive integers. Returns the updated configuration.

### Calculate optimal packing

```
POST /api/v1/calculate
Content-Type: application/json

{ "items": 12001 }
```

```json
{
  "total_items": 12250,
  "packs": [
    { "size": 5000, "count": 2 },
    { "size": 2000, "count": 1 },
    { "size": 250,  "count": 1 }
  ]
}
```

## Development

```bash
# Run tests with race detector
make test

# Run tests with coverage report (generates coverage.html)
make test-cover

# Lint (golangci-lint)
make lint

# Format source files
make fmt

# Tidy dependencies
make tidy
```

## CI/CD & Deployment

The repository has three GitHub Actions workflows:

| Workflow | Trigger | Steps |
|---|---|---|
| `lint.yml` | Pull request | golangci-lint |
| `build.yml` | Pull request | go mod tidy, go build, go test -cover |
| `deploy.yml` | Push to `main` | go mod tidy, build, test, docker build, ECR push, ECS deploy |

On every merge to `main` the image is built, tagged with the commit SHA, pushed to **AWS ECR**, and deployed to **AWS ECS Fargate** with a rolling update. The workflow waits for service stability before completing.

Required repository secrets: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_ACCOUNT_ID`, `AWS_REGION`, `PORT`, `ALLOW_ORIGINS`.

The application is live at **http://pack-calculator-frontend-alb-1042378341.eu-central-1.elb.amazonaws.com**
