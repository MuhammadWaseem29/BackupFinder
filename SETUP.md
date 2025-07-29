# BackupFinder - Automated Setup Guide

BackupFinder now includes automated Go installation and setup functionality specifically designed for Linux systems. The tool can automatically detect, install, and configure Go if it's not present on your system.

## 🚀 Quick Start

### For Linux Systems (Recommended)

1. **Automated Installation**:
   ```bash
   # Clone the repository
   git clone https://github.com/MuhammadWaseem29/BackupFinder.git
   cd BackupFinder
   
   # Run the automated installer
   chmod +x install.sh
   ./install.sh
   ```

2. **Basic Setup**:
   ```bash
   # Alternative basic setup
   chmod +x setup.sh
   ./setup.sh
   ```

### For All Systems

1. **Check System Requirements**:
   ```bash
   # If you already have Go installed, build and check
   go build -o backupfinder ./cmd/backupfinder
   ./backupfinder system-check
   ```

2. **Interactive Setup Guide**:
   ```bash
   ./backupfinder setup
   ```

3. **Quick Start Guide**:
   ```bash
   ./backupfinder quick-start
   ```

## 🔧 Installation Scripts

### install.sh (Advanced)
- Comprehensive installation script with multiple fallback options
- Supports package managers (apt, yum, dnf, pacman, zypper, brew)
- Falls back to official Go installation if package manager fails
- Automatic environment configuration
- Built-in verification and testing
- System architecture detection

**Features:**
- ✅ Multi-platform support (Linux, macOS)
- ✅ Multiple installation methods
- ✅ Version compatibility checking
- ✅ Automatic PATH configuration
- ✅ GOPATH setup
- ✅ Build verification
- ✅ Health check testing

**Usage:**
```bash
./install.sh                    # Full automated installation
./install.sh --help            # Show help
./install.sh --check-only      # Only check Go installation
./install.sh --build-only      # Only build (assumes Go is installed)
./install.sh --force-install   # Force Go installation
```

### setup.sh (Basic)
- Simple and straightforward installation
- Downloads Go 1.22.3 from official source
- Basic environment setup
- Suitable for most Linux systems

**Usage:**
```bash
./setup.sh
```

## 🖥️ Built-in Commands

BackupFinder includes several built-in commands to help with setup and diagnostics:

### System Diagnostics
```bash
./backupfinder system-check
```
**Output:**
- Operating system and architecture
- Go installation status
- Version compatibility check
- Environment variables (GOPATH, GOROOT)
- Platform-specific installation instructions

### Interactive Setup
```bash
./backupfinder setup
```
**Features:**
- Step-by-step setup guidance
- Automatic system detection
- Build and test instructions
- Troubleshooting tips

### Quick Start Guide
```bash
./backupfinder quick-start
```
**Provides:**
- System readiness check
- Build instructions
- Basic usage examples
- Next steps

### Health Check
```bash
./backupfinder health-check
```
**Verifies:**
- Embedded assets
- Pattern files
- System requirements

## 📋 Manual Installation

If the automated scripts don't work for your system, you can install Go manually:

### Linux
```bash
# Download Go
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz

# Remove existing installation
sudo rm -rf /usr/local/go

# Extract
sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

### macOS
```bash
# Using Homebrew
brew install go

# Or download installer from: https://go.dev/dl/
```

### Windows
```bash
# Using Chocolatey
choco install golang

# Using Scoop
scoop install go

# Or download installer from: https://go.dev/dl/
```

## 🏗️ Building BackupFinder

Once Go is installed:

```bash
# Build
go build -o backupfinder ./cmd/backupfinder

# Make executable (Linux/macOS)
chmod +x backupfinder

# Test
./backupfinder health-check
./backupfinder version
```

## 📖 Usage Examples

```bash
# Basic usage with wordlist
./backupfinder -u example.com -w

# Custom extensions
./backupfinder -u example.com -e txt,bak,old,backup

# Multiple targets
./backupfinder -l targets.txt -w

# JSON output
./backupfinder -u example.com -w --json -o results.json

# Silent mode
./backupfinder -u example.com -w --silent

# Custom concurrency
./backupfinder -u example.com -w -c 20
```

## 🔍 Troubleshooting

### Go Not Found After Installation
```bash
# Reload shell configuration
source ~/.bashrc

# Or restart terminal

# Verify PATH
echo $PATH | grep go
```

### Permission Issues
```bash
# Make scripts executable
chmod +x install.sh setup.sh

# If sudo is needed for Go installation
sudo ./install.sh
```

### Build Failures
```bash
# Check Go installation
go version
go env

# Verify you're in the right directory
ls -la | grep go.mod

# Clean and rebuild
go clean
go build -o backupfinder ./cmd/backupfinder
```

### Version Compatibility
```bash
# Check minimum requirements
./backupfinder system-check

# Upgrade Go if needed (Linux)
./install.sh --force-install
```

## 🌟 Features

- **Automated Go Detection**: Checks if Go is already installed
- **Version Validation**: Ensures Go version meets minimum requirements (1.18+)
- **Multi-Platform Support**: Works on Linux, macOS, and Windows
- **Multiple Installation Methods**: Package managers, official downloads, manual setup
- **Environment Configuration**: Automatic PATH and GOPATH setup
- **Build Verification**: Tests the build process
- **Health Checks**: Comprehensive system diagnostics
- **Interactive Guidance**: Step-by-step setup instructions

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Test your changes with the setup scripts
4. Submit a pull request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙋‍♂️ Support

If you encounter issues with the installation scripts:

1. Run `./backupfinder system-check` for diagnostics
2. Check the troubleshooting section above
3. Open an issue on GitHub with system details

---

**Note**: The automated installation scripts are designed primarily for Linux systems. For other operating systems, use the manual installation instructions or built-in diagnostic commands.
