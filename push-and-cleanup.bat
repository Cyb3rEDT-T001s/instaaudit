@echo off
setlocal enabledelayedexpansion

echo InstaAudit GitHub Push and Cleanup Script
echo ==========================================

:: Check if Git is installed
git --version >nul 2>&1
if errorlevel 1 (
    echo Error: Git is not installed or not in PATH
    echo Please install Git from https://git-scm.com/
    pause
    exit /b 1
)

:: Initialize Git repository if needed
if not exist .git (
    echo Initializing Git repository...
    git init
    git branch -M main
)

:: Add all files
echo Adding files to Git...
git add .

:: Commit changes
set /p commit_msg="Enter commit message (or press Enter for default): "
if "!commit_msg!"=="" set commit_msg=Update InstaAudit project with Termux support

echo Committing changes...
git commit -m "!commit_msg!"

:: Check if remote exists
git remote get-url origin >nul 2>&1
if errorlevel 1 (
    echo No remote repository found.
    set /p repo_url="Enter your GitHub repository URL: "
    git remote add origin !repo_url!
)

:: Push to GitHub
echo Uploading to GitHub...
git push -u origin main

if errorlevel 1 (
    echo Upload failed. This might be due to:
    echo 1. Authentication issues - make sure you're logged in to Git
    echo 2. Repository doesn't exist on GitHub
    echo 3. Network connectivity issues
    echo.
    echo To fix authentication, run: git config --global user.name "Your Name"
    echo And: git config --global user.email "your.email@example.com"
    echo.
    echo Script will NOT be deleted due to upload failure.
    pause
    exit /b 1
)

echo.
echo ========================================
echo SUCCESS! Project uploaded to GitHub
echo ========================================
echo.
echo Files uploaded:
echo - InstaAudit source code
echo - Termux installation scripts
echo - Documentation and guides
echo - Build scripts for all platforms
echo.
echo The upload script will self-delete in 3 seconds...
timeout /t 3 /nobreak >nul

:: Create a temporary batch file to delete this script
echo @echo off > temp_cleanup.bat
echo timeout /t 1 /nobreak ^>nul >> temp_cleanup.bat
echo del "%~nx0" >> temp_cleanup.bat
echo del "%%~f0" >> temp_cleanup.bat

:: Run the cleanup script and exit
start /b temp_cleanup.bat
exit