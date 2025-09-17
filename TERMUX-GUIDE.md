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
```

## Usage Examples

### Basic Scan
```bash
./instaaudit -H 192.168.1.1 -p common
```

### Scan Specific Ports
```bash
./instaaudit -H example.com -p 80,443,22,21
```

### Full Scan with All Formats
```bash
./instaaudit -H target.com -p common -f all -o scan_results
```

### Quick Scan (Skip Heavy Tests)
```bash
./instaaudit -H 10.0.0.1 --skip-exploits --skip-recon
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

### Go Build Fails
```bash
# Clear module cache
go clean -modcache
go mod download
go mod tidy
```

### Permission Denied
```bash
# Fix permissions
chmod +x instaaudit
chmod +x install-termux.sh
```

### Network Issues
```bash
# Test basic connectivity
ping google.com
nslookup example.com
```

### Storage Access
```bash
# Allow Termux storage access
termux-setup-storage
```

## Advanced Usage

### Custom Config
Create `config.json`:
```json
{
  "timeout": "2s",
  "output_path": "/sdcard/instaaudit/reports",
  "output_format": "json",
  "skip_exploits": false,
  "aggressive": false
}
```

### Automated Scans
Create a scan script:
```bash
#!/data/data/com.termux/files/usr/bin/bash
./instaaudit -H $1 -p common -f all -o "/sdcard/scans/$(date +%Y%m%d_%H%M%S)"
```

## Security Notes

- InstaAudit is for authorized testing only
- Respect network policies and laws
- Some features may trigger security alerts
- Use responsibly on your own networks

## Getting Help

- Check `./instaaudit --help` for all options
- Read `COMMANDS.md` for detailed usage
- Review `UNDERSTANDING-RESULTS.md` for result interpretation