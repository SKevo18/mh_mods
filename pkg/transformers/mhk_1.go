package transformers

import (
	"errors"
	"os"
)

var xorKey = []byte{0x0C, 0x38, 0x4E, 0x41, 0x0C, 0x2B, 0x70, 0xB2, 0xD4, 0x04, 0x4C, 0x20, 0x6F}

// Packs MHK1 data file from given `rootFolder` into `dataFileLocation`.
// The input folder should contain the zip file that will be packed.
func packMhk1(dataFileLocation string, inputFolder string) error {
	// TODO: pack files in root folder into a zip file
	// tricky, because the ZIP file itself has some withcraftery in it

	// Read the input folder as ZIP
	dataBytes, err := readFile(inputFolder + "/mhke.zip")
	if err != nil {
		return err
	}

	// Apply XOR key to the data
	dataBytes = xorData(dataBytes, xorKey)

	// Open the data file for writing
	outputFile, err := openFileForWriting(dataFileLocation)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Write the packed data to the data file
	_, err = outputFile.Write(dataBytes)
	if err != nil {
		return err
	}

	return nil
}

// Unpacks MHK1 data file from `dataFileLocation` into `rootFolder`.
// The output folder will contain the unpacked zip file (`mhke.zip`)
func unpackMhk1(dataFileLocation string, outputFolder string) error {
	// Obtain the data file in byte form
	dataBytes, err := readFile(dataFileLocation)
	if err != nil {
		return err
	}

	// Apply XOR key to the data
	dataBytes = xorData(dataBytes, xorKey)

	// Create the output folder
	err = os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		return err
	}

	// Open the data file for writing
	outputFile, err := openFileForWriting(outputFolder + "/mhke.zip")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Write the packed data to the data file
	_, err = outputFile.Write(dataBytes)
	if err != nil {
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
