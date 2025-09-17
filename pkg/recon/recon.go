package recon

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// ReconResult represents reconnaissance information
type ReconResult struct {
	Host          string            `json:"host"`
	IPAddresses   []string          `json:"ip_addresses"`
	Subdomains    []string          `json:"subdomains"`
	Technologies  []string          `json:"technologies"`
	Headers       map[string]string `json:"headers"`
	OSFingerprint string            `json:"os_fingerprint"`
}

// DNSRecon performs DNS reconnaissance
func DNSRecon(host string) *ReconResult {
	result := &ReconResult{
		Host:         host,
		IPAddresses:  []string{},
		Subdomains:   []string{},
		Technologies: []string{},
		Headers:      make(map[string]string),
	}

	// Resolve IP addresses
	ips, err := net.LookupIP(host)
	if err == nil {
		for _, ip := range ips {
			result.IPAddresses = append(result.IPAddresses, ip.String())
		}
	}

	// Check for common subdomains
	commonSubdomains := []string{
		"www", "mail", "ftp", "admin", "api", "dev", "test", "staging",
		"blog", "shop", "support", "help", "docs", "cdn", "static",
	}

	for _, subdomain := range commonSubdomains {
		fullDomain := subdomain + "." + host
		_, err := net.LookupIP(fullDomain)
		if err == nil {
			result.Subdomains = append(result.Subdomains, fullDomain)
		}
	}

	return result
}

// HTTPRecon performs HTTP-based reconnaissance
func HTTPRecon(host string, port int) map[string]string {
	headers := make(map[string]string)

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, fmt.Sprintf("%d", port)), 5*time.Second)
	if err != nil {
		return headers
	}
	defer conn.Close()

	// Send HTTP request
	request := fmt.Sprintf("GET / HTTP/1.1\r\nHost: %s\r\nUser-Agent: InstaAudit/1.0\r\n\r\n", host)
	conn.Write([]byte(request))

	// Read response
	buffer := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return headers
	}

	response := string(buffer[:n])
	lines := strings.Split(response, "\r\n")

	// Parse headers
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				headers[key] = value
			}
		}
	}

	return headers
}

// DetectTechnologies identifies web technologies from headers and content
func DetectTechnologies(headers map[string]string, content string) []string {
	var technologies []string

	// Check headers for technology indicators
	if server, exists := headers["Server"]; exists {
		server = strings.ToLower(server)
		if strings.Contains(server, "apache") {
			technologies = append(technologies, "Apache")
		}
		if strings.Contains(server, "nginx") {
			technologies = append(technologies, "Nginx")
		}
		if strings.Contains(server, "iis") {
			technologies = append(technologies, "IIS")
		}
		if strings.Contains(server, "cloudflare") {
			technologies = append(technologies, "Cloudflare")
		}
	}

	if xPoweredBy, exists := headers["X-Powered-By"]; exists {
		xPoweredBy = strings.ToLower(xPoweredBy)
		if strings.Contains(xPoweredBy, "php") {
			technologies = append(technologies, "PHP")
		}
		if strings.Contains(xPoweredBy, "asp.net") {
			technologies = append(technologies, "ASP.NET")
		}
	}

	// Check content for technology indicators
	content = strings.ToLower(content)
	if strings.Contains(content, "wordpress") {
		technologies = append(technologies, "WordPress")
	}
	if strings.Contains(content, "drupal") {
		technologies = append(technologies, "Drupal")
	}
	if strings.Contains(content, "joomla") {
		technologies = append(technologies, "Joomla")
	}

	return technologies
}

// OSFingerprinting attempts basic OS detection
func OSFingerprinting(host string) string {
	// This is a simplified OS fingerprinting - real tools use TCP/IP stack analysis

	// Try to connect and analyze TTL values
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:80", host), 3*time.Second)
	if err != nil {
		// Try ping-like approach with different ports
		for _, port := range []int{22, 23, 80, 443} {
			conn, err = net.DialTimeout("tcp", net.JoinHostPort(host, fmt.Sprintf("%d", port)), 1*time.Second)
			if err == nil {
				conn.Close()
				break
			}
		}
		if err != nil {
			return "Unknown"
		}
	} else {
		conn.Close()
	}

	// In a real implementation, you would:
	// 1. Analyze TCP window sizes
	// 2. Check TCP options
	// 3. Analyze ICMP responses
	// 4. Check service banners

	// For demo purposes, return generic detection
	return "Linux/Unix-like (estimated)"
}

// PerformReconnaissance conducts comprehensive reconnaissance
func PerformReconnaissance(host string, openPorts []int) *ReconResult {
	// Start with DNS reconnaissance
	result := DNSRecon(host)

	// Perform OS fingerprinting
	result.OSFingerprint = OSFingerprinting(host)

	// Check HTTP services for additional information
	for _, port := range openPorts {
		if port == 80 || port == 8080 || port == 443 || port == 8443 {
			headers := HTTPRecon(host, port)

			// Merge headers (in case multiple HTTP services)
			for k, v := range headers {
				result.Headers[fmt.Sprintf("Port%d-%s", port, k)] = v
			}

			// Get content for technology detection
			content := ""
			if contentType, exists := headers["Content-Type"]; exists {
				content = contentType
			}

			technologies := DetectTechnologies(headers, content)
			result.Technologies = append(result.Technologies, technologies...)
		}
	}

	// Remove duplicate technologies
	result.Technologies = removeDuplicates(result.Technologies)

	return result
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}