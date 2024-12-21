package transformers

/*
Dino and Aliens contains multiple .dat files, that contain multiple files.
Each DAT file contains a file name terminated with 0x0, followed by the file size,
and then the file data. The file data is XORed with the base file name used as key.

Note: the original files in the DAT are not packed in alphabetical order, but this tool still produces a format that the game can read.
*/

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	cp "github.com/otiai10/copy"
)

var paths = []string{
	"animations.dat",
	"geomobjs.dat",
	"characters.dat",
	"texts.dat",
	"textures.dat",
	filepath.Join("levels", "demos.dat"),
	filepath.Join("levels", "levels.dat"),
	filepath.Join("levels", "preview.dat"),
}

func packDinoAliens(rootPath string, outputFolder string) error {
	copyNonDatFiles(rootPath, outputFolder)

	// pack DAT files
	for _, path := range paths {
		originalPath := filepath.Join(rootPath, path)
		packedPath := filepath.Join(outputFolder, path)
		if _, err := os.Stat(originalPath); os.IsNotExist(err) {
			fmt.Printf("Skipping %s because it doesn't exist - but note that the game won't boot correctly without all correct files!\n", originalPath)
			continue
		}

		if err := processPack(originalPath, packedPath); err != nil {
			return err
		}
	}
	return nil
}

// This function unpacks all DAT files in place and copies the rest of the files,
// retaining the original file structure.
func unpackDinoAliens(rootPath string, outputFolder string) error {
	copyNonDatFiles(rootPath, outputFolder)

	// unpack DAT files
	for _, path := range paths {
		originalPath := filepath.Join(rootPath, path)
		unpackPath := filepath.Join(outputFolder, path)
		if _, err := os.Stat(originalPath); os.IsNotExist(err) {
			continue
		}

		if err := processUnpack(originalPath, unpackPath); err != nil {
			return err
		}
	}
	return nil
}

func copyNonDatFiles(rootPath string, outputFolder string) error {
	if err := cp.Copy(rootPath, outputFolder, cp.Options{
		Skip: func(_ os.FileInfo, src string, _ string) (bool, error) {
			relativePath, err := filepath.Rel(rootPath, src)
			if err != nil {
				return false, err
			}
			return slices.Contains(paths, relativePath), nil
		},
	}); err != nil {
		return err
	}
	return nil
}

func processUnpack(inputFile string, outputPath string) error {
	// open
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	fullSize, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	for offset := int64(0); offset < fullSize; {
		// read packed file name
		fileName, err := readUntilNullByte(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// read packed file packedSize
		var packedSize int32
		if err := binary.Read(file, binary.LittleEndian, &packedSize); err != nil {
			return err
		}

		// save current offset
		currentOffset, err := file.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}

		// create out file
		outputFilePath := filepath.Join(outputPath, string(fileName))
		if err := os.MkdirAll(filepath.Dir(outputFilePath), os.ModePerm); err != nil {
			return err
		}
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		// XOR data
		outBuffer := make([]byte, packedSize)
		if err := binary.Read(file, binary.LittleEndian, &outBuffer); err != nil {
			return err
		}
		for i := range outBuffer {
			outBuffer[i] ^= fileName[i%len(fileName)]
		}

		// write
		if _, err := outputFile.Write(outBuffer); err != nil {
			return err
		}

		// jump to next file
		offset = currentOffset + int64(packedSize)
		if _, err := file.Seek(offset, io.SeekStart); err != nil {
			return err
		}
	}

	return nil
}

func processPack(inputDatFolder string, outputPath string) error {
	// make parent folder
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return err
	}

	// create output DAT file
	outputDatFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputDatFile.Close()

	// read files in unpacked DAT
	err = filepath.Walk(inputDatFolder, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// read file data
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		// XOR file data with file name
		fileName := filepath.Base(filePath)
		for i := range fileData {
			fileData[i] ^= fileName[i%len(fileName)]
		}

		// write file name with null termination
		if _, err := outputDatFile.Write(append([]byte(fileName), 0x0)); err != nil {
			return err
		}

		// write file size
		packedSize := int32(len(fileData))
		if err := binary.Write(outputDatFile, binary.LittleEndian, packedSize); err != nil {
			return err
		}

		// write XORed file data
		if _, err := outputDatFile.Write(fileData); err != nil {
			return err
		}

		return nil
	})

	return err
}

func transformDinoAliens(action string, args []string) error {
	switch action {
	case "pack":
		return packDinoAliens(args[0], args[1])
	case "unpack":
		return unpackDinoAliens(args[0], args[1])
	default:
		return errors.New("invalid action")
	}
}
