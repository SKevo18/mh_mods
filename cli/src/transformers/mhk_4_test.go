package transformers

import (
	"path/filepath"
	"testing"
)

func Test_PackMhk4(t *testing.T) {
	ok, err := comparePacked(t, "mhk_4", mhkFixtures, filepath.Join(fixturePackedFolder, "mhk", "sample_mhk_4"))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if !ok {
		t.Fatalf("Packed file does not match expected data.")
	}
}

func Test_UnpackMhk4(t *testing.T) {
	ok, err := compareUnpacked(t, "mhk_4", mhkFixtures, filepath.Join(fixturePackedFolder, "mhk", "sample_mhk_4"))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if !ok {
		t.Fatalf("Unpacked files do not match expected data.")
	}
}
