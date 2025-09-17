package report

import (
	"fmt"
	"html/template"
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
		fmt.Fprintf(file, "🚨 SPECIFIC SECURITY ISSUES FOUND:\n")
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
			fmt.Fprintf(file, "Port %d: ", port)
			
			// Simple port explanations
			switch port {
			case 22:
				fmt.Fprintf(file, "SSH (Secure Shell) - Remote access to the system\n")
				fmt.Fprintf(file, "   Risk: Medium - Could allow unauthorized access if weak passwords\n")
			case 80:
				fmt.Fprintf(file, "HTTP - Web server (unencrypted)\n")
				fmt.Fprintf(file, "   Risk: Medium - Should redirect to HTTPS for security\n")
			case 443:
				fmt.Fprintf(file, "HTTPS - Secure web server\n")
				fmt.Fprintf(file, "   Risk: Low - This is secure web traffic\n")
			case 3306:
				fmt.Fprintf(file, "MySQL Database\n")
				fmt.Fprintf(file, "   Risk: CRITICAL - Database should NEVER be open to internet!\n")
			case 5432:
				fmt.Fprintf(file, "PostgreSQL Database\n")
				fmt.Fprintf(file, "   Risk: CRITICAL - Database should NEVER be open to internet!\n")
			case 3389:
				fmt.Fprintf(file, "Windows Remote Desktop\n")
				fmt.Fprintf(file, "   Risk: CRITICAL - Remote access should NEVER be open to internet!\n")
			default:
				fmt.Fprintf(file, "Unknown service\n")
				fmt.Fprintf(file, "   Risk: Unknown - investigate what service is running\n")
			}
			fmt.Fprintf(file, "\n")
		}
	}

	return nil
}

// generateEducationalHTML creates a visual, beginner-friendly HTML report
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
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>InstaAudit Educational Report</h1>
            <h2>Security Scan Results for {{.ScanResult.Host}}</h2>
            <p>Scanned on {{.Timestamp.Format "January 2, 2006 at 3:04 PM"}}</p>
        </div>

        <div class="trust-warning">
            <h3>Important: Trust & Verification</h3>
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
            <h2>Overall Security Assessment</h2>
            <div class="risk-{{.Summary.RiskLevel | lower}}">
                <h3>Security Level: {{.Summary.RiskLevel}}</h3>
                <p><strong>What this means:</strong></p>
                {{if eq .Summary.RiskLevel "Critical"}}
                    <p class="danger">IMMEDIATE ACTION REQUIRED! Critical security vulnerabilities detected.</p>
                {{else if eq .Summary.RiskLevel "High"}}
                    <p class="warning">SERIOUS ISSUES FOUND! High-risk security problems that should be fixed within days.</p>
                {{else if eq .Summary.RiskLevel "Medium"}}
                    <p class="warning">SOME CONCERNS: Medium-risk issues that should be addressed within weeks.</p>
                {{else}}
                    <p class="good-news">LOOKING GOOD! No major security issues detected, just minor improvements possible.</p>
                {{end}}
            </div>
        </div>

        <div class="section">
            <h2>Security Issues Breakdown</h2>
            <div class="grid">
                <div class="card risk-critical">
                    <h3>Critical Issues: {{.Summary.CriticalIssues}}</h3>
                    <p><strong>Fix immediately!</strong> These are like leaving your house keys in the door.</p>
                </div>
                <div class="card risk-high">
                    <h3>High Risk Issues: {{.Summary.HighRiskIssues}}</h3>
                    <p><strong>Fix within days!</strong> These are like leaving your front door unlocked.</p>
                </div>
                <div class="card risk-medium">
                    <h3>Medium Risk Issues: {{.Summary.MediumRiskIssues}}</h3>
                    <p><strong>Fix within weeks.</strong> These are like leaving a window unlocked.</p>
                </div>
                <div class="card risk-low">
                    <h3>Low Risk Issues: {{.Summary.LowRiskIssues}}</h3>
                    <p><strong>Fix when convenient.</strong> Minor security improvements.</p>
                </div>
            </div>
        </div>

        {{if .Summary.RiskDetails}}
        <div class="section">
            <h2>Specific Security Issues Found</h2>
            {{range $index, $detail := .Summary.RiskDetails}}
            <div class="risk-high">
                <h4>Issue {{add $index 1}}: {{$detail}}</h4>
                <div class="verification-box">
                    <strong>How to verify this finding:</strong>
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

        <div class="section">
            <h2>Learning Resources</h2>
            <div class="grid">
                <div class="card">
                    <h3>Understanding Results</h3>
                    <p>Read <span class="code">UNDERSTANDING-RESULTS.md</span> for detailed explanations of security concepts.</p>
                </div>
                <div class="card">
                    <h3>Verification Guide</h3>
                    <p>Learn how to verify findings in <span class="code">VERIFICATION-GUIDE.md</span></p>
                </div>
                <div class="card">
                    <h3>Trust & Security</h3>
                    <p>Read <span class="code">TRUST-AND-VERIFICATION.md</span> to understand how to trust security tools.</p>
                </div>
            </div>
        </div>

        <div class="section">
            <h2>Next Steps</h2>
            <ol>
                <li><strong>Verify critical findings</strong> using the verification scripts</li>
                <li><strong>Prioritize fixes</strong> starting with Critical and High risk issues</li>
                <li><strong>Learn more</strong> about the security concepts mentioned</li>
                <li><strong>Re-scan</strong> after making fixes to confirm improvements</li>
                <li><strong>Ask for help</strong> if you're unsure about any findings</li>
            </ol>
        </div>

        <div class="trust-warning">
            <h3>Remember: Security is a Journey, Not a Destination</h3>
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
	}).Parse(htmlTemplate))

	return tmpl.Execute(file, auditReport)
}

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