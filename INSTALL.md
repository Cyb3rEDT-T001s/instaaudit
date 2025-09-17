# InstaAudit Installation Guide

## Platform-Specific Instructions

### 🪟 Windows

1. **Install Go:**
   - Download from: https://golang.org/dl/
   - Choose: `go1.21.x.windows-amd64.msi`
   - Install with default settings

2. **Build InstaAudit:**
   ```cmd
   build.bat
   ```

3. **Run:**
   ```cmd
   instaaudit.exe -H target.com
   ```

### 🐧 Linux

#### Ubuntu/Debian:
```bash
# Install Go
sudo apt update
sudo apt install golang-go

# Build InstaAudit
chmod +x build.sh
./build.sh

# Run
./instaaudit -H target.com
```

#### CentOS/RHEL/Fedora:
```bash
# Install Go
sudo yum install golang  # or dnf install golang

# Build InstaAudit
chmod +x build.sh
./build.sh

# Run
./instaaudit -H target.com
```

#### Arch Linux:
```bash
# Install Go
sudo pacman -S go

# Build InstaAudit
chmod +x build.sh
./build.sh

# Run
./instaaudit -H target.com
```

### 🍎 macOS

#### Using Homebrew:
```bash
# Install Go
brew install go

# Build InstaAudit
chmod +x build.sh
./build.sh

# Run
./instaaudit -H target.com
```

#### Manual Installation:
```bash
# Download Go from https://golang.org/dl/
# Choose: go1.21.x.darwin-amd64.pkg
# Install the package

# Build InstaAudit
chmod +x build.sh
./build.sh

# Run
./instaaudit -H target.com
```

## 🐳 Docker Installation

```bash
# Build Docker image
chmod +x docker-build.sh
./docker-build.sh

# Run scan
docker run --rm instaaudit:latest -H target.com

# Save reports
docker run --rm -v $(pwd):/reports instaaudit:latest -H target.com -o /reports/audit
```

## 🔧 Using Make (Linux/macOS)

```bash
# Build for current platform
make build

# Build for all platforms
make all-platforms

# Install system-wide
make install

# Clean build files
make clean

# Run demo
make demo
```

## 📦 One-Command Install (Linux/macOS)

```bash
# Automatic installation
curl -sSL https://raw.githubusercontent.com/Cyb3rEDT-T001s/instaaudit/main/install-unix.sh | bash
```

## ✅ Verification

After installation, verify InstaAudit works:

```bash
# Show help
./instaaudit --help

# Test scan (safe)
./instaaudit -H google.com -p "80,443"

# Check version
./instaaudit --version
```

## 🚀 Quick Examples

```bash
# Basic security scan
./instaaudit -H target.com

# Comprehensive scan with HTML report
./instaaudit -H target.com -A -f html

# Database security audit
./instaaudit -H db-server.com -p "3306,5432,27017,6379"

# Web application security
./instaaudit -H webapp.com -p "80,443,8080,8443" -f html

# All report formats
./instaaudit -H target.com -f all
```

## 🛠️ Troubleshooting

### Go Not Found
```bash
# Check if Go is installed
go version

# Add Go to PATH (if needed)
export PATH=$PATH:/usr/local/go/bin
```

### Permission Denied
```bash
# Make scripts executable
chmod +x build.sh install-unix.sh docker-build.sh

# Make binary executable
chmod +x instaaudit
```

### Build Errors
```bash
# Clean and rebuild
make clean
go clean -cache
go mod tidy
make build
```

---

**InstaAudit supports Windows, Linux, macOS, and Docker!** 🌍