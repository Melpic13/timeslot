.PHONY: all build test test-race test-cover coverage-check bench lint fmt vet tidy-check fuzz clean help quality

COVERAGE_MIN ?= 90.0

all: quality build

build:
	go build ./...

test:
	go test -v ./...

test-race:
	go test -v -race ./...

test-cover:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

coverage-check: test-cover
	@total=$$(go tool cover -func=coverage.out | awk '/^total:/{gsub("%", "", $$3); print $$3}'); \
	echo "Total coverage: $$total% (minimum: $(COVERAGE_MIN)%)"; \
	awk -v total="$$total" -v min="$(COVERAGE_MIN)" 'BEGIN { exit (total + 0 < min + 0) ? 1 : 0 }'

bench:
	go test -bench=. -benchmem ./...

lint:
	golangci-lint run

fmt:
	gofmt -s -w .
	@if command -v goimports >/dev/null 2>&1; then goimports -w .; else echo "goimports not installed; skipping goimports"; fi

vet:
	go vet ./...

tidy-check:
	go mod tidy
	git diff --exit-code -- go.mod go.sum

fuzz:
	go test -fuzz=FuzzParseTimeRange -fuzztime=30s ./availability/...
	go test -fuzz=FuzzParseRule -fuzztime=30s ./recurrence/...
	go test -fuzz=FuzzParseICS -fuzztime=30s ./ical/...

quality: vet test-race coverage-check

clean:
	rm -f coverage.out coverage.html benchmark.txt
	go clean -testcache

help:
	@echo "Available targets:"
	@echo "  build          - Build all packages"
	@echo "  test           - Run unit tests"
	@echo "  test-race      - Run tests with race detector"
	@echo "  test-cover     - Run tests with coverage reports"
	@echo "  coverage-check - Enforce minimum coverage (COVERAGE_MIN)"
	@echo "  bench          - Run benchmarks"
	@echo "  lint           - Run golangci-lint"
	@echo "  fmt            - Format source files"
	@echo "  vet            - Run go vet"
	@echo "  tidy-check     - Verify go mod tidy has no diff"
	@echo "  fuzz           - Run package fuzz tests"
	@echo "  quality        - Vet + race tests + coverage gate"
	@echo "  clean          - Clean build artifacts"
