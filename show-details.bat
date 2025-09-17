@echo off
setlocal enabledelayedexpansion

set REPORT=%1
if "%REPORT%"=="" set REPORT=audit_report.json

echo InstaAudit Detailed Results Viewer
echo ==================================

if not exist "%REPORT%" (
    echo Error: Report file not found: %REPORT%
    echo Usage: %0 [report.json]
    pause
    exit /b 1
)

echo Reading report: %REPORT%
echo.

:: Check if PowerShell is available for JSON parsing
powershell -Command "Get-Host" >nul 2>&1
if errorlevel 1 (
    echo Warning: PowerShell not available - showing basic info
    goto :basic_parsing
)

:: Use PowerShell to parse JSON
echo RISK LEVEL:
powershell -Command "(Get-Content '%REPORT%' | ConvertFrom-Json).summary.risk_level" 2>nul
echo.

echo RISK BREAKDOWN:
powershell -Command "$json = Get-Content '%REPORT%' | ConvertFrom-Json; Write-Host 'Critical Issues:' $json.summary.critical_issues; Write-Host 'High Risk Issues:' $json.summary.high_risk_issues; Write-Host 'Medium Risk Issues:' $json.summary.medium_risk_issues; Write-Host 'Low Risk Issues:' $json.summary.low_risk_issues" 2>nul
echo.

echo DETAILED ISSUES:
powershell -Command "$json = Get-Content '%REPORT%' | ConvertFrom-Json; $json.summary.risk_details | Select-Object -First 10 | ForEach-Object { Write-Host '  •' $_ }" 2>nul
echo.

echo OPEN PORTS:
powershell -Command "$json = Get-Content '%REPORT%' | ConvertFrom-Json; $json.scan_result.open_ports | ForEach-Object { Write-Host '  • Port' $_ }" 2>nul
echo.

echo SERVICES FOUND:
powershell -Command "$json = Get-Content '%REPORT%' | ConvertFrom-Json; $json.audit_result.services | ForEach-Object { Write-Host '  • Port' $_.port':' $_.service $_.version }" 2>nul
echo.

goto :recommendations

:basic_parsing
echo BASIC INFO (Limited parsing without PowerShell):
echo.

:: Basic text parsing
findstr /C:"risk_level" "%REPORT%" | head -1
findstr /C:"open_ports" "%REPORT%" | head -1
findstr /C:"vulnerabilities_found" "%REPORT%" | head -1
echo.

:recommendations
echo RECOMMENDATIONS:
echo 1. Review all high and critical issues immediately
echo 2. Check the HTML report for detailed explanations  
echo 3. Use verification tools to confirm findings
echo 4. Prioritize database and system security issues
echo.

echo AVAILABLE REPORTS:
dir audit_report.* 2>nul
echo.

echo To verify results, run:
echo   verify-results.bat [target_host] %REPORT%
echo.

pause