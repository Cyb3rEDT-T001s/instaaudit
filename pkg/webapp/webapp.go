package webapp

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// WebAppResult represents web application security check results
type WebAppResult struct {
	URL         string            `json:"url"`
	Host        string            `json:"host"`
	Port        int               `json:"port"`
	Accessible  bool              `json:"accessible"`
	Headers     map[string]string `json:"headers"`
	Warnings    []string          `json:"warnings"`
	TLSInfo     *TLSInfo          `json:"tls_info,omitempty"`
	Severity    string            `json:"severity"`
}

// TLSInfo contains TLS/SSL security information
type TLSInfo struct {
	Version     string   `json:"version"`
	CipherSuite string   `json:"cipher_suite"`
	Warnings    []string `json:"warnings"`
}

// CheckHTTPSecurity performs comprehensive HTTP/HTTPS security checks
func CheckHTTPSecurity(host string, port int) *WebAppResult {
	result := &WebAppResult{
		Host:     host,
		Port:     port,
		Headers:  make(map[string]string),
		Warnings: []string{},
		Severity: "Low",
	}

	// Determine protocol
	protocol := "http"
	if port == 443 || port == 8443 {
		protocol = "https"
	}
	
	result.URL = fmt.Sprintf("%s://%s:%d", protocol, host, port)

	// Create HTTP client with custom transport
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Get(result.URL)
	if err != nil {
		return result
	}
	defer resp.Body.Close()

	result.Accessible = true

	// Extract headers
	for name, values := range resp.Header {
		result.Headers[name] = strings.Join(values, ", ")
	}

	// Check for TLS if HTTPS
	if protocol == "https" && resp.TLS != nil {
		result.TLSInfo = checkTLSSecurity(resp.TLS)
		result.Warnings = append(result.Warnings, result.TLSInfo.Warnings...)
	}

	// Check security headers
	checkSecurityHeaders(result)

	// Check server information disclosure
	checkServerDisclosure(result)

	return result
}

// checkTLSSecurity analyzes TLS configuration
func checkTLSSecurity(tlsState *tls.ConnectionState) *TLSInfo {
	info := &TLSInfo{
		Warnings: []string{},
	}

	// Check TLS version
	switch tlsState.Version {
	case tls.VersionSSL30:
		info.Version = "SSL 3.0"
		info.Warnings = append(info.Warnings, "SSL 3.0 is deprecated and insecure")
	case tls.VersionTLS10:
		info.Version = "TLS 1.0"
		info.Warnings = append(info.Warnings, "TLS 1.0 is deprecated")
	case tls.VersionTLS11:
		info.Version = "TLS 1.1"
		info.Warnings = append(info.Warnings, "TLS 1.1 is deprecated")
	case tls.VersionTLS12:
		info.Version = "TLS 1.2"
	case tls.VersionTLS13:
		info.Version = "TLS 1.3"
	default:
		info.Version = "Unknown"
		info.Warnings = append(info.Warnings, "Unknown TLS version")
	}

	// Check cipher suite
	info.CipherSuite = tls.CipherSuiteName(tlsState.CipherSuite)
	if strings.Contains(strings.ToLower(info.CipherSuite), "rc4") {
		info.Warnings = append(info.Warnings, "Weak cipher suite detected (RC4)")
	}

	// Check certificate
	for _, cert := range tlsState.PeerCertificates {
		if time.Now().After(cert.NotAfter) {
			info.Warnings = append(info.Warnings, "Certificate has expired")
		}
		if time.Now().Add(30*24*time.Hour).After(cert.NotAfter) {
			info.Warnings = append(info.Warnings, "Certificate expires within 30 days")
		}
		if cert.Issuer.CommonName == cert.Subject.CommonName {
			info.Warnings = append(info.Warnings, "Self-signed certificate detected")
		}
	}

	return info
}

