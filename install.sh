#!/bin/bash

# Advanced BackupFinder Setup Script
# Supports multiple installation methods and environment detection

set -e

# Configuration
GO_VERSION="1.22.3"
MIN_GO_VERSION="1.18"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_NAME="backupfinder"

# Colors
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly PURPLE='\033[0;35m'
readonly CYAN='\033[0;36m'
readonly NC='\033[0m'

# Logging functions
log_info() { echo -e "${BLUE}[INFO]${NC} $*"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $*"; }
log_warning() { echo -e "${YELLOW}[WARNING]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }
log_debug() { echo -e "${PURPLE}[DEBUG]${NC} $*"; }

# Banner
show_banner() {
    echo -e "${CYAN}"
    echo "=================================================="
    echo "          BackupFinder Setup Script v2.0         "
    echo "          Automated Go Installation & Build       "
    echo "=================================================="
    echo -e "${NC}"
}

# System detection
detect_system() {
    local os arch
    os=$(uname -s)
    arch=$(uname -m)
    
    case "$os" in
        Linux*)
            case "$arch" in
                x86_64) echo "linux-amd64" ;;
                i386|i686) echo "linux-386" ;;
                armv6l) echo "linux-armv6l" ;;
                armv7l) echo "linux-armv6l" ;;
                aarch64|arm64) echo "linux-arm64" ;;
                *) log_error "Unsupported architecture: $arch"; exit 1 ;;
            esac
            ;;
        Darwin*)
            case "$arch" in
                x86_64) echo "darwin-amd64" ;;
                arm64) echo "darwin-arm64" ;;
                *) log_error "Unsupported macOS architecture: $arch"; exit 1 ;;
            esac
            ;;
        *)
            log_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
}

# Version comparison
version_ge() {
    [ "$(printf '%s\n' "$1" "$2" | sort -V | head -n1)" = "$2" ]
}

# Check Go installation
check_go() {
    if command -v go >/dev/null 2>&1; then
        local current_version
        current_version=$(go version | grep -oE 'go[0-9]+\.[0-9]+\.[0-9]+' | sed 's/go//')
        
        log_info "Found Go version: $current_version"
        
        if version_ge "$current_version" "$MIN_GO_VERSION"; then
            log_success "Go version is sufficient (>= $MIN_GO_VERSION)"
            return 0
        else
            log_warning "Go version $current_version is too old (need >= $MIN_GO_VERSION)"
            return 1
        fi
    else
        log_info "Go is not installed"
        return 1
    fi
}

# Install Go using package manager
install_go_package_manager() {
    log_info "Attempting to install Go using package manager..."
    
    if command -v apt-get >/dev/null 2>&1; then
        # Ubuntu/Debian
        log_info "Detected APT package manager"
        if sudo apt-get update && sudo apt-get install -y golang-go; then
            log_success "Go installed via APT"
            return 0
        fi
    elif command -v yum >/dev/null 2>&1; then
        # RHEL/CentOS
        log_info "Detected YUM package manager"
        if sudo yum install -y golang; then
            log_success "Go installed via YUM"
            return 0
        fi
    elif command -v dnf >/dev/null 2>&1; then
        # Fedora
        log_info "Detected DNF package manager"
        if sudo dnf install -y golang; then
            log_success "Go installed via DNF"
            return 0
        fi
    elif command -v pacman >/dev/null 2>&1; then
        # Arch Linux
        log_info "Detected Pacman package manager"
        if sudo pacman -S --noconfirm go; then
            log_success "Go installed via Pacman"
            return 0
        fi
    elif command -v zypper >/dev/null 2>&1; then
        # openSUSE
        log_info "Detected Zypper package manager"
        if sudo zypper install -y go; then
            log_success "Go installed via Zypper"
            return 0
        fi
    elif command -v brew >/dev/null 2>&1; then
        # macOS with Homebrew
        log_info "Detected Homebrew package manager"
        if brew install go; then
            log_success "Go installed via Homebrew"
            return 0
        fi
    fi
    
    log_warning "Package manager installation failed or not available"
    return 1
}

# Install Go from official source
install_go_official() {
    local platform temp_dir
    platform=$(detect_system)
    temp_dir=$(mktemp -d)
    
    log_info "Installing Go $GO_VERSION for $platform from official source..."
    
    cd "$temp_dir"
    
    # Download
    local download_url="https://go.dev/dl/go${GO_VERSION}.${platform}.tar.gz"
    log_info "Downloading from: $download_url"
    
    if command -v wget >/dev/null 2>&1; then
        wget -q --show-progress "$download_url" || {
            log_error "Download failed"
            return 1
        }
    elif command -v curl >/dev/null 2>&1; then
        curl -L -o "go${GO_VERSION}.${platform}.tar.gz" "$download_url" || {
            log_error "Download failed"
            return 1
        }
    else
        log_error "Neither wget nor curl is available"
        return 1
    fi
    
    # Remove existing installation
    log_info "Removing existing Go installation..."
    sudo rm -rf /usr/local/go
    
    # Extract
    log_info "Extracting Go..."
    sudo tar -C /usr/local -xzf "go${GO_VERSION}.${platform}.tar.gz" || {
        log_error "Extraction failed"
        return 1
    }
    
    # Cleanup
    cd - >/dev/null
    rm -rf "$temp_dir"
    
    log_success "Go $GO_VERSION installed successfully"
    return 0
}

