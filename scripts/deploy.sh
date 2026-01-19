#!/bin/bash
# Estate Index: Deploy to static host
# Supports multiple deployment targets

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BUILD_DIR="$PROJECT_ROOT/public"

echo "Estate Index: Deploy"
echo "===================="
echo ""

# Check build directory exists
if [ ! -d "$BUILD_DIR" ]; then
  echo "Error: Build directory not found at $BUILD_DIR"
  echo "Run ./scripts/build.sh first"
  exit 1
fi

# Parse deployment target
TARGET="${1:-help}"

case $TARGET in
  netlify)
    echo "Deploying to Netlify..."
    if command -v netlify &> /dev/null; then
      netlify deploy --prod --dir="$BUILD_DIR"
      echo "✓ Deployed to Netlify"
    else
      echo "Error: netlify-cli not found. Install with: npm install -g netlify-cli"
      exit 1
    fi
    ;;
  
  vercel)
    echo "Deploying to Vercel..."
    if command -v vercel &> /dev/null; then
      vercel --prod
      echo "✓ Deployed to Vercel"
    else
      echo "Error: vercel CLI not found. Install with: npm install -g vercel"
      exit 1
    fi
    ;;
  
  s3)
    echo "Deploying to AWS S3..."
    if command -v aws &> /dev/null; then
      BUCKET="${2:-estateindex-static}"
      aws s3 sync "$BUILD_DIR" "s3://$BUCKET" --delete
      echo "✓ Deployed to S3 bucket: $BUCKET"
    else
      echo "Error: AWS CLI not found. Install with: pip install awscli"
      exit 1
    fi
    ;;
  
  local)
    echo "Local deployment (copying files)..."
    DEST="${2:-.}"
    mkdir -p "$DEST"
    cp -r "$BUILD_DIR"/* "$DEST/"
    echo "✓ Files copied to: $DEST"
    ;;
  
  *)
    echo "Usage: $0 <target> [options]"
    echo ""
    echo "Supported targets:"
    echo "  netlify          Deploy to Netlify (requires netlify-cli)"
    echo "  vercel           Deploy to Vercel (requires vercel CLI)"
    echo "  s3 <bucket>      Deploy to AWS S3 (requires aws-cli)"
    echo "  local <path>     Copy to local directory"
    echo ""
    echo "Examples:"
    echo "  $0 netlify"
    echo "  $0 s3 my-bucket"
    echo "  $0 local /var/www/html"
    exit 0
    ;;
esac

echo ""
echo "Deployment complete!"
