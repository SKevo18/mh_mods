package transformers

import (
	"encoding/binary"
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
	zipData, err := os.ReadFile(zipLocation)
	if err != nil {
		return err
	}
	checksum, err := xorData(zipData, dataFileLocation)
	if err != nil {
		return err
	}

	// append checksum to data file
	log.Printf("Appending checksum `%d` to `%s`...\n", checksum, dataFileLocation)
	if err := appendChecksum(dataFileLocation, checksum); err != nil {
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
	data, err := os.ReadFile(dataFileLocation)
	if err != nil {
		return err
	}
	_, err = xorData(data[:len(data)-4], zipLocation) // minus checksum at end
	if err != nil {
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
func xorData(dataBytes []byte, outputPath string) (uint32, error) {
	// xor and uint32 checksum
	checksum := uint32(0)
	keyLength := len(xorKey)
	for i := range dataBytes {
		dataBytes[i] ^= xorKey[i%keyLength]

		signedByte := int8(dataBytes[i])
		if signedByte < 0 {
			checksum -= uint32(signedByte * -1)
		} else {
			checksum += uint32(signedByte)
		}
	}

	// write
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return 0, err
	}
	if err := os.WriteFile(outputPath, dataBytes, os.ModePerm); err != nil {
		return 0, err
	}

	return checksum, nil
}

// Appends checksum to the end of a file.
func appendChecksum(dataFileLocation string, checksum uint32) error {
	dataFile, err := os.OpenFile(dataFileLocation, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	dataFile, err = os.OpenFile(dataFileLocation, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, checksum)
	if _, err = dataFile.Write(buf); err != nil {
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
		return errors.New("invalid action")
	}
}
