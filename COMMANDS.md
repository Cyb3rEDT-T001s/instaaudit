# 📋 InstaAudit Command Reference Guide

## 🚀 Basic Usage

```bash
instaaudit -H <target> [options]
```

## 📖 Command Structure

### Required Parameters
```bash
-H, --host string      Target host to audit (REQUIRED)
```

### Optional Parameters
```bash
-p, --ports string     Ports to scan
-t, --timeout int      Connection timeout in seconds (default: 2)
-f, --format string    Output format (default: json)
-o, --output string    Output file path (default: ./audit_report)
-A, --aggressive       Enable comprehensive scanning
-c, --config string    Configuration file path
    --skip-exploits    Skip exploit testing
    --skip-recon       Skip reconnaissance
-h, --help            Show help
```

## 🎯 Port Specification Options

### Common Ports
```bash
# Use predefined common ports (21,22,23,25,53,80,110,143,443,993,995,3306,3389,5432,5900,6379,8080,8443,27017)
instaaudit -H target.com
instaaudit -H target.com -p "common"
```

### Specific Ports
```bash
# Single port
instaaudit -H target.com -p "80"

# Multiple ports
instaaudit -H target.com -p "80,443,22"

# Port ranges
instaaudit -H target.com -p "80-90"

# Mixed specification
instaaudit -H target.com -p "22,80-90,443,8000-8080"
```

### Service-Specific Port Groups
```bash
# Web services
instaaudit -H target.com -p "80,443,8080,8443"

# Email services
instaaudit -H target.com -p "25,110,143,993,995"

# Database services
instaaudit -H target.com -p "3306,5432,27017,6379"

# Remote access services
instaaudit -H target.com -p "22,23,3389,5900"
```

## 📊 Output Format Options

### Single Format
```bash
# JSON (default) - Machine readable
instaaudit -H target.com -f json

# HTML - Beautiful web dashboard
instaaudit -H target.com -f html

# CSV - Spreadsheet format
instaaudit -H target.com -f csv

# TXT - Plain text summary
instaaudit -H target.com -f txt
```

### Multiple Formats
```bash
# All formats at once
instaaudit -H target.com -f all

# Generates: audit_report.json, audit_report.html, audit_report.csv, audit_report.txt
```

## ⚡ Scanning Modes

### Basic Scan
```bash
# Standard security scan
instaaudit -H target.com
```

### Aggressive Scan
```bash
# Comprehensive security assessment
instaaudit -H target.com -A

# Aggressive with specific ports
instaaudit -H target.com -A -p "80,443,3306,5432"

# Aggressive with HTML report
instaaudit -H target.com -A -f html
```

### Stealth Scan
```bash
# Skip reconnaissance (lower footprint)
instaaudit -H target.com --skip-recon

# Skip exploit testing (faster)
instaaudit -H target.com --skip-exploits

# Minimal scan
instaaudit -H target.com --skip-recon --skip-exploits -t 5
```

## 🕐 Timeout Options

```bash
# Fast scan (1 second timeout)
instaaudit -H target.com -t 1

# Standard scan (2 seconds - default)
instaaudit -H target.com -t 2

# Slow/reliable scan (5 seconds)
instaaudit -H target.com -t 5

# Very thorough scan (10 seconds)
instaaudit -H target.com -t 10
```

## 📁 Output File Options

### Default Output
```bash
# Creates: audit_report.json (or specified format)
instaaudit -H target.com
```

### Custom Output Names
```bash
# Custom filename
instaaudit -H target.com -o security_assessment

# Creates: security_assessment.json (or specified format)

# With path
instaaudit -H target.com -o reports/daily_scan

# Creates: reports/daily_scan.json
```

## 🎯 Complete Command Examples

### 🌐 Web Application Security
```bash
# Basic web security scan
instaaudit -H webapp.com -p "80,443" -f html

# Comprehensive web audit
instaaudit -H webapp.com -p "80,443,8080,8443" -A -f all -o web_security_audit

# Quick web check
instaaudit -H webapp.com -p "80,443" --skip-recon -f html
```

### 🗄️ Database Security Assessment
```bash
# Database security scan
instaaudit -H db-server.com -p "3306,5432,27017,6379" -A -f html

# MySQL specific
instaaudit -H mysql-server.com -p "3306" -A -f all -o mysql_audit

# All databases with timeout
instaaudit -H db-cluster.com -p "3306,5432,27017,6379" -A -t 5 -f html
```

### 🏢 Enterprise Network Audit
```bash
# Comprehensive enterprise scan
instaaudit -H corporate-server.com -A -f all -o enterprise_audit

# Infrastructure assessment
instaaudit -H infrastructure.com -p "22,23,80,443,3389,5900" -A -f html

# Complete security audit
instaaudit -H target.com -A -t 5 -f all -o complete_security_assessment
```

