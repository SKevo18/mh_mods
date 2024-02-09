package transformers

import (
	"errors"
	"fmt"
	"os"
)

// Opens existing data file for reading in byte mode.
// If the file does not exist, an error is returned.
func openDataFileForReading(dataFileLocation string) (*os.File, error) {
	file, err := os.Open(dataFileLocation)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error opening data file: %s", err))
	}

	return file, nil
}

// Reads the data file into a byte array.
// Returns an error if the file cannot be read.
func readDataFile(dataFileLocation string) ([]byte, error) {
	file, err := openDataFileForReading(dataFileLocation)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dataFileInfo, err := file.Stat()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading data file: %s", err))
	}

	dataBytes := make([]byte, dataFileInfo.Size())
	_, err = file.Read(dataBytes)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading data file: %s", err))
	}

	return dataBytes, nil

}

// Opens new data file for writing in byte mode.
// If the file does not exist, it is created.
// If the file already exists, it is truncated to zero length.
func openDataFileForWriting(dataFileLocation string) (*os.File, error) {
	file, err := os.Create(dataFileLocation)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error creating data file: %s", err))
	}

	return file, nil
}

// A simple XOR function that applies a key to the data (used in MHK Extra).
func xorData(data []byte, key []byte) []byte {
	keyLength := len(key)
	for i := range data {
		data[i] ^= key[i%keyLength]
	}
	return data
}

// Dynamically determines the appropriate (un)packing function based on the game ID.
// Action can either be `pack` or `unpack`.
func Transform(action string, gameId string, dataFileLocation string, rootFolder string) error {
	var transformFunction func(string, string, string) error

	switch gameId {
	case "mhk_extra":
		transformFunction = transformMhk1
	case "mhk_2":
		transformFunction = transformMhk2
	case "mhk_3":
		transformFunction = transformMhk3
	case "mhk_4":
		transformFunction = transformMhk4
	default:
		return errors.New("Invalid game ID")
	}

	var err error
	err = transformFunction(action, dataFileLocation, rootFolder)

	if err != nil {
		return errors.New(fmt.Sprintf("Error transforming data: %s", err))
	}

	return nil
}
