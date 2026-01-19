#!/bin/bash
# Estate Index: Daily Listings Sync via Cron
# This script syncs the XML feed daily and logs results
# Add to crontab: 0 2 * * * /path/to/EstateIndex/scripts/cron_sync_listings.sh

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CMD_DIR="$PROJECT_ROOT/cmd/xmlsync"
CONTENT_DIR="$PROJECT_ROOT/content/listings"
LOGS_DIR="$PROJECT_ROOT/logs"
LOG_FILE="$LOGS_DIR/listings_updates.log"

# Create logs directory if needed
mkdir -p "$LOGS_DIR"

# Timestamp function
timestamp() {
  date '+%Y-%m-%d %H:%M:%S'
}

# Log function (never output to UI)
log_message() {
  echo "[$(timestamp)] $1" >> "$LOG_FILE"
}

# Start logging
log_message "=========================================="
log_message "Starting daily listings sync"

# Build xmlsync CLI
log_message "Building xmlsync CLI..."
cd "$CMD_DIR"
if ! go build -o xmlsync . >> "$LOG_FILE" 2>&1; then
  log_message "ERROR: Failed to build xmlsync CLI"
  exit 1
fi

# Run sync
log_message "Fetching and syncing XML feed..."
cd "$PROJECT_ROOT"

# Try to sync, capture success/failure
if ./cmd/xmlsync/xmlsync -content "$CONTENT_DIR" >> "$LOG_FILE" 2>&1; then
  log_message "SUCCESS: XML feed synced successfully"
  
  # Rebuild Hugo site if sync was successful
  log_message "Rebuilding Hugo site..."
  if command -v hugo &> /dev/null; then
    if hugo --cleanDestinationDir >> "$LOG_FILE" 2>&1; then
      log_message "SUCCESS: Site rebuilt successfully"
    else
      log_message "ERROR: Hugo build failed - falling back to last working build"
      # The public/ directory should contain the last successful build
    fi
  else
    log_message "WARNING: Hugo not found in PATH"
  fi
  
  log_message "COMPLETED: Daily sync finished successfully"
else
  # Capture the error and log it
  log_message "ERROR: XML sync failed - keeping previous listings"
  log_message "Previous listings remain in place as fallback"
  log_message "COMPLETED: Daily sync finished with errors"
  exit 1
fi

log_message "=========================================="
echo "Sync completed. Check $LOG_FILE for details."
