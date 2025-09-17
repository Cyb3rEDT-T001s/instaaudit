package report

import (
	"fmt"
	"os"
	"strings"
)

// EducationalReport generates beginner-friendly explanations
func GenerateEducationalReport(auditReport *Report, outputPath string) error {
	file, err := os.Create(outputPath + "_explained.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	// Header
	fmt.Fprintf(file, "🎓 InstaAudit Educational Report\n")
	fmt.Fprintf(file, "================================\n\n")
	fmt.Fprintf(file, "Target: %s\n", auditReport.ScanResult.Host)
	fmt.Fprintf(file, "Scan Date: %s\n\n", auditReport.Timestamp.Format("2006-01-02 15:04:05"))

	// Executive Summary for Beginners
	fmt.Fprintf(file, "📊 WHAT WAS FOUND (Simple Explanation):\n")
	fmt.Fprintf(file, "=======================================\n\n")

	riskEmoji := getRiskEmoji(auditReport.Summary.RiskLevel)
	fmt.Fprintf(file, "%s Overall Security Level: %s\n\n", riskEmoji, auditReport.Summary.RiskLevel)

	// Explain risk levels
	fmt.Fprintf(file, "Risk Level Meanings:\n")
	fmt.Fprintf(file, "🟢 Low = Good security, minor improvements possible\n")
	fmt.Fprintf(file, "🟡 Medium = Some security issues, should be fixed\n")
	fmt.Fprintf(file, "🟠 High = Serious problems, fix soon\n")
	fmt.Fprintf(file, "🔴 Critical = Immediate danger, fix right now!\n\n")

	// Explain what was scanned
	fmt.Fprintf(file, "🔍 WHAT WE CHECKED:\n")
	fmt.Fprintf(file, "==================\n")
	fmt.Fprintf(file, "• Scanned %d network 'doors' (ports) on %s\n", auditReport.Summary.TotalPorts, auditReport.ScanResult.Host)
	fmt.Fprintf(file, "• Found %d doors that are open\n", auditReport.Summary.OpenPorts)
	fmt.Fprintf(file, "• Identified %d services (programs) running\n", auditReport.Summary.ServicesFound)
	fmt.Fprintf(file, "• Discovered %d security issues\n\n", auditReport.Summary.Vulnerabilities+auditReport.Summary.ExploitableServices+auditReport.Summary.Misconfigurations)

	// Explain open ports in simple terms
	if len(auditReport.ScanResult.OpenPorts) > 0 {
		fmt.Fprintf(file, "🚪 OPEN DOORS (PORTS) FOUND:\n")
		fmt.Fprintf(file, "============================\n\n")

		for _, port := range auditReport.ScanResult.OpenPorts {
			service := "Unknown Service"
			for _, svc := range auditReport.AuditResult.Services {
				if svc.Port == port {
					service = svc.Service
					break
				}
			}

			explanation := getPortExplanationSimple(port, service)
			fmt.Fprintf(file, "Port %d - %s\n", port, explanation.ServiceName)
			fmt.Fprintf(file, "   What it does: %s\n", explanation.Description)
			fmt.Fprintf(file, "   Risk level: %s\n", explanation.RiskLevel)
			fmt.Fprintf(file, "   Should it be open? %s\n", explanation.ShouldBe)
			fmt.Fprintf(file, "   How to check: %s\n\n", explanation.VerifyWith)
		}
	}

	// Explain security issues found
	if auditReport.Summary.Vulnerabilities > 0 || auditReport.Summary.ExploitableServices > 0 {
		fmt.Fprintf(file, "🚨 SECURITY PROBLEMS FOUND:\n")
		fmt.Fprintf(file, "===========================\n\n")

		// Vulnerabilities
		for _, vuln := range auditReport.AuditResult.Vulnerabilities {
			fmt.Fprintf(file, "Problem: %s\n", vuln.Description)
			fmt.Fprintf(file, "Severity: %s (Score: %.1f/10)\n", vuln.Severity, vuln.Score)
			fmt.Fprintf(file, "What this means: This is a known security weakness that attackers could exploit\n")
			fmt.Fprintf(file, "What to do: Update the software or apply security patches\n\n")
		}

		// Exploit results
		for _, exploit := range auditReport.AuditResult.ExploitResults {
			if exploit.Success {
				fmt.Fprintf(file, "Problem: %s on port %d\n", exploit.ExploitName, exploit.Port)
				fmt.Fprintf(file, "Service: %s\n", exploit.Service)
				fmt.Fprintf(file, "Risk: %s\n", exploit.Severity)
				fmt.Fprintf(file, "Details: %s\n", exploit.Details)
				fmt.Fprintf(file, "What this means: This service has a security weakness\n")
				fmt.Fprintf(file, "What to do: Fix the configuration or update the software\n\n")
			}
		}
	}

	// Database security issues
	if auditReport.Summary.DatabaseIssues > 0 {
		fmt.Fprintf(file, "🗄️ DATABASE SECURITY ISSUES:\n")
		fmt.Fprintf(file, "============================\n\n")

		for _, db := range auditReport.AuditResult.DatabaseResults {
			if db.Accessible {
				fmt.Fprintf(file, "Database Type: %s (Port %d)\n", db.Service, db.Port)
				fmt.Fprintf(file, "Risk Level: %s\n", db.Severity)
				fmt.Fprintf(file, "What this means: Your database can be accessed from the internet\n")
				fmt.Fprintf(file, "Why this is bad: Attackers could steal or delete your data\n")
				fmt.Fprintf(file, "Problems found:\n")
				for _, warning := range db.Warnings {
					fmt.Fprintf(file, "   • %s\n", warning)
				}
				fmt.Fprintf(file, "What to do: Use a firewall to block internet access to your database\n\n")
			}
		}
	}

	// Web application issues
	if auditReport.Summary.WebAppIssues > 0 {
		fmt.Fprintf(file, "🌐 WEBSITE SECURITY ISSUES:\n")
		fmt.Fprintf(file, "===========================\n\n")

		for _, web := range auditReport.AuditResult.WebAppResults {
			if len(web.Warnings) > 0 {
				fmt.Fprintf(file, "Website: %s\n", web.URL)
				fmt.Fprintf(file, "Risk Level: %s\n", web.Severity)
				fmt.Fprintf(file, "Problems found:\n")
				for _, warning := range web.Warnings {
					fmt.Fprintf(file, "   • %s\n", warning)
					
					// Add explanations for common issues
					if strings.Contains(warning, "Missing HSTS") {
						fmt.Fprintf(file, "     Explanation: Website doesn't force secure connections\n")
					} else if strings.Contains(warning, "Missing Content Security Policy") {
						fmt.Fprintf(file, "     Explanation: Website lacks protection against malicious scripts\n")
					} else if strings.Contains(warning, "Admin panel accessible") {
						fmt.Fprintf(file, "     Explanation: Management interface is publicly accessible\n")
					}
				}
				fmt.Fprintf(file, "\n")
			}
		}
	}

	// Recommendations section
	fmt.Fprintf(file, "💡 WHAT YOU SHOULD DO:\n")
	fmt.Fprintf(file, "======================\n\n")

	switch auditReport.Summary.RiskLevel {
	case "Critical":
		fmt.Fprintf(file, "🚨 URGENT - Act immediately:\n")
		fmt.Fprintf(file, "1. Disconnect dangerous services from the internet\n")
		fmt.Fprintf(file, "2. Change all default passwords\n")
		fmt.Fprintf(file, "3. Install security updates\n")
		fmt.Fprintf(file, "4. Set up a firewall\n")
		fmt.Fprintf(file, "5. Get professional help if needed\n\n")
	case "High":
		fmt.Fprintf(file, "⚠️ Important - Fix within a week:\n")
		fmt.Fprintf(file, "1. Update software to latest versions\n")
		fmt.Fprintf(file, "2. Configure proper security settings\n")
		fmt.Fprintf(file, "3. Restrict network access where possible\n")
		fmt.Fprintf(file, "4. Monitor for suspicious activity\n\n")
	case "Medium":
		fmt.Fprintf(file, "🔶 Moderate - Fix within a month:\n")
		fmt.Fprintf(file, "1. Implement missing security features\n")
		fmt.Fprintf(file, "2. Review and update configurations\n")
		fmt.Fprintf(file, "3. Plan security improvements\n")
		fmt.Fprintf(file, "4. Regular security monitoring\n\n")
	default:
		fmt.Fprintf(file, "✅ Good - Keep monitoring:\n")
		fmt.Fprintf(file, "1. Maintain current security measures\n")
		fmt.Fprintf(file, "2. Keep software updated\n")
		fmt.Fprintf(file, "3. Regular security scans\n")
		fmt.Fprintf(file, "4. Stay informed about new threats\n\n")
	}

	// Expert verification section
	fmt.Fprintf(file, "🔬 FOR EXPERTS - HOW TO VERIFY THESE RESULTS:\n")
	fmt.Fprintf(file, "=============================================\n\n")

	fmt.Fprintf(file, "Manual verification commands:\n\n")

	// Port verification
	fmt.Fprintf(file, "Verify open ports:\n")
	fmt.Fprintf(file, "  nmap -F %s\n", auditReport.ScanResult.Host)
	fmt.Fprintf(file, "  nmap -sV %s\n\n", auditReport.ScanResult.Host)

	// Service-specific verification
	for _, port := range auditReport.ScanResult.OpenPorts {
		switch port {
		case 22:
			fmt.Fprintf(file, "Verify SSH (port 22):\n")
			fmt.Fprintf(file, "  ssh -V %s\n", auditReport.ScanResult.Host)
			fmt.Fprintf(file, "  ssh %s -o PreferredAuthentications=none\n\n", auditReport.ScanResult.Host)
		case 80:
			fmt.Fprintf(file, "Verify HTTP (port 80):\n")
			fmt.Fprintf(file, "  curl -I http://%s\n", auditReport.ScanResult.Host)
			fmt.Fprintf(file, "  curl http://%s/robots.txt\n\n", auditReport.ScanResult.Host)
		case 443:
			fmt.Fprintf(file, "Verify HTTPS (port 443):\n")
			fmt.Fprintf(file, "  curl -I https://%s\n", auditReport.ScanResult.Host)
			fmt.Fprintf(file, "  openssl s_client -connect %s:443\n\n", auditReport.ScanResult.Host)
		case 3306:
			fmt.Fprintf(file, "Verify MySQL (port 3306):\n")
			fmt.Fprintf(file, "  mysql -h %s -u root -p\n", auditReport.ScanResult.Host)
			fmt.Fprintf(file, "  nmap -sV -p 3306 %s\n\n", auditReport.ScanResult.Host)
		}
	}

	fmt.Fprintf(file, "Cross-reference with other tools:\n")
	fmt.Fprintf(file, "  • Nmap: https://nmap.org/\n")
	fmt.Fprintf(file, "  • SSL Labs: https://www.ssllabs.com/ssltest/\n")
	fmt.Fprintf(file, "  • Security Headers: https://securityheaders.com/\n\n")

	fmt.Fprintf(file, "📚 Learn more:\n")
	fmt.Fprintf(file, "  • Port numbers: https://www.speedguide.net/ports.php\n")
	fmt.Fprintf(file, "  • Web security: https://owasp.org/\n")
	fmt.Fprintf(file, "  • Network security: https://www.sans.org/\n\n")

	fmt.Fprintf(file, "Generated by InstaAudit - Professional Security Auditing Tool\n")
	fmt.Fprintf(file, "https://github.com/Cyb3rEDT-T001s/instaaudit\n")

	return nil
}

// Simple port explanation structure
type SimplePortExplanation struct {
	ServiceName string
	Description string
	RiskLevel   string
	ShouldBe    string
	VerifyWith  string
}

// getPortExplanationSimple provides simple explanations for beginners
func getPortExplanationSimple(port int, service string) SimplePortExplanation {
	explanations := map[int]SimplePortExplanation{
		22:    {"SSH (Secure Remote Login)", "Allows secure remote access to the computer", "Medium", "Only if you need remote access", "ssh user@target.com"},
		25:    {"SMTP (Email Sending)", "Sends email messages", "Medium", "Only for email servers", "telnet target.com 25"},
		53:    {"DNS (Website Name Lookup)", "Translates website names to numbers", "Low", "Usually safe", "nslookup google.com target.com"},
		80:    {"HTTP (Website)", "Serves websites without encryption", "Medium", "Should redirect to HTTPS", "curl -I http://target.com"},
		110:   {"POP3 (Email Download)", "Downloads email from server", "Medium", "Use secure version instead", "telnet target.com 110"},
		143:   {"IMAP (Email Access)", "Access email on server", "Medium", "Use secure version instead", "telnet target.com 143"},
		443:   {"HTTPS (Secure Website)", "Serves websites with encryption", "Low", "Good - secure websites", "curl -I https://target.com"},
		993:   {"IMAPS (Secure Email)", "Encrypted email access", "Low", "Good - secure email", "openssl s_client -connect target.com:993"},
		995:   {"POP3S (Secure Email)", "Encrypted email download", "Low", "Good - secure email", "openssl s_client -connect target.com:995"},
		3306:  {"MySQL Database", "Stores application data", "Critical", "Should NOT be public", "mysql -h target.com -u root -p"},
		3389:  {"Remote Desktop", "Windows remote control", "Critical", "Should NOT be public", "Remote Desktop Connection"},
		5432:  {"PostgreSQL Database", "Advanced database server", "Critical", "Should NOT be public", "psql -h target.com -U postgres"},
		6379:  {"Redis Database", "Fast memory database", "Critical", "Should NOT be public", "redis-cli -h target.com"},
		27017: {"MongoDB Database", "Document database", "Critical", "Should NOT be public", "mongo target.com:27017"},
	}

	if explanation, exists := explanations[port]; exists {
		return explanation
	}

	return SimplePortExplanation{
		ServiceName: fmt.Sprintf("Unknown Service (%s)", service),
		Description: "A program is running but we don't know what it does",
		RiskLevel:   "Medium",
		ShouldBe:    "Investigate what this service does",
		VerifyWith:  fmt.Sprintf("nmap -sV -p %d target.com", port),
	}
}

// getRiskEmoji returns emoji for risk levels
func getRiskEmoji(riskLevel string) string {
	switch riskLevel {
	case "Critical":
		return "🔴"
	case "High":
		return "🟠"
	case "Medium":
		return "🟡"
	case "Low":
		return "🟢"
	default:
		return "⚪"
	}
}