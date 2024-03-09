#!/bin/bash
# This script compares two directories and lists the files that are different.
# Useful for packing only modded files, if you modded them in place of the original files.

if [ "$#" -ne 2 ]; then
    echo "Compares two directories and lists files that are different"
    echo "Usage: $0 <modified_dir> <original_dir>"
    exit 1
fi

dir1=$1
dir2=$2

changed_files=0
only_in_dir1=0
only_in_dir2=0
total_files=0

while IFS= read -r filepath; do
    total_files=$((total_files + 1))
    filepath_in_dir2="${filepath/$dir1/$dir2}"
    
    if [ -e "$filepath_in_dir2" ]; then
        if ! diff -q "$filepath" "$filepath_in_dir2" > /dev/null; then
            echo "Changed: $filepath <=> $filepath_in_dir2"
            changed_files=$((changed_files + 1))
        fi
    else
        echo "Only in $dir1: $filepath"
        only_in_dir1=$((only_in_dir1 + 1))
    fi
done < <(find "$dir1" -type f)

while IFS= read -r filepath; do
    filepath_in_dir1="${filepath/$dir2/$dir1}"
    if [ ! -e "$filepath_in_dir1" ]; then
        echo "Only in $dir2: $filepath"
        only_in_dir2=$((only_in_dir2 + 1))
    fi
done < <(find "$dir2" -type f)

echo ""
echo "Summary:"
echo "Total files checked: $total_files"
echo "Changed files: $changed_files"
echo "Files only in $dir1: $only_in_dir1"
echo "Files only in $dir2: $only_in_dir2"
