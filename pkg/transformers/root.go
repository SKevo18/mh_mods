package transformers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// A simple struct to represent a file entry.
type FileEntry struct {
	// The path of the file relative to the root folder.
	FilePath string
	// The size of the file in bytes.
	FileSize int64
	// Content position, in the data file
	ContentOffset int64
}

// Walks through files in `rootFolder` and returns array of relative file entries.
func walkFiles(rootFolder string) ([]FileEntry, error) {
	var files []FileEntry

	err := filepath.Walk(rootFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// append file
		if !info.IsDir() {
			relPath, err := filepath.Rel(rootFolder, path)
			if err != nil {
				return err
			}

			files = append(files, FileEntry{
				FilePath: relPath,
				FileSize: info.Size(),
			})
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking through files: %s", err)
	}

	return files, nil
}

// Dynamically determines the appropriate transform function based on the game ID.
// Action can either be `pack` or `unpack`.
func Transform(action string, gameId string, dataFileLocation string, rootFolder string) error {
	var transformFunction func(string, string, string) error

	switch gameId {
	case "mhk_extra", "mhk_1", "mhke":
		transformFunction = transformMhk1
	case "mhk_2", "schatzjaeger": // should also work on Schatzjäger (Jump and Run), but untested
		transformFunction = transformMhk2
	case "mhk_3":
		transformFunction = transformMhk3
	case "mhk_4", "mhk_thunder":
		transformFunction = transformMhk4
	default:
		return errors.New("invalid game ID! Please, use one of the following: `mhk_extra`, `mhk_1`, `mhk_2`, `schatzjaeger`, `mhk_3`, `mhk_4`, `mhk_thunder`")
	}

	err := transformFunction(action, dataFileLocation, rootFolder)
	if err != nil {
		return fmt.Errorf("error transforming data: %s", err)
	}

	return nil
}
