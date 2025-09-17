#!/bin/bash

echo "========================================"
echo "InstaAudit Git Repository Update"
echo "========================================"
echo

# Check if Git is installed
if ! command -v git &> /dev/null; then
    echo "Error: Git is not installed"
    echo "Please install Git first"
    exit 1
fi

# Check if we're in a Git repository
if [ ! -d .git ]; then
    echo "Error: Not in a Git repository"
    echo "Run 'git init' first or clone the repository"
    exit 1
fi

# Show current status
echo "Current Git Status:"
echo "----------------------"
git status --short
echo

# Clean up any build artifacts and reports (preserve Docker files)
echo "Cleaning up unnecessary files..."
rm -f instaaudit instaaudit.exe
rm -f audit_report*
rm -f *_educational.*
rm -f *_explained.*
rm -f temp_cleanup.bat
echo "Preserved Docker files (Dockerfile, docker-compose.yml, .dockerignore)"

# Add all relevant files
echo "Adding files to Git..."
git add .
git add -A

# Show what will be committed
echo
echo "Files to be committed:"
git diff --cached --name-status
echo

# Get commit message
read -p "Enter commit message (or press Enter for auto-message): " commit_msg
if [ -z "$commit_msg" ]; then
    commit_msg="Update InstaAudit: Enhanced Termux support, Docker deployment, educational features, and verification tools"
fi

# Commit changes
echo "Committing changes..."
git commit -m "$commit_msg"

if [ $? -ne 0 ]; then
    echo "Warning: No changes to commit or commit failed"
    exit 1
fi

# Check if remote exists
if ! git remote get-url origin &> /dev/null; then
    echo "No remote repository configured."
    read -p "Enter your GitHub repository URL (https://github.com/username/repo.git): " repo_url
    if [ -z "$repo_url" ]; then
        echo "Error: No repository URL provided"
        exit 1
    fi
    git remote add origin "$repo_url"
    echo "Remote repository added: $repo_url"
fi

# Show remote info
echo
echo "Remote repository:"
git remote -v
echo

# Push to GitHub
echo "Pushing to GitHub..."
git push -u origin main

if [ $? -ne 0 ]; then
    echo
    echo "Push failed. Common solutions:"
    echo
    echo "1. Authentication Issues:"
    echo "   - Make sure you're logged into Git"
    echo "   - Use: git config --global user.name \"Your Name\""
    echo "   - Use: git config --global user.email \"your.email@example.com\""
    echo "   - For GitHub, you may need a Personal Access Token"
    echo
    echo "2. Repository Issues:"
    echo "   - Make sure the repository exists on GitHub"
    echo "   - Check if you have write permissions"
    echo
    echo "3. Branch Issues:"
    echo "   - Try: git push -u origin master (if main branch is master)"
    echo "   - Or: git branch -M main (to rename current branch to main)"
    echo
    exit 1
fi

# Success message
echo
echo "========================================"
echo "SUCCESS! Repository Updated"
echo "========================================"
echo
echo "Update Summary:"
echo "------------------"
echo "• Cleaned up unnecessary files (preserved Docker files)"
echo "• Committed changes with message: \"$commit_msg\""
echo "• Pushed to GitHub successfully"
echo
echo "Your repository now includes:"
echo "• Enhanced Termux support and mobile use cases"
echo "• Docker deployment (Dockerfile, docker-compose.yml)"
echo "• Educational reports and verification tools"
echo "• Trust and verification documentation"
echo "• Cross-platform compatibility"
echo "• Mobile security testing features"
echo
echo "View your repository at:"
git remote get-url origin
echo
echo "Next Steps:"
echo "• Test Docker deployment: docker build -t instaaudit ."
echo "• Share repository with others"
echo "• Create releases/tags for versions"
echo "• Test installation on different platforms"
echo

# Self-cleanup
echo "Cleaning up deployment script..."
sleep 2
rm "$0"