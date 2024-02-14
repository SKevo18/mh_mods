#!/bin/bash
# In order to use this script, you need to build two files containing a list of
# paths to files, one per line. It is used to compare whether all the files that
# are expected to be in the unpacked directory are actually there.
#
# Context: MHK 4 has five files in the data index that are present twice, this script
# aided in identifying which were the duplicated files.

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <expected_file> <actual_file>"
    echo "Both files should contain a list of paths to files, one per line."
    exit 1
fi

expected_file="$1"
actual_file="$2"

# -23 suppresses lines that are unique to the second file and lines
# that are common, leaving only lines unique to the first file
comm -23 <(sort "$expected_file") <(sort "$actual_file")
