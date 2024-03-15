package transformers

import (
	"archive/zip"
	"compress/flate"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Unzip a file into a folder.
func unzipFile(zipFilePath string, outputFolder string) error {
	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		return fmt.Errorf("error creating output folder: %s", err)
	}

	// unzip
	zipReader, err := zip.OpenReader(zipFilePath)
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

// Zips a folder into a zip file.
func zipFolder(folderPath string, zipFilePath string) error {
	// create output
	outFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// create zip writer
	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// register compressor
	zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})

	err = filepath.Walk(folderPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// open file to zip
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// use base path, to keep files at root
		zipPath := filepath.Base(filePath)
		zipFile, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}

		// write file to zip
		_, err = io.Copy(zipFile, file)
		return err
	})

	return err
}
