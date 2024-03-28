package util_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"mhmods/src/util"
)

const (
	originalPath = "fixture/original.txt"
	modPath      = "fixture/mod_%d.txt"
	mergedPath   = "fixture/merged.txt"
	expectedPath = "fixture/expected.txt"
)

func getModPaths(nMods int) []string {
	paths := make([]string, nMods)
	for i := 0; i < nMods; i++ {
		paths[i] = fmt.Sprintf(modPath, i+1)
	}
	return paths
}

func TestMergeModFiles(t *testing.T) {
	modPaths := getModPaths(3)

	// merge 3 dummy mods
	err := util.MergeModFiles(originalPath, modPaths, mergedPath)
	if err != nil {
		t.Error(err)
	}

	// compare merged to expected
	mergedContent, err := os.ReadFile(mergedPath)
	if err != nil {
		t.Error(err)
	}

	expectedContent, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, string(expectedContent), string(mergedContent))
}
