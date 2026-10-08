#!/usr/bin/env bash

# Check args
if [ "$#" -lt 2 ]; then
    echo "Usage: $0 <maps_folder> <images_folder> [src_folder]"
    exit 1
fi

MAPS_DIR="$1"
IMAGES_DIR="$2"
SRC_DIR="${3:-client/gui/src}"

# Check deps
if ! command -v jq &> /dev/null; then
    echo "Error: 'jq' isn't installed."
    exit 1
fi

# Get used image files from JSON maps
USED_IMAGES=$(find "$MAPS_DIR" -type f -name "*.json" -exec jq -r '
  (.. | .image? // empty)
' {} + | while IFS= read -r img_path; do
    if [ -n "$img_path" ]; then
        basename "$img_path"
    fi
done | sort -u)

# Delete unused images
while IFS= read -r -d '' img_file; do
    filename=$(basename "$img_file")
    name_without_ext="${filename%.*}"
    
    used_in_maps=0
    used_in_src=0
    
    # Check if file is used in maps
    if echo "$USED_IMAGES" | grep -qx "$filename"; then
        used_in_maps=1
    fi
    
    # Check if file is used in src
    if [ -d "$SRC_DIR" ]; then
        if grep -rE "[\"\`/]${name_without_ext}[\.\"\`]" "$SRC_DIR" | grep -vE "fmt\.Print|log\.|slog\.|println?\(" | grep -q . ; then
            used_in_src=1
        fi
    fi
    
    # Delete if not used
    if [ "$used_in_maps" -eq 0 ] && [ "$used_in_src" -eq 0 ]; then
        rm "$img_file"
    fi
done < <(find "$IMAGES_DIR" -type f \( -name "*.png" -o -name "*.jpg" -o -name "*.jpeg" -o -name "*.webp" \) -print0)

# Delete empty directories
find "$IMAGES_DIR" -type d -empty -delete
