package transformers

import (
	"errors"
)

func packMhk3(dataFileLocation string, rootFolder string) error {
	return errors.New("not implemented")
}

func unpackMhk3(dataFileLocation string, rootFolder string) error {
	return errors.New("not implemented")
}

func transformMhk3(action string, dataFileLocation string, rootFolder string) error {
	switch action {
	case "pack":
		return packMhk3(dataFileLocation, rootFolder)
	case "unpack":
		return unpackMhk3(dataFileLocation, rootFolder)
	default:
		return errors.New("Invalid action!")
	}
}
