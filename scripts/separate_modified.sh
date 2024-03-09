#!/bin/bash
# This script leaves the original files unchanged and copies only
# files that have been modified or are extra to a new directory
#
# Also see `compare.sh` for "dry-run" version of this script.

if [ "$#" -ne 3 ]; then
    echo "Separates modified files from original files."
    echo "Also see 'compare.sh' for 'dry-run' version of this script."
    echo "Usage: $0 <modified_dir> <original_dir> <output_dir>"
    exit 1
fi

dir1=$1
dir2=$2
dir3=$3

mkdir -p "$dir3"

copy_to_dir3() {
    local src_file=$1
    local dest_file=${src_file/$dir1/$dir3}
    mkdir -p "$(dirname "$dest_file")"
    cp "$src_file" "$dest_file"
}

while IFS= read -r filepath; do
    filepath_in_dir2="${filepath/$dir1/$dir2}"
    
    if [ -e "$filepath_in_dir2" ]; then
        if ! diff -q "$filepath" "$filepath_in_dir2" > /dev/null; then
            copy_to_dir3 "$filepath"
        fi
    else
        copy_to_dir3 "$filepath"
    fi
done < <(find "$dir1" -type f)

echo "Files have been copied to: $dir3"
