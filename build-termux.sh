#!/data/data/com.termux/files/usr/bin/bash

echo "Building InstaAudit for Termux..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Run: pkg install golang"
    exit 1
fi

# Set Go environment
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Build for Android ARM
echo "Building for Android ARM64..."
GOOS=android GOARCH=arm64 go build -o instaaudit-android cmd/main.go

# Build regular version
echo "Building regular version..."
go build -o instaaudit cmd/main.go

# Make executable
chmod +x instaaudit*

echo "Build complete!"
echo "Run with: ./instaaudit --help"