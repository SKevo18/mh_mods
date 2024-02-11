package transformers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Packs MHK2 data file from given `inputFolder` into `dataFileLocation`.
func packMhk2(dataFileLocation string, inputPath string) error {
    log.Printf("Creating `%s`...", dataFileLocation)
    outFile, err := os.Create(dataFileLocation)
    if err != nil {
        return err
    }
    defer outFile.Close()

    files, err := walkFiles(inputPath)
    if err != nil {
        return err
    }
    log.Printf("Packing %d files...", len(files))

    // header
    log.Println("Generating header...")
    header := generateHeader("Moorhuhn", uint32(len(files)))
    if _, err = outFile.Write(header); err != nil {
        return err
    }

    // save file entry data
    log.Println("Saving file entries data...")
    offset := int64(0x40) + int64(len(files)*0x80)
    for _, file := range files {
        fileEntry := make([]byte, 0x80)
        relativePath, _ := filepath.Rel(inputPath, file.Filename)
        pathWithPrefix := strings.ReplaceAll(relativePath, "/", "\\")
        copy(fileEntry, pathWithPrefix)

        binary.LittleEndian.PutUint64(fileEntry[0x68:], uint64(offset))
        binary.LittleEndian.PutUint64(fileEntry[0x6C:], uint64(file.Filesize))

        offset += file.Filesize + (file.Filesize % 0x100)
        if _, err := outFile.Write(fileEntry); err != nil {
            return err
        }
    }

    log.Println("Writing file data...")
    for _, file := range files {
        fileData, err := os.ReadFile(filepath.Join(inputPath, file.Filename))
        if err != nil {
            return err
        }

		// encrypt, if necessary
        if filepath.Ext(file.Filename) == ".txt" {
            log.Printf("Encrypting `%s`...", file.Filename)
            encryptConfig(fileData)
        }

		// write data
        paddingLength := 0x100 - (len(fileData) % 0x100)
        if paddingLength == 0x100 {
            paddingLength = 0 // no padding needed if `fileData` is already aligned
        }
        padding := make([]byte, paddingLength)

		log.Printf("Writing `%s`...", file.Filename)
        if _, err := outFile.Write(append(fileData, padding...)); err != nil {
            return err
        }
    }

    return nil
}

// Unpacks MHK2 data file from `dataFileLocation` into `outputFolder`.
func unpackMhk2(dataFileLocation string, outputPath string) error {
	log.Printf("Unpacking data file `%s`...", dataFileLocation)
	dataFile, err := os.Open(dataFileLocation)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	// read header
	header := make([]byte, 0x40)
	if _, err := dataFile.Read(header); err != nil {
		return err
	}
	dataFileName, numFiles := readHeader([0x40]byte(header))
	log.Printf("Data file name: `%s`, files: %d", dataFileName, numFiles)

	// read file entries
	log.Printf("Reading %d file entries...", numFiles)
	fileEntry := make([]byte, 0x80)
	for i := uint32(0); i < numFiles; i++ {
		if _, err := dataFile.Read(fileEntry); err != nil {
			return err
		}

		// read file position
		indexOffset, err := dataFile.Seek(0x0, io.SeekCurrent)
		if err != nil {
			return err
		}

		// read file entry
		fileName := getFilename(fileEntry)
		fileSize, err := getFileLength(fileEntry)
		if err != nil {
			return err
		}

		// where the file data begins
		dataPosition, err := getPosition(fileEntry)
		if err != nil {
			return err
		}
		log.Printf("Unpacking packed file `%s`...", fileName)

		// seek to file data
		if _, err := dataFile.Seek(int64(dataPosition), io.SeekStart); err != nil {
			return err
		}

		// read data
		fileData := make([]byte, fileSize)
		if _, err := dataFile.Read(fileData); err != nil {
			return err
		}

		// decrypt, if necessary
		if filepath.Ext(fileName) == ".txt" {
			decryptConfig(fileData)
		}

		// create output directory
		outputFilePath := filepath.Join(outputPath, strings.ReplaceAll(fileName, "\\", "/")) // the path in the data file uses backslashes
		os.MkdirAll(filepath.Dir(outputFilePath), os.ModePerm)

		// write decoded output
		os.WriteFile(outputFilePath, fileData, os.ModePerm)

		// seek back to file entry
		if _, err := dataFile.Seek(indexOffset, io.SeekStart); err != nil {
			return err
		}
	}

	return nil
}

// Extracts the file position from a file entry block.
func getPosition(fileEntry []byte) (uint32, error) {
	if len(fileEntry) < 104 {
		return 0, errors.New("Input data must have at least 104 bytes")
	}
	return binary.LittleEndian.Uint32(fileEntry[0x68:0x6C]), nil
}

// Extracts the file length from a file entry block.
func getFileLength(fileEntry []byte) (uint32, error) {
	if len(fileEntry) < 108 {
		return 0, errors.New("Input data must have at least 108 bytes")
	}
	return binary.LittleEndian.Uint32(fileEntry[0x6C:0x70]), nil
}

// Extracts the filename from a file entry block.
func getFilename(fileEntry []byte) string {
	i := bytes.IndexByte(fileEntry, 0x0)
	if i == -(0x1) {
		i = len(fileEntry)
	}
	return string(fileEntry[:i])
}

// Generates header for MHK2 data file.
func generateHeader(name string, numFiles uint32) []byte {
	header := make([]byte, 0x40)
	copy(header, name)

	binary.LittleEndian.PutUint32(header[0x20:], numFiles)
	binary.LittleEndian.PutUint32(header[0x24:], 0x100)

	return header
}

// Reads header from MHK2 data file. Returns data file name and number of files.
func readHeader(header [0x40]byte) (string, uint32) {
	name := strings.Trim(string(header[:32]), "\x00")
	numFiles := binary.LittleEndian.Uint32(header[0x20:0x24])

	return name, numFiles
}

// Encrypts MHK2 `.txt` config files.
func encryptConfig(data []byte) {
	key := uint(0x1234)

	for i := range data {
		uVar2 := data[i] & 0x55
		cVar1 := data[i] & 0xAA
		cVar1 >>= 0x1
		uVar2 <<= 0x1
		data[i] = (uVar2 ^ cVar1) ^ byte(key & 0xFF)
		key = (key * 0x3) + (0x2 & 0xffff)
	}
}

// Decrypts MHK2 `.txt` config files.
func decryptConfig(data []byte) []byte {
	var uVar3 uint16 = 0x1234

	for i := range data {
		uVar2 := (uVar3 & 0xFF) ^ uint16(data[i])
		cVar1 := uVar2 & 0xFF
		uVar3 = ((uVar3 * 0x3) + 0x2) & 0xFFFF

		data[i] = byte(((uVar2 >> 0x1) ^ (cVar1 << 0x1)) & 0xFF)
		data[i] &= 0x55
		data[i] ^= byte((cVar1 << 0x1) & 0xFF)
	}

	return data
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
