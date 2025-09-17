package education

import "fmt"

// PortExplanation provides educational information about ports
type PortExplanation struct {
	Port        int    `json:"port"`
	ServiceName string `json:"service_name"`
	Description string `json:"description"`
	RiskLevel   string `json:"risk_level"`
	WhyOpen     string `json:"why_open"`
	ShouldBe    string `json:"should_be"`
	VerifyWith  string `json:"verify_with"`
}

// GetPortExplanation returns educational information about a specific port
func GetPortExplanation(port int, service string) *PortExplanation {
	explanations := map[int]*PortExplanation{
		21: {
			Port:        21,
			ServiceName: "FTP (File Transfer Protocol)",
			Description: "Used for transferring files between computers. Sends passwords in plain text.",
			RiskLevel:   "High",
			WhyOpen:     "File sharing, website uploads, legacy systems",
			ShouldBe:    "Closed - Use SFTP (port 22) instead for secure file transfer",
			VerifyWith:  "ftp target.com (try anonymous login)",
		},
		22: {
			Port:        22,
			ServiceName: "SSH (Secure Shell)",
			Description: "Secure remote login and file transfer. Encrypted connection.",
			RiskLevel:   "Medium",
			WhyOpen:     "Remote server administration, secure file transfer",
			ShouldBe:    "Open only if remote access needed. Use key-based authentication.",
			VerifyWith:  "ssh user@target.com",
		},
		23: {
			Port:        23,
			ServiceName: "Telnet",
			Description: "Unencrypted remote login. Sends passwords in plain text.",
			RiskLevel:   "Critical",
			WhyOpen:     "Legacy systems, old network equipment",
			ShouldBe:    "Closed - Use SSH (port 22) instead",
			VerifyWith:  "telnet target.com 23",
		},
		25: {
			Port:        25,
			ServiceName: "SMTP (Email Sending)",
			Description: "Sends email messages between servers.",
			RiskLevel:   "Medium",
			WhyOpen:     "Email server, sending notifications",
			ShouldBe:    "Open only for email servers. Should require authentication.",
			VerifyWith:  "telnet target.com 25 (try EHLO command)",
		},
		53: {
			Port:        53,
			ServiceName: "DNS (Domain Name System)",
			Description: "Translates website names to IP addresses.",
			RiskLevel:   "Low",
			WhyOpen:     "DNS server, domain name resolution",
			ShouldBe:    "Open for DNS servers. Should be configured securely.",
			VerifyWith:  "nslookup google.com target.com",
		},
		80: {
			Port:        80,
			ServiceName: "HTTP (Web Server)",
			Description: "Serves websites without encryption. Data sent in plain text.",
			RiskLevel:   "Medium",
			WhyOpen:     "Website hosting, web applications",
			ShouldBe:    "Should redirect to HTTPS (port 443) for security",
			VerifyWith:  "curl -I http://target.com",
		},
		110: {
			Port:        110,
			ServiceName: "POP3 (Email Retrieval)",
			Description: "Downloads email from server. Unencrypted.",
			RiskLevel:   "Medium",
			WhyOpen:     "Email server for client access",
			ShouldBe:    "Use POP3S (port 995) instead for encryption",
			VerifyWith:  "telnet target.com 110",
		},
		143: {
			Port:        143,
			ServiceName: "IMAP (Email Access)",
			Description: "Access email on server. Unencrypted.",
			RiskLevel:   "Medium",
			WhyOpen:     "Email server for client access",
			ShouldBe:    "Use IMAPS (port 993) instead for encryption",
			VerifyWith:  "telnet target.com 143",
		},
		443: {
			Port:        443,
			ServiceName: "HTTPS (Secure Web Server)",
			Description: "Serves websites with SSL/TLS encryption. Secure.",
			RiskLevel:   "Low",
			WhyOpen:     "Secure website hosting, web applications",
			ShouldBe:    "Good - encrypted web traffic",
			VerifyWith:  "curl -I https://target.com",
		},
		993: {
			Port:        993,
			ServiceName: "IMAPS (Secure Email Access)",
			Description: "Encrypted email access using SSL/TLS.",
			RiskLevel:   "Low",
			WhyOpen:     "Secure email server access",
			ShouldBe:    "Good - encrypted email access",
			VerifyWith:  "openssl s_client -connect target.com:993",
		},
		995: {
			Port:        995,
			ServiceName: "POP3S (Secure Email Retrieval)",
			Description: "Encrypted email download using SSL/TLS.",
			RiskLevel:   "Low",
			WhyOpen:     "Secure email server access",
			ShouldBe:    "Good - encrypted email retrieval",
			VerifyWith:  "openssl s_client -connect target.com:995",
		},
		3306: {
			Port:        3306,
			ServiceName: "MySQL Database",
			Description: "Database server for storing application data.",
			RiskLevel:   "Critical",
			WhyOpen:     "Database server, web applications",
			ShouldBe:    "Should NOT be accessible from internet. Use firewall.",
			VerifyWith:  "mysql -h target.com -u root -p",
		},
		3389: {
			Port:        3389,
			ServiceName: "RDP (Remote Desktop)",
			Description: "Windows remote desktop access.",
			RiskLevel:   "Critical",
			WhyOpen:     "Remote Windows administration",
			ShouldBe:    "Should NOT be accessible from internet. Use VPN.",
			VerifyWith:  "rdesktop target.com (or Windows Remote Desktop)",
		},
		5432: {
			Port:        5432,
			ServiceName: "PostgreSQL Database",
			Description: "Advanced database server for applications.",
			RiskLevel:   "Critical",
			WhyOpen:     "Database server, web applications",
			ShouldBe:    "Should NOT be accessible from internet. Use firewall.",
			VerifyWith:  "psql -h target.com -U postgres",
		},
		5900: {
			Port:        5900,
			ServiceName: "VNC (Remote Desktop)",
			Description: "Cross-platform remote desktop access.",
			RiskLevel:   "Critical",
			WhyOpen:     "Remote desktop administration",
			ShouldBe:    "Should NOT be accessible from internet. Use VPN.",
			VerifyWith:  "vncviewer target.com",
		},
		6379: {
			Port:        6379,
			ServiceName: "Redis Database",
			Description: "In-memory database and cache server.",
			RiskLevel:   "Critical",
			WhyOpen:     "Caching, session storage, real-time applications",
			ShouldBe:    "Should NOT be accessible from internet. Use firewall.",
			VerifyWith:  "redis-cli -h target.com",
		},
		8080: {
			Port:        8080,
			ServiceName: "HTTP Alternative",
			Description: "Alternative web server port, often for development.",
			RiskLevel:   "Medium",
			WhyOpen:     "Development servers, alternative web services",
			ShouldBe:    "Should use HTTPS and proper authentication",
			VerifyWith:  "curl -I http://target.com:8080",
		},
		8443: {
			Port:        8443,
			ServiceName: "HTTPS Alternative",
			Description: "Alternative secure web server port.",
			RiskLevel:   "Low",
			WhyOpen:     "Alternative HTTPS services, admin panels",
			ShouldBe:    "Good if properly configured with valid SSL",
			VerifyWith:  "curl -I https://target.com:8443",
		},
		27017: {
			Port:        27017,
			ServiceName: "MongoDB Database",
			Description: "NoSQL document database server.",
			RiskLevel:   "Critical",
			WhyOpen:     "Document database, web applications, APIs",
			ShouldBe:    "Should NOT be accessible from internet. Use authentication.",
			VerifyWith:  "mongo target.com:27017",
		},
	}

	if explanation, exists := explanations[port]; exists {
		return explanation
	}

	// Default explanation for unknown ports
	return &PortExplanation{
		Port:        port,
		ServiceName: fmt.Sprintf("Unknown Service (%s)", service),
		Description: "This port is running a service that wasn't identified.",
		RiskLevel:   "Medium",
		WhyOpen:     "Unknown - could be custom application or misconfigured service",
		ShouldBe:    "Investigate what service is running and if it should be public",
		VerifyWith:  fmt.Sprintf("nmap -sV -p %d target.com", port),
	}
}

