# InstaAudit on Termux (Android)

## Quick Setup

1. **Install Termux** from F-Droid (recommended) or Google Play Store

2. **Run the installation script:**
   ```bash
   chmod +x install-termux.sh
   ./install-termux.sh
   ```

3. **Test the installation:**
   ```bash
   ./instaaudit --help
   ```

## Manual Installation

If the script doesn't work, follow these steps:

### Step 1: Update Termux
```bash
pkg update && pkg upgrade
```

### Step 2: Install Dependencies
```bash
pkg install golang git make
```

### Step 3: Set Up Go Environment
```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
```

### Step 4: Build InstaAudit
```bash
go mod tidy
go build -o instaaudit cmd/main.go
chmod +x instaaudit

# Add to PATH for easier access
cp instaaudit $PREFIX/bin/
```

## Termux Use Cases & Examples

### 🏠 **Home Network Security (Most Common)**
```bash
# Scan your home router
instaaudit -H 192.168.1.1 -p common

# Check all devices on your network
instaaudit -H 192.168.1.0/24 -p 22,80,443

# Quick router security check
instaaudit -H 192.168.1.1 -p 80,443,22,23,8080 --skip-exploits
```

### 📱 **Mobile Penetration Testing**
```bash
# Test public WiFi security
instaaudit -H gateway-ip -p common -t 1

# Scan hotel/cafe networks (with permission)
instaaudit -H 10.0.0.1 -p 80,443,22,23 --skip-recon

# Quick device discovery
instaaudit -H 192.168.43.0/24 -p 22,80,443,8080
```

### 🎓 **Learning & Education**
```bash
# Safe learning on your own devices
instaaudit -H localhost -p common

# Test your own website/server
instaaudit -H your-domain.com -p 80,443,22

# Educational scan with explanations
instaaudit -H example.com -p 80,443 -f html
```

### 🔍 **Bug Bounty & Research (Authorized Only)**
```bash
# Initial reconnaissance
instaaudit -H target.com -p common --skip-exploits

# Web application focus
instaaudit -H webapp.com -p 80,443,8080,8443 -A

# Quick service enumeration
instaaudit -H target.com -p 1-1000 -t 1
```

### 🏢 **Corporate/Work Networks (With Permission)**
```bash
# Internal network assessment
instaaudit -H 10.0.0.0/24 -p common

# Server security check
instaaudit -H server.company.com -p 22,80,443,3389

# Database security audit
instaaudit -H db-server -p 3306,5432,1433,27017
```

### 📊 **IoT Device Testing**
```bash
# Smart home devices
instaaudit -H 192.168.1.100 -p 80,443,23,8080,9000

# IP cameras and NVRs
instaaudit -H camera-ip -p 80,554,8000,8080

# Router and access points
instaaudit -H 192.168.1.1 -p 80,443,22,23,8080,8443
```

## 📱 **Termux-Specific Execution Methods**

### Method 1: Direct Command (Recommended)
```bash
# After installation, use globally
instaaudit -H 192.168.1.1 -p common
```

### Method 2: Termux Runner Script
```bash
# Use the helper script
./run-termux.sh -H target.com -p 80,443,22
./run-termux.sh --help
```

### Method 3: Full Path Execution
```bash
# If global command doesn't work
$(pwd)/instaaudit -H example.com -p common
$HOME/instaaudit/instaaudit -H target.com -p 80,443
```

### Method 4: Bash Wrapper
```bash
# For compatibility issues
bash -c "./instaaudit -H target.com -p common"
```

## Termux Limitations

### Network Restrictions
- Some network operations may require root access
- Raw sockets might not work without root
- Some system-level checks are limited

### Workarounds
- Use TCP connect scans instead of SYN scans
- Focus on application-level testing
- Use external tools when needed

### Performance Tips
- Use smaller port ranges for faster scans
- Skip heavy operations with `--skip-exploits`
- Use shorter timeouts: `-t 1`

## Troubleshooting

### Command Execution Issues

**Problem: `./instaaudit` doesn't work**
```bash
# Solution 1: Use direct command (if installed globally)
instaaudit --help

# Solution 2: Use full path
$(pwd)/instaaudit --help

# Solution 3: Use bash
bash -c "./instaaudit --help"

# Solution 4: Check if it's in PATH
which instaaudit
echo $PATH
```

