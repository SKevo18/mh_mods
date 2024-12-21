package transformers

import (
	"path/filepath"
	"testing"
)

func Test_PackMhk3(t *testing.T) {
	ok, err := comparePacked(t, "mhk_3", mhkFixtures, filepath.Join(fixturePackedFolder, "mhk", "sample_mhk_3"))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if !ok {
		t.Fatalf("Packed file does not match expected data.")
	}
}

func Test_UnpackMhk3(t *testing.T) {
	ok, err := compareUnpacked(t, "mhk_3", mhkFixtures, filepath.Join(fixturePackedFolder, "mhk", "sample_mhk_3"))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if !ok {
		t.Fatalf("Unpacked files do not match expected data.")
	}
}
