package util

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"unicode/utf8"
)

type Lines map[int]string

// Merges files recursively based on the original directory structure
func MergeModFilesRecursively(originalDir string, modDirs []string, destDir string) error {
	return filepath.Walk(originalDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ensure it's a file
		if info.IsDir() {
			return nil
		}

		// paths
		relPath, err := filepath.Rel(originalDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, relPath)

		// copy file, so it's always present regardless if modified or not
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		copyFile(path, destPath)

		// can only merge UTF-8 files
		if isUTF8File(path) {
			// collect existing modified file paths
			existingModPaths := []string{}
			for _, modDir := range modDirs {
				modPath := filepath.Join(modDir, relPath)
				if _, err := os.Stat(modPath); !os.IsNotExist(err) {
					existingModPaths = append(existingModPaths, modPath)
				}
			}

			if len(existingModPaths) > 0 {
				// if there are modded versions, merge them into the original file
				return MergeModFiles(path, existingModPaths, destPath)
			}
		}
		return nil
	})
}

// Merges mod files into the original file lines, always preferring newer changes
func MergeModFiles(originalPath string, modPaths []string, destPath string) error {
	// original lines
	originalLines, err := readFileLines(originalPath)
	if err != nil {
		return err
	}

	// merge
	modifiedLines := make(Lines)
	for _, modPath := range modPaths {
		modLines, err := readFileLines(modPath)
		if err != nil {
			return err
		}
		for lineNum, modLine := range modLines {
			if originalLine, exists := originalLines[lineNum]; !exists || (exists && modLine != originalLine) {
				modifiedLines[lineNum] = modLine
			}
		}
	}

	// write merged
	return writeMergedContent(destPath, originalLines, modifiedLines)
}

func writeMergedContent(destPath string, originalLines, modifiedLines Lines) error {
	// open
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()
	writer := bufio.NewWriter(destFile)

	// find max line number
	maxLine := maxKey(originalLines)
	if maxModLine := maxKey(modifiedLines); maxModLine > maxLine {
		maxLine = maxModLine
	}

	// write
	for i := 1; i <= maxLine; i++ {
		if line, modified := modifiedLines[i]; modified {
			_, err = writer.WriteString(line + "\n")
		} else if line, exists := originalLines[i]; exists {
			_, err = writer.WriteString(line + "\n")
		}
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func readFileLines(filePath string) (Lines, error) {
	// open
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// read
	lines := make(Lines)
	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		lines[lineNum] = scanner.Text()
		lineNum++
	}
	return lines, nil
}

// Simply copies a file from `src` to `dst`
func copyFile(src string, dst string) error {
	// open `src`
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// create `dest`
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// copy
	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

// Determines if 10 KB of the file are valid UTF-8 characters
func isUTF8File(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 1024*10)
	_, err = file.Read(buf)
	return err == nil && utf8.Valid(buf)
}

// Returns maximum line number in array of lines
func maxKey(lines Lines) int {
	max := 0
	for k := range lines {
		if k > max {
			max = k
		}
	}
	return max
}
