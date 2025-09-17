#!/bin/bash

# Manual Verification Guide for InstaAudit Results
# Helps beginners verify findings step-by-step

echo "🔍 InstaAudit Manual Verification Helper"
echo "======================================="

if [ $# -eq 0 ]; then
    echo "Usage: $0 <target_host>"
    echo "Example: $0 example.com"
    exit 1
fi

TARGET=$1
echo "Target: $TARGET"
echo ""

echo "🎓 This script will guide you through manually verifying InstaAudit results."
echo "Each test shows you HOW to verify findings yourself."
echo ""

# Basic connectivity test
echo "1️⃣  BASIC CONNECTIVITY TEST"
echo "============================"
echo "Testing if we can reach $TARGET..."
if ping -c 3 $TARGET >/dev/null 2>&1; then
    echo "✅ $TARGET is reachable"
else
    echo "❌ $TARGET is not reachable - check the hostname/IP"
    echo "   This might explain why InstaAudit had issues"
fi
echo ""

# DNS resolution
echo "2️⃣  DNS RESOLUTION TEST"
echo "======================"
echo "Checking if $TARGET resolves to an IP address..."
IP=$(nslookup $TARGET 2>/dev/null | grep -A1 "Name:" | tail -1 | awk '{print $2}')
if [ ! -z "$IP" ]; then
    echo "✅ $TARGET resolves to: $IP"
else
    echo "❌ DNS resolution failed for $TARGET"
fi
echo ""

# Common port tests
echo "3️⃣  COMMON PORT VERIFICATION"
echo "============================"
echo "Testing common ports that InstaAudit typically scans..."

test_port() {
    local port=$1
    local service=$2
    local risk=$3
    
    echo "Testing port $port ($service)..."
    if timeout 3 bash -c "</dev/tcp/$TARGET/$port" 2>/dev/null; then
        echo "✅ Port $port is OPEN - $service is accessible"
        echo "   Risk Level: $risk"
        case $port in
            80)
                echo "   💡 Verify: curl -I http://$TARGET"
                ;;
            443)
                echo "   💡 Verify: curl -I https://$TARGET"
                ;;
            22)
                echo "   💡 Verify: ssh $TARGET (should show SSH banner)"
                ;;
            3306)
                echo "   ⚠️  WARNING: MySQL database might be exposed!"
                echo "   💡 Verify: mysql -h $TARGET -u root -p"
                ;;
            5432)
                echo "   ⚠️  WARNING: PostgreSQL database might be exposed!"
                echo "   💡 Verify: psql -h $TARGET -U postgres"
                ;;
        esac
    else
        echo "❌ Port $port is CLOSED or filtered"
    fi
    echo ""
}

test_port 80 "HTTP Web Server" "Medium"
test_port 443 "HTTPS Web Server" "Low"
test_port 22 "SSH Remote Access" "Medium"
test_port 21 "FTP File Transfer" "High"
test_port 23 "Telnet (Insecure)" "High"
test_port 3306 "MySQL Database" "CRITICAL"
test_port 5432 "PostgreSQL Database" "CRITICAL"
test_port 3389 "Windows RDP" "CRITICAL"

# Web server specific tests
echo "4️⃣  WEB SERVER VERIFICATION"
echo "==========================="
echo "If InstaAudit found web services, let's verify them..."

if timeout 5 curl -I -s http://$TARGET >/dev/null 2>&1; then
    echo "✅ HTTP web server confirmed"
    echo "Headers:"
    curl -I -s http://$TARGET | head -5
    echo ""
    echo "💡 Security check - look for these headers:"
    echo "   - Strict-Transport-Security (forces HTTPS)"
    echo "   - X-Frame-Options (prevents clickjacking)"
    echo "   - X-Content-Type-Options (prevents MIME attacks)"
    echo ""
else
    echo "❌ No HTTP web server found"
fi

if timeout 5 curl -I -s https://$TARGET >/dev/null 2>&1; then
    echo "✅ HTTPS web server confirmed"
    echo "SSL Certificate check:"
    echo | openssl s_client -connect $TARGET:443 -servername $TARGET 2>/dev/null | grep -E "(subject=|issuer=|Verify return code)"
    echo ""
else
    echo "❌ No HTTPS web server found"
fi

# Database verification
echo "5️⃣  DATABASE EXPOSURE CHECK"
echo "==========================="
echo "⚠️  CRITICAL: Checking if databases are exposed to the internet..."
echo "(Databases should NEVER be accessible from the internet!)"
echo ""

check_database() {
    local port=$1
    local name=$2
    local test_cmd=$3
    
    echo "Checking $name on port $port..."
    if timeout 3 bash -c "</dev/tcp/$TARGET/$port" 2>/dev/null; then
        echo "🚨 CRITICAL: $name database appears to be accessible!"
        echo "   This is a SERIOUS security risk!"
        echo "   💡 Manual test: $test_cmd"
        echo "   🔧 Fix: Block port $port from internet access"
    else
        echo "✅ $name database is not accessible (good!)"
    fi
    echo ""
}

check_database 3306 "MySQL" "mysql -h $TARGET -u root -p"
check_database 5432 "PostgreSQL" "psql -h $TARGET -U postgres"
check_database 27017 "MongoDB" "mongo $TARGET:27017"
check_database 6379 "Redis" "redis-cli -h $TARGET"

# Summary and recommendations
echo "6️⃣  VERIFICATION SUMMARY"
echo "======================="
echo "✅ Manual verification complete!"
echo ""
echo "🎯 What to do next:"
echo "1. Compare these results with your InstaAudit report"
echo "2. If results differ significantly, investigate why"
echo "3. Focus on any CRITICAL issues (exposed databases)"
echo "4. Use online tools for additional verification:"
echo "   - SSL Labs: https://www.ssllabs.com/ssltest/"
echo "   - Security Headers: https://securityheaders.com/"
echo ""
echo "🔍 For more detailed verification:"
echo "   ./verify-results.sh $TARGET audit_report.json"
echo ""
echo "📚 Learn more:"
echo "   - Read VERIFICATION-GUIDE.md"
echo "   - Read TRUST-AND-VERIFICATION.md"
echo "   - Check UNDERSTANDING-RESULTS.md"
echo ""
echo "🛡️  Remember: Always verify critical security findings!"
echo "    Don't rely on just one tool - use multiple sources!"