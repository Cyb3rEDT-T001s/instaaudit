package database

import (
	"database/sql"
	"fmt"
	"net"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// DatabaseResult represents database security check results
type DatabaseResult struct {
	Service     string   `json:"service"`
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	Accessible  bool     `json:"accessible"`
	Warnings    []string `json:"warnings"`
	Severity    string   `json:"severity"`
}

// CheckMySQL performs comprehensive MySQL security checks
func CheckMySQL(host string, port int) *DatabaseResult {
	result := &DatabaseResult{
		Service:  "MySQL",
		Host:     host,
		Port:     port,
		Warnings: []string{},
		Severity: "Low",
	}

	// Test common default credentials
	defaultCreds := [][]string{
		{"root", ""},
		{"root", "root"},
		{"root", "password"},
		{"admin", "admin"},
		{"mysql", "mysql"},
		{"test", "test"},
	}

	for _, cred := range defaultCreds {
		username, password := cred[0], cred[1]
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", username, password, host, port)
		
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			continue
		}
		
		db.SetConnMaxLifetime(3 * time.Second)
		if err := db.Ping(); err == nil {
			result.Accessible = true
			result.Warnings = append(result.Warnings, 
				fmt.Sprintf("MySQL accessible with default credentials: %s:%s", username, password))
			result.Severity = "Critical"
			
			// Additional checks if we can connect
			checkMySQLSecurity(db, result)
			db.Close()
			break
		}
		db.Close()
	}

	// Test for unauthenticated access
	if !result.Accessible {
		dsn := fmt.Sprintf("@tcp(%s:%d)/", host, port)
		db, err := sql.Open("mysql", dsn)
		if err == nil {
			db.SetConnMaxLifetime(3 * time.Second)
			if err := db.Ping(); err == nil {
				result.Accessible = true
				result.Warnings = append(result.Warnings, "MySQL allows unauthenticated access")
				result.Severity = "Critical"
			}
			db.Close()
		}
	}

	return result
}

// checkMySQLSecurity performs additional security checks when connected
func checkMySQLSecurity(db *sql.DB, result *DatabaseResult) {
	// Check for anonymous users
	var anonymousUsers int
	err := db.QueryRow("SELECT COUNT(*) FROM mysql.user WHERE user = ''").Scan(&anonymousUsers)
	if err == nil && anonymousUsers > 0 {
		result.Warnings = append(result.Warnings, "Anonymous users found in MySQL")
		result.Severity = "High"
	}

	// Check for users with empty passwords
	var emptyPassUsers int
	err = db.QueryRow("SELECT COUNT(*) FROM mysql.user WHERE authentication_string = '' OR password = ''").Scan(&emptyPassUsers)
	if err == nil && emptyPassUsers > 0 {
		result.Warnings = append(result.Warnings, "Users with empty passwords found")
		result.Severity = "Critical"
	}

	// Check MySQL version
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err == nil {
		if strings.Contains(version, "5.5") || strings.Contains(version, "5.6") {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Outdated MySQL version: %s", version))
			result.Severity = "Medium"
		}
	}
}

// CheckPostgreSQL performs PostgreSQL security checks
func CheckPostgreSQL(host string, port int) *DatabaseResult {
	result := &DatabaseResult{
		Service:  "PostgreSQL",
		Host:     host,
		Port:     port,
		Warnings: []string{},
		Severity: "Low",
	}

	// Test common default credentials
	defaultCreds := [][]string{
		{"postgres", ""},
		{"postgres", "postgres"},
		{"postgres", "password"},
		{"admin", "admin"},
	}

	for _, cred := range defaultCreds {
		username, password := cred[0], cred[1]
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable connect_timeout=3", 
			host, port, username, password)
		
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			continue
		}
		
		if err := db.Ping(); err == nil {
			result.Accessible = true
			result.Warnings = append(result.Warnings, 
				fmt.Sprintf("PostgreSQL accessible with default credentials: %s:%s", username, password))
			result.Severity = "Critical"
			
			checkPostgreSQLSecurity(db, result)
			db.Close()
			break
		}
		db.Close()
	}

	return result
}

// checkPostgreSQLSecurity performs additional PostgreSQL security checks
func checkPostgreSQLSecurity(db *sql.DB, result *DatabaseResult) {
	// Check PostgreSQL version
	var version string
	err := db.QueryRow("SELECT version()").Scan(&version)
	if err == nil {
		if strings.Contains(version, "9.") {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Outdated PostgreSQL version: %s", version))
			result.Severity = "Medium"
		}
	}

	// Check for superuser accounts
	var superUsers int
	err = db.QueryRow("SELECT COUNT(*) FROM pg_user WHERE usesuper = true").Scan(&superUsers)
	if err == nil && superUsers > 1 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Multiple superuser accounts found: %d", superUsers))
		result.Severity = "Medium"
	}
}

// CheckMongoDB performs MongoDB security checks
func CheckMongoDB(host string, port int) *DatabaseResult {
	result := &DatabaseResult{
		Service:  "MongoDB",
		Host:     host,
		Port:     port,
		Warnings: []string{},
		Severity: "Low",
	}

	// Test for unauthenticated access
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, fmt.Sprintf("%d", port)), 3*time.Second)
	if err != nil {
		return result
	}
	defer conn.Close()

	// Send MongoDB handshake
	handshake := []byte{
		0x3a, 0x00, 0x00, 0x00, // message length
		0x01, 0x00, 0x00, 0x00, // request id
		0x00, 0x00, 0x00, 0x00, // response to
		0xd4, 0x07, 0x00, 0x00, // opcode (query)
	}

	conn.Write(handshake)
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	n, err := conn.Read(buffer)

	if err == nil && n > 0 {
		result.Accessible = true
		result.Warnings = append(result.Warnings, "MongoDB appears to allow unauthenticated connections")
		result.Severity = "Critical"
	}

	return result
}

// CheckRedis performs Redis security checks
func CheckRedis(host string, port int) *DatabaseResult {
	result := &DatabaseResult{
		Service:  "Redis",
		Host:     host,
		Port:     port,
		Warnings: []string{},
		Severity: "Low",
	}

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, fmt.Sprintf("%d", port)), 3*time.Second)
	if err != nil {
		return result
	}
	defer conn.Close()

	// Test Redis INFO command (should work without auth if misconfigured)
	conn.Write([]byte("*1\r\n$4\r\nINFO\r\n"))
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	n, err := conn.Read(buffer)

	if err == nil && n > 0 {
		response := string(buffer[:n])
		if strings.Contains(response, "redis_version") {
			result.Accessible = true
			result.Warnings = append(result.Warnings, "Redis accessible without authentication")
			result.Severity = "Critical"

			// Check for dangerous commands
			if strings.Contains(response, "flushall") || strings.Contains(response, "flushdb") {
				result.Warnings = append(result.Warnings, "Dangerous Redis commands are enabled")
			}
		}
	}

	return result
}

// RunDatabaseChecks performs all database security checks
func RunDatabaseChecks(host string, openPorts []int) []*DatabaseResult {
	var results []*DatabaseResult

	for _, port := range openPorts {
		switch port {
		case 3306:
			results = append(results, CheckMySQL(host, port))
		case 5432:
			results = append(results, CheckPostgreSQL(host, port))
		case 27017, 27018, 27019:
			results = append(results, CheckMongoDB(host, port))
		case 6379:
			results = append(results, CheckRedis(host, port))
		}
	}

	return results
}