#!/data/data/com.termux/files/usr/bin/bash

echo "InstaAudit Termux Installation Script"
echo "===================================="

# Update Termux packages
echo "Updating Termux packages..."
pkg update -y && pkg upgrade -y

# Install required packages
echo "Installing required packages..."
pkg install -y golang git make

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go installation failed"
    exit 1
fi

echo "Go version: $(go version)"

# Set up Go environment
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc

# Create Go workspace
mkdir -p $GOPATH/src $GOPATH/bin $GOPATH/pkg

# Initialize Go module if not exists
if [ ! -f "go.mod" ]; then
    echo "Initializing Go module..."
    go mod init github.com/Cyb3rEDT-T001s/instaaudit
fi

# Download dependencies
echo "Downloading dependencies..."
go mod tidy

# Build the application
echo "Building InstaAudit for Termux..."
go build -o instaaudit cmd/main.go

# Make executable
chmod +x instaaudit

# Install to Termux bin directory for global access
echo "Installing to system PATH..."
cp instaaudit $PREFIX/bin/
echo "InstaAudit installed to $PREFIX/bin/"

echo ""
echo "Installation complete!"
echo ""
echo "🎯 How to run InstaAudit in Termux:"
echo ""
echo "Method 1 (Recommended - Global command):"
echo "  instaaudit --help"
echo "  instaaudit -H 192.168.1.1 -p common"
echo ""
echo "Method 2 (If global doesn't work):"
echo "  \$(pwd)/instaaudit --help"
echo "  bash -c \"./instaaudit -H target.com -p 80,443\""
echo ""
echo "Method 3 (Full path):"
echo "  $HOME/instaaudit/instaaudit --help"
echo ""
echo "🔧 Troubleshooting:"
echo "- If './instaaudit' doesn't work, use 'instaaudit' directly"
echo "- If permission denied, run: chmod +x instaaudit"
echo "- If command not found, use full path or \$(pwd)/instaaudit"
echo ""
echo "Note: Some network features may be limited on Android/Termux"