package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SystemResult represents system-level security check results
type SystemResult struct {
	CheckType string   `json:"check_type"`
	Findings  []string `json:"findings"`
	Severity  string   `json:"severity"`
}

// CheckSUIDBinaries identifies SUID/SGID binaries
func CheckSUIDBinaries(rootDirs []string) *SystemResult {
	result := &SystemResult{
		CheckType: "SUID/SGID Binaries",
		Findings:  []string{},
		Severity:  "Low",
	}

	// Known risky SUID binaries
	riskyBinaries := []string{
		"nmap", "vim", "find", "bash", "sh", "more", "less", "nano", "cp", "mv",
		"awk", "man", "wget", "curl", "tar", "zip", "unzip", "python", "perl", "ruby",
	}

	for _, rootDir := range rootDirs {
		filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Continue on errors
			}

			mode := info.Mode()
			if mode&os.ModeSetuid != 0 {
				result.Findings = append(result.Findings, fmt.Sprintf("SUID binary: %s", path))
				
				// Check if it's a risky binary
				basename := filepath.Base(path)
				for _, risky := range riskyBinaries {
					if strings.Contains(basename, risky) {
						result.Findings = append(result.Findings, 
							fmt.Sprintf("RISKY SUID binary detected: %s", path))
						result.Severity = "High"
					}
				}
			}

			if mode&os.ModeSetgid != 0 {
				result.Findings = append(result.Findings, fmt.Sprintf("SGID binary: %s", path))
			}

			return nil
		})
	}

	if len(result.Findings) > 10 {
		result.Severity = "Medium"
	}

	return result
}

// CheckWorldWritableFiles identifies world-writable files and directories
func CheckWorldWritableFiles(rootDirs []string) *SystemResult {
	result := &SystemResult{
		CheckType: "World-Writable Files",
		Findings:  []string{},
		Severity:  "Low",
	}

	// Critical directories that should never be world-writable
	criticalDirs := []string{"/etc", "/bin", "/sbin", "/usr/bin", "/usr/sbin", "/boot"}

	for _, rootDir := range rootDirs {
		filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.Mode().Perm()&0002 != 0 {
				result.Findings = append(result.Findings, fmt.Sprintf("World-writable: %s", path))
				
				// Check if it's in a critical directory
				for _, criticalDir := range criticalDirs {
					if strings.HasPrefix(path, criticalDir) {
						result.Findings = append(result.Findings, 
							fmt.Sprintf("CRITICAL: World-writable file in system directory: %s", path))
						result.Severity = "Critical"
					}
				}
			}

			return nil
		})
	}

	if len(result.Findings) > 5 && result.Severity != "Critical" {
		result.Severity = "Medium"
	}

	return result
}

// CheckFilePermissions checks critical system file permissions
func CheckFilePermissions() *SystemResult {
	result := &SystemResult{
		CheckType: "Critical File Permissions",
		Findings:  []string{},
		Severity:  "Low",
	}

	// Critical files to check
	criticalFiles := map[string]os.FileMode{
		"/etc/passwd":  0644,
		"/etc/shadow":  0640,
		"/etc/group":   0644,
		"/etc/gshadow": 0640,
		"/etc/sudoers": 0440,
		"/etc/ssh/sshd_config": 0644,
	}

	for filePath, expectedMode := range criticalFiles {
		info, err := os.Stat(filePath)
		if err != nil {
			continue // File doesn't exist or can't access
		}

		actualMode := info.Mode().Perm()
		if actualMode != expectedMode {
			result.Findings = append(result.Findings, 
				fmt.Sprintf("Incorrect permissions on %s: %o (expected %o)", 
					filePath, actualMode, expectedMode))
			
			// Check for overly permissive files
			if actualMode&0022 != 0 && (filePath == "/etc/shadow" || filePath == "/etc/gshadow") {
				result.Severity = "Critical"
			} else if actualMode > expectedMode {
				result.Severity = "High"
			}
		}
	}

	return result
}

// CheckProcesses identifies suspicious or risky running processes
func CheckProcesses() *SystemResult {
	result := &SystemResult{
		CheckType: "Process Analysis",
		Findings:  []string{},
		Severity:  "Low",
	}

	// This is a simplified version - in a real implementation,
	// you would parse /proc/*/cmdline and /proc/*/stat
	
	// Check for common risky processes
	riskyProcesses := []string{
		"telnetd", "rsh", "rlogin", "ftp", "tftp", "finger", "rexec",
	}

	// Read /proc directory (Linux-specific)
	procDir := "/proc"
	entries, err := os.ReadDir(procDir)
	if err != nil {
		result.Findings = append(result.Findings, "Cannot access /proc directory (not Linux?)")
		return result
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Try to read cmdline
		cmdlinePath := filepath.Join(procDir, entry.Name(), "cmdline")
		cmdlineData, err := os.ReadFile(cmdlinePath)
		if err != nil {
			continue
		}

		cmdline := string(cmdlineData)
		for _, risky := range riskyProcesses {
			if strings.Contains(cmdline, risky) {
				result.Findings = append(result.Findings, 
					fmt.Sprintf("Risky process detected: %s (PID: %s)", risky, entry.Name()))
				result.Severity = "Medium"
			}
		}
	}

	return result
}

// CheckSystemConfiguration performs various system configuration checks
func CheckSystemConfiguration() *SystemResult {
	result := &SystemResult{
		CheckType: "System Configuration",
		Findings:  []string{},
		Severity:  "Low",
	}

	// Check if running as root
	if os.Geteuid() == 0 {
		result.Findings = append(result.Findings, "Running with root privileges")
		result.Severity = "Medium"
	}

	// Check umask (Windows doesn't have umask, so skip this check)
	// This is a Unix/Linux specific check
	result.Findings = append(result.Findings, "Umask check skipped (Windows system)")

	// Check for core dumps enabled
	corePattern := "/proc/sys/kernel/core_pattern"
	if data, err := os.ReadFile(corePattern); err == nil {
		pattern := strings.TrimSpace(string(data))
		if pattern != "" && pattern != "core" {
			result.Findings = append(result.Findings, 
				fmt.Sprintf("Core dumps may be enabled: %s", pattern))
		}
	}

	return result
}

// RunSystemChecks performs all system-level security checks
func RunSystemChecks() []*SystemResult {
	var results []*SystemResult

	// Define directories to scan (adjust based on OS)
	scanDirs := []string{"/bin", "/sbin", "/usr/bin", "/usr/sbin", "/usr/local/bin"}
	
	// Only scan if directories exist (Linux/Unix systems)
	var existingDirs []string
	for _, dir := range scanDirs {
		if _, err := os.Stat(dir); err == nil {
			existingDirs = append(existingDirs, dir)
		}
	}

	if len(existingDirs) > 0 {
		results = append(results, CheckSUIDBinaries(existingDirs))
		results = append(results, CheckWorldWritableFiles(existingDirs))
	}

	results = append(results, CheckFilePermissions())
	results = append(results, CheckProcesses())
	results = append(results, CheckSystemConfiguration())

	return results
}