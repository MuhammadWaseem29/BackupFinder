#!/bin/bash

# BackupFinder Setup Script
# This script checks for Go installation and sets up the environment

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Go version to install
GO_VERSION="1.22.3"
GO_ARCH="linux-amd64"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}       BackupFinder Setup Script       ${NC}"
echo -e "${BLUE}========================================${NC}"
echo

# Function to print colored output
print_status() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_info() {
    echo -e "${BLUE}[i]${NC} $1"
}

# Function to check if Go is installed
check_go_installation() {
    if command -v go >/dev/null 2>&1; then
        GO_CURRENT_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        print_status "Go is already installed (version: $GO_CURRENT_VERSION)"
        
        # Check if version is sufficient (1.18+)
        if [ "$(printf '%s\n' "1.18" "$GO_CURRENT_VERSION" | sort -V | head -n1)" = "1.18" ]; then
            print_status "Go version is sufficient for BackupFinder"
            return 0
        else
            print_warning "Go version is too old. Minimum required: 1.18"
            return 1
        fi
    else
        print_warning "Go is not installed"
        return 1
    fi
}

# Function to install Go
install_go() {
    print_info "Installing Go $GO_VERSION..."
    
    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    # Download Go
    print_info "Downloading Go $GO_VERSION for $GO_ARCH..."
    if ! wget -q --show-progress "https://go.dev/dl/go${GO_VERSION}.${GO_ARCH}.tar.gz"; then
        print_error "Failed to download Go. Please check your internet connection."
        exit 1
    fi
    
    # Remove existing Go installation
    print_info "Removing existing Go installation (if any)..."
    sudo rm -rf /usr/local/go
    
    # Extract Go
    print_info "Extracting Go to /usr/local/go..."
    if ! sudo tar -C /usr/local -xzf "go${GO_VERSION}.${GO_ARCH}.tar.gz"; then
        print_error "Failed to extract Go archive"
        exit 1
    fi
    
    # Clean up
    cd - >/dev/null
    rm -rf "$TEMP_DIR"
    
    print_status "Go $GO_VERSION installed successfully"
}

# Function to setup Go environment
setup_go_environment() {
    print_info "Setting up Go environment..."
    
    # Check if PATH is already set
    if echo "$PATH" | grep -q "/usr/local/go/bin"; then
        print_status "Go binary path already in PATH"
    else
        # Add to bashrc
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        print_status "Added Go binary path to ~/.bashrc"
    fi
    
    # Check if GOPATH is set
    if echo "$PATH" | grep -q "$HOME/go/bin"; then
        print_status "GOPATH bin already in PATH"
    else
        # Add GOPATH to bashrc
        echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
        print_status "Added GOPATH bin to ~/.bashrc"
    fi
    
    # Source bashrc for current session
    export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
    
    print_status "Go environment configured"
}

# Function to verify Go installation
verify_go_installation() {
    print_info "Verifying Go installation..."
    
    if command -v go >/dev/null 2>&1; then
        GO_VERSION_INSTALLED=$(go version)
        print_status "Go verification successful: $GO_VERSION_INSTALLED"
        
        # Show Go environment
        print_info "Go environment:"
        echo "  GOROOT: $(go env GOROOT)"
        echo "  GOPATH: $(go env GOPATH)"
        echo "  GOOS: $(go env GOOS)"
        echo "  GOARCH: $(go env GOARCH)"
        
        return 0
    else
        print_error "Go installation verification failed"
        return 1
    fi
}

# Function to build BackupFinder
build_backupfinder() {
    print_info "Building BackupFinder..."
    
    # Check if we're in the right directory
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Please run this script from the BackupFinder root directory."
        exit 1
    fi
    
    # Build the application
    if go build -o backupfinder ./cmd/backupfinder; then
        print_status "BackupFinder built successfully"
        
        # Make it executable
        chmod +x backupfinder
        
        # Show binary info
        BINARY_SIZE=$(du -h backupfinder | cut -f1)
        print_info "Binary size: $BINARY_SIZE"
        print_info "Binary location: $(pwd)/backupfinder"
        
        return 0
    else
        print_error "Failed to build BackupFinder"
        return 1
    fi
}

# Function to run BackupFinder
run_backupfinder() {
    print_info "Running BackupFinder health check..."
    
    if ./backupfinder health-check; then
        print_status "BackupFinder is working correctly!"
        echo
        print_info "To use BackupFinder, run: ./backupfinder --help"
        print_info "Example: ./backupfinder -u example.com -w"
    else
        print_error "BackupFinder health check failed"
        return 1
    fi
}

# Main execution flow
main() {
    # Check if running on Linux
    if [ "$(uname)" != "Linux" ]; then
        print_error "This script is designed for Linux systems only"
        exit 1
    fi
    
    # Check for required tools
    if ! command -v wget >/dev/null 2>&1; then
        print_error "wget is required but not installed. Please install it first."
        exit 1
    fi
    
    if ! command -v sudo >/dev/null 2>&1; then
        print_error "sudo is required but not available."
        exit 1
    fi
    
    # Check Go installation
    if ! check_go_installation; then
        print_info "Installing Go..."
        install_go
        setup_go_environment
        
        # Verify installation
        if ! verify_go_installation; then
            print_error "Go installation failed"
            exit 1
        fi
    fi
    
    # Build BackupFinder
    if build_backupfinder; then
        print_status "Setup completed successfully!"
        echo
        
        # Run BackupFinder
        run_backupfinder
    else
        print_error "Setup failed during build process"
        exit 1
    fi
}

# Run main function
main "$@"
