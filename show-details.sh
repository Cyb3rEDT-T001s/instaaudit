#!/bin/bash

# Show detailed InstaAudit results
# Usage: ./show-details.sh [report.json]

REPORT=${1:-"audit_report.json"}

echo "🔍 InstaAudit Detailed Results Viewer"
echo "====================================="

if [ ! -f "$REPORT" ]; then
    echo "❌ Report file not found: $REPORT"
    echo "Usage: $0 [report.json]"
    exit 1
fi

echo "📊 Reading report: $REPORT"
echo ""

# Extract and display risk details using jq if available
if command -v jq &> /dev/null; then
    echo "🚨 RISK LEVEL:"
    jq -r '.summary.risk_level' "$REPORT" 2>/dev/null || echo "Unable to parse risk level"
    echo ""
    
    echo "📈 RISK BREAKDOWN:"
    echo "Critical Issues: $(jq -r '.summary.critical_issues // 0' "$REPORT" 2>/dev/null)"
    echo "High Risk Issues: $(jq -r '.summary.high_risk_issues // 0' "$REPORT" 2>/dev/null)"
    echo "Medium Risk Issues: $(jq -r '.summary.medium_risk_issues // 0' "$REPORT" 2>/dev/null)"
    echo "Low Risk Issues: $(jq -r '.summary.low_risk_issues // 0' "$REPORT" 2>/dev/null)"
    echo ""
    
    echo "🔍 DETAILED ISSUES:"
    jq -r '.summary.risk_details[]?' "$REPORT" 2>/dev/null | head -20 | while read -r line; do
        echo "  • $line"
    done
    echo ""
    
    echo "🔓 OPEN PORTS:"
    jq -r '.scan_result.open_ports[]?' "$REPORT" 2>/dev/null | while read -r port; do
        echo "  • Port $port"
    done
    echo ""
    
    echo "🛡️  SERVICES FOUND:"
    jq -r '.audit_result.services[]? | "  • Port \(.port): \(.service) \(.version // "")"' "$REPORT" 2>/dev/null
    echo ""
    
else
    echo "⚠️  jq not installed - showing basic info"
    echo ""
    
    # Basic grep-based parsing
    echo "🔓 OPEN PORTS:"
    grep -o '"open_ports":\[[^]]*\]' "$REPORT" | sed 's/.*\[\(.*\)\].*/\1/' | tr ',' '\n' | while read -r port; do
        echo "  • Port $port"
    done
    echo ""
    
    echo "📊 SUMMARY COUNTS:"
    grep -o '"vulnerabilities_found":[0-9]*' "$REPORT" | cut -d':' -f2 | head -1 | xargs -I {} echo "  • Vulnerabilities: {}"
    grep -o '"database_issues":[0-9]*' "$REPORT" | cut -d':' -f2 | head -1 | xargs -I {} echo "  • Database Issues: {}"
    grep -o '"webapp_issues":[0-9]*' "$REPORT" | cut -d':' -f2 | head -1 | xargs -I {} echo "  • Web App Issues: {}"
    grep -o '"system_issues":[0-9]*' "$REPORT" | cut -d':' -f2 | head -1 | xargs -I {} echo "  • System Issues: {}"
    echo ""
fi

echo "💡 RECOMMENDATIONS:"
echo "1. Review all high and critical issues immediately"
echo "2. Check the HTML report for detailed explanations"
echo "3. Use verification tools to confirm findings"
echo "4. Prioritize database and system security issues"
echo ""

echo "📄 AVAILABLE REPORTS:"
ls -la audit_report.* 2>/dev/null | while read -r line; do
    echo "  $line"
done

echo ""
echo "🔧 To verify results, run:"
echo "  ./verify-results.sh $(jq -r '.scan_result.host // "target"' "$REPORT" 2>/dev/null || echo "target") $REPORT"