// GetSecurityRecommendation provides security advice based on findings
func GetSecurityRecommendation(riskLevel string, findings []string) string {
	switch riskLevel {
	case "Critical":
		return `🚨 IMMEDIATE ACTION REQUIRED:
• This poses a serious security risk
• Could lead to data breach or system compromise
• Fix within 24 hours
• Consider taking the service offline until fixed
• Implement network firewall rules
• Review access logs for suspicious activity`

	case "High":
		return `⚠️ HIGH PRIORITY - Fix within 1 week:
• Significant security vulnerability
• Could be exploited by attackers
• Implement proper authentication
• Use encryption where possible
• Restrict network access
• Monitor for unusual activity`

	case "Medium":
		return `🔶 MEDIUM PRIORITY - Fix within 1 month:
• Security improvement needed
• Reduces overall security posture
• Update software versions
• Implement security best practices
• Review configuration settings
• Consider additional security measures`

	case "Low":
		return `✅ LOW PRIORITY - Monitor and improve:
• Good security practices in place
• Minor improvements possible
• Keep software updated
• Regular security reviews
• Monitor for changes
• Maintain current security measures`

	default:
		return `📋 REVIEW RECOMMENDED:
• Assess the security implications
• Verify the configuration is intentional
• Implement appropriate security measures
• Regular monitoring and updates`
	}
}

// GetVerificationSteps provides manual verification commands
func GetVerificationSteps(port int, service string) []string {
	steps := []string{}

	switch port {
	case 22:
		steps = append(steps, "ssh -V target.com  # Check SSH version")
		steps = append(steps, "ssh target.com -o PreferredAuthentications=none  # Test auth")
	case 80:
		steps = append(steps, "curl -I http://target.com  # Check headers")
		steps = append(steps, "curl http://target.com/robots.txt  # Check robots.txt")
	case 443:
		steps = append(steps, "curl -I https://target.com  # Check HTTPS headers")
		steps = append(steps, "openssl s_client -connect target.com:443  # Check SSL")
	case 3306:
		steps = append(steps, "mysql -h target.com -u root -p  # Test MySQL access")
		steps = append(steps, "nmap -sV -p 3306 target.com  # Check version")
	case 5432:
		steps = append(steps, "psql -h target.com -U postgres  # Test PostgreSQL")
		steps = append(steps, "nmap -sV -p 5432 target.com  # Check version")
	default:
		steps = append(steps, fmt.Sprintf("nmap -sV -p %d target.com  # Identify service", port))
		steps = append(steps, fmt.Sprintf("telnet target.com %d  # Test connection", port))
	}

	return steps
}