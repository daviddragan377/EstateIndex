#!/bin/bash
# Estate Index: Development Setup Guide
# Run this script to verify all dependencies are installed

set -e

echo "Estate Index: Development Environment Setup"
echo "============================================"
echo ""

# Check Hugo
echo "Checking Hugo..."
if command -v hugo &> /dev/null; then
    HUGO_VERSION=$(hugo version | cut -d' ' -f5)
    echo "✓ Hugo $HUGO_VERSION installed"
else
    echo "✗ Hugo not found. Install from: https://gohugo.io/installation/"
    exit 1
fi

# Check Go
echo "Checking Go..."
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo "✓ Go $GO_VERSION installed"
else
    echo "✗ Go not found. Install from: https://golang.org/doc/install"
    exit 1
fi

# Check Node.js
echo "Checking Node.js..."
if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    echo "✓ Node.js $NODE_VERSION installed"
else
    echo "✗ Node.js not found. Install from: https://nodejs.org/"
    exit 1
fi

# Check npm
echo "Checking npm..."
if command -v npm &> /dev/null; then
    NPM_VERSION=$(npm --version)
    echo "✓ npm $NPM_VERSION installed"
else
    echo "✗ npm not found. Install from: https://nodejs.org/"
    exit 1
fi

echo ""
echo "Installing Node dependencies..."
npm install

echo ""
echo "Building xmlsync CLI..."
cd cmd/xmlsync
go build -o xmlsync .
cd ../..
echo "✓ xmlsync CLI built"

echo ""
echo "============================================"
echo "✅ Development environment is ready!"
echo ""
echo "Next steps:"
echo "  1. Build the site:      ./scripts/build.sh"
echo "  2. View locally:        hugo server"
echo "  3. Deploy to Netlify:   ./scripts/deploy.sh netlify"
echo ""
echo "For more info, see README.md"
