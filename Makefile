.PHONY: all build run test test-cover fmt lint tidy docker-build

BINARY_NAME := pack-calculator
MAIN_PATH   := ./cmd/server
IMAGE       ?= $(BINARY_NAME)-backend

all: lint test build

build:
	go build -ldflags="-w -s" -o bin/$(BINARY_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)

test:
	go test -v -race ./...

test-cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

fmt:
	gofmt -w .

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

docker-build:
	docker build -t $(IMAGE) .
