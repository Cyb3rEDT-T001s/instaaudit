#!/bin/bash

echo "========================================"
echo "InstaAudit - Security Auditing Tool"
echo "Cross-Platform Build Script"
echo "========================================"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed!"
    echo ""
    echo "Install Go:"
    echo "  Linux:   sudo apt install golang-go  (or download from golang.org)"
    echo "  macOS:   brew install go  (or download from golang.org)"
    echo "  Manual:  https://golang.org/dl/"
    exit 1
fi

echo "✅ Go found!"
go version
echo ""

# Detect OS
OS=$(uname -s)
ARCH=$(uname -m)

case $OS in
    "Linux")
        BINARY_NAME="instaaudit-linux"
        ;;
    "Darwin")
        BINARY_NAME="instaaudit-macos"
        ;;
    *)
        BINARY_NAME="instaaudit"
        ;;
esac

echo "🔨 Building InstaAudit for $OS ($ARCH)..."
go mod tidy

if go build -o "$BINARY_NAME" cmd/main.go; then
    echo "✅ Build successful!"
    echo ""
    echo "Executable: ./$BINARY_NAME"
    echo ""
    echo "Usage:"
    echo "  ./$BINARY_NAME -H target.com"
    echo "  ./$BINARY_NAME -H target.com -A -f html"
    echo ""
    echo "Make executable (if needed):"
    echo "  chmod +x $BINARY_NAME"
    echo ""
else
    echo "❌ Build failed!"
    exit 1
fi