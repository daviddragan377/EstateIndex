#!/bin/bash
# Estate Index: Watch mode for development
# Automatically rebuilds and syncs on changes

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "Estate Index: Watch Mode"
echo "========================"
echo ""
echo "Watching for changes... (Press Ctrl+C to stop)"
echo ""

# Check for file watcher
if command -v fswatch &> /dev/null; then
  fswatch -r "$PROJECT_ROOT/content" "$PROJECT_ROOT/assets" "$PROJECT_ROOT/layouts" \
    --exclude='\.DS_Store' \
    --exclude='\.git' | while read; do
    echo "[$(date '+%H:%M:%S')] Changes detected - rebuilding..."
    "$SCRIPT_DIR/build.sh" --skip-sync 2>&1 | grep -E '(Error|✓|Step)' || true
  done
elif command -v entr &> /dev/null; then
  find "$PROJECT_ROOT/content" "$PROJECT_ROOT/assets" "$PROJECT_ROOT/layouts" -type f | \
    entr bash -c "echo '[$(date '+%H:%M:%S')] Changes detected - rebuilding...' && $SCRIPT_DIR/build.sh --skip-sync 2>&1 | grep -E '(Error|✓|Step)' || true"
else
  echo "Error: fswatch or entr not found"
  echo "Install one with:"
  echo "  macOS: brew install fswatch"
  echo "  Ubuntu: apt-get install inotify-tools && npm install -g entr"
  exit 1
fi
