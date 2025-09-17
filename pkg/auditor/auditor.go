package auditor

import (
	"fmt"
	"strings"
	"time"

	"github.com/Cyb3rEDT-T001s/instaaudit/pkg/database"
	"github.com/Cyb3rEDT-T001s/instaaudit/pkg/exploits"
	"github.com/Cyb3rEDT-T001s/instaaudit/pkg/recon"
	"github.com/Cyb3rEDT-T001s/instaaudit/pkg/system"
	"github.com/Cyb3rEDT-T001s/instaaudit/pkg/webapp"
)

// Vulnerability represents a security vulnerability
type Vulnerability struct {
	CVE         string  `json:"cve"`
	Description string  `json:"description"`
	Severity    string  `json:"severity"`
	Score       float64 `json:"score"`
}

// ServiceInfo represents information about a detected service
type ServiceInfo struct {
	Port    int    `json:"port"`
	Service string `json:"service"`
	Version string `json:"version,omitempty"`
}

// AuditResult contains the results of security auditing
type AuditResult struct {
	Host            string                     `json:"host"`
	Services        []ServiceInfo              `json:"services"`
	Vulnerabilities []Vulnerability            `json:"vulnerabilities"`
	Misconfigs      []string                   `json:"misconfigurations"`
	ExploitResults  []*exploits.ExploitResult  `json:"exploit_results"`
	ReconData       *recon.ReconResult         `json:"reconnaissance"`
	DatabaseResults []*database.DatabaseResult `json:"database_results"`
	WebAppResults   []*webapp.WebAppResult     `json:"webapp_results"`
	SystemResults   []*system.SystemResult     `json:"system_results"`
}

// IdentifyService attempts to identify the service running on a port
func IdentifyService(port int) string {
	serviceMap := map[int]string{
		21:    "FTP",
		22:    "SSH",
		23:    "Telnet",
		25:    "SMTP",
		53:    "DNS",
		80:    "HTTP",
		110:   "POP3",
		143:   "IMAP",
		443:   "HTTPS",
		993:   "IMAPS",
		995:   "POP3S",
		3306:  "MySQL",
		3389:  "RDP",
		5432:  "PostgreSQL",
		5900:  "VNC",
		6379:  "Redis",
		8080:  "HTTP-Alt",
		8443:  "HTTPS-Alt",
		27017: "MongoDB",
	}

	if service, exists := serviceMap[port]; exists {
		return service
	}
	return "Unknown"
}

// CheckBasicMisconfigurations performs basic security checks
func CheckBasicMisconfigurations(services []ServiceInfo) []string {
	var misconfigs []string

	for _, service := range services {
		switch service.Service {
		case "FTP":
			misconfigs = append(misconfigs, "FTP service detected - consider using SFTP instead")
		case "Telnet":
			misconfigs = append(misconfigs, "Telnet service detected - unencrypted protocol, use SSH instead")
		case "HTTP":
			if service.Port == 80 {
				misconfigs = append(misconfigs, "HTTP service on port 80 - consider redirecting to HTTPS")
			}
		case "SSH":
			misconfigs = append(misconfigs, "SSH service detected - ensure key-based authentication is enabled")
		case "MySQL", "PostgreSQL", "MongoDB", "Redis":
			misconfigs = append(misconfigs, fmt.Sprintf("%s database detected - ensure proper authentication and access controls", service.Service))
		}
	}

	return misconfigs
}

// RealVulnerabilityCheck performs actual vulnerability assessment based on service detection
func RealVulnerabilityCheck(service, version string) []Vulnerability {
	var vulnerabilities []Vulnerability
	
	// Only report real vulnerabilities based on actual service analysis
	// This function now only returns vulnerabilities when real issues are detected
	
	// Example: Check for known vulnerable versions (you can expand this)
	if strings.Contains(strings.ToLower(version), "apache/2.2") {
		vulnerabilities = append(vulnerabilities, Vulnerability{
			CVE:         "CVE-2017-15715",
			Description: "Apache HTTP Server 2.2.x vulnerability - Expression injection in mod_rewrite",
			Severity:    "High",
			Score:       8.1,
		})
	}
	
	if strings.Contains(strings.ToLower(version), "nginx/1.1") {
		vulnerabilities = append(vulnerabilities, Vulnerability{
			CVE:         "CVE-2013-2028",
			Description: "Nginx 1.1.x buffer overflow vulnerability",
			Severity:    "High", 
			Score:       7.5,
		})
	}
	
	// Add more real vulnerability checks here based on actual service versions
	// For now, return empty array to avoid false positives
	return vulnerabilities
}

// PerformAudit conducts a comprehensive security audit
func PerformAudit(host string, openPorts []int) *AuditResult {
	var services []ServiceInfo
	var allVulns []Vulnerability
	var allExploitResults []*exploits.ExploitResult

	// Perform reconnaissance first
	reconData := recon.PerformReconnaissance(host, openPorts)

	// Identify services and perform exploit checks
	for _, port := range openPorts {
		service := IdentifyService(port)
		
		// Try to get more detailed service info using banner grabbing
		if banner, err := exploits.BannerGrabber(host, port, 3*time.Second); err == nil && banner != "" {
			// Extract version info from banner if possible
			service = fmt.Sprintf("%s (%s)", service, banner[:min(50, len(banner))])
		}
		
		serviceInfo := ServiceInfo{
			Port:    port,
			Service: service,
		}
		services = append(services, serviceInfo)

		// Check for real vulnerabilities based on service detection
		vulns := RealVulnerabilityCheck(IdentifyService(port), service)
		allVulns = append(allVulns, vulns...)

		// Run exploit checks
		exploitResults := exploits.RunExploitChecks(host, port, IdentifyService(port))
		allExploitResults = append(allExploitResults, exploitResults...)
	}

	// Check for misconfigurations
	misconfigs := CheckBasicMisconfigurations(services)

	// Perform database security checks
	databaseResults := database.RunDatabaseChecks(host, openPorts)

	// Perform web application security checks
	webAppResults := webapp.RunWebAppChecks(host, openPorts)

	// Perform system-level security checks
	systemResults := system.RunSystemChecks()

	return &AuditResult{
		Host:            host,
		Services:        services,
		Vulnerabilities: allVulns,
		Misconfigs:      misconfigs,
		ExploitResults:  allExploitResults,
		ReconData:       reconData,
		DatabaseResults: databaseResults,
		WebAppResults:   webAppResults,
		SystemResults:   systemResults,
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}