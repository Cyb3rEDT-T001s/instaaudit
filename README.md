# InstaAudit - Professional Security Auditing Tool 🛡️

A comprehensive, educational security auditing tool built in Go that makes cybersecurity accessible to everyone - from beginners to experts.

## ✨ Key Features

### 🔍 **Comprehensive Security Assessment**
- **Port Scanning**: Fast concurrent TCP port scanning with service identification
- **Database Security**: MySQL, PostgreSQL, MongoDB, Redis vulnerability testing
- **Web Application Security**: SSL/TLS analysis, security headers, vulnerability detection
- **System Security**: File permissions, SUID/SGID analysis, configuration auditing
- **Reconnaissance**: DNS enumeration, subdomain discovery, technology detection

### 🎓 **Educational & Beginner-Friendly**
- **Explains everything in simple terms** - Perfect for people new to cybersecurity
- **Educational reports** that teach you about ports, services, and security risks
- **Risk level explanations** with clear guidance on what to fix first
- **Cross-verification guides** for security professionals

### 📊 **Professional Reporting**
- **HTML Dashboard**: Beautiful web interface with charts and risk analysis
- **Educational Report**: Beginner-friendly explanations of all findings
- **JSON/CSV/TXT**: Technical data for integration and analysis
- **Expert Verification**: Commands to manually verify all findings

### 🌍 **Cross-Platform Support**
- **Windows**: Native .exe with simple build script
- **Linux**: Native binary with auto-installer
- **macOS**: Native binary with Homebrew support
- **Docker**: Containerized deployment for any platform

## 🚀 Quick Start

### 📥 Installation

#### **Windows Users:**
```cmd
# 1. Install Go from: https://golang.org/dl/
# 2. Clone and build
git clone https://github.com/Cyb3rEDT-T001s/instaaudit.git
cd instaaudit
build.bat
```

#### **Linux/macOS Users:**
```bash
# One-command install
curl -sSL https://raw.githubusercontent.com/Cyb3rEDT-T001s/instaaudit/main/install-unix.sh | bash

# Or manual install
git clone https://github.com/Cyb3rEDT-T001s/instaaudit.git
cd instaaudit
chmod +x install-unix.sh
./install-unix.sh
```

#### **Termux (Android) Users:**
```bash
# Install from Termux
git clone https://github.com/Cyb3rEDT-T001s/instaaudit.git
cd instaaudit
chmod +x install-termux.sh
./install-termux.sh
```

#### **Docker Users:**
```bash
git clone https://github.com/Cyb3rEDT-T001s/instaaudit.git
cd instaaudit

# Build Docker image
chmod +x docker-build.sh
./docker-build.sh

# Or use Docker Compose
docker-compose build
```

### 🎯 Usage Examples

#### **For Beginners:**
```bash
# Basic security scan (generates educational report)
./instaaudit -H target.com

# Get beautiful HTML report that explains everything
./instaaudit -H target.com -f html
```

#### **For Security Professionals:**
```bash
# Comprehensive security audit
./instaaudit -H target.com -A -f all

# Database security focus
./instaaudit -H db-server.com -p "3306,5432,27017,6379" -A

# Web application security
./instaaudit -H webapp.com -p "80,443,8080,8443" -f html
```

### 📊 What You Get

After running a scan, InstaAudit generates:

1. **📋 Educational Report** (`audit_report_explained.txt`)
   - Simple explanations of what each finding means
   - Risk levels explained in plain English
   - Step-by-step recommendations

2. **🌐 HTML Dashboard** (`audit_report.html`)
   - Beautiful web interface with charts
   - Color-coded risk levels
   - Professional security assessment

3. **📄 Technical Reports** (JSON, CSV, TXT)
   - Machine-readable data for integration
   - Detailed technical findings
   - Expert verification commands

### Cross-Platform Builds

Build for multiple platforms:
```bash
# Build for all platforms
make all-platforms

# Specific platforms
make linux    # Creates instaaudit-linux
make windows  # Creates instaaudit.exe
make macos    # Creates instaaudit-macos
```

### Docker Support

Run InstaAudit in Docker:
```bash
# Build Docker image
chmod +x docker-build.sh
./docker-build.sh

# Run scan in container
docker run --rm instaaudit:latest -H target.com

# Save reports to host
docker run --rm -v $(pwd):/reports instaaudit:latest -H target.com -o /reports/audit
```

### Package Managers

Install Go using package managers:

**Linux:**
```bash
# Ubuntu/Debian
sudo apt update && sudo apt install golang-go

# CentOS/RHEL/Fedora
sudo yum install golang  # or dnf install golang

# Arch Linux
sudo pacman -S go

# Alpine Linux
sudo apk add go
```

**macOS:**
```bash
# Homebrew
brew install go

# MacPorts
sudo port install go
```

## 🔧 Command Options

```
Required:
  -H, --host string      Target host to audit

Optional:
  -p, --ports string     Ports to scan (e.g., "80,443,22" or "common" for common ports)
  -t, --timeout int      Connection timeout in seconds (default: 2)
  -f, --format string    Output format: json, csv, html, txt, all (default: json)
  -o, --output string    Output file path without extension (default: "./audit_report")
  -A, --aggressive       Enable comprehensive scanning (more thorough)
      --skip-exploits    Skip exploit testing (faster scan)
      --skip-recon       Skip reconnaissance (stealth mode)
  -c, --config string    Configuration file path
  -h, --help            Show help and usage examples
```

### 📋 Usage Examples

```bash
# Basic scan with educational report
./instaaudit -H target.com

# Comprehensive scan with all reports
./instaaudit -H target.com -A -f all

# Quick web security check
./instaaudit -H webapp.com -p "80,443"

# Database security audit
./instaaudit -H db-server.com -p "3306,5432,27017,6379"

# Stealth scan (minimal footprint)
./instaaudit -H target.com --skip-recon -t 5

# Generate only HTML report
./instaaudit -H target.com -f html -o security_assessment
```

