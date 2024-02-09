package transformers

import (
	"errors"
)

func packMhk2(dataFileLocation string, rootFolder string) error {
	return errors.New("not implemented")
}

func unpackMhk2(dataFileLocation string, rootFolder string) error {
	return errors.New("not implemented")
}

func transformMhk2(action string, dataFileLocation string, rootFolder string) error {
	switch action {
	case "pack":
		return packMhk2(dataFileLocation, rootFolder)
	case "unpack":
		return unpackMhk2(dataFileLocation, rootFolder)
	default:
		return errors.New("Invalid action!")
	}
}
