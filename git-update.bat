@echo off
setlocal enabledelayedexpansion

echo ========================================
echo InstaAudit Git Repository Update
echo ========================================
echo.

:: Check if Git is installed
git --version >nul 2>&1
if errorlevel 1 (
    echo Error: Git is not installed or not in PATH
    echo Please install Git from https://git-scm.com/
    pause
    exit /b 1
)

:: Check if we're in a Git repository
if not exist .git (
    echo Error: Not in a Git repository
    echo Run 'git init' first or clone the repository
    pause
    exit /b 1
)

:: Show current status
echo Current Git Status:
echo ----------------------
git status --short
echo.

:: Clean up any build artifacts and reports (preserve Docker files)
echo Cleaning up unnecessary files...
if exist instaaudit del instaaudit
if exist instaaudit.exe del instaaudit.exe
if exist audit_report* del audit_report*
if exist *_educational.* del *_educational.*
if exist *_explained.* del *_explained.*
if exist temp_cleanup.bat del temp_cleanup.bat
echo Preserved Docker files (Dockerfile, docker-compose.yml, .dockerignore)

:: Add all relevant files
echo Adding files to Git...
git add .
git add -A

:: Show what will be committed
echo.
echo Files to be committed:
git diff --cached --name-status
echo.

:: Get commit message
set /p commit_msg="Enter commit message (or press Enter for auto-message): "
if "!commit_msg!"=="" (
    set commit_msg=Update InstaAudit: Enhanced Termux support, Docker deployment, educational features, and verification tools
)

:: Commit changes
echo Committing changes...
git commit -m "!commit_msg!"

if errorlevel 1 (
    echo Warning: No changes to commit or commit failed
    pause
    exit /b 1
)

:: Check if remote exists
git remote get-url origin >nul 2>&1
if errorlevel 1 (
    echo No remote repository configured.
    set /p repo_url="Enter your GitHub repository URL (https://github.com/username/repo.git): "
    if "!repo_url!"=="" (
        echo Error: No repository URL provided
        pause
        exit /b 1
    )
    git remote add origin !repo_url!
    echo Remote repository added: !repo_url!
)

:: Show remote info
echo.
echo Remote repository:
git remote -v
echo.

:: Push to GitHub
echo Pushing to GitHub...
git push -u origin main

if errorlevel 1 (
    echo.
    echo Push failed. Common solutions:
    echo.
    echo 1. Authentication Issues:
    echo    - Make sure you're logged into Git
    echo    - Use: git config --global user.name "Your Name"
    echo    - Use: git config --global user.email "your.email@example.com"
    echo    - For GitHub, you may need a Personal Access Token
    echo.
    echo 2. Repository Issues:
    echo    - Make sure the repository exists on GitHub
    echo    - Check if you have write permissions
    echo.
    echo 3. Branch Issues:
    echo    - Try: git push -u origin master (if main branch is master)
    echo    - Or: git branch -M main (to rename current branch to main)
    echo.
    pause
    exit /b 1
)

:: Success message
echo.
echo ========================================
echo SUCCESS! Repository Updated
echo ========================================
echo.
echo Update Summary:
echo ------------------
echo • Cleaned up unnecessary files (preserved Docker files)
echo • Committed changes with message: "!commit_msg!"
echo • Pushed to GitHub successfully
echo.
echo Your repository now includes:
echo • Enhanced Termux support and mobile use cases
echo • Docker deployment (Dockerfile, docker-compose.yml)
echo • Educational reports and verification tools
echo • Trust and verification documentation
echo • Cross-platform compatibility
echo • Mobile security testing features
echo.
echo View your repository at:
git remote get-url origin
echo.
echo Next Steps:
echo • Test Docker deployment: docker build -t instaaudit .
echo • Share repository with others
echo • Create releases/tags for versions
echo • Test installation on different platforms
echo.

:: Self-cleanup
echo Cleaning up deployment script...
timeout /t 3 /nobreak >nul
del "%~f0"