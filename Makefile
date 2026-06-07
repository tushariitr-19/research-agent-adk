.PHONY: build run cli test test-integration clean dev

# Build the binary
build:
	go build -o research-agent ./cmd/

# Run with web UI
run: build
	GOOGLE_API_KEY=$(GOOGLE_API_KEY) ./research-agent web api webui

# Run with CLI only
cli: build
	GOOGLE_API_KEY=$(GOOGLE_API_KEY) ./research-agent

# Unit tests only (no network)
test:
	go test ./tests/... -v

# Integration tests (requires network)
test-integration:
	go test ./tests/... -v -tags integration

# All tests
test-all: test test-integration

# Clean built binary
clean:
	rm -f research-agent

# Run without building (development)
dev:
	GOOGLE_API_KEY=$(GOOGLE_API_KEY) go run ./cmd/ web api webui