// checkSecurityHeaders validates HTTP security headers
func checkSecurityHeaders(result *WebAppResult) {
	requiredHeaders := map[string]string{
		"Content-Security-Policy": "Missing Content Security Policy",
		"X-Frame-Options":         "Missing X-Frame-Options (clickjacking protection)",
		"X-Content-Type-Options":  "Missing X-Content-Type-Options",
		"Strict-Transport-Security": "Missing HSTS header",
		"X-XSS-Protection":        "Missing XSS Protection header",
	}

	for header, warning := range requiredHeaders {
		if _, exists := result.Headers[header]; !exists {
			result.Warnings = append(result.Warnings, warning)
			result.Severity = "Medium"
		}
	}

	// Check for insecure header values
	if frameOptions, exists := result.Headers["X-Frame-Options"]; exists {
		if strings.ToLower(frameOptions) == "allowall" {
			result.Warnings = append(result.Warnings, "X-Frame-Options set to ALLOWALL (insecure)")
			result.Severity = "High"
		}
	}

	// Check for information disclosure headers
	disclosureHeaders := []string{"Server", "X-Powered-By", "X-AspNet-Version"}
	for _, header := range disclosureHeaders {
		if value, exists := result.Headers[header]; exists {
			result.Warnings = append(result.Warnings, 
				fmt.Sprintf("Information disclosure via %s header: %s", header, value))
		}
	}
}

// checkServerDisclosure checks for server information disclosure
func checkServerDisclosure(result *WebAppResult) {
	if server, exists := result.Headers["Server"]; exists {
		server = strings.ToLower(server)
		
		// Check for version disclosure
		if strings.Contains(server, "/") {
			result.Warnings = append(result.Warnings, "Server version disclosed in headers")
		}

		// Check for specific server types
		if strings.Contains(server, "apache") {
			if strings.Contains(server, "2.2") {
				result.Warnings = append(result.Warnings, "Outdated Apache version detected")
				result.Severity = "Medium"
			}
		}
		
		if strings.Contains(server, "nginx") {
			if strings.Contains(server, "1.1") {
				result.Warnings = append(result.Warnings, "Outdated Nginx version detected")
				result.Severity = "Medium"
			}
		}
	}
}

// CheckCommonVulnerabilities tests for common web vulnerabilities
func CheckCommonVulnerabilities(host string, port int) []string {
	var warnings []string
	
	protocol := "http"
	if port == 443 || port == 8443 {
		protocol = "https"
	}
	
	baseURL := fmt.Sprintf("%s://%s:%d", protocol, host, port)
	
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Test for directory traversal
	traversalPaths := []string{
		"../../../etc/passwd",
		"..\\..\\..\\windows\\system32\\drivers\\etc\\hosts",
		"....//....//....//etc/passwd",
	}

	for _, path := range traversalPaths {
		resp, err := client.Get(baseURL + "/" + path)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				warnings = append(warnings, fmt.Sprintf("Potential directory traversal vulnerability: %s", path))
				break
			}
		}
	}

	// Test for admin panels
	adminPaths := []string{"/admin", "/administrator", "/wp-admin", "/phpmyadmin", "/admin.php"}
	for _, path := range adminPaths {
		resp, err := client.Get(baseURL + path)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				warnings = append(warnings, fmt.Sprintf("Admin panel accessible: %s", path))
			}
		}
	}

	return warnings
}

// RunWebAppChecks performs all web application security checks
func RunWebAppChecks(host string, openPorts []int) []*WebAppResult {
	var results []*WebAppResult

	webPorts := []int{80, 443, 8080, 8443, 8000, 8888, 9090}
	
	for _, port := range openPorts {
		for _, webPort := range webPorts {
			if port == webPort {
				result := CheckHTTPSecurity(host, port)
				
				// Add common vulnerability checks
				vulnWarnings := CheckCommonVulnerabilities(host, port)
				result.Warnings = append(result.Warnings, vulnWarnings...)
				
				if len(vulnWarnings) > 0 {
					result.Severity = "High"
				}
				
				results = append(results, result)
				break
			}
		}
	}

	return results
}