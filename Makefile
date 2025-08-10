.PHONY: all build test test-race test-cover bench lint fmt vet fuzz clean help

all: lint test build

build:
	go build ./...

test:
	go test -v ./...

test-race:
	go test -v -race ./...

test-cover:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

bench:
	go test -bench=. -benchmem ./...

lint:
	golangci-lint run

fmt:
	gofmt -s -w .
	goimports -w .

vet:
	go vet ./...

fuzz:
	go test -fuzz=FuzzParseRule -fuzztime=1m ./recurrence/...
	go test -fuzz=FuzzParseICS -fuzztime=1m ./ical/...

clean:
	rm -f coverage.out coverage.html benchmark.txt
	go clean -testcache

help:
	@echo "Available targets:"
	@echo "  build      - Build the library"
	@echo "  test       - Run tests"
	@echo "  test-race  - Run tests with race detector"
	@echo "  test-cover - Run tests with coverage"
	@echo "  bench      - Run benchmarks"
	@echo "  lint       - Run linters"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  fuzz       - Run fuzz tests"
	@echo "  clean      - Clean build artifacts"
