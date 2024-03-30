package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SKevo18/gopatch"
	cp "github.com/otiai10/copy"
)

func PatchModFiles(rootDir string, outputDir string, patchFilePaths []string) error {
	patchLines, err := gopatch.ReadPatchFiles(patchFilePaths)
	if err != nil {
		return err
	}

	if len(patchLines) == 0 {
		return fmt.Errorf("no valid patch files found in mod paths for root `%s`", rootDir)
	}
	if err := gopatch.PatchDir(rootDir, outputDir, patchLines); err != nil {
		return err
	}

	return nil
}

// Copies mod files from `modRootPaths` (looksup "source" directory here)
// to `outputDir` and returns a list of patch files found in the mod roots.
func CopyModFiles(modRootPaths []string, outputDir string) ([]string, error) {
	patchFilePaths := []string{}
	for _, modPath := range modRootPaths {
		patchFile := filepath.Join(modPath, "patch.txt")
		if _, err := os.Stat(patchFile); !os.IsNotExist(err) {
			patchFilePaths = append(patchFilePaths, patchFile)
		}

		sourceDir := filepath.Join(modPath, "source")
		if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
			continue
		}

		if err := cp.Copy(sourceDir, outputDir); err != nil {
			return nil, fmt.Errorf("fatal error while copying mods: %s", err)
		}
	}

	return patchFilePaths, nil
}
