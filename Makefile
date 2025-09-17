# InstaAudit - Cross-Platform Makefile

.PHONY: all build clean install test help linux windows macos

# Default target
all: build

# Build for current platform
build:
	@echo "Building InstaAudit for current platform..."
	go mod tidy
	go build -o instaaudit cmd/main.go
	@echo "✅ Build complete: ./instaaudit"

# Cross-platform builds
linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o instaaudit-linux cmd/main.go
	@echo "✅ Linux build complete: ./instaaudit-linux"

windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o instaaudit.exe cmd/main.go
	@echo "✅ Windows build complete: ./instaaudit.exe"

macos:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o instaaudit-macos cmd/main.go
	@echo "✅ macOS build complete: ./instaaudit-macos"

# Build for all platforms
all-platforms: linux windows macos
	@echo "✅ All platform builds complete!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f instaaudit instaaudit-* *.exe
	rm -f audit_report.*
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Install to system (Linux/macOS)
install: build
	@echo "Installing InstaAudit to /usr/local/bin..."
	sudo cp instaaudit /usr/local/bin/
	sudo chmod +x /usr/local/bin/instaaudit
	@echo "✅ InstaAudit installed system-wide"

# Quick test scan
demo:
	@echo "Running demo scan..."
	./instaaudit -H google.com -p "80,443"

# Help
help:
	@echo "InstaAudit Build Commands:"
	@echo ""
	@echo "  make build         - Build for current platform"
	@echo "  make linux         - Build for Linux"
	@echo "  make windows       - Build for Windows" 
	@echo "  make macos         - Build for macOS"
	@echo "  make all-platforms - Build for all platforms"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make install       - Install system-wide (Linux/macOS)"
	@echo "  make test          - Run tests"
	@echo "  make demo          - Run demo scan"
	@echo "  make help          - Show this help"