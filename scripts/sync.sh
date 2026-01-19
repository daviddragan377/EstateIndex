#!/bin/bash
# Estate Index: Sync XML feed only
# Useful for scheduled syncs without rebuilding the entire site

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CONTENT_DIR="$PROJECT_ROOT/content/listings"
CMD_DIR="$PROJECT_ROOT/cmd/xmlsync"

echo "Estate Index: XML Sync"
echo "======================"
echo ""

# Parse arguments
DRY_RUN=false
while [[ $# -gt 0 ]]; do
  case $1 in
    --dry-run)
      DRY_RUN=true
      shift
      ;;
    --help)
      echo "Usage: $0 [OPTIONS]"
      echo ""
      echo "Options:"
      echo "  --dry-run   Show what would be done without writing files"
      echo "  --help      Show this help message"
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Build xmlsync if not already built
if [ ! -f "$CMD_DIR/xmlsync" ]; then
  echo "Building xmlsync CLI..."
  cd "$CMD_DIR"
  go build -o xmlsync .
  echo ""
fi

# Run sync
cd "$CMD_DIR"
SYNC_ARGS="-content $CONTENT_DIR"
if [ "$DRY_RUN" = true ]; then
  SYNC_ARGS="$SYNC_ARGS --dry-run"
fi

./xmlsync $SYNC_ARGS

echo ""
echo "Sync complete!"
