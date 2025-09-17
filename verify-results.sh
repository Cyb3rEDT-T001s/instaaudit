#!/bin/bash

# InstaAudit Results Verification Script
# Cross-checks InstaAudit findings with other tools

echo "🔍 InstaAudit Results Verification Tool"
echo "======================================"

if [ $# -eq 0 ]; then
    echo "Usage: $0 <target_host> [instaaudit_report.json]"
    echo "Example: $0 example.com audit_report.json"
    exit 1
fi

TARGET=$1
REPORT=${2:-"audit_report.json"}

echo "Target: $TARGET"
echo "Report: $REPORT"
echo ""

# Check if required tools are installed
check_tool() {
    if ! command -v $1 &> /dev/null; then
        echo "⚠️  $1 not installed - skipping $2 verification"
        return 1
    fi
    return 0
}

echo "🔧 Checking available verification tools..."
NMAP_AVAILABLE=false
CURL_AVAILABLE=false
DIG_AVAILABLE=false

if check_tool "nmap" "port scanning"; then
    NMAP_AVAILABLE=true
fi

if check_tool "curl" "web services"; then
    CURL_AVAILABLE=true
fi

if check_tool "dig" "DNS resolution"; then
    DIG_AVAILABLE=true
fi

echo ""

# Extract open ports from InstaAudit report if available
if [ -f "$REPORT" ]; then
    echo "📊 Extracting findings from InstaAudit report..."
    
    # Extract open ports (basic JSON parsing)
    OPEN_PORTS=$(grep -o '"port":[0-9]*' "$REPORT" | cut -d':' -f2 | sort -n | uniq | tr '\n' ',' | sed 's/,$//')
    
    if [ ! -z "$OPEN_PORTS" ]; then
        echo "InstaAudit found open ports: $OPEN_PORTS"
    else
        echo "No open ports found in report or report format not recognized"
    fi
    echo ""
else
    echo "⚠️  InstaAudit report not found: $REPORT"
    echo "Proceeding with basic verification scan..."
    echo ""
fi

# DNS Resolution Check
if [ "$DIG_AVAILABLE" = true ]; then
    echo "🌐 DNS Resolution Verification"
    echo "-----------------------------"
    dig +short $TARGET
    echo ""
fi

# Port Scanning Verification
if [ "$NMAP_AVAILABLE" = true ]; then
    echo "🔍 Port Scanning Verification"
    echo "----------------------------"
    
    if [ ! -z "$OPEN_PORTS" ]; then
        echo "Verifying specific ports: $OPEN_PORTS"
        nmap -p $OPEN_PORTS $TARGET
    else
        echo "Running common ports scan..."
        nmap -F $TARGET
    fi
    echo ""
    
    echo "🔬 Service Detection Verification"
    echo "--------------------------------"
    if [ ! -z "$OPEN_PORTS" ]; then
        nmap -sV -p $OPEN_PORTS $TARGET
    else
        nmap -sV -F $TARGET
    fi
    echo ""
fi

# Web Service Verification
if [ "$CURL_AVAILABLE" = true ]; then
    echo "🌐 Web Service Verification"
    echo "--------------------------"
    
    # Check HTTP
    echo "Testing HTTP (port 80):"
    timeout 5 curl -I -s http://$TARGET 2>/dev/null | head -5 || echo "HTTP not accessible"
    echo ""
    
    # Check HTTPS
    echo "Testing HTTPS (port 443):"
    timeout 5 curl -I -s https://$TARGET 2>/dev/null | head -5 || echo "HTTPS not accessible"
    echo ""
    
    # Check security headers
    echo "Security Headers Check:"
    timeout 5 curl -I -s https://$TARGET 2>/dev/null | grep -i -E "(strict-transport|x-frame|x-content|x-xss)" || echo "No security headers found"
    echo ""
fi

# Database Service Quick Check
echo "🗄️  Database Service Quick Check"
echo "-------------------------------"

check_db_port() {
    local port=$1
    local service=$2
    
    if timeout 3 bash -c "</dev/tcp/$TARGET/$port" 2>/dev/null; then
        echo "⚠️  $service (port $port) is accessible!"
        return 0
    else
        echo "✅ $service (port $port) is not accessible"
        return 1
    fi
}

check_db_port 3306 "MySQL"
check_db_port 5432 "PostgreSQL" 
check_db_port 27017 "MongoDB"
check_db_port 6379 "Redis"
check_db_port 1433 "SQL Server"
echo ""

# SSL/TLS Quick Check
if command -v openssl &> /dev/null; then
    echo "🔒 SSL/TLS Quick Check"
    echo "---------------------"
    echo "Testing SSL certificate..."
    timeout 10 openssl s_client -connect $TARGET:443 -servername $TARGET </dev/null 2>/dev/null | grep -E "(subject=|issuer=|Verify return code)" || echo "SSL check failed or not available"
    echo ""
fi

# Summary and Recommendations
echo "📋 Verification Summary"
echo "======================"
echo ""
echo "✅ Cross-verification completed for: $TARGET"
echo ""
echo "🔍 What to do next:"
echo "1. Compare results with your InstaAudit report"
echo "2. Investigate any discrepancies"
echo "3. Use specialized tools for detailed analysis"
echo "4. Document all findings"
echo ""
echo "⚠️  Remember:"
echo "- Different tools may show different results due to timing"
echo "- Some services may be filtered by firewalls"
echo "- Always verify critical findings manually"
echo ""
echo "🛠️  For deeper analysis, consider:"
echo "- nmap -A $TARGET (aggressive scan)"
echo "- nikto -h http://$TARGET (web vulnerabilities)"
echo "- testssl.sh $TARGET (SSL analysis)"
echo "- dirb http://$TARGET (directory enumeration)"