## 📊 Example Output

### 🎓 **For Beginners - Educational Report:**
```
🎓 InstaAudit Educational Report
================================

🚪 OPEN DOORS (PORTS) FOUND:

Port 80 - HTTP (Website)
   What it does: Serves websites without encryption
   Risk level: Medium
   Should it be open? Should redirect to HTTPS
   How to check: curl -I http://target.com

Port 3306 - MySQL Database
   What it does: Stores application data
   Risk level: Critical
   Should it be open? Should NOT be public
   How to check: mysql -h target.com -u root -p

💡 WHAT YOU SHOULD DO:
⚠️ Important - Fix within a week:
1. Use a firewall to block database access from internet
2. Redirect HTTP traffic to HTTPS
3. Update software to latest versions
```

### 🔬 **For Experts - Technical Report:**
```json
{
  "summary": {
    "risk_level": "High",
    "open_ports": 5,
    "vulnerabilities_found": 2,
    "database_issues": 1
  },
  "verification_commands": {
    "port_scan": "nmap -sV target.com",
    "web_security": "curl -I https://target.com",
    "database_check": "mysql -h target.com -u root -p"
  }
}
```

### 🌐 **HTML Dashboard Features:**
- **Risk Level Dashboard** with color-coded findings
- **Interactive Charts** showing security metrics  
- **Detailed Explanations** for each finding
- **Verification Commands** for expert review
- **Executive Summary** for management reporting

## 🛡️ Security Features

### 🔍 **Network Security**
- **Port Scanning**: Fast concurrent TCP scanning with service identification
- **Banner Grabbing**: Service version detection and fingerprinting
- **SSL/TLS Analysis**: Certificate validation, cipher analysis, protocol testing
- **Service Detection**: Automatic identification of running services

### 🗄️ **Database Security** 
- **MySQL/PostgreSQL/MongoDB/Redis**: Comprehensive security testing
- **Default Credentials**: Automated testing for weak authentication
- **Access Control**: Network accessibility validation
- **Version Analysis**: Known vulnerability detection

### 🌐 **Web Application Security**
- **Security Headers**: HSTS, CSP, X-Frame-Options, XSS Protection analysis
- **SSL/TLS Configuration**: Certificate validation and cipher strength testing
- **Admin Panel Detection**: Exposed management interfaces
- **Common Vulnerabilities**: Directory traversal, information disclosure

### 💻 **System Security**
- **File Permissions**: Critical system file auditing
- **SUID/SGID Analysis**: Privileged binary detection
- **Process Monitoring**: Running service analysis
- **Configuration Review**: Security hardening assessment

### 🕵️ **Reconnaissance**
- **DNS Enumeration**: Subdomain discovery and DNS analysis
- **Technology Detection**: Web framework and software identification
- **Information Gathering**: Passive reconnaissance techniques

## 🎯 Who Can Use InstaAudit?

### 👨‍🎓 **Beginners & Students**
- **Learn cybersecurity** with educational explanations
- **Understand security concepts** through practical examples
- **Get started** with hands-on security testing
- **Build knowledge** with guided learning materials

### 👨‍💼 **IT Professionals**
- **Security auditing** for internal systems
- **Compliance checking** for regulatory requirements
- **Risk assessment** for business systems
- **Network monitoring** and security validation

### 🔒 **Security Experts**
- **Penetration testing** with comprehensive scanning
- **Vulnerability assessment** with expert verification
- **Security consulting** with professional reports
- **Red team operations** with stealth scanning options

### 🏢 **Organizations**
- **Internal security audits** for compliance
- **Vendor security assessment** before partnerships
- **Regular security monitoring** of infrastructure
- **Security awareness training** with educational reports

## 📚 Learning Resources

### **Included Documentation:**
- 📖 **UNDERSTANDING-RESULTS.md** - Complete beginner's guide to security
- 🔍 **VERIFICATION-GUIDE.md** - Expert cross-checking and validation
- 📋 **INSTALL.md** - Platform-specific installation instructions

### **Educational Features:**
- **Port explanations** in simple terms
- **Risk level guidance** with clear priorities
- **Security recommendations** with actionable steps
- **Verification commands** for manual testing

## ⚖️ Legal & Ethical Use

### **✅ Authorized Use:**
- Systems you own or manage
- Authorized penetration testing
- Security research with permission
- Educational purposes on test systems

### **❌ Unauthorized Use:**
- Scanning systems without permission
- Testing third-party systems without consent
- Malicious security testing
- Violating terms of service

**⚠️ Important**: Always ensure you have explicit permission before scanning any systems. Unauthorized scanning may violate local laws and regulations.

## 🤝 Contributing

We welcome contributions from the cybersecurity community:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### **Areas for Contribution:**
- Additional security checks
- New vulnerability detection
- Educational content improvements
- Cross-platform compatibility
- Performance optimizations

## 📞 Support & Community

- 🐛 **Issues**: [GitHub Issues](https://github.com/Cyb3rEDT-T001s/instaaudit/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/Cyb3rEDT-T001s/instaaudit/discussions)
- 📖 **Documentation**: [Wiki](https://github.com/Cyb3rEDT-T001s/instaaudit/wiki)
- 🔗 **Main Toolkit**: [Cyb3rEDT-T001s](https://github.com/Cyb3rEDT-T001s)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**🎉 InstaAudit - Making Cybersecurity Accessible to Everyone!** 

*From complete beginners learning their first port scan to security experts conducting professional assessments - InstaAudit bridges the knowledge gap with educational, comprehensive security auditing.* 🛡️🎓