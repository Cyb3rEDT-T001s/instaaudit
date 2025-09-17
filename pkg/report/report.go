package report

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/Cyb3rEDT-T001s/instaaudit/pkg/auditor"
	"github.com/Cyb3rEDT-T001s/instaaudit/pkg/scanner"
)

// Report represents the complete audit report
type Report struct {
	Timestamp   time.Time            `json:"timestamp"`
	ScanResult  *scanner.ScanResult  `json:"scan_result"`
	AuditResult *auditor.AuditResult `json:"audit_result"`
	Summary     *Summary             `json:"summary"`
}

// Summary provides a high-level overview of the audit
type Summary struct {
	TotalPorts         int    `json:"total_ports_scanned"`
	OpenPorts          int    `json:"open_ports"`
	ServicesFound      int    `json:"services_found"`
	Vulnerabilities    int    `json:"vulnerabilities_found"`
	Misconfigurations  int    `json:"misconfigurations_found"`
	ExploitableServices int   `json:"exploitable_services"`
	SubdomainsFound    int    `json:"subdomains_found"`
	TechnologiesFound  int    `json:"technologies_found"`
	DatabaseIssues     int    `json:"database_issues"`
	WebAppIssues       int    `json:"webapp_issues"`
	SystemIssues       int    `json:"system_issues"`
	RiskLevel          string `json:"risk_level"`
}

// GenerateReport creates a comprehensive report
func GenerateReport(scanResult *scanner.ScanResult, auditResult *auditor.AuditResult) *Report {
	exploitableCount := 0
	for _, exploit := range auditResult.ExploitResults {
		if exploit.Success {
			exploitableCount++
		}
	}

	subdomainCount := 0
	techCount := 0
	if auditResult.ReconData != nil {
		subdomainCount = len(auditResult.ReconData.Subdomains)
		techCount = len(auditResult.ReconData.Technologies)
	}

	// Count database issues
	dbIssues := 0
	for _, db := range auditResult.DatabaseResults {
		if db.Accessible {
			dbIssues++
		}
	}

	// Count web app issues
	webIssues := 0
	for _, web := range auditResult.WebAppResults {
		if len(web.Warnings) > 0 {
			webIssues++
		}
	}

	// Count system issues
	sysIssues := 0
	for _, sys := range auditResult.SystemResults {
		if len(sys.Findings) > 0 {
			sysIssues++
		}
	}

	summary := &Summary{
		TotalPorts:         len(scanResult.Results),
		OpenPorts:          len(scanResult.OpenPorts),
		ServicesFound:      len(auditResult.Services),
		Vulnerabilities:    len(auditResult.Vulnerabilities),
		Misconfigurations:  len(auditResult.Misconfigs),
		ExploitableServices: exploitableCount,
		SubdomainsFound:    subdomainCount,
		TechnologiesFound:  techCount,
		DatabaseIssues:     dbIssues,
		WebAppIssues:       webIssues,
		SystemIssues:       sysIssues,
		RiskLevel:          calculateRiskLevel(auditResult),
	}

	return &Report{
		Timestamp:   time.Now(),
		ScanResult:  scanResult,
		AuditResult: auditResult,
		Summary:     summary,
	}
}

// calculateRiskLevel determines the overall risk level
func calculateRiskLevel(audit *auditor.AuditResult) string {
	criticalCount := 0
	highRiskCount := 0
	mediumRiskCount := 0

	// Check vulnerabilities
	for _, vuln := range audit.Vulnerabilities {
		if vuln.Score >= 9.0 {
			criticalCount++
		} else if vuln.Score >= 7.0 {
			highRiskCount++
		} else if vuln.Score >= 4.0 {
			mediumRiskCount++
		}
	}

	// Check exploit results
	for _, exploit := range audit.ExploitResults {
		if exploit.Success {
			switch exploit.Severity {
			case "Critical":
				criticalCount++
			case "High":
				highRiskCount++
			case "Medium":
				mediumRiskCount++
			}
		}
	}

	// Check database results
	for _, db := range audit.DatabaseResults {
		if db.Accessible && db.Severity == "Critical" {
			criticalCount++
		} else if db.Accessible && db.Severity == "High" {
			highRiskCount++
		}
	}

	// Check web app results
	for _, web := range audit.WebAppResults {
		if web.Severity == "Critical" {
			criticalCount++
		} else if web.Severity == "High" {
			highRiskCount++
		} else if web.Severity == "Medium" {
			mediumRiskCount++
		}
	}

	if criticalCount > 0 {
		return "Critical"
	} else if highRiskCount > 0 {
		return "High"
	} else if mediumRiskCount > 0 || len(audit.Misconfigs) > 2 {
		return "Medium"
	}
	return "Low"
}

