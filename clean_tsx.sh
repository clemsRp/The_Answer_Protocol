#!/usr/bin/env bash

# Check args
if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <maps_folder> <tilesets_folder> <root_folder>"
    exit 1
fi

MAPS_DIR="$1"
TILESETS_DIR="$2"
ROOT_DIR="$3"

# Check deps
if ! command -v jq &> /dev/null; then
    echo "Error: 'jq' isn't installed."
    exit 1
fi

# Get used tsx files
USED_TSX=$(find "$MAPS_DIR" -type f -name "*.json" -exec jq -r '.tilesets[]?.source // empty' {} + | xargs -n1 basename | sort -u)

# Delete unused tsx files
while IFS= read -r -d '' tsx_file; do
    filename=$(basename "$tsx_file")
    
    # Check if file is used
    if ! echo "$USED_TSX" | grep -qx "$filename"; then
        rm "$tsx_file"
    fi
done < <(find "$ROOT_DIR" -type f -name "*.tsx" -print0)
