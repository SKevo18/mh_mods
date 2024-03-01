package transformers

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Unzip a file into a folder.
func unzipFile(zipFile string, outputFolder string) error {
	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		return fmt.Errorf("error creating output folder: %s", err)
	}

	// unzip
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("error opening ZIP file: %s", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		filePath := filepath.Join(outputFolder, file.Name)

		// read
		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("error opening file: %s", err)
		}
		defer fileReader.Close()

		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory: %s", err)
		}

		// write
		outFile, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("error creating file: %s", err)
		}
		defer outFile.Close()

		if _, err = io.Copy(outFile, fileReader); err != nil {
			return fmt.Errorf("error copying file: %s", err)
		}
	}

	return nil
}

// Zips a folder into a file.
func zipFolder(folder string, zipFile string) error {
	// open
	outFile, err := os.Create(zipFile)
	if err != nil {
		return fmt.Errorf("error creating ZIP file: %s", err)
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
		return fmt.Errorf("error while walking through files: %s", err)
	}

	return nil
}
