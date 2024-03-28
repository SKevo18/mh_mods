package transformers_test

import (
	"os"

	"mhmods/src/transformers"
)

const (
	fixtureAssetsFolder = "fixture/assets"
	fixturePackedFolder   = "fixture/packed"
)

func compareUnpacked(gameId string) (bool, error) {
	packedPath := fixturePackedFolder + "/sample_" + gameId

	// create temp unpacked dir
	tempDir := os.TempDir() + "/mhmods_test_" + gameId
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return false, err
	}
	defer os.RemoveAll(tempDir)

	// unpack
	if err := transformers.Transform("unpack", gameId, packedPath, tempDir); err != nil {
		return false, err
	}

	// compare
	ok, err := compareDirectories(fixtureAssetsFolder, tempDir)
	if err != nil {
		return false, err
	}

	return ok, nil
}


func comparePacked(gameId string) (bool, error) {
	packedPath := fixturePackedFolder + "/sample_" + gameId

	// create temp packed file
	temp, err := os.CreateTemp("", "mhmods_test_packed_"+gameId)
	if err != nil {
		return false, err
	}
	defer os.Remove(temp.Name())

	// pack
	if err := transformers.Transform("pack", gameId, temp.Name(), fixtureAssetsFolder); err != nil {
		return false, err
	}

	// compare
	ok, err := compareFiles(packedPath, temp.Name())
	if err != nil {
		return false, err
	}

	return ok, nil
}
