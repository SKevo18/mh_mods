package transformers

import (
	"errors"
)

func packMhk4(dataFileLocation string, rootFolder string) error {
	return errors.New("not implemented")
}

func unpackMhk4(dataFileLocation string, rootFolder string) error {
	return errors.New("not implemented")
}

func transformMhk4(action string, dataFileLocation string, rootFolder string) error {
	switch action {
	case "pack":
		return packMhk4(dataFileLocation, rootFolder)
	case "unpack":
		return unpackMhk4(dataFileLocation, rootFolder)
	default:
		return errors.New("invalid action")
	}
}
