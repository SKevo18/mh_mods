package cmd

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"mhmods/pkg/transformers"

	"github.com/spf13/cobra"
)

func PackmodCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "packmod <game ID> <original data file> <output modded data file> <mod paths>...",
		Short: "Pack all mod paths into a single data file",
		Long:  `Packs all mod paths into a single data file for a specific game.`,
		Args:  cobra.MinimumNArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			gameID := args[0]
			originalDataFile := args[1]
			outputDataFile := args[2]
			modPaths := args[3:]

			// create temp dir
			tempDir, err := os.MkdirTemp("", "mhmods_temp")
			if err != nil {
				log.Fatalf("Fatal error: %s", err)
			}
			defer os.RemoveAll(tempDir)

			// unpack, merge and repack
			if err := transformers.Transform("unpack", gameID, originalDataFile, tempDir); err != nil {
				log.Fatalf("Fatal error during unpacking: %s", err)
			}
			if err := mergeRecursively(tempDir, modPaths...); err != nil {
				log.Fatalf("Fatal error during merging: %s", err)
			}
			if err := transformers.Transform("pack", gameID, outputDataFile, tempDir); err != nil {
				log.Fatalf("Fatal error during repacking: %s", err)
			}

			log.Printf("Packed modded data file: %s (paths: %v)", outputDataFile, modPaths)
		},
	}
}

// Merges directories and their file contents recursively.
func mergeRecursively(dest string, srcs ...string) error {
	// key		= relative path
	// value 	= set of unique lines
	allUniqueLines := make(map[string]map[string]struct{})

	for _, src := range srcs {
		err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(src, path)
			if err != nil {
				return err
			}

			lines, err := readLines(path)
			if err != nil {
				return err
			}

			// setup struct
			if allUniqueLines[relPath] == nil {
				allUniqueLines[relPath] = make(map[string]struct{})
			}
			for _, line := range lines {
				allUniqueLines[relPath][line] = struct{}{}
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	// create/update files with merged content
	for relPath, linesSet := range allUniqueLines {
		var lines []string
		for line := range linesSet {
			lines = append(lines, line)
		}

		destPath := filepath.Join(dest, relPath)
		if err := writeLines(destPath, lines); err != nil {
			return err
		}
	}

	return nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(path string, lines []string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
