#!/data/data/com.termux/files/usr/bin/bash

# Termux InstaAudit Runner Script
# Handles different execution methods for Termux

echo "🔍 InstaAudit Termux Runner"
echo "=========================="

# Check if instaaudit binary exists
if [ ! -f "instaaudit" ]; then
    echo "❌ InstaAudit binary not found in current directory"
    echo "Run the installation script first: ./install-termux.sh"
    exit 1
fi

# Make sure it's executable
chmod +x instaaudit

echo "🎯 Available execution methods:"
echo "1. Global command (if installed to PATH)"
echo "2. Current directory execution"
echo "3. Bash wrapper execution"
echo ""

# Try different execution methods
echo "Testing execution methods..."

# Method 1: Global command
if command -v instaaudit &> /dev/null; then
    echo "✅ Method 1: Global command works"
    echo "   Usage: instaaudit -H target.com -p common"
    EXEC_METHOD="instaaudit"
# Method 2: Current directory
elif [ -x "./instaaudit" ]; then
    echo "✅ Method 2: Current directory execution works"
    echo "   Usage: \$(pwd)/instaaudit -H target.com -p common"
    EXEC_METHOD="$(pwd)/instaaudit"
# Method 3: Bash wrapper
else
    echo "✅ Method 3: Using bash wrapper"
    echo "   Usage: bash -c \"./instaaudit -H target.com -p common\""
    EXEC_METHOD="bash -c ./instaaudit"
fi

echo ""

# If arguments provided, run InstaAudit
if [ $# -gt 0 ]; then
    echo "🚀 Running InstaAudit with arguments: $@"
    echo "Using execution method: $EXEC_METHOD"
    echo ""
    
    if [ "$EXEC_METHOD" = "instaaudit" ]; then
        instaaudit "$@"
    elif [ "$EXEC_METHOD" = "$(pwd)/instaaudit" ]; then
        $(pwd)/instaaudit "$@"
    else
        bash -c "./instaaudit $*"
    fi
else
    echo "📋 Termux Use Case Examples:"
    echo ""
    echo "🏠 Home Network Security:"
    echo "  $0 -H 192.168.1.1 -p common                    # Scan router"
    echo "  $0 -H 192.168.1.100 -p 80,443,22,23           # Scan IoT device"
    echo ""
    echo "📱 Mobile Security Testing:"
    echo "  $0 -H \$(ip route | grep default | awk '{print \$3}') -p common  # Check WiFi gateway"
    echo "  $0 -H hotel-wifi-gateway -p 22,23,80,443 -t 1  # Quick hotel WiFi check"
    echo ""
    echo "🎓 Learning & Education:"
    echo "  $0 -H localhost -p common                      # Safe local testing"
    echo "  $0 -H your-website.com -p 80,443 -f html      # Your own site audit"
    echo ""
    echo "🔍 Website Security Check:"
    echo "  $0 -H example.com -p 80,443,8080              # Web services scan"
    echo "  $0 -H target.com --skip-exploits --skip-recon # Quick web check"
    echo ""
    echo "⚡ Quick Mobile Scans:"
    echo "  $0 -H target.com -p 80,443 -t 1               # Fast scan (saves battery)"
    echo "  $0 -H 192.168.1.0/24 -p 22,80,443             # Network discovery"
    echo ""
    echo "📊 IoT Device Security:"
    echo "  $0 -H smart-tv-ip -p 80,443,8080,9000         # Smart TV check"
    echo "  $0 -H camera-ip -p 80,554,8000                # IP camera audit"
    echo ""
    echo "🛠️  Advanced Options:"
    echo "  $0 --help                                      # Show all options"
    echo "  $0 -H target.com -A                           # Aggressive scan"
    echo "  $0 -H target.com -f all -o /sdcard/scan       # All formats to SD card"
    echo ""
    echo "Or run directly using:"
    if [ "$EXEC_METHOD" = "instaaudit" ]; then
        echo "  instaaudit -H target.com -p common"
    elif [ "$EXEC_METHOD" = "$(pwd)/instaaudit" ]; then
        echo "  \$(pwd)/instaaudit -H target.com -p common"
    else
        echo "  bash -c \"./instaaudit -H target.com -p common\""
    fi
    echo ""
    echo "📚 Learn more:"
    echo "  cat TERMUX-GUIDE.md                           # Complete Termux guide"
    echo "  cat TRUST-AND-VERIFICATION.md                 # How to verify results"
fi