**Problem: Permission denied**
```bash
# Fix permissions
chmod +x instaaudit
ls -la instaaudit  # Should show -rwxr-xr-x

# If still doesn't work, copy to bin
cp instaaudit $PREFIX/bin/
```

**Problem: Command not found**
```bash
# Check if binary exists
ls -la instaaudit

# Add current directory to PATH temporarily
export PATH=$PATH:$(pwd)

# Or use absolute path
$HOME/your-project-folder/instaaudit --help
```

### Go Build Fails
```bash
# Clear module cache
go clean -modcache
go mod download
go mod tidy

# Check Go installation
go version
which go
```

### Network Issues
```bash
# Test basic connectivity
ping google.com
nslookup example.com

# Check if you have network permissions
curl -I https://google.com
```

### Storage Access
```bash
# Allow Termux storage access
termux-setup-storage

# Check storage permissions
ls /sdcard/
```

### Termux-Specific Issues

**Problem: Shared library errors**
```bash
# Update Termux packages
pkg update && pkg upgrade

# Reinstall Go if needed
pkg reinstall golang
```

**Problem: Binary won't execute**
```bash
# Check architecture
uname -m

# Rebuild for correct architecture
GOOS=android GOARCH=arm64 go build -o instaaudit cmd/main.go
```

## 🚀 **Mobile Security Scenarios**

### **Scenario 1: Hotel WiFi Security Check**
```bash
# Check if hotel network is secure
instaaudit -H $(ip route | grep default | awk '{print $3}') -p 22,23,80,443

# Look for open services that shouldn't be there
# Red flags: SSH (22), Telnet (23), databases
```

### **Scenario 2: Home Network Audit**
```bash
# Find all devices on your network
nmap -sn 192.168.1.0/24  # Discovery scan first

# Then audit each device
instaaudit -H 192.168.1.1 -p common    # Router
instaaudit -H 192.168.1.100 -p common  # Smart TV
instaaudit -H 192.168.1.150 -p common  # IoT device
```

### **Scenario 3: Website Security Check**
```bash
# Check your own website security
instaaudit -H yoursite.com -p 80,443 -f html

# Focus on web security
instaaudit -H webapp.com -p 80,443,8080,8443 --skip-recon
```

### **Scenario 4: Learning Lab Setup**
```bash
# Create a learning environment
# Scan intentionally vulnerable apps (with permission)
instaaudit -H dvwa.local -p 80,443
instaaudit -H metasploitable -p common
```

## 📱 **Termux Optimization Tips**

### **Battery & Performance**
```bash
# Quick scans for mobile
instaaudit -H target.com -p 80,443,22 -t 1

# Skip heavy operations
instaaudit -H target.com --skip-exploits --skip-recon

# Limit port range
instaaudit -H target.com -p 1-1000
```

### **Storage Management**
```bash
# Save to external storage
termux-setup-storage
instaaudit -H target.com -o /sdcard/security-scans/scan1

# Compress old reports
gzip /sdcard/security-scans/*.json
```

### **Network Considerations**
```bash
# Use mobile data carefully (data usage)
instaaudit -H target.com -p 80,443 -t 2

# WiFi-only scans for larger ranges
instaaudit -H 192.168.1.0/24 -p common
```

## 🔧 **Advanced Termux Usage**

### **Custom Configuration**
```bash
# Create mobile-optimized config
cat > mobile-config.json << EOF
{
  "timeout": "1s",
  "output_path": "/sdcard/instaaudit/reports",
  "output_format": "html",
  "skip_exploits": true,
  "aggressive": false
}
EOF

# Use custom config
instaaudit -H target.com -c mobile-config.json
```

### **Automated Mobile Scanning**
```bash
# Create scan automation script
cat > mobile-scan.sh << 'EOF'
#!/data/data/com.termux/files/usr/bin/bash

TARGET=$1
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OUTPUT_DIR="/sdcard/security-scans"

mkdir -p "$OUTPUT_DIR"

echo "📱 Mobile Security Scan: $TARGET"
echo "Timestamp: $TIMESTAMP"

# Quick scan optimized for mobile
instaaudit -H "$TARGET" -p common -t 1 \
  -o "$OUTPUT_DIR/scan_${TIMESTAMP}" \
  -f html --skip-exploits

echo "✅ Scan complete: $OUTPUT_DIR/scan_${TIMESTAMP}.html"
EOF

chmod +x mobile-scan.sh
```

