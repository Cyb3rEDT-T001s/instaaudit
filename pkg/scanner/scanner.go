package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// ScanResult represents the result of a port scan
type ScanResult struct {
	Host      string       `json:"host"`
	Results   map[int]bool `json:"results"`
	OpenPorts []int        `json:"open_ports"`
}

// ScanPort checks if a specific port is open on a host
func ScanPort(host string, port int, timeout time.Duration) bool {
	target := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// ScanPorts performs concurrent port scanning
func ScanPorts(host string, ports []int, timeout time.Duration) *ScanResult {
	results := make(map[int]bool)
	var openPorts []int
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Limit concurrent goroutines to avoid overwhelming the target
	semaphore := make(chan struct{}, 100)

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			open := ScanPort(host, p, timeout)
			mutex.Lock()
			results[p] = open
			if open {
				openPorts = append(openPorts, p)
			}
			mutex.Unlock()
		}(port)
	}
	wg.Wait()

	return &ScanResult{
		Host:      host,
		Results:   results,
		OpenPorts: openPorts,
	}
}

// GetCommonPorts returns a list of commonly scanned ports
func GetCommonPorts() []int {
	return []int{
		21, 22, 23, 25, 53, 80, 110, 111, 135, 139, 143, 443, 993, 995,
		1723, 3306, 3389, 5432, 5900, 6379, 8080, 8443, 8888, 9090, 27017,
	}
}