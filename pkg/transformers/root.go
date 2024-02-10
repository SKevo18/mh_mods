package transformers

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// A simple struct to represent a file entry.
type FileEntry struct {
	// The path of the file relative to the root folder.
	Filename string
	// The size of the file in bytes.
	Filesize int64
}

// Walks through files in `rootFolder` and returns array of file entries.
func walkFiles(rootFolder string) ([]FileEntry, error) {
	var files []FileEntry

	err := filepath.Walk(rootFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(rootFolder, path)
			if err != nil {
				return err
			}

			files = append(files, FileEntry{
				Filename: relPath,
				Filesize: info.Size(),
			})
		}

		return nil
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error walking through files: %s", err))
	}

	return files, nil
}

// A simple XOR function that applies a key to the data (used in MHK Extra).
func xorData(data []byte, key []byte) []byte {
	keyLength := len(key)
	for i := range data {
		data[i] ^= key[i%keyLength]
	}
	return data
}

// Unzip a file into a folder.
func unzipFile(zipFile string, outputFolder string) error {
	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		return errors.New(fmt.Sprintf("error creating output folder: %s", err))
	}

	// Unzip the file
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return errors.New(fmt.Sprintf("error opening ZIP file: %s", err))
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		filePath := filepath.Join(outputFolder, file.Name)

		// read
		fileReader, err := file.Open()
		if err != nil {
			return errors.New(fmt.Sprintf("error opening file: %s", err))
		}
		defer fileReader.Close()

		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return errors.New(fmt.Sprintf("error creating directory: %s", err))
		}

		// write
		outFile, err := os.Create(filePath)
		if err != nil {
			return errors.New(fmt.Sprintf("error creating file: %s", err))
		}
		defer outFile.Close()

		if _, err = io.Copy(outFile, fileReader); err != nil {
			return errors.New(fmt.Sprintf("error copying file: %s", err))
		}
	}

	return nil
}

// Zips a folder into a file.
func zipFolder(folder string, zipFile string) error {
	outFile, err := os.Create(zipFile)
	if err != nil {
		return errors.New(fmt.Sprintf("error creating ZIP file: %s", err))
	}
	defer outFile.Close()
	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// only files can be zipped
		if info.IsDir() {
			return nil
		}

		// header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// read path is absolute, convert to relative
		header.Name = filepath.ToSlash(filepath.Base(path))
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errors.New(fmt.Sprintf("error while walking through files: %s", err))
	}

	return nil
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
