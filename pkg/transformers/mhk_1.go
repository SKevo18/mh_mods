package transformers

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

// XOR key used to pack/unpack MHK1 data files.
var xorKey = []byte{0x0C, 0x38, 0x4E, 0x41, 0x0C, 0x2B, 0x70, 0xB2, 0xD4, 0x04, 0x4C, 0x20, 0x6F}

// Packs MHK1 data file from given `inputPath` (path to unpacked root of `mhke.zip`)
// into `dataFileLocation` (path to result data file).
func packMhk1(dataFileLocation string, inputPath string) error {
	zipLocation := dataFileLocation + ".zip"

	// temp zip
	log.Printf("Zipping `%s`...\n", inputPath)
	if err := zipFolder(inputPath, zipLocation); err != nil {
		return err
	}

	// xor
	log.Printf("Applying XOR to `%s`...\n", zipLocation)
	if err := xorMhk1File(zipLocation, dataFileLocation); err != nil {
		return err
	}

	// remove temp zip
	if err := os.Remove(zipLocation); err != nil {
		return err
	}

	return nil
}

// Unpacks MHK1 data file from `dataFileLocation` into `outputPath`.
func unpackMhk1(dataFileLocation string, outputPath string) error {
	zipLocation := outputPath + ".zip"

	// xor
	log.Printf("Applying XOR to `%s`...\n", dataFileLocation)
	if err := xorMhk1File(dataFileLocation, zipLocation); err != nil {
		return err
	}

	// unzip
	log.Printf("Unzipping `%s`...\n", zipLocation)
	if err := unzipFile(zipLocation, outputPath); err != nil {
		return err
	}

	// remove temp zip
	if err := os.Remove(zipLocation); err != nil {
		return err
	}

	return nil
}

// Applies XOR operation on MHK 1 file at `toXorPath` and writes result to `outputPath`.
// XOR is symmetric, so this function can be used for both packing and unpacking.
func xorMhk1File(toXorPath string, outputPath string) error {
	// read
	dataBytes, err := os.ReadFile(toXorPath)
	if err != nil {
		return err
	}

	// xor
	dataBytes = xorData(dataBytes, xorKey)

	// write
	if err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return err
	}
	if err = os.WriteFile(outputPath, dataBytes, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func transformMhk1(action string, dataFileLocation string, rootFolder string) error {
	switch action {
	case "pack":
		return packMhk1(dataFileLocation, rootFolder)
	case "unpack":
		return unpackMhk1(dataFileLocation, rootFolder)
	default:
		return errors.New("Invalid action!")
	}
}
