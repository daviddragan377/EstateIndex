#!/bin/bash
# Estate Index: Sync XML feed and build static site
# This script orchestrates the complete build process

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CONTENT_DIR="$PROJECT_ROOT/content/listings"
CMD_DIR="$PROJECT_ROOT/cmd/xmlsync"
BUILD_DIR="$PROJECT_ROOT/public"

# Load configuration from .env
source "$SCRIPT_DIR/config.sh" || true

# Set defaults if config.sh didn't load them
BASE_URL="${BASE_URL:-https://estateindex.com/}"
HUGO_ENV="${HUGO_ENV:-development}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Estate Index Build System${NC}"
echo "========================================"
echo ""

# Parse arguments
DRY_RUN=false
SKIP_SYNC=false
while [[ $# -gt 0 ]]; do
  case $1 in
    --dry-run)
      DRY_RUN=true
      shift
      ;;
    --skip-sync)
      SKIP_SYNC=true
      shift
      ;;
    --help)
      echo "Usage: $0 [OPTIONS]"
      echo ""
      echo "Options:"
      echo "  --dry-run      Show what would be done without writing files"
      echo "  --skip-sync    Skip XML sync, only build site"
      echo "  --help         Show this help message"
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Step 1: Build xmlsync CLI
echo -e "${BLUE}Step 1: Building xmlsync CLI...${NC}"
cd "$CMD_DIR"
if go build -o xmlsync . 2>&1; then
  echo -e "${GREEN}✓ xmlsync CLI built successfully${NC}"
else
  echo -e "${RED}✗ Failed to build xmlsync CLI${NC}"
  exit 1
fi
echo ""

# Step 2: Sync XML feed
if [ "$SKIP_SYNC" = false ]; then
  echo -e "${BLUE}Step 2: Syncing XML feed...${NC}"
  
  SYNC_ARGS="-content $CONTENT_DIR"
  if [ "$DRY_RUN" = true ]; then
    SYNC_ARGS="$SYNC_ARGS --dry-run"
  fi
  
  if ./xmlsync $SYNC_ARGS; then
    echo -e "${GREEN}✓ XML sync completed${NC}"
  else
    echo -e "${RED}✗ XML sync failed${NC}"
    exit 1
  fi
else
  echo -e "${BLUE}Step 2: Skipping XML sync (--skip-sync)${NC}"
fi
echo ""

# Step 3: Build site with Hugo
if [ "$DRY_RUN" = false ]; then
  echo -e "${BLUE}Step 3: Building site with Hugo...${NC}"
  echo -e "  Base URL: $BASE_URL"
  cd "$PROJECT_ROOT"
  
  if command -v hugo &> /dev/null; then
    if hugo -b "$BASE_URL" -e "$HUGO_ENV" --cleanDestinationDir 2>&1; then
      echo -e "${GREEN}✓ Site built successfully${NC}"
      echo -e "  Output: $BUILD_DIR"
    else
      echo -e "${RED}✗ Hugo build failed${NC}"
      exit 1
    fi
  else
    echo -e "${RED}✗ Hugo not found in PATH${NC}"
    exit 1
  fi
else
  echo -e "${BLUE}Step 3: Skipping Hugo build (dry-run mode)${NC}"
fi
echo ""

# Step 4: Summary
echo -e "${GREEN}========================================"
echo "Build Complete!"
echo "========================================"

if [ "$DRY_RUN" = false ] && [ -d "$BUILD_DIR" ]; then
  FILE_COUNT=$(find "$BUILD_DIR" -type f | wc -l)
  echo -e "Site files: $FILE_COUNT"
  echo -e "Ready to deploy: $BUILD_DIR${NC}"
else
  echo -e "Dry run mode - no files written${NC}"
fi
