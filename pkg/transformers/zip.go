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

// Zips a folder into a zip file.
func zipFolder(folder string, zipFile string) error {
	// create
	outFile, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// zip writer
	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// walk
	filepath.Walk(folder, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// create zip file
		zipFile, err := zipWriter.Create(filePath)
		if err != nil {
			return err
		}

		// read
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// write
		buf := make([]byte, 1024 * 1024) // 1 MB
		for {
			n, err := file.Read(buf)
			if err != nil && err != io.EOF {
				return err
			}
			if n == 0 {
				break
			}

			_, err = zipFile.Write(buf[:n])
			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}
