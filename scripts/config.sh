#!/bin/bash
# Load environment variables from .env file if it exists
# This is sourced by build scripts to set up the environment

if [ -f .env ]; then
  export $(cat .env | grep -v '^#' | xargs)
else
  echo "⚠️  .env file not found. Using defaults from .env.example"
  echo "   Copy .env.example to .env and customize for your deployment."
fi

# Export defaults if not set
export BASE_URL="${BASE_URL:-https://estateindex.com/}"
export HUGO_ENV="${HUGO_ENV:-development}"

# Ensure trailing slash on BASE_URL
if [[ ! "$BASE_URL" =~ /$ ]]; then
  export BASE_URL="${BASE_URL}/"
fi

echo "✓ Configuration loaded"
echo "  Base URL: $BASE_URL"
echo "  Environment: $HUGO_ENV"