### 🔍 Reconnaissance Focus
```bash
# Reconnaissance heavy scan
instaaudit -H target.com -A -f html -o recon_report

# Skip exploits, focus on discovery
instaaudit -H target.com --skip-exploits -f json -o discovery_scan
```

### ⚡ Performance Optimized
```bash
# Fast scan
instaaudit -H target.com -t 1 --skip-recon --skip-exploits -f json

# Balanced scan
instaaudit -H target.com -t 3 -f html

# Thorough scan
instaaudit -H target.com -A -t 10 -f all -o thorough_audit
```

## 🔧 Configuration File Usage

### Create Configuration File
```json
// config.json
{
  "timeout": "5s",
  "max_workers": 200,
  "output_path": "./reports/security_audit",
  "output_format": "html",
  "ports": [21, 22, 23, 25, 53, 80, 110, 143, 443, 993, 995, 3306, 3389, 5432, 5900, 6379, 8080, 8443, 27017],
  "skip_exploits": false,
  "skip_recon": false,
  "aggressive": true
}
```

### Use Configuration File
```bash
# Use config file
instaaudit -H target.com -c config.json

# Override config settings
instaaudit -H target.com -c config.json -f all -o custom_output
```

## 🎓 Educational Commands

### For Learning
```bash
# Basic educational scan
instaaudit -H google.com -f html

# Learn about specific services
instaaudit -H target.com -p "80,443" -f html

# Comprehensive learning scan
instaaudit -H safe-target.com -A -f all
```

### Safe Practice Targets
```bash
# Safe targets for learning (won't cause issues)
instaaudit -H google.com -p "80,443"
instaaudit -H github.com -p "80,443,22"
instaaudit -H stackoverflow.com -p "80,443"
```

## 🔬 Expert Commands

### Professional Assessment
```bash
# Full professional audit
instaaudit -H target.com -A -f all -t 5 -o professional_audit

# Compliance scan
instaaudit -H internal-server.com -A -f csv -o compliance_$(date +%Y%m%d)

# Penetration testing reconnaissance
instaaudit -H target.com -A -t 10 -f json -o pentest_recon
```

### Batch Scanning
```bash
# Multiple targets (use in script)
for target in server1.com server2.com server3.com; do
    instaaudit -H $target -A -f html -o audit_$target
done
```

## 🚨 Emergency/Incident Response

### Quick Security Check
```bash
# Fast incident response scan
instaaudit -H compromised-server.com --skip-recon -t 1 -f txt

# Emergency database check
instaaudit -H db-server.com -p "3306,5432,27017,6379" -A -f json

# Quick web security check
instaaudit -H webapp.com -p "80,443" --skip-exploits -f html
```

## 📊 Report Generation Commands

### View Reports
```bash
# Open HTML report in browser
start audit_report.html          # Windows
open audit_report.html           # macOS
xdg-open audit_report.html       # Linux

# View JSON report
type audit_report.json           # Windows
cat audit_report.json            # Linux/macOS

# View educational report
type audit_report_explained.txt  # Windows
cat audit_report_explained.txt   # Linux/macOS
```

## 🛠️ Troubleshooting Commands

### Verify Installation
```bash
# Check if InstaAudit is working
instaaudit --help

# Test with safe target
instaaudit -H google.com -p "80,443" -f txt
```

### Debug Issues
```bash
# Verbose output (if implemented)
instaaudit -H target.com -v

# Test specific port
instaaudit -H target.com -p "80" -t 10 -f txt

# Minimal test
instaaudit -H 8.8.8.8 -p "53" -f txt
```

## 🎯 Common Use Cases

### Daily Security Monitoring
```bash
instaaudit -H internal-server.com -f json -o daily_$(date +%Y%m%d)
```

### Weekly Comprehensive Audit
```bash
instaaudit -H all-servers.com -A -f all -o weekly_audit_$(date +%Y%m%d)
```

### Pre-Deployment Security Check
```bash
instaaudit -H new-server.com -A -f html -o pre_deployment_check
```

### Compliance Reporting
```bash
instaaudit -H production-server.com -A -f csv -o compliance_report
```

### Vendor Security Assessment
```bash
instaaudit -H vendor-system.com -A -f all -o vendor_assessment
```

## 🔍 Verification Commands

### Cross-Check with Other Tools
```bash
# After InstaAudit scan, verify with:
nmap -sV target.com
curl -I https://target.com
openssl s_client -connect target.com:443
```

## ⚠️ Important Notes

### Legal Usage
- Only scan systems you own or have permission to test
- Unauthorized scanning may violate laws
- Use responsibly for security testing and education

### Performance Considerations
- Use appropriate timeouts for network conditions
- Consider target system load when choosing scan intensity
- Use stealth options for sensitive environments

### Best Practices
- Always verify findings with multiple tools
- Document all scans with timestamps
- Keep reports secure and confidential
- Follow responsible disclosure for vulnerabilities

---

**📚 This guide covers all InstaAudit commands and options. Use it as a reference for both learning and professional security auditing!** 🛡️