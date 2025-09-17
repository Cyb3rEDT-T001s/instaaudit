#!/bin/bash

echo "🐳 InstaAudit Docker Build Script"
echo "================================="

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed"
    echo "Please install Docker first: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker daemon is running
if ! docker info &> /dev/null; then
    echo "❌ Docker daemon is not running"
    echo "Please start Docker first"
    exit 1
fi

# Set image name and tag
IMAGE_NAME="instaaudit"
TAG=${1:-"latest"}
FULL_IMAGE_NAME="$IMAGE_NAME:$TAG"

echo "📦 Building Docker image: $FULL_IMAGE_NAME"
echo ""

# Build the Docker image
echo "🔨 Building image..."
docker build -t "$FULL_IMAGE_NAME" .

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Docker image built successfully!"
    echo ""
    echo "📊 Image Information:"
    docker images | grep "$IMAGE_NAME"
    echo ""
    echo "🚀 Usage Examples:"
    echo ""
    echo "Basic scan:"
    echo "  docker run --rm $FULL_IMAGE_NAME -H target.com -p common"
    echo ""
    echo "Scan with output (mount volume):"
    echo "  docker run --rm -v \$(pwd)/reports:/home/instaaudit/reports $FULL_IMAGE_NAME -H target.com -p common -o reports/scan"
    echo ""
    echo "Interactive mode:"
    echo "  docker run --rm -it $FULL_IMAGE_NAME /bin/sh"
    echo ""
    echo "Help:"
    echo "  docker run --rm $FULL_IMAGE_NAME --help"
    echo ""
    echo "🔧 Advanced Usage:"
    echo ""
    echo "Network scanning (host network):"
    echo "  docker run --rm --network host $FULL_IMAGE_NAME -H 192.168.1.1 -p common"
    echo ""
    echo "Custom config:"
    echo "  docker run --rm -v \$(pwd)/config.json:/home/instaaudit/config.json $FULL_IMAGE_NAME -H target.com -c config.json"
    echo ""
else
    echo ""
    echo "❌ Docker build failed!"
    echo ""
    echo "🔍 Troubleshooting:"
    echo "• Check if Dockerfile exists"
    echo "• Ensure all source files are present"
    echo "• Check Docker daemon status"
    echo "• Review build logs above"
    exit 1
fi