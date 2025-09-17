#!/bin/bash

echo "========================================"
echo "InstaAudit - Docker Build"
echo "========================================"

# Build Docker image
echo "🐳 Building Docker image..."
docker build -t instaaudit:latest .

if [ $? -eq 0 ]; then
    echo "✅ Docker image built successfully!"
    echo ""
    echo "Usage:"
    echo "  docker run --rm instaaudit:latest -H target.com"
    echo "  docker run --rm instaaudit:latest -H target.com -A -f json"
    echo ""
    echo "To save reports, mount a volume:"
    echo "  docker run --rm -v \$(pwd):/reports instaaudit:latest -H target.com -o /reports/audit"
    echo ""
else
    echo "❌ Docker build failed!"
    exit 1
fi