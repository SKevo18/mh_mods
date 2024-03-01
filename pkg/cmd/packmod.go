package cmd

import (
	"bufio"
	"fmt"
	"io"
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
			if err := mergeRecursively(tempDir, true, modPaths...); err != nil {
				log.Fatalf("Fatal error during merging: %s", err)
			}
			if err := transformers.Transform("pack", gameID, outputDataFile, tempDir); err != nil {
				log.Fatalf("Fatal error during repacking: %s", err)
			}

			log.Printf("Packed modded data file: %s (paths: %v)", outputDataFile, modPaths)
		},
	}
}

// Merges directories and file contents recursively
func mergeRecursively(dest string, ignoreConflicts bool, srcs ...string) error {
	for _, src := range srcs {
		err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// get relative path, to create the same structure in the destination
			relativePath, err := filepath.Rel(src, path)
			if err != nil {
				return err
			}
			destPath := filepath.Join(dest, relativePath)

			// copy
			if info.IsDir() {
				return os.MkdirAll(destPath, info.Mode())
			}

			// copy destPath file into temporary file to use for merge
			tempOriginal, err := os.CreateTemp("", fmt.Sprintf("mhmods_temp_%s", filepath.Base(destPath)))
			if err != nil {
				return err
			}
			defer os.Remove(tempOriginal.Name())

			original, err := os.Open(destPath)
			if err != nil {
				return err
			}
			defer original.Close()

			io.Copy(tempOriginal, original)

			// merge
			return mergeFileContents(destPath, path, tempOriginal.Name(), ignoreConflicts)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Merges the contents of two files into the destination file line by line.
func mergeFileContents(destFile string, srcFile1 string, srcFile2 string, ignoreConflicts bool) error {
	// open files
	file1, err := os.Open(srcFile1)
	if err != nil {
		return err
	}
	defer file1.Close()
	file2, err := os.Open(srcFile2)
	if err != nil {
		return err
	}
	defer file2.Close()

	// create destination file
	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer dest.Close()

	// scan
	scanner1 := bufio.NewScanner(file1)
	scanner2 := bufio.NewScanner(file2)

	for scanner1.Scan() && scanner2.Scan() {
		line1 := scanner1.Text()
		line2 := scanner2.Text()

		if err := scanner1.Err(); err != nil {
			return err
		}
		if err := scanner2.Err(); err != nil {
			return err
		}

		// equal, write first
		if line1 == line2 {
			_, err := dest.WriteString(line1 + "\n")
			if err != nil {
				return err
			}
		} else {
			// ignore conflict, choose first
			if ignoreConflicts {
				_, err := dest.WriteString(line1 + "\n")
				if err != nil {
					return err
				}
			} else {
				// error, if not ignored (currently unused)
				return fmt.Errorf("conflict detected and `ignoreConflicts` is false: `%s` vs `%s`", line1, line2)
			}
		}
	}

	return nil
}
