#!/bin/bash

# build_docker_image.sh - Flexible Docker image builder
# This script detects the OS, builds the appropriate binary, and creates a Docker image

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="lazy-rabbit-secretary"
DOCKER_IMAGE_NAME="${PROJECT_NAME}"
DOCKER_TAG="${1:-latest}"  # Use first argument as tag, default to 'latest'
BUILD_DIR="build"

# Default target environment (production)
DEFAULT_TARGET_OS="linux"
DEFAULT_TARGET_ARCH="amd64"
DEFAULT_TARGET_PLATFORM="${DEFAULT_TARGET_OS}-${DEFAULT_TARGET_ARCH}"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    # Convert architecture names
    case $arch in
        x86_64)
            arch="amd64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    
    echo "${os}-${arch}"
}

# Function to build binary for current platform
build_binary() {
    local platform=$1
    local os=$(echo $platform | cut -d'-' -f1)
    local arch=$(echo $platform | cut -d'-' -f2)
    
    print_status "Building binary for platform: $platform"
    
    case $os in
        darwin)
            print_status "Building for macOS ($arch)..."
            make build-darwin
            ;;
        linux)
            print_status "Building for Linux ($arch)..."
            make build-linux
            ;;
        *)
            print_error "Unsupported OS: $os"
            exit 1
            ;;
    esac
    
    print_success "Binary build completed for $platform"
}

# Function to prepare binary for Docker (default: linux-amd64)
prepare_docker_binary() {
    local platform=${1:-"$DEFAULT_TARGET_PLATFORM"}  # Default to configured target platform
    local os=$(echo $platform | cut -d'-' -f1)
    local arch=$(echo $platform | cut -d'-' -f2)
    
    # For Docker, we use Linux static binary (default: amd64)
    print_status "Preparing Linux static binary for Docker (target: $platform)..."
    
    # Build static Linux binary
    make build-static
    
    # Check if target binary exists
    local target_binary="${BUILD_DIR}/${PROJECT_NAME}-${platform}-static"
    
    if [ ! -f "$target_binary" ]; then
        print_error "Target static binary not found: $target_binary"
        exit 1
    fi
    
    print_success "Docker binary prepared for target environment:"
    print_success "  - Target: $platform"
    print_success "  - Binary: $target_binary"
}

# Function to build Docker image
build_docker_image() {
    local image_name="${DOCKER_IMAGE_NAME}:${DOCKER_TAG}"
    
    print_status "Building Docker image: $image_name"
    
    # Check if Dockerfile exists
    if [ ! -f "Dockerfile" ]; then
        print_error "Dockerfile not found in current directory"
        exit 1
    fi
    
    # Check if Docker buildx is available for multi-arch builds
    if docker buildx version >/dev/null 2>&1; then
        print_status "Building Docker image for linux/amd64 (target environment)..."
        
        # Create and use buildx builder if it doesn't exist
        docker buildx create --name multiarch --use 2>/dev/null || docker buildx use multiarch
        
        # Build for AMD64 Linux (target environment)
        docker buildx build \
            --platform linux/amd64 \
            --tag "$image_name" \
            --load \
            .
    else
        print_warning "Docker buildx not available, building for current platform only"
        
        # Build for current platform only
        docker build -t "$image_name" .
    fi
    
    if [ $? -eq 0 ]; then
        print_success "Docker image built successfully: $image_name"
        
        # Show image info
        print_status "Docker image information:"
        docker images "$image_name" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"
    else
        print_error "Docker image build failed"
        exit 1
    fi
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [DOCKER_TAG]"
    echo ""
    echo "Arguments:"
    echo "  DOCKER_TAG    Docker image tag (default: latest)"
    echo ""
    echo "Examples:"
    echo "  $0                    # Build with tag 'latest'"
    echo "  $0 v1.0.0            # Build with tag 'v1.0.0'"
    echo "  $0 dev               # Build with tag 'dev'"
    echo ""
    echo "This script will:"
    echo "  1. Detect your current OS and architecture"
    echo "  2. Build the appropriate binary for your platform"
    echo "  3. Build Linux static binary for Docker (default: amd64)"
    echo "  4. Create a Docker image for target environment (default: linux/amd64)"
    echo ""
    echo "Default target environment:"
    echo "  - OS: ${DEFAULT_TARGET_OS}"
    echo "  - Architecture: ${DEFAULT_TARGET_ARCH}"
    echo "  - Platform: ${DEFAULT_TARGET_PLATFORM}"
    echo "  - Use case: Intel/AMD Linux servers (production target)"
}

# Main execution
main() {
    # Check if help is requested
    if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
        show_usage
        exit 0
    fi
    
    print_status "Starting Docker image build process..."
    print_status "Target Docker image: ${DOCKER_IMAGE_NAME}:${DOCKER_TAG}"
    
    # Detect current platform
    local current_platform=$(detect_platform)
    print_status "Detected platform: $current_platform"
    
    # Check if Makefile exists
    if [ ! -f "Makefile" ]; then
        print_error "Makefile not found in current directory"
        exit 1
    fi
    
    # Build binary for current platform (for development/testing)
    build_binary "$current_platform"
    
    # Prepare Linux static binary for Docker (uses default: linux-amd64)
    prepare_docker_binary
    
    # Build Docker image
    build_docker_image
    
    print_success "Docker image build process completed!"
    print_status "You can now run the container with:"
    print_status "  docker run -p 8080:8080 ${DOCKER_IMAGE_NAME}:${DOCKER_TAG}"
}

# Run main function
main "$@"
