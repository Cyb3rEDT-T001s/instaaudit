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

# Create symlink for global access
mkdir -p $HOME/.local/bin
cp instaaudit $HOME/.local/bin/
echo 'export PATH=$PATH:$HOME/.local/bin' >> ~/.bashrc

echo ""
echo "Installation complete!"
echo ""
echo "To run InstaAudit:"
echo "  ./instaaudit --help"
echo ""
echo "Or from anywhere:"
echo "  source ~/.bashrc"
echo "  instaaudit --help"
echo ""
echo "Example usage:"
echo "  ./instaaudit -H 192.168.1.1 -p common"
echo ""
echo "Note: Some features may be limited on Android/Termux"