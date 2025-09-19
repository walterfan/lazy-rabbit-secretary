#!/bin/bash

# Build script for lazy-rabbit-secretary
# This script builds the Go binary locally for Linux deployment

set -e

echo "🚀 Building lazy-rabbit-secretary for Linux deployment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="lazy-rabbit-secretary"
TARGET_OS="linux"
TARGET_ARCH="amd64"

# Build flags for static linking (important for Alpine)
BUILD_FLAGS="-a -ldflags '-linkmode external -extldflags \"-static\"'"

echo -e "${BLUE}📋 Build Configuration:${NC}"
echo -e "  Binary Name: ${YELLOW}${BINARY_NAME}${NC}"
echo -e "  Target OS: ${YELLOW}${TARGET_OS}${NC}"
echo -e "  Target Arch: ${YELLOW}${TARGET_ARCH}${NC}"
echo -e "  CGO Enabled: ${YELLOW}1${NC} (for database drivers)"
echo ""

# Check if go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed or not in PATH${NC}"
    exit 1
fi

# Show Go version
echo -e "${BLUE}🔧 Go Version:${NC}"
go version
echo ""

# Clean previous build
if [ -f "$BINARY_NAME" ]; then
    echo -e "${YELLOW}🧹 Cleaning previous build...${NC}"
    rm "$BINARY_NAME"
fi

# Install dependencies if needed
echo -e "${BLUE}📦 Checking dependencies...${NC}"
if [ ! -d "vendor" ]; then
    echo -e "${YELLOW}📥 Running go mod vendor...${NC}"
    go mod vendor
fi

# Build the binary
echo -e "${BLUE}🔨 Building binary for ${TARGET_OS}/${TARGET_ARCH}...${NC}"

# Set build environment
export CGO_ENABLED=1
export GOOS=$TARGET_OS
export GOARCH=$TARGET_ARCH

# For cross-compilation from macOS to Linux, we need to use Docker or disable CGO
# Check if we're cross-compiling
if [[ "$OSTYPE" == "darwin"* ]] && [[ "$TARGET_OS" == "linux" ]]; then
    echo -e "${YELLOW}⚠️  Cross-compiling from macOS to Linux detected${NC}"
    echo -e "${BLUE}🐳 Using Docker for CGO cross-compilation...${NC}"
    
    # Use Docker to build with proper Linux CGO environment
    # Use multi-platform Docker image for cross-compilation
    docker run --rm \
        --platform linux/amd64 \
        -v "$(pwd)":/app \
        -w /app \
        -e CGO_ENABLED=1 \
        -e GOOS=linux \
        -e GOARCH=amd64 \
        golang:1.24-alpine \
        sh -c "
            apk add --no-cache gcc musl-dev sqlite-dev postgresql-dev pkgconfig && \
            go build -a -ldflags '-linkmode external -extldflags \"-static\"' -o $BINARY_NAME .
        "
else
    # Native build or same OS cross-compilation
    eval "go build $BUILD_FLAGS -o $BINARY_NAME ."
fi

# Check if build was successful
if [ $? -eq 0 ] && [ -f "$BINARY_NAME" ]; then
    echo -e "${GREEN}✅ Build successful!${NC}"
    
    # Show binary info
    echo -e "${BLUE}📊 Binary Information:${NC}"
    ls -lh "$BINARY_NAME"
    file "$BINARY_NAME"
    echo ""
    
    echo -e "${GREEN}🎉 Ready for Docker deployment!${NC}"
    echo -e "${BLUE}Next steps:${NC}"
    echo -e "  1. ${YELLOW}docker build -t lazy-rabbit-secretary .${NC}"
    echo -e "  2. ${YELLOW}docker-compose up -d${NC}"
    echo ""
else
    echo -e "${RED}❌ Build failed!${NC}"
    exit 1
fi