// SaveAsJSON saves the report in JSON format
func SaveAsJSON(report *Report, filepath string) error {
	file, err := os.Create(filepath + ".json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

// SaveAsCSV saves the report in CSV format
func SaveAsCSV(report *Report, filepath string) error {
	file, err := os.Create(filepath + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{"Category", "Item", "Severity", "Details"}
	writer.Write(headers)

	// Write port data
	for port, isOpen := range report.ScanResult.Results {
		if isOpen {
			service := "Unknown"
			for _, svc := range report.AuditResult.Services {
				if svc.Port == port {
					service = svc.Service
					break
				}
			}
			record := []string{"Open Port", fmt.Sprintf("Port %d", port), "Info", service}
			writer.Write(record)
		}
	}

	// Write vulnerabilities
	for _, vuln := range report.AuditResult.Vulnerabilities {
		record := []string{"Vulnerability", vuln.CVE, vuln.Severity, vuln.Description}
		writer.Write(record)
	}

	// Write exploit results
	for _, exploit := range report.AuditResult.ExploitResults {
		if exploit.Success {
			record := []string{"Exploit", exploit.ExploitName, exploit.Severity, exploit.Details}
			writer.Write(record)
		}
	}

	return nil
}

// SaveAsHTML saves the report in HTML format
func SaveAsHTML(report *Report, filepath string) error {
	return GenerateHTMLReport(report, filepath)
}

// SaveAsText saves the report in plain text format
func SaveAsText(report *Report, filepath string) error {
	file, err := os.Create(filepath + ".txt")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write report header
	fmt.Fprintf(file, "InstaAudit Enhanced Security Report\n")
	fmt.Fprintf(file, "===================================\n\n")
	fmt.Fprintf(file, "Host: %s\n", report.ScanResult.Host)
	fmt.Fprintf(file, "Timestamp: %s\n\n", report.Timestamp.Format(time.RFC3339))

	// Write summary
	fmt.Fprintf(file, "Executive Summary:\n")
	fmt.Fprintf(file, "-----------------\n")
	fmt.Fprintf(file, "Total Ports Scanned: %d\n", report.Summary.TotalPorts)
	fmt.Fprintf(file, "Open Ports: %d\n", report.Summary.OpenPorts)
	fmt.Fprintf(file, "Services Found: %d\n", report.Summary.ServicesFound)
	fmt.Fprintf(file, "Vulnerabilities: %d\n", report.Summary.Vulnerabilities)
	fmt.Fprintf(file, "Misconfigurations: %d\n", report.Summary.Misconfigurations)
	fmt.Fprintf(file, "Exploitable Services: %d\n", report.Summary.ExploitableServices)
	fmt.Fprintf(file, "Database Issues: %d\n", report.Summary.DatabaseIssues)
	fmt.Fprintf(file, "Web App Issues: %d\n", report.Summary.WebAppIssues)
	fmt.Fprintf(file, "System Issues: %d\n", report.Summary.SystemIssues)
	fmt.Fprintf(file, "Risk Level: %s\n\n", report.Summary.RiskLevel)

	// Write detailed sections
	writeDetailedSections(file, report)

	return nil
}

// writeDetailedSections writes all detailed sections to the text report
func writeDetailedSections(file *os.File, report *Report) {
	// Open ports section
	if len(report.ScanResult.OpenPorts) > 0 {
		fmt.Fprintf(file, "Open Ports & Services:\n")
		fmt.Fprintf(file, "---------------------\n")
		for _, port := range report.ScanResult.OpenPorts {
			service := "Unknown"
			for _, svc := range report.AuditResult.Services {
				if svc.Port == port {
					service = svc.Service
					break
				}
			}
			fmt.Fprintf(file, "Port %d: %s\n", port, service)
		}
		fmt.Fprintf(file, "\n")
	}

	// Database issues
	if len(report.AuditResult.DatabaseResults) > 0 {
		fmt.Fprintf(file, "Database Security Issues:\n")
		fmt.Fprintf(file, "------------------------\n")
		for _, db := range report.AuditResult.DatabaseResults {
			if db.Accessible {
				fmt.Fprintf(file, "%s on port %d (%s):\n", db.Service, db.Port, db.Severity)
				for _, warning := range db.Warnings {
					fmt.Fprintf(file, "  - %s\n", warning)
				}
				fmt.Fprintf(file, "\n")
			}
		}
	}

	// Web application issues
	if len(report.AuditResult.WebAppResults) > 0 {
		fmt.Fprintf(file, "Web Application Security Issues:\n")
		fmt.Fprintf(file, "-------------------------------\n")
		for _, web := range report.AuditResult.WebAppResults {
			if len(web.Warnings) > 0 {
				fmt.Fprintf(file, "%s (%s):\n", web.URL, web.Severity)
				for _, warning := range web.Warnings {
					fmt.Fprintf(file, "  - %s\n", warning)
				}
				fmt.Fprintf(file, "\n")
			}
		}
	}

	// System security issues
	if len(report.AuditResult.SystemResults) > 0 {
		fmt.Fprintf(file, "System Security Issues:\n")
		fmt.Fprintf(file, "----------------------\n")
		for _, sys := range report.AuditResult.SystemResults {
			if len(sys.Findings) > 0 {
				fmt.Fprintf(file, "%s (%s):\n", sys.CheckType, sys.Severity)
				for _, finding := range sys.Findings {
					fmt.Fprintf(file, "  - %s\n", finding)
				}
				fmt.Fprintf(file, "\n")
			}
		}
	}
}

// HTML Report Generation Functions

// HTMLReportData contains all data needed for HTML report generation
type HTMLReportData struct {
	Report      *Report
	GeneratedAt string
	Summary     *HTMLSummary
	Sections    []HTMLSection
}

// HTMLSummary contains summary statistics for the HTML report
type HTMLSummary struct {
	RiskLevel           string
	RiskColor          string
	TotalFindings      int
	CriticalFindings   int
	HighFindings       int
	MediumFindings     int
	LowFindings        int
}

// HTMLSection represents a section in the HTML report
type HTMLSection struct {
	Title    string
	Items    []HTMLItem
	HasItems bool
}

// HTMLItem represents an individual finding or result
type HTMLItem struct {
	Title       string
	Description string
	Severity    string
	SeverityColor string
	Details     []string
}

// GenerateHTMLReport creates a comprehensive HTML security report
func GenerateHTMLReport(auditReport *Report, outputPath string) error {
	data := prepareHTMLData(auditReport)
	
	tmpl := getHTMLTemplate()
	
	file, err := os.Create(outputPath + ".html")
	if err != nil {
		return fmt.Errorf("failed to create HTML file: %v", err)
	}
	defer file.Close()

	t, err := template.New("report").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	return t.Execute(file, data)
}

// prepareHTMLData converts the audit report into HTML-friendly format
func prepareHTMLData(auditReport *Report) *HTMLReportData {
	data := &HTMLReportData{
		Report:      auditReport,
		GeneratedAt: time.Now().Format("January 2, 2006 at 3:04 PM"),
		Summary:     createHTMLSummary(auditReport),
		Sections:    []HTMLSection{},
	}

	// Add port scan results section
	if len(auditReport.ScanResult.OpenPorts) > 0 {
		portSection := HTMLSection{
			Title:    "Open Ports & Services",
			Items:    []HTMLItem{},
			HasItems: true,
		}

		for _, port := range auditReport.ScanResult.OpenPorts {
			service := "Unknown"
			for _, svc := range auditReport.AuditResult.Services {
				if svc.Port == port {
					service = svc.Service
					break
				}
			}

			item := HTMLItem{
				Title:         fmt.Sprintf("Port %d", port),
				Description:   fmt.Sprintf("Service: %s", service),
				Severity:      "Info",
				SeverityColor: "#17a2b8",
				Details:       []string{},
			}
			portSection.Items = append(portSection.Items, item)
		}
		data.Sections = append(data.Sections, portSection)
	}

	// Add vulnerabilities section
	if len(auditReport.AuditResult.Vulnerabilities) > 0 {
		vulnSection := HTMLSection{
			Title:    "Vulnerabilities",
			Items:    []HTMLItem{},
			HasItems: true,
		}

		for _, vuln := range auditReport.AuditResult.Vulnerabilities {
			item := HTMLItem{
				Title:         vuln.CVE,
				Description:   vuln.Description,
				Severity:      vuln.Severity,
				SeverityColor: getSeverityColor(vuln.Severity),
				Details:       []string{fmt.Sprintf("CVSS Score: %.1f", vuln.Score)},
			}
			vulnSection.Items = append(vulnSection.Items, item)
		}
		data.Sections = append(data.Sections, vulnSection)
	}

	// Add database results section
	if len(auditReport.AuditResult.DatabaseResults) > 0 {
		dbSection := HTMLSection{
			Title:    "Database Security Issues",
			Items:    []HTMLItem{},
			HasItems: false,
		}

		for _, db := range auditReport.AuditResult.DatabaseResults {
			if db.Accessible {
				item := HTMLItem{
					Title:         fmt.Sprintf("%s Database", db.Service),
					Description:   fmt.Sprintf("Port %d - %s", db.Port, db.Severity),
					Severity:      db.Severity,
					SeverityColor: getSeverityColor(db.Severity),
					Details:       db.Warnings,
				}
				dbSection.Items = append(dbSection.Items, item)
				dbSection.HasItems = true
			}
		}
		
		if dbSection.HasItems {
			data.Sections = append(data.Sections, dbSection)
		}
	}

	return data
}

// createHTMLSummary generates summary statistics for the HTML report
func createHTMLSummary(auditReport *Report) *HTMLSummary {
	summary := &HTMLSummary{
		RiskLevel: auditReport.Summary.RiskLevel,
		RiskColor: getRiskColor(auditReport.Summary.RiskLevel),
	}

	// Count findings by severity
	for _, vuln := range auditReport.AuditResult.Vulnerabilities {
		summary.TotalFindings++
		switch vuln.Severity {
		case "Critical":
			summary.CriticalFindings++
		case "High":
			summary.HighFindings++
		case "Medium":
			summary.MediumFindings++
		case "Low":
			summary.LowFindings++
		}
	}

	for _, exploit := range auditReport.AuditResult.ExploitResults {
		if exploit.Success {
			summary.TotalFindings++
			switch exploit.Severity {
			case "Critical":
				summary.CriticalFindings++
			case "High":
				summary.HighFindings++
			case "Medium":
				summary.MediumFindings++
			case "Low":
				summary.LowFindings++
			}
		}
	}

	// Add database issues
	for _, db := range auditReport.AuditResult.DatabaseResults {
		if db.Accessible {
			summary.TotalFindings++
			switch db.Severity {
			case "Critical":
				summary.CriticalFindings++
			case "High":
				summary.HighFindings++
			case "Medium":
				summary.MediumFindings++
			case "Low":
				summary.LowFindings++
			}
		}
	}

	return summary
}

// getSeverityColor returns the appropriate color for a severity level
func getSeverityColor(severity string) string {
	switch severity {
	case "Critical":
		return "#dc3545" // Red
	case "High":
		return "#fd7e14" // Orange
	case "Medium":
		return "#ffc107" // Yellow
	case "Low":
		return "#28a745" // Green
	default:
		return "#17a2b8" // Blue (Info)
	}
}

// getRiskColor returns the appropriate color for a risk level
func getRiskColor(riskLevel string) string {
	switch riskLevel {
	case "Critical":
		return "#dc3545"
	case "High":
		return "#fd7e14"
	case "Medium":
		return "#ffc107"
	case "Low":
		return "#28a745"
	default:
		return "#6c757d" // Gray
	}
}

// getHTMLTemplate returns the HTML template for the security report
func getHTMLTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>InstaAudit Security Report - {{.Report.ScanResult.Host}}</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; margin: 0; padding: 20px; background-color: #f8f9fa; }
        .container { max-width: 1200px; margin: 0 auto; background: white; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; border-radius: 8px 8px 0 0; }
        .header h1 { margin: 0; font-size: 2.5em; }
        .header p { margin: 10px 0 0 0; opacity: 0.9; }
        .summary { padding: 30px; background: #f8f9fa; border-bottom: 1px solid #dee2e6; }
        .summary-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-top: 20px; }
        .summary-card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); text-align: center; }
        .summary-card h3 { margin: 0 0 10px 0; font-size: 2em; }
        .summary-card p { margin: 0; color: #6c757d; }
        .content { padding: 30px; }
        .section { margin-bottom: 40px; }
        .section h2 { color: #495057; border-bottom: 2px solid #e9ecef; padding-bottom: 10px; }
        .finding { background: white; border: 1px solid #dee2e6; border-radius: 8px; margin-bottom: 15px; overflow: hidden; }
        .finding-header { padding: 15px 20px; background: #f8f9fa; border-bottom: 1px solid #dee2e6; display: flex; justify-content: space-between; align-items: center; }
        .finding-title { font-weight: bold; flex-grow: 1; }
        .severity-badge { padding: 4px 12px; border-radius: 12px; color: white; font-size: 0.85em; font-weight: bold; }
        .finding-body { padding: 20px; }
        .finding-details { margin-top: 10px; }
        .finding-details ul { margin: 10px 0; padding-left: 20px; }
        .stats-row { display: flex; justify-content: space-around; margin-top: 20px; }
        .stat-item { text-align: center; }
        .stat-number { font-size: 2em; font-weight: bold; margin-bottom: 5px; }
        .footer { text-align: center; padding: 20px; color: #6c757d; border-top: 1px solid #dee2e6; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔒 InstaAudit Security Report</h1>
            <p>Target: {{.Report.ScanResult.Host}} | Generated: {{.GeneratedAt}}</p>
        </div>

        <div class="summary">
            <h2>Executive Summary</h2>
            <div class="stats-row">
                <div class="stat-item">
                    <div class="stat-number">{{.Report.Summary.OpenPorts}}</div>
                    <div>Open Ports</div>
                </div>
                <div class="stat-item">
                    <div class="stat-number">{{.Summary.TotalFindings}}</div>
                    <div>Total Findings</div>
                </div>
                <div class="stat-item">
                    <div class="stat-number">{{.Report.Summary.ServicesFound}}</div>
                    <div>Services</div>
                </div>
                <div class="stat-item">
                    <div class="stat-number" style="color: {{.Summary.RiskColor}};">{{.Summary.RiskLevel}}</div>
                    <div>Risk Level</div>
                </div>
            </div>
            
            <div class="summary-grid">
                <div class="summary-card">
                    <h3 style="color: #dc3545;">{{.Summary.CriticalFindings}}</h3>
                    <p>Critical Issues</p>
                </div>
                <div class="summary-card">
                    <h3 style="color: #fd7e14;">{{.Summary.HighFindings}}</h3>
                    <p>High Risk Issues</p>
                </div>
                <div class="summary-card">
                    <h3 style="color: #ffc107;">{{.Summary.MediumFindings}}</h3>
                    <p>Medium Risk Issues</p>
                </div>
                <div class="summary-card">
                    <h3 style="color: #28a745;">{{.Summary.LowFindings}}</h3>
                    <p>Low Risk Issues</p>
                </div>
            </div>
        </div>

        <div class="content">
            {{range .Sections}}
            {{if .HasItems}}
            <div class="section">
                <h2>{{.Title}}</h2>
                {{range .Items}}
                <div class="finding">
                    <div class="finding-header">
                        <div class="finding-title">{{.Title}}</div>
                        <span class="severity-badge" style="background-color: {{.SeverityColor}};">{{.Severity}}</span>
                    </div>
                    <div class="finding-body">
                        <p>{{.Description}}</p>
                        {{if .Details}}
                        <div class="finding-details">
                            <strong>Details:</strong>
                            <ul>
                                {{range .Details}}
                                <li>{{.}}</li>
                                {{end}}
                            </ul>
                        </div>
                        {{end}}
                    </div>
                </div>
                {{end}}
            </div>
            {{end}}
            {{end}}
        </div>

        <div class="footer">
            <p>Generated by InstaAudit v1.0 | {{.Report.Timestamp.Format "2006-01-02 15:04:05 UTC"}}</p>
        </div>
    </div>
</body>
</html>`
}