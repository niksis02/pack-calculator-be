# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependencies separately for faster rebuilds
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /bin/pack-calculator ./cmd/server

# ---- Runtime stage ----
FROM scratch

# TLS certificates (needed for any outbound HTTPS calls)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /bin/pack-calculator /pack-calculator

EXPOSE 3000

ENTRYPOINT ["/pack-calculator"]
