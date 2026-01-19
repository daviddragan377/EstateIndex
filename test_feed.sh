#!/bin/bash
set -e

echo "Building xmlsync tool..."
cd /workspaces/EstateIndex
go build -o cmd/xmlsync/xmlsync ./cmd/xmlsync/

echo ""
echo "Running xmlsync..."
./cmd/xmlsync/xmlsync -content ./content/listings

echo ""
echo "Checking generated listings..."
ls -la content/listings/ | head -20
