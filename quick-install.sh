#!/bin/bash

# BackupFinder Quick Installer
# This script makes setup scripts executable and runs the appropriate installer

set -e

echo "🚀 BackupFinder Quick Installer"
echo "==============================="
echo

# Check if we're on Linux
if [ "$(uname)" != "Linux" ]; then
    echo "❌ This automated installer is designed for Linux systems only."
    echo "🔧 For other systems:"
    echo "   1. Install Go manually: https://go.dev/dl/"
    echo "   2. Build: go build -o backupfinder ./cmd/backupfinder"
    echo "   3. Run: ./backupfinder system-check"
    exit 1
fi

# Make scripts executable
echo "🔧 Making scripts executable..."
chmod +x setup.sh install.sh 2>/dev/null || true

# Check if Go is already installed
if command -v go >/dev/null 2>&1; then
    echo "✅ Go is already installed: $(go version)"
    echo "🏗️  Building BackupFinder..."
    
    if go build -o backupfinder ./cmd/backupfinder; then
        chmod +x backupfinder
        echo "✅ Build successful!"
        
        # Run health check
        if ./backupfinder health-check >/dev/null 2>&1; then
            echo "✅ BackupFinder is working correctly!"
            echo
            echo "🎯 Usage examples:"
            echo "   ./backupfinder -u example.com -w"
            echo "   ./backupfinder --help"
        else
            echo "⚠️  Build completed but health check failed"
            echo "   Run: ./backupfinder system-check"
        fi
    else
        echo "❌ Build failed. Running system diagnostics..."
        echo "   This will help identify the issue:"
        if [ -x "./backupfinder" ]; then
            ./backupfinder system-check
        else
            echo "   Manual check: go version && go env"
        fi
    fi
else
    echo "🔧 Go is not installed. Running automated installer..."
    echo
    
    # Run the comprehensive installer
    if [ -x "./install.sh" ]; then
        ./install.sh
    elif [ -x "./setup.sh" ]; then
        ./setup.sh
    else
        echo "❌ Setup scripts not found or not executable"
        echo "🔧 Manual installation:"
        echo "   chmod +x install.sh && ./install.sh"
        exit 1
    fi
fi

echo
echo "🎉 Setup complete! BackupFinder is ready to use."
