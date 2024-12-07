package transformers_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// Compares two files.
func compareFiles(file1, file2 string) (bool, error) {
	content1, err := os.ReadFile(file1)
	if err != nil {
		return false, fmt.Errorf("failed to read file %s: %w", file1, err)
	}
	content2, err := os.ReadFile(file2)
	if err != nil {
		return false, fmt.Errorf("failed to read file %s: %w", file2, err)
	}

	return reflect.DeepEqual(content1, content2), nil
}

// Compare two directories and their contents.
func compareDirectories(dir1, dir2 string) (bool, error) {
	files1, err := os.ReadDir(dir1)
	if err != nil {
		return false, fmt.Errorf("failed to read directory %s: %w", dir1, err)
	}
	files2, err := os.ReadDir(dir2)
	if err != nil {
		return false, fmt.Errorf("failed to read directory %s: %w", dir2, err)
	}

	if len(files1) != len(files2) {
		return false, nil
	}

	for i := range files1 {
		fileName1 := files1[i].Name()
		fileName2 := files2[i].Name()
		if fileName1 != fileName2 {
			return false, nil
		}

		filePath1 := filepath.Join(dir1, fileName1)
		filePath2 := filepath.Join(dir2, fileName2)
		if files1[i].IsDir() != files2[i].IsDir() {
			return false, nil
		}

		if files1[i].IsDir() {
			same, err := compareDirectories(filePath1, filePath2)
			if err != nil {
				return false, err
			}
			if !same {
				return false, nil
			}
		} else {
			content1, err := os.ReadFile(filePath1)
			if err != nil {
				return false, fmt.Errorf("failed to read file %s: %w", filePath1, err)
			}

			content2, err := os.ReadFile(filePath2)
			if err != nil {
				return false, fmt.Errorf("failed to read file %s: %w", filePath2, err)
			}

			if !reflect.DeepEqual(content1, content2) {
				return false, nil
			}
		}
	}

	return true, nil
}
