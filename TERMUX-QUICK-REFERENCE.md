# 📱 InstaAudit Termux Quick Reference

## 🚀 **Installation**
```bash
git clone https://github.com/Cyb3rEDT-T001s/instaaudit.git
cd instaaudit
chmod +x install-termux.sh
./install-termux.sh
```

## ⚡ **Quick Commands**

### **Basic Usage**
```bash
instaaudit -H target.com -p common          # Standard scan
instaaudit -H 192.168.1.1 -p 80,443,22     # Specific ports
instaaudit -H target.com --help             # Show help
```

### **Mobile-Optimized Scans**
```bash
instaaudit -H target.com -p 80,443 -t 1     # Fast scan (1s timeout)
instaaudit -H target.com --skip-exploits    # Skip heavy tests
instaaudit -H target.com -p common -f html  # HTML report only
```

## 🏠 **Common Use Cases**

| Scenario | Command | Purpose |
|----------|---------|---------|
| **Home Router** | `instaaudit -H 192.168.1.1 -p common` | Check router security |
| **WiFi Gateway** | `instaaudit -H $(ip route \| grep default \| awk '{print $3}') -p 22,23,80,443` | Public WiFi safety |
| **Smart Device** | `instaaudit -H 192.168.1.100 -p 80,443,8080` | IoT security check |
| **Website** | `instaaudit -H yoursite.com -p 80,443` | Web security audit |
| **Local Test** | `instaaudit -H localhost -p common` | Safe learning |
| **Network Scan** | `instaaudit -H 192.168.1.0/24 -p 22,80,443` | Find devices |

## 🔧 **Execution Methods**

| Method | Command | When to Use |
|--------|---------|-------------|
| **Global** | `instaaudit -H target.com -p common` | After installation |
| **Runner Script** | `./run-termux.sh -H target.com -p common` | If global fails |
| **Full Path** | `$(pwd)/instaaudit -H target.com -p common` | Path issues |
| **Bash Wrapper** | `bash -c "./instaaudit -H target.com -p common"` | Permission issues |

## 📊 **Output Options**

```bash
# Save to external storage
instaaudit -H target.com -o /sdcard/scans/scan1

# Multiple formats
instaaudit -H target.com -f all -o scan_results

# Educational report (auto-generated)
# Creates: scan_results_educational.html
```

## 🎯 **Port Ranges**

| Range | Usage | Example |
|-------|-------|---------|
| `common` | Standard services | `instaaudit -H target.com -p common` |
| `80,443,22` | Web + SSH | `instaaudit -H target.com -p 80,443,22` |
| `1-1000` | First 1000 ports | `instaaudit -H target.com -p 1-1000` |
| `80,443,8080,8443` | Web services | `instaaudit -H target.com -p 80,443,8080,8443` |

## 🚨 **Security Priorities**

### **Critical Issues (Fix Immediately)**
- Exposed databases (ports 3306, 5432, 27017)
- Remote access (ports 22, 3389, 5900)
- Default credentials

### **High Priority (Fix Soon)**
- Unencrypted web (port 80 without HTTPS redirect)
- Missing security headers
- Unnecessary services

### **Medium Priority (Fix When Possible)**
- Outdated software versions
- Weak SSL configurations
- Information disclosure

## 🔍 **Verification Commands**

```bash
# Cross-check results
./verify-results.sh target.com audit_report.json

# Manual verification
./manual-verify.sh target.com

# Show detailed findings
./show-details.sh audit_report.json
```

## 📱 **Mobile Tips**

### **Battery Optimization**
- Use `-t 1` for faster scans
- Add `--skip-exploits` to reduce CPU usage
- Limit port ranges: `-p 80,443,22` instead of `common`

### **Data Usage**
- Scan local networks on WiFi
- Use specific ports for remote scans
- Avoid large port ranges on mobile data

### **Storage Management**
```bash
# Setup external storage
termux-setup-storage

# Save to SD card
instaaudit -H target.com -o /sdcard/security-scans/scan1

# Compress old reports
gzip /sdcard/security-scans/*.json
```

## 🛠️ **Troubleshooting**

| Problem | Solution |
|---------|----------|
| `./instaaudit` doesn't work | Use `instaaudit` directly or `$(pwd)/instaaudit` |
| Permission denied | `chmod +x instaaudit` |
| Command not found | Check if in PATH: `which instaaudit` |
| Network errors | Test connectivity: `ping target.com` |
| Go build fails | `pkg update && pkg upgrade golang` |

## 📚 **Learning Resources**

| File | Purpose |
|------|---------|
| `TERMUX-GUIDE.md` | Complete Termux guide |
| `TRUST-AND-VERIFICATION.md` | How to verify results |
| `UNDERSTANDING-RESULTS.md` | Interpret scan results |
| `VERIFICATION-GUIDE.md` | Cross-check findings |
| `audit_report_educational.html` | Visual learning report |

## ⚖️ **Legal & Ethical**

### **✅ Always OK**
- Your own devices and networks
- Authorized penetration testing
- Bug bounty programs (follow rules)

### **❌ Never OK**
- Unauthorized network scanning
- Accessing systems without permission
- Ignoring "no testing" policies

### **🤔 Ask First**
- Corporate networks
- Public WiFi networks
- Shared hosting environments

## 🎓 **Learning Path**

1. **Start Safe**: `instaaudit -H localhost -p common`
2. **Home Network**: `instaaudit -H 192.168.1.1 -p common`
3. **Your Website**: `instaaudit -H yoursite.com -p 80,443`
4. **Verify Results**: `./verify-results.sh target.com audit_report.json`
5. **Learn Concepts**: Read educational reports
6. **Practice More**: Try different scenarios with permission

---

**Remember**: InstaAudit is a learning tool. Always verify critical findings and use responsibly! 🛡️