### **Network Discovery & Scanning**
```bash
# Find network range
ip route | grep -E "192\.168|10\.|172\."

# Discover active hosts
nmap -sn 192.168.1.0/24 | grep "Nmap scan report"

# Batch scan discovered hosts
for ip in $(nmap -sn 192.168.1.0/24 | grep "Nmap scan report" | awk '{print $5}'); do
    echo "Scanning $ip..."
    instaaudit -H $ip -p 22,80,443 -t 1 -o "/sdcard/scans/$ip"
done
```

## 📋 **Mobile Security Checklists**

### **Home Network Security Checklist**
```bash
# 1. Router Security
instaaudit -H 192.168.1.1 -p 22,23,80,443,8080
# ❌ SSH/Telnet should be disabled
# ✅ HTTPS admin interface only

# 2. IoT Device Check
instaaudit -H smart-device-ip -p common
# ❌ No default passwords
# ❌ No unnecessary services

# 3. Network Isolation
# Check if IoT devices can reach main network
```

### **Public WiFi Safety Check**
```bash
# 1. Gateway Security
instaaudit -H $(ip route | grep default | awk '{print $3}') -p common

# 2. Look for red flags:
# ❌ Open databases (3306, 5432, 27017)
# ❌ File shares (445, 139, 21)
# ❌ Remote access (22, 3389, 5900)
```

### **Website Security Audit**
```bash
# 1. Basic Web Security
instaaudit -H yoursite.com -p 80,443

# 2. Check for:
# ✅ HTTPS redirect (port 80 → 443)
# ✅ Security headers
# ❌ Admin interfaces exposed
```

## 🎯 **Real-World Termux Scenarios**

### **Scenario: Coffee Shop Security**
```bash
# You're at a coffee shop and want to check network security
# 1. Find the gateway
GATEWAY=$(ip route | grep default | awk '{print $3}')

# 2. Quick security check
instaaudit -H $GATEWAY -p 22,23,80,443,21,445 -t 1

# 3. Red flags to look for:
# - SSH/Telnet access to router
# - File sharing services
# - Unencrypted admin interfaces
```

### **Scenario: Smart Home Setup**
```bash
# Setting up IoT devices securely
# 1. Scan each new device
instaaudit -H 192.168.1.100 -p common  # Smart TV
instaaudit -H 192.168.1.101 -p common  # Security camera
instaaudit -H 192.168.1.102 -p common  # Smart speaker

# 2. Check for common IoT issues:
# - Default credentials
# - Unnecessary services
# - Unencrypted communications
```

### **Scenario: Bug Bounty Research**
```bash
# Mobile bug bounty hunting (authorized targets only)
# 1. Initial reconnaissance
instaaudit -H target.com -p 80,443,8080,8443 --skip-exploits

# 2. Service enumeration
instaaudit -H target.com -p 1-10000 -t 1

# 3. Focus on findings
instaaudit -H target.com -p found-ports -A
```

## 🔒 **Security & Legal Notes**

### **Legal Considerations**
- ✅ **Your own networks** - Always OK
- ✅ **Authorized testing** - With written permission
- ✅ **Bug bounty programs** - Follow program rules
- ❌ **Unauthorized scanning** - Illegal in most places
- ❌ **Public networks** - Without permission

### **Ethical Guidelines**
- Only scan networks you own or have permission to test
- Respect rate limits and don't overload systems
- Report vulnerabilities responsibly
- Don't access or modify data without authorization
- Follow local laws and regulations

### **Mobile-Specific Considerations**
- **Data usage** - Scans can consume mobile data
- **Battery drain** - Intensive scans drain battery
- **Network detection** - Some networks detect scanning
- **Legal jurisdiction** - Laws vary by location

### **Best Practices**
- Start with your own devices and networks
- Use educational mode to learn concepts
- Verify findings with multiple tools
- Document your testing authorization
- Keep learning about cybersecurity ethics

## Getting Help

- Check `./instaaudit --help` for all options
- Read `COMMANDS.md` for detailed usage
- Review `UNDERSTANDING-RESULTS.md` for result interpretation