# Setup environment
setup_environment() {
    log_info "Setting up Go environment..."
    
    local shell_rc
    case "$SHELL" in
        */bash) shell_rc="$HOME/.bashrc" ;;
        */zsh) shell_rc="$HOME/.zshrc" ;;
        */fish) shell_rc="$HOME/.config/fish/config.fish" ;;
        *) shell_rc="$HOME/.profile" ;;
    esac
    
    # Create shell config if it doesn't exist
    touch "$shell_rc"
    
    # Add Go paths
    local go_path_export='export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin'
    
    if ! grep -q "/usr/local/go/bin" "$shell_rc"; then
        echo "$go_path_export" >> "$shell_rc"
        log_success "Added Go paths to $shell_rc"
    else
        log_info "Go paths already present in $shell_rc"
    fi
    
    # Export for current session
    export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
    
    # Set GOPATH if not set
    if [ -z "$GOPATH" ]; then
        export GOPATH="$HOME/go"
        echo 'export GOPATH=$HOME/go' >> "$shell_rc"
        log_success "Set GOPATH to $HOME/go"
    fi
    
    # Create Go workspace directories
    mkdir -p "$HOME/go/bin" "$HOME/go/src" "$HOME/go/pkg"
    
    log_success "Environment setup complete"
}

# Verify installation
verify_installation() {
    log_info "Verifying Go installation..."
    
    # Source shell config for verification
    export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
    
    if command -v go >/dev/null 2>&1; then
        local version goroot gopath
        version=$(go version)
        goroot=$(go env GOROOT 2>/dev/null || echo "Not set")
        gopath=$(go env GOPATH 2>/dev/null || echo "Not set")
        
        log_success "Go installation verified: $version"
        log_info "GOROOT: $goroot"
        log_info "GOPATH: $gopath"
        
        return 0
    else
        log_error "Go verification failed"
        return 1
    fi
}

# Build BackupFinder
build_tool() {
    log_info "Building BackupFinder..."
    
    cd "$SCRIPT_DIR"
    
    if [ ! -f "go.mod" ]; then
        log_error "go.mod not found. Please run this script from the BackupFinder root directory."
        return 1
    fi
    
    # Clean previous builds
    rm -f "$BINARY_NAME"
    
    # Build with optimizations
    log_info "Compiling with optimizations..."
    go build -ldflags="-w -s" -o "$BINARY_NAME" ./cmd/backupfinder || {
        log_error "Build failed"
        return 1
    }
    
    chmod +x "$BINARY_NAME"
    
    local size
    size=$(du -h "$BINARY_NAME" | cut -f1)
    log_success "BackupFinder built successfully (size: $size)"
    
    return 0
}

# Test the tool
test_tool() {
    log_info "Testing BackupFinder..."
    
    cd "$SCRIPT_DIR"
    
    if [ ! -x "$BINARY_NAME" ]; then
        log_error "BackupFinder binary not found or not executable"
        return 1
    fi
    
    # Run health check
    if ./"$BINARY_NAME" health-check >/dev/null 2>&1; then
        log_success "BackupFinder health check passed"
        
        # Show usage
        echo
        log_info "BackupFinder is ready to use!"
        echo "Usage examples:"
        echo "  ./$BINARY_NAME --help"
        echo "  ./$BINARY_NAME -u example.com -w"
        echo "  ./$BINARY_NAME -u example.com -e txt,bak,old"
        echo "  ./$BINARY_NAME version"
        
        return 0
    else
        log_error "BackupFinder health check failed"
        return 1
    fi
}

# Main function
main() {
    show_banner
    
    # Check if Go is already installed and sufficient
    if check_go; then
        log_info "Skipping Go installation"
    else
        log_info "Installing Go..."
        
        # Try package manager first
        if ! install_go_package_manager; then
            log_info "Falling back to official installation method..."
            if ! install_go_official; then
                log_error "Failed to install Go"
                exit 1
            fi
        fi
        
        # Setup environment
        setup_environment
        
        # Verify installation
        if ! verify_installation; then
            log_error "Go installation verification failed"
            exit 1
        fi
    fi
    
    # Build the tool
    if ! build_tool; then
        log_error "Failed to build BackupFinder"
        exit 1
    fi
    
    # Test the tool
    if ! test_tool; then
        log_error "BackupFinder testing failed"
        exit 1
    fi
    
    echo
    log_success "Setup completed successfully! 🎉"
    echo
    echo "Note: You may need to restart your shell or run 'source ~/.bashrc' to use Go in new terminal sessions."
}

# Handle script arguments
case "${1:-}" in
    --help|-h)
        echo "BackupFinder Setup Script"
        echo "Usage: $0 [options]"
        echo ""
        echo "Options:"
        echo "  --help, -h          Show this help message"
        echo "  --check-only        Only check Go installation, don't install"
        echo "  --build-only        Only build BackupFinder, don't install Go"
        echo "  --force-install     Force Go installation even if already present"
        echo ""
        exit 0
        ;;
    --check-only)
        check_go
        exit $?
        ;;
    --build-only)
        show_banner
        if check_go; then
            build_tool && test_tool
        else
            log_error "Go is not properly installed. Run without --build-only first."
            exit 1
        fi
        exit $?
        ;;
    --force-install)
        show_banner
        log_info "Force installing Go..."
        install_go_official && setup_environment && verify_installation && build_tool && test_tool
        exit $?
        ;;
    "")
        main
        ;;
    *)
        log_error "Unknown option: $1"
        echo "Use --help for usage information"
        exit 1
        ;;
esac
