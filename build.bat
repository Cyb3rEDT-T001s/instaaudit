@echo off
echo ========================================
echo InstaAudit - Security Auditing Tool
echo Build Script
echo ========================================

REM Check if Go is installed
where go >nul 2>&1
if %errorlevel% neq 0 (
    if exist "C:\Program Files\Go\bin\go.exe" (
        set "PATH=%PATH%;C:\Program Files\Go\bin"
    ) else if exist "C:\Go\bin\go.exe" (
        set "PATH=%PATH%;C:\Go\bin"
    ) else (
        echo ❌ Go not installed! 
        echo Download from: https://golang.org/dl/
        pause
        exit /b 1
    )
)

echo ✅ Go found!
go version
echo.

echo 🔨 Building InstaAudit...
go mod tidy
go build -o instaaudit.exe cmd/main.go

if %errorlevel% neq 0 (
    echo ❌ Build failed!
    pause
    exit /b 1
)

echo ✅ Build successful!
echo.
echo Usage:
echo   instaaudit.exe -H target.com
echo   instaaudit.exe -H target.com -A -f html
echo.
pause