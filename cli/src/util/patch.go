package util

import (
	"github.com/SKevo18/gopatch"
)

func PatchModFiles(rootDir string, outputDir string, patchFilePaths []string) error {
	patchLines, err := gopatch.ReadPatchFiles(patchFilePaths)
	if err != nil {
		return err
	}

	if err := gopatch.PatchDir(rootDir, outputDir, patchLines); err != nil {
		return err
	}

	return nil
}
