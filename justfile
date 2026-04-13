# Build the zet binary
build:
    go build -o zet ./cmd/zet

# Install to ~/.local/bin
install: build
    mv zet ~/.local/bin/zet

# Download dependencies
deps:
    go mod download

# Run tests
test:
    go test ./...

# Clean build artifacts
clean:
    rm -f zet
