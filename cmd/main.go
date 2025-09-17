package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/Cyb3rEDT-T001s/instaaudit/pkg/auditor"
	"github.com/Cyb3rEDT-T001s/instaaudit/pkg/config"
	"github.com/Cyb3rEDT-T001s/instaaudit/pkg/report"
	"github.com/Cyb3rEDT-T001s/instaaudit/pkg/scanner"
)

var (
	host         string
	ports        string
	timeout      int
	outputPath   string
	outputFormat string
	configFile   string
	skipExploits bool
	skipRecon    bool
	aggressive   bool
)

var rootCmd = &cobra.Command{
	Use:   "instaaudit",
	Short: "InstaAudit - Advanced Security Auditing Tool",
	Long: `InstaAudit is a comprehensive, fast security auditing tool that performs:
- Port scanning with service identification
- Vulnerability assessment and exploit testing
- Database security auditing
- Web application security testing
- System-level security checks
- Reconnaissance and information gathering
- Multi-format report generation (JSON, CSV, HTML, TXT)`,
	Run: runAudit,
}

func init() {
	rootCmd.Flags().StringVarP(&host, "host", "H", "", "Target host to audit (required)")
	rootCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan (e.g., '80,443,22' or 'common' for common ports)")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 2, "Connection timeout in seconds")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "./audit_report", "Output file path (without extension)")
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "Output format (json, csv, html, txt, all)")
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file path")
	rootCmd.Flags().BoolVar(&skipExploits, "skip-exploits", false, "Skip exploit testing")
	rootCmd.Flags().BoolVar(&skipRecon, "skip-recon", false, "Skip reconnaissance")
	rootCmd.Flags().BoolVarP(&aggressive, "aggressive", "A", false, "Enable aggressive scanning (more thorough but slower)")

	rootCmd.MarkFlagRequired("host")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runAudit(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Printf("Warning: Could not load config file: %v", err)
		cfg = config.DefaultConfig()
	}

	// Override config with command line flags
	if timeout > 0 {
		cfg.Timeout = time.Duration(timeout) * time.Second
	}
	if outputPath != "./audit_report" {
		cfg.OutputPath = outputPath
	}
	if outputFormat != "json" {
		cfg.OutputFormat = outputFormat
	}
	cfg.SkipExploits = skipExploits
	cfg.SkipRecon = skipRecon
	cfg.Aggressive = aggressive

	// Parse ports
	var targetPorts []int
	if ports == "" || ports == "common" {
		targetPorts = scanner.GetCommonPorts()
	} else {
		targetPorts = parsePorts(ports)
	}

	fmt.Printf("🔍 Starting InstaAudit Enhanced Security Assessment\n")
	fmt.Printf("Target: %s\n", host)
	fmt.Printf("Ports: %d total\n", len(targetPorts))
	fmt.Printf("Timeout: %v\n", cfg.Timeout)
	if cfg.Aggressive {
		fmt.Printf("Mode: Aggressive (comprehensive checks enabled)\n")
	}
	fmt.Println(strings.Repeat("=", 60))

	// Perform port scan
	fmt.Printf("🔎 Scanning %d ports...\n", len(targetPorts))
	scanResult := scanner.ScanPorts(host, targetPorts, cfg.Timeout)
	fmt.Printf("✅ Port scan complete. Found %d open ports: %v\n", len(scanResult.OpenPorts), scanResult.OpenPorts)

	// Perform comprehensive security audit
	fmt.Println("🛡️  Performing comprehensive security audit...")
	
	if cfg.SkipExploits {
		fmt.Println("⏭️  Skipping exploit tests...")
	}
	if cfg.SkipRecon {
		fmt.Println("⏭️  Skipping reconnaissance...")
	}
	if cfg.Aggressive {
		fmt.Println("⚡ Aggressive mode: Running enhanced security checks...")
	}

	auditResult := auditor.PerformAudit(host, scanResult.OpenPorts)

	// Generate comprehensive report
	fmt.Println("📊 Generating comprehensive report...")
	auditReport := report.GenerateReport(scanResult, auditResult)

	// Save report in requested format(s)
	if err := saveReport(auditReport, cfg.OutputPath, cfg.OutputFormat); err != nil {
		log.Fatalf("Failed to save report: %v", err)
	}

	// Always generate educational report for beginners
	fmt.Println("📚 Generating educational report for beginners...")
	if err := report.GenerateEducationalReport(auditReport, cfg.OutputPath); err != nil {
		log.Printf("Warning: Could not generate educational report: %v", err)
	}

	// Print executive summary
	printExecutiveSummary(auditReport)
}

