package transformers

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func packMhk4(dataFileLocation string, rootFolder string) error {
	return errors.New("not implemented")
}

func unpackMhk4(dataFileLocation string, outputDirectory string) error {
	dataFile, err := os.Open(dataFileLocation)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	// get file size
	fileInfo, err := dataFile.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()
	log.Printf("Data file size: %d", fileSize)

	// go to beginning of files
	var indexOffset uint32
	if _, err := dataFile.Seek(0x0C, io.SeekStart); err != nil {
		return err
	}

	binary.Read(dataFile, binary.LittleEndian, &indexOffset)
	dataFile.Seek(int64(indexOffset), io.SeekStart)

	// read file entries
	log.Printf("Reading file entries...")
	for {
		// if at end of file
		if currentOffset, _ := dataFile.Seek(0x0, io.SeekCurrent); currentOffset >= fileSize {
			break
		}

		// read filename
		var filenameLength uint8
		if err := binary.Read(dataFile, binary.LittleEndian, &filenameLength); err != nil {
			return err
		}
		filenameBytes := make([]byte, filenameLength)
		if _, err := dataFile.Read(filenameBytes); err != nil {
			return err
		}

		// make parent dir for file
		filename := strings.ReplaceAll(string(filenameBytes), ":", "")
		filename = strings.ReplaceAll(filename, "\\", "/")
		outputDirectory := filepath.Join(outputDirectory, filepath.Dir(filename))
		if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
			return err
		}
		outputPath := filepath.Join(outputDirectory, filepath.Base(filename))

		dataFile.Seek(0x1, io.SeekCurrent) // skip byte after filename

		// get file offset and length
		var fileOffset, fileLength uint32
		binary.Read(dataFile, binary.LittleEndian, &fileOffset)
		binary.Read(dataFile, binary.LittleEndian, &fileLength)
		nextOffset, _ := dataFile.Seek(0x4, io.SeekCurrent)

		// write file
		dataFile.Seek(int64(fileOffset), io.SeekStart)
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		io.CopyN(outputFile, dataFile, int64(fileLength))
		outputFile.Close()

		// go to next file
		dataFile.Seek(nextOffset, io.SeekStart)
		log.Printf("Filename: %s; Offset: %d; Size: %d", filename, fileOffset, fileLength)
	}

	return nil
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
