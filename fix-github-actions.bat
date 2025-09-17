@echo off
echo ========================================
echo Fix GitHub Actions - Update to v4
echo ========================================

echo 🔧 The GitHub Actions workflow has been updated to use:
echo   ✅ actions/upload-artifact@v4 (was v3)
echo   ✅ actions/setup-go@v5 (was v4)
echo.

echo 📦 Pushing fix to GitHub...

git init
git branch -M main
git add .
git commit -m "🔧 Fix GitHub Actions - Update to latest versions

✅ Updated actions/upload-artifact from v3 to v4
✅ Updated actions/setup-go from v4 to v5
✅ Resolves deprecated action warnings"

git remote add origin https://github.com/Cyb3rEDT-T001s/instaaudit.git
git push -u origin main --force

if %errorlevel% equ 0 (
    echo ✅ GitHub Actions fixed!
    echo.
    echo 🎉 InstaAudit will now build without warnings!
    echo 📊 Repository: https://github.com/Cyb3rEDT-T001s/instaaudit
    echo.
    echo 🧹 Cleaning up...
    del fix-github-actions.bat
) else (
    echo ❌ Push failed - check authentication
)

pause