package transformers

import (
	"path/filepath"
	"testing"
)

var (
	correctPacked = filepath.Join(fixturePackedFolder, "sample_dino_aliens")
	correctUnpacked = filepath.Join(fixtureAssetsFolder, "dino_aliens", "Data")
)

func Test_PackDinoAliens(t *testing.T) {
	tempPacked := filepath.Join(t.TempDir(), "sample_dino_aliens")
	if err := packDinoAliens(correctUnpacked, tempPacked); err != nil {
		t.Fatalf("Error: %v", err)
	}

	if ok, err := compareDirectories(tempPacked, correctPacked); err != nil {
		t.Fatalf("Error: %v", err)
	} else if !ok {
		t.Fatalf("Directories %s and %s do not match", tempPacked, correctPacked)
	}
}

func Test_UnpackDinoAliens(t *testing.T) {
	tempUnpacked := t.TempDir()
	if err := unpackDinoAliens(correctPacked, tempUnpacked); err != nil {
		t.Fatalf("Error: %v", err)
	}

	if ok, err := compareDirectories(tempUnpacked, correctUnpacked); err != nil {
		t.Fatalf("Error: %v", err)
	} else if !ok {
		t.Fatalf("Directories %s and %s do not match", tempUnpacked, correctUnpacked)
	}
}