func parsePorts(portStr string) []int {
	var ports []int
	parts := strings.Split(portStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			// Handle range (e.g., "80-90")
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) == 2 {
				start, err1 := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
				end, err2 := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
				if err1 == nil && err2 == nil {
					for i := start; i <= end; i++ {
						ports = append(ports, i)
					}
				}
			}
		} else {
			// Handle single port
			if port, err := strconv.Atoi(part); err == nil {
				ports = append(ports, port)
			}
		}
	}

	return ports
}

func saveReport(auditReport *report.Report, outputPath, format string) error {
	switch format {
	case "json":
		return report.SaveAsJSON(auditReport, outputPath)
	case "csv":
		return report.SaveAsCSV(auditReport, outputPath)
	case "html":
		return report.SaveAsHTML(auditReport, outputPath)
	case "txt":
		return report.SaveAsText(auditReport, outputPath)
	case "all":
		if err := report.SaveAsJSON(auditReport, outputPath); err != nil {
			return err
		}
		if err := report.SaveAsCSV(auditReport, outputPath); err != nil {
			return err
		}
		if err := report.SaveAsHTML(auditReport, outputPath); err != nil {
			return err
		}
		return report.SaveAsText(auditReport, outputPath)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

func printExecutiveSummary(auditReport *report.Report) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("🔒 INSTAAUDIT EXECUTIVE SECURITY SUMMARY")
	fmt.Println(strings.Repeat("=", 70))
	
	// Risk level with color coding
	riskEmoji := getRiskEmoji(auditReport.Summary.RiskLevel)
	fmt.Printf("🎯 Target: %s\n", auditReport.ScanResult.Host)
	fmt.Printf("📅 Scan Date: %s\n", auditReport.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("%s Overall Risk Level: %s\n", riskEmoji, auditReport.Summary.RiskLevel)
	
	fmt.Println("\n📊 SCAN STATISTICS:")
	fmt.Printf("   • Total Ports Scanned: %d\n", auditReport.Summary.TotalPorts)
	fmt.Printf("   • Open Ports Found: %d\n", auditReport.Summary.OpenPorts)
	fmt.Printf("   • Services Identified: %d\n", auditReport.Summary.ServicesFound)
	fmt.Printf("   • Technologies Detected: %d\n", auditReport.Summary.TechnologiesFound)
	fmt.Printf("   • Subdomains Found: %d\n", auditReport.Summary.SubdomainsFound)

	fmt.Println("\n🚨 SECURITY FINDINGS:")
	fmt.Printf("   • Vulnerabilities: %d\n", auditReport.Summary.Vulnerabilities)
	fmt.Printf("   • Exploitable Services: %d\n", auditReport.Summary.ExploitableServices)
	fmt.Printf("   • Database Issues: %d\n", auditReport.Summary.DatabaseIssues)
	fmt.Printf("   • Web App Issues: %d\n", auditReport.Summary.WebAppIssues)
	fmt.Printf("   • System Issues: %d\n", auditReport.Summary.SystemIssues)
	fmt.Printf("   • Misconfigurations: %d\n", auditReport.Summary.Misconfigurations)

	// Show critical findings
	if auditReport.Summary.RiskLevel == "Critical" || auditReport.Summary.RiskLevel == "High" {
		fmt.Println("\n⚠️  CRITICAL ISSUES DETECTED:")
		
		// Database issues
		for _, db := range auditReport.AuditResult.DatabaseResults {
			if db.Accessible && (db.Severity == "Critical" || db.Severity == "High") {
				fmt.Printf("   🔴 %s database accessible on port %d\n", db.Service, db.Port)
			}
		}
		
		// Successful exploits
		for _, exploit := range auditReport.AuditResult.ExploitResults {
			if exploit.Success && (exploit.Severity == "Critical" || exploit.Severity == "High") {
				fmt.Printf("   🔴 %s vulnerability on port %d\n", exploit.ExploitName, exploit.Port)
			}
		}
	}

	// Show open ports
	if len(auditReport.ScanResult.OpenPorts) > 0 {
		fmt.Printf("\n🔓 Open Ports: %v\n", auditReport.ScanResult.OpenPorts)
	}

	// Show detected technologies
	if auditReport.AuditResult.ReconData != nil && len(auditReport.AuditResult.ReconData.Technologies) > 0 {
		fmt.Printf("🔧 Technologies: %v\n", auditReport.AuditResult.ReconData.Technologies)
	}

	fmt.Printf("\n📄 Reports saved to: %s.*\n", outputPath)
	fmt.Println(strings.Repeat("=", 70))
	
	// Recommendations
	if auditReport.Summary.RiskLevel == "Critical" || auditReport.Summary.RiskLevel == "High" {
		fmt.Println("🚨 IMMEDIATE ACTION REQUIRED:")
		fmt.Println("   • Review and secure database access")
		fmt.Println("   • Patch identified vulnerabilities")
		fmt.Println("   • Implement proper access controls")
		fmt.Println("   • Consider network segmentation")
	}
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