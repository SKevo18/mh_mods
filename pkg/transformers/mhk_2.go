package transformers

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Reads header from MHK2 data file.
func readHeader(header []byte) (string, uint32) {
	return string(header[:0x20]), binary.LittleEndian.Uint32(header[0x20:])
}

// Generates header for MHK2 data file.
func generateHeader(name string, numFiles uint32) []byte {
	header := make([]byte, 0x40)
	copy(header, name)

	binary.LittleEndian.PutUint32(header[0x20:], numFiles)
	binary.LittleEndian.PutUint32(header[0x24:], 0x100)

	return header
}

// Encrypts MHK2 `.txt` config files.
func encryptConfig(data []byte) {
	key := uint(0x1234)
	for i := range data {
		uVar2 := data[i] & 0x55
		cVar1 := data[i] & 0xAA
		cVar1 >>= 1
		uVar2 <<= 1
		data[i] = (uVar2 ^ cVar1) ^ byte(key&0xFF)
		key = (key * 3) + 2&0xffff
	}
}

// Decrypts MHK2 `.txt` config files.
func decryptConfig(data []byte) {
	encryptConfig(data)
	encryptConfig(data)
}

// Packs MHK2 data file from given `inputFolder` into `dataFileLocation`.
func packMhk2(dataFileLocation string, inputPath string) error {
	outFile, err := os.Create(dataFileLocation)
	if err != nil {
		return err
	}
	defer outFile.Close()

	files, err := walkFiles(inputPath)
	if err != nil {
		return err
	}

	// header
	header := generateHeader("Moorhuhn", uint32(len(files)))
	_, err = outFile.Write(header)

	// save entries
	offset := int64(0x40 + (len(files) * 0x80))
	for _, file := range files {
		byteFileEntry := make([]byte, 0x80)
		relativePath, _ := filepath.Rel(inputPath, file.Filename)
		copy(byteFileEntry, strings.ReplaceAll(relativePath, "/", "\\"))

		binary.LittleEndian.PutUint64(byteFileEntry[0x68:], uint64(offset))
		binary.LittleEndian.PutUint64(byteFileEntry[0x6C:], uint64(file.Filesize))

		offset += file.Filesize + (file.Filesize % 0x100)
		outFile.Write(byteFileEntry)
	}

	// write data
	for _, file := range files {
		fileData, err := os.ReadFile(inputPath + "/" + file.Filename)
		if err != nil {
			return err
		}

		if filepath.Ext(file.Filename) == ".txt" {
			encryptConfig(fileData)
		}

		padding := make([]byte, file.Filesize%0x100)
		outFile.Write(append(fileData, padding...))
	}

	return nil
}

// Unpacks MHK2 data file from `dataFileLocation` into `outputFolder`.
func unpackMhk2(dataFileLocation string, outputPath string) error {
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}
	inFile, err := os.Open(dataFileLocation)
	if err != nil {
		return err
	}
	defer inFile.Close()

	// header
	header := make([]byte, 0x40)
	_, err = inFile.Read(header)
	if err != nil {
		return err
	}
	_, numFiles := readHeader(header)

	for i := uint32(0); i < numFiles; i++ {
		// read file entry
		fileEntry := make([]byte, 0x80)
		if _, err = inFile.Read(fileEntry); err != nil {
			return err
		}
		filename := string(fileEntry[:0x68])
		filesize := binary.LittleEndian.Uint64(fileEntry[0x6C:])

		// read data
		fileData := make([]byte, filesize)
		if _, err = inFile.Read(fileData); err != nil {
			return err
		}

		// decrypt config
		if filepath.Ext(filename) == ".txt" {
			decryptConfig(fileData)
		}

		// write
		path := outputPath + "/" + strings.ReplaceAll(filename, "\\", "/")
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}
		if err = os.WriteFile(path, fileData, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// Generic function to pack or unpack MHK2 data files.
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
