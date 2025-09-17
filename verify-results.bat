@echo off
setlocal enabledelayedexpansion

echo InstaAudit Results Verification Tool (Windows)
echo =============================================

if "%1"=="" (
    echo Usage: %0 ^<target_host^> [instaaudit_report.json]
    echo Example: %0 example.com audit_report.json
    pause
    exit /b 1
)

set TARGET=%1
set REPORT=%2
if "%REPORT%"=="" set REPORT=audit_report.json

echo Target: %TARGET%
echo Report: %REPORT%
echo.

:: Check available tools
echo Checking available verification tools...

:: Check for nmap
nmap --version >nul 2>&1
if errorlevel 1 (
    echo Warning: Nmap not installed - download from https://nmap.org/
    set NMAP_AVAILABLE=false
) else (
    echo ✅ Nmap available
    set NMAP_AVAILABLE=true
)

:: Check for curl
curl --version >nul 2>&1
if errorlevel 1 (
    echo Warning: Curl not available - some web checks will be skipped
    set CURL_AVAILABLE=false
) else (
    echo ✅ Curl available
    set CURL_AVAILABLE=true
)

:: Check for PowerShell
powershell -Command "Get-Host" >nul 2>&1
if errorlevel 1 (
    echo Warning: PowerShell not available
    set PS_AVAILABLE=false
) else (
    echo ✅ PowerShell available
    set PS_AVAILABLE=true
)

echo.

:: DNS Resolution Check
echo DNS Resolution Verification
echo ---------------------------
nslookup %TARGET%
echo.

:: Port Scanning Verification
if "%NMAP_AVAILABLE%"=="true" (
    echo Port Scanning Verification
    echo --------------------------
    echo Running common ports scan...
    nmap -F %TARGET%
    echo.
    
    echo Service Detection Verification
    echo ------------------------------
    nmap -sV -F %TARGET%
    echo.
)

:: Web Service Verification
if "%CURL_AVAILABLE%"=="true" (
    echo Web Service Verification
    echo ------------------------
    
    echo Testing HTTP port 80:
    curl -I -s --connect-timeout 5 http://%TARGET% 2>nul | findstr /i "HTTP"
    if errorlevel 1 echo HTTP not accessible
    echo.
    
    echo Testing HTTPS port 443:
    curl -I -s --connect-timeout 5 https://%TARGET% 2>nul | findstr /i "HTTP"
    if errorlevel 1 echo HTTPS not accessible
    echo.
)

:: Database Service Quick Check using PowerShell
if "%PS_AVAILABLE%"=="true" (
    echo Database Service Quick Check
    echo ---------------------------
    
    echo Checking MySQL port 3306...
    powershell -Command "try { $tcp = New-Object System.Net.Sockets.TcpClient; $tcp.ConnectAsync('%TARGET%', 3306).Wait(3000); if($tcp.Connected) { Write-Host 'WARNING: MySQL port 3306 is accessible!' -ForegroundColor Red } else { Write-Host 'OK: MySQL port 3306 not accessible' -ForegroundColor Green } } catch { Write-Host 'OK: MySQL port 3306 not accessible' -ForegroundColor Green } finally { if($tcp) { $tcp.Close() } }"
    
    echo Checking PostgreSQL port 5432...
    powershell -Command "try { $tcp = New-Object System.Net.Sockets.TcpClient; $tcp.ConnectAsync('%TARGET%', 5432).Wait(3000); if($tcp.Connected) { Write-Host 'WARNING: PostgreSQL port 5432 is accessible!' -ForegroundColor Red } else { Write-Host 'OK: PostgreSQL port 5432 not accessible' -ForegroundColor Green } } catch { Write-Host 'OK: PostgreSQL port 5432 not accessible' -ForegroundColor Green } finally { if($tcp) { $tcp.Close() } }"
    
    echo Checking MongoDB port 27017...
    powershell -Command "try { $tcp = New-Object System.Net.Sockets.TcpClient; $tcp.ConnectAsync('%TARGET%', 27017).Wait(3000); if($tcp.Connected) { Write-Host 'WARNING: MongoDB port 27017 is accessible!' -ForegroundColor Red } else { Write-Host 'OK: MongoDB port 27017 not accessible' -ForegroundColor Green } } catch { Write-Host 'OK: MongoDB port 27017 not accessible' -ForegroundColor Green } finally { if($tcp) { $tcp.Close() } }"
    
    echo.
)

:: Summary
echo Verification Summary
echo ===================
echo.
echo ✅ Cross-verification completed for: %TARGET%
echo.
echo What to do next:
echo 1. Compare results with your InstaAudit report
echo 2. Investigate any discrepancies  
echo 3. Use specialized tools for detailed analysis
echo 4. Document all findings
echo.
echo Remember:
echo - Different tools may show different results due to timing
echo - Some services may be filtered by firewalls
echo - Always verify critical findings manually
echo.
echo For deeper analysis, consider:
echo - nmap -A %TARGET% (aggressive scan)
echo - Online tools like SSL Labs, Security Headers
echo - Specialized vulnerability scanners
echo.

pause