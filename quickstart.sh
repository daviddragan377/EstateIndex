#!/bin/bash
# Quick Start Guide for Estate Index XML Sync

echo "🏗️  Estate Index - XML Sync Quick Start"
echo "======================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.19+ first."
    exit 1
fi

echo "✅ Go is installed"
echo ""

# Compile the sync tool
echo "📦 Building XML sync tool..."
cd cmd/xmlsync || exit 1

if go build -o xmlsync main.go; then
    echo "✅ Build successful"
else
    echo "❌ Build failed"
    exit 1
fi

cd ../.. || exit 1
echo ""

# Offer options
echo "Select an option:"
echo ""
echo "1) Run XML sync (parse feed and create listing files)"
echo "2) Run XML sync in dry-run mode (preview changes)"
echo "3) Run Hugo development server"
echo "4) Build Hugo for production"
echo "5) Exit"
echo ""

read -p "Enter your choice (1-5): " choice

case $choice in
    1)
        echo ""
        echo "🔄 Syncing listings from XML feed..."
        ./cmd/xmlsync/xmlsync
        echo ""
        echo "✅ Sync complete!"
        echo ""
        echo "Next steps:"
        echo "1. Run 'hugo' to build the static site"
        echo "2. Or run 'hugo server' to preview"
        ;;
    2)
        echo ""
        echo "🔍 Dry-run mode (no files will be written)..."
        ./cmd/xmlsync/xmlsync -dry-run
        ;;
    3)
        echo ""
        echo "🚀 Starting Hugo development server..."
        echo "Open browser to: http://localhost:1313"
        echo ""
        hugo server -D
        ;;
    4)
        echo ""
        echo "🏗️  Building Hugo for production..."
        hugo -e production
        echo ""
        echo "✅ Build complete! Generated site in public/ directory"
        ;;
    5)
        echo "Goodbye!"
        exit 0
        ;;
    *)
        echo "Invalid choice"
        exit 1
        ;;
esac
