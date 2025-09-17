package report

import (
	"fmt"
	"os"
	"strings"
)

// EducationalReport generates beginner-friendly explanations
func GenerateEducationalReport(auditReport *Report, outputPath string) error {
	// Create both text and HTML educational reports
	if err := generateEducationalText(auditReport, outputPath); err != nil {
		return err
	}
	return generateEducationalHTML(auditReport, outputPath)
}

// generateEducationalText creates a simple text explanation
func generateEducationalText(auditReport *Report, outputPath string) error {
	file, err := os.Create(outputPath + "_educational.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	// Header with trust information
	fmt.Fprintf(file, "🎓 InstaAudit Educational Report\n")
	fmt.Fprintf(file, "================================\n\n")
	fmt.Fprintf(file, "🛡️  TRUST & VERIFICATION:\n")
	fmt.Fprintf(file, "This report shows what InstaAudit found, but you should ALWAYS verify important findings!\n")
	fmt.Fprintf(file, "• Use './verify-results.sh %s' to cross-check with other tools\n", auditReport.ScanResult.Host)
	fmt.Fprintf(file, "• Test critical findings manually (see VERIFICATION-GUIDE.md)\n")
	fmt.Fprintf(file, "• When in doubt, ask an expert or use multiple tools\n\n")

	fmt.Fprintf(file, "Target: %s\n", auditReport.ScanResult.Host)
	fmt.Fprintf(file, "Scan Date: %s\n\n", auditReport.Timestamp.Format("2006-01-02 15:04:05"))

	// Executive Summary for Beginners
	fmt.Fprintf(file, "📊 WHAT WAS FOUND (Simple Explanation):\n")
	fmt.Fprintf(file, "=======================================\n\n")

	riskEmoji := getRiskEmoji(auditReport.Summary.RiskLevel)
	fmt.Fprintf(file, "%s Overall Security Level: %s\n\n", riskEmoji, auditReport.Summary.RiskLevel)

	// Explain risk levels with real-world context
	fmt.Fprintf(file, "🎯 What Risk Levels Mean:\n")
	fmt.Fprintf(file, "🟢 Low = Like a house with good locks - minor improvements possible\n")
	fmt.Fprintf(file, "🟡 Medium = Like leaving a window unlocked - should be fixed\n")
	fmt.Fprintf(file, "🟠 High = Like leaving your front door open - fix soon!\n")
	fmt.Fprintf(file, "🔴 Critical = Like leaving your house keys in the door - fix NOW!\n\n")

	// Detailed risk breakdown
	fmt.Fprintf(file, "📈 DETAILED SECURITY BREAKDOWN:\n")
	fmt.Fprintf(file, "==============================\n")
	fmt.Fprintf(file, "Critical Issues: %d (Fix immediately!)\n", auditReport.Summary.CriticalIssues)
	fmt.Fprintf(file, "High Risk Issues: %d (Fix within days)\n", auditReport.Summary.HighRiskIssues)
	fmt.Fprintf(file, "Medium Risk Issues: %d (Fix within weeks)\n", auditReport.Summary.MediumRiskIssues)
	fmt.Fprintf(file, "Low Risk Issues: %d (Fix when convenient)\n\n", auditReport.Summary.LowRiskIssues)

	// Explain what was scanned
	fmt.Fprintf(file, "🔍 WHAT WE CHECKED:\n")
	fmt.Fprintf(file, "==================\n")
	fmt.Fprintf(file, "Think of your computer/server like a building with doors (ports):\n")
	fmt.Fprintf(file, "• Scanned %d network 'doors' (ports) on %s\n", auditReport.Summary.TotalPorts, auditReport.ScanResult.Host)
	fmt.Fprintf(file, "• Found %d doors that are open (accessible from internet)\n", auditReport.Summary.OpenPorts)
	fmt.Fprintf(file, "• Identified %d services (programs answering the doors)\n", auditReport.Summary.ServicesFound)
	fmt.Fprintf(file, "• Found %d potential security problems\n\n", len(auditReport.Summary.RiskDetails))

	// Show specific security issues found
	if len(auditReport.Summary.RiskDetails) > 0 {
		fmt.Fprintf(file, "�  SPECIFIC SECURITY ISSUES FOUND:\n")
		fmt.Fprintf(file, "=================================\n\n")
		for i, detail := range auditReport.Summary.RiskDetails {
			fmt.Fprintf(file, "%d. %s\n", i+1, detail)
			fmt.Fprintf(file, "   ➤ How to verify: Use './verify-results.sh %s' or test manually\n\n", auditReport.ScanResult.Host)
		}
	}

	// Explain open ports in simple terms
	if len(auditReport.ScanResult.OpenPorts) > 0 {
		fmt.Fprintf(file, "🚪 OPEN DOORS (PORTS) FOUND:\n")
		fmt.Fprintf(file, "============================\n\n")
		fmt.Fprintf(file, "These are 'doors' that are open on %s:\n\n", auditReport.ScanResult.Host)

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
// gene
rateEducationalHTML creates a visual, beginner-friendly HTML report
func generateEducationalHTML(auditReport *Report, outputPath string) error {
	file, err := os.Create(outputPath + "_educational.html")
	if err != nil {
		return err
	}
	defer file.Close()

	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>InstaAudit Educational Report - {{.ScanResult.Host}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 20px; border-radius: 10px; margin-bottom: 20px; }
        .trust-warning { background: #fff3cd; border: 1px solid #ffeaa7; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .risk-critical { background: #ffebee; border-left: 5px solid #f44336; padding: 10px; margin: 10px 0; }
        .risk-high { background: #fff3e0; border-left: 5px solid #ff9800; padding: 10px; margin: 10px 0; }
        .risk-medium { background: #fffde7; border-left: 5px solid #ffeb3b; padding: 10px; margin: 10px 0; }
        .risk-low { background: #e8f5e8; border-left: 5px solid #4caf50; padding: 10px; margin: 10px 0; }
        .port-explanation { background: #f8f9fa; border: 1px solid #dee2e6; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .verification-box { background: #e3f2fd; border: 1px solid #2196f3; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .good-news { color: #4caf50; font-weight: bold; }
        .warning { color: #ff9800; font-weight: bold; }
        .danger { color: #f44336; font-weight: bold; }
        .code { background: #f4f4f4; padding: 2px 5px; border-radius: 3px; font-family: monospace; }
        .section { margin: 20px 0; }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
        .card { background: white; border: 1px solid #ddd; border-radius: 8px; padding: 15px; }
        .emoji { font-size: 1.2em; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎓 InstaAudit Educational Report</h1>
            <h2>Security Scan Results for {{.ScanResult.Host}}</h2>
            <p>Scanned on {{.Timestamp.Format "January 2, 2006 at 3:04 PM"}}</p>
        </div>

        <div class="trust-warning">
            <h3>🛡️ Important: Trust & Verification</h3>
            <p><strong>Always verify important security findings!</strong> This report shows what InstaAudit detected, but you should:</p>
            <ul>
                <li>Cross-check critical findings with other tools</li>
                <li>Test database exposure claims manually</li>
                <li>Verify SSL/certificate issues with online tools</li>
                <li>When in doubt, consult a security expert</li>
            </ul>
            <p><strong>Verification Command:</strong> <span class="code">./verify-results.sh {{.ScanResult.Host}}</span></p>
        </div>

        <div class="section">
            <h2>📊 Overall Security Assessment</h2>
            <div class="risk-{{.Summary.RiskLevel | lower}}">
                <h3>{{getRiskEmoji .Summary.RiskLevel}} Security Level: {{.Summary.RiskLevel}}</h3>
                <p><strong>What this means:</strong></p>
                {{if eq .Summary.RiskLevel "Critical"}}
                    <p class="danger">🚨 <strong>IMMEDIATE ACTION REQUIRED!</strong> Critical security vulnerabilities detected that could allow hackers easy access to your system.</p>
                {{else if eq .Summary.RiskLevel "High"}}
                    <p class="warning">⚠️ <strong>SERIOUS ISSUES FOUND!</strong> High-risk security problems that should be fixed within days.</p>
                {{else if eq .Summary.RiskLevel "Medium"}}
                    <p class="warning">🟡 <strong>SOME CONCERNS:</strong> Medium-risk issues that should be addressed within weeks.</p>
                {{else}}
                    <p class="good-news">✅ <strong>LOOKING GOOD!</strong> No major security issues detected, just minor improvements possible.</p>
                {{end}}
            </div>
        </div>

        <div class="section">
            <h2>📈 Security Issues Breakdown</h2>
            <div class="grid">
                <div class="card risk-critical">
                    <h3>🔴 Critical Issues: {{.Summary.CriticalIssues}}</h3>
                    <p><strong>Fix immediately!</strong> These are like leaving your house keys in the door.</p>
                </div>
                <div class="card risk-high">
                    <h3>🟠 High Risk Issues: {{.Summary.HighRiskIssues}}</h3>
                    <p><strong>Fix within days!</strong> These are like leaving your front door unlocked.</p>
                </div>
                <div class="card risk-medium">
                    <h3>🟡 Medium Risk Issues: {{.Summary.MediumRiskIssues}}</h3>
                    <p><strong>Fix within weeks.</strong> These are like leaving a window unlocked.</p>
                </div>
                <div class="card risk-low">
                    <h3>🟢 Low Risk Issues: {{.Summary.LowRiskIssues}}</h3>
                    <p><strong>Fix when convenient.</strong> Minor security improvements.</p>
                </div>
            </div>
        </div>

        {{if .Summary.RiskDetails}}
        <div class="section">
            <h2>🚨 Specific Security Issues Found</h2>
            {{range $index, $detail := .Summary.RiskDetails}}
            <div class="risk-high">
                <h4>Issue {{add $index 1}}: {{$detail}}</h4>
                <div class="verification-box">
                    <strong>🔍 How to verify this finding:</strong>
                    <p>Don't just trust this result! Here's how to check it yourself:</p>
                    <ul>
                        <li>Use the verification script: <span class="code">./verify-results.sh {{$.ScanResult.Host}}</span></li>
                        <li>Test manually with appropriate tools</li>
                        <li>Check with online security scanners</li>
                        <li>Consult the VERIFICATION-GUIDE.md for detailed steps</li>
                    </ul>
                </div>
            </div>
            {{end}}
        </div>
        {{end}}

        {{if .ScanResult.OpenPorts}}
        <div class="section">
            <h2>🚪 Open Ports (Network "Doors") Found</h2>
            <p><strong>What are ports?</strong> Think of ports like doors on a building. Each door has a number and serves a specific purpose.</p>
            
            {{range .ScanResult.OpenPorts}}
            <div class="port-explanation">
                <h4>🚪 Port {{.}} is OPEN</h4>
                {{$portInfo := getPortExplanation .}}
                <p><strong>What this port is for:</strong> {{$portInfo.Purpose}}</p>
                <p><strong>Is this normal?</strong> {{$portInfo.Normal}}</p>
                <p><strong>Security concern:</strong> {{$portInfo.Risk}}</p>
                
                <div class="verification-box">
                    <strong>🔍 Verify this yourself:</strong>
                    <p>Test if port {{.}} is really open: <span class="code">nc -zv {{$.ScanResult.Host}} {{.}}</span></p>
                    {{if eq . 80}}
                    <p>For web servers: <span class="code">curl -I http://{{$.ScanResult.Host}}</span></p>
                    {{else if eq . 443}}
                    <p>For HTTPS: <span class="code">curl -I https://{{$.ScanResult.Host}}</span></p>
                    {{end}}
                </div>
            </div>
            {{end}}
        </div>
        {{end}}

        <div class="section">
            <h2>🎓 Learning Resources</h2>
            <div class="grid">
                <div class="card">
                    <h3>📚 Understanding Results</h3>
                    <p>Read <span class="code">UNDERSTANDING-RESULTS.md</span> for detailed explanations of security concepts.</p>
                </div>
                <div class="card">
                    <h3>🔍 Verification Guide</h3>
                    <p>Learn how to verify findings in <span class="code">VERIFICATION-GUIDE.md</span></p>
                </div>
                <div class="card">
                    <h3>🛡️ Trust & Security</h3>
                    <p>Read <span class="code">TRUST-AND-VERIFICATION.md</span> to understand how to trust security tools.</p>
                </div>
            </div>
        </div>

        <div class="section">
            <h2>🚀 Next Steps</h2>
            <ol>
                <li><strong>Verify critical findings</strong> using the verification scripts</li>
                <li><strong>Prioritize fixes</strong> starting with Critical and High risk issues</li>
                <li><strong>Learn more</strong> about the security concepts mentioned</li>
                <li><strong>Re-scan</strong> after making fixes to confirm improvements</li>
                <li><strong>Ask for help</strong> if you're unsure about any findings</li>
            </ol>
        </div>

        <div class="trust-warning">
            <h3>🎯 Remember: Security is a Journey, Not a Destination</h3>
            <p>This scan is a snapshot in time. Security requires ongoing attention, learning, and verification. 
            Always question results, verify findings, and continue learning about cybersecurity!</p>
        </div>
    </div>
</body>
</html>`

	// Execute template with report data
	tmpl := template.Must(template.New("educational").Funcs(template.FuncMap{
		"lower": strings.ToLower,
		"add": func(a, b int) int { return a + b },
		"getRiskEmoji": func(risk string) string {
			switch risk {
			case "Critical": return "🔴"
			case "High": return "🟠"
			case "Medium": return "🟡"
			case "Low": return "🟢"
			default: return "⚪"
			}
		},
		"getPortExplanation": func(port int) map[string]string {
			explanations := map[int]map[string]string{
				22:   {"Purpose": "SSH (Secure Shell) - Remote access to the system", "Normal": "Normal for servers, concerning for personal computers", "Risk": "Medium - Could allow unauthorized access if weak passwords"},
				80:   {"Purpose": "HTTP - Web server (unencrypted)", "Normal": "Normal for websites", "Risk": "Medium - Should redirect to HTTPS for security"},
				443:  {"Purpose": "HTTPS - Secure web server", "Normal": "Normal and good for websites", "Risk": "Low - This is secure web traffic"},
				3306: {"Purpose": "MySQL Database", "Normal": "Should NEVER be open to internet", "Risk": "CRITICAL - Database exposed to hackers!"},
				5432: {"Purpose": "PostgreSQL Database", "Normal": "Should NEVER be open to internet", "Risk": "CRITICAL - Database exposed to hackers!"},
				3389: {"Purpose": "Windows Remote Desktop", "Normal": "Should NEVER be open to internet", "Risk": "CRITICAL - Remote access exposed!"},
			}
			
			if exp, exists := explanations[port]; exists {
				return exp
			}
			return map[string]string{
				"Purpose": fmt.Sprintf("Service running on port %d", port),
				"Normal": "Depends on what service this is",
				"Risk": "Unknown - investigate what service is running",
			}
		},
	}).Parse(htmlTemplate))

	return tmpl.Execute(file, auditReport)
}