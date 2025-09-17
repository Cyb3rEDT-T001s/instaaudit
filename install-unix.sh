#!/bin/bash

echo "========================================"
echo "InstaAudit - Installation Script"
echo "Linux & macOS Support"
echo "========================================"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed!"
    echo ""
    echo "Installing Go..."
    
    # Detect OS
    OS=$(uname -s)
    case $OS in
        "Linux")
            echo "Detected Linux - Installing Go..."
            if command -v apt &> /dev/null; then
                sudo apt update
                sudo apt install -y golang-go
            elif command -v yum &> /dev/null; then
                sudo yum install -y golang
            elif command -v pacman &> /dev/null; then
                sudo pacman -S go
            else
                echo "Please install Go manually from: https://golang.org/dl/"
                exit 1
            fi
            ;;
        "Darwin")
            echo "Detected macOS - Installing Go..."
            if command -v brew &> /dev/null; then
                brew install go
            else
                echo "Please install Homebrew first: https://brew.sh/"
                echo "Or install Go manually from: https://golang.org/dl/"
                exit 1
            fi
            ;;
        *)
            echo "Unsupported OS. Please install Go manually from: https://golang.org/dl/"
            exit 1
            ;;
    esac
fi

echo "✅ Go is installed!"
go version
echo ""

# Build InstaAudit
echo "🔨 Building InstaAudit..."
go mod tidy

if go build -o instaaudit cmd/main.go; then
    echo "✅ Build successful!"
    
    # Make executable
    chmod +x instaaudit
    
    # Optionally install system-wide
    echo ""
    read -p "Install InstaAudit system-wide? (y/n): " -n 1 -r
    echo ""
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Installing to /usr/local/bin..."
        sudo cp instaaudit /usr/local/bin/
        sudo chmod +x /usr/local/bin/instaaudit
        echo "✅ InstaAudit installed system-wide!"
        echo "You can now run: instaaudit -H target.com"
    else
        echo "✅ InstaAudit ready to use!"
        echo "Run: ./instaaudit -H target.com"
    fi
    
    echo ""
    echo "🎯 Quick Start:"
    echo "  ./instaaudit -H google.com"
    echo "  ./instaaudit -H target.com -A -f html"
    echo ""
    
else
    echo "❌ Build failed!"
    exit 1
fi