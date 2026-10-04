#!/usr/bin/env bash

# Check args
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <maps_folder> <images_folder>"
    exit 1
fi

MAPS_DIR="$1"
IMAGES_DIR="$2"

# Check deps
if ! command -v jq &> /dev/null; then
    echo "Error: 'jq' isn't installed."
    exit 1
fi

# Get used image files from JSON maps
USED_IMAGES=$(find "$MAPS_DIR" -type f -name "*.json" -exec jq -r '
  (.. | .image? // empty)
' {} + | xargs -n1 basename | sort -u)

# Delete unused images
while IFS= read -r -d '' img_file; do
    filename=$(basename "$img_file")
    
    # Check if file is used
    if ! echo "$USED_IMAGES" | grep -qx "$filename"; then
        rm "$img_file"
    fi
done < <(find "$IMAGES_DIR" -type f \( -name "*.png" -o -name "*.jpg" -o -name "*.jpeg" -o -name "*.webp" \) -print0)
