#!/bin/sh
# Run this script with data files placed inside `../data` directory
# The files should be named by their game IDs. The game ID is determined
# by the file name before first dot. Everything after the first dot is ignored,
# and used to determine the version of the data file, or other metadata
# (sometimes, the game has multiple versions, so there can be multiple data files for same game ID).
# 
# Example: `../data/mhk_2.en.dat` for English version of Moorhuhn Kart 2
# Or `../data/mhk_3.sar` for a version of Moorhuhn Kart 3
#
# This script requires the `idlemod` binary to be placed inside `../build` directory.

cd "$(dirname "$0")/../"

echo "This will unpack all data files under the 'data' directory into 'unpacked' directory."
echo "Type 'y' to continue:"
read -r confirm
if [ "$confirm" != "y" ]; then
    echo "Aborted."
    exit 1
fi

echo "Unpacking all data files..."
for file in data/*; do
    game_id=$(basename "$file" | cut -d '.' -f 1)
    if [ "$game_id" = "README" ]; then
        continue
    fi

    echo "Unpacking $file for game $game_id..."
    build/idlemod unpack "$game_id" "$file" "./unpacked/$game_id"
done

echo "Done."

