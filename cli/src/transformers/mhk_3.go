// "MHK 3 is almost the same format [as MHK 4], just the index is stored before the data.
// And the file offset is relative to the beginning of the data blob, not the sar file offset."
// - pyramidensurfer
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

// Packs MHK 3 data files.
func packMhk3(dataFileLocation string, inputFolder string) error {
	dataFile, err := os.Create(dataFileLocation)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	// read input folder
	fileEntries, err := walkFiles(inputFolder)
	if err != nil {
		return err
	}

	// write header
	dataFile.Seek(0x0, io.SeekStart)
	dataFile.Write([]byte("SARC2\x00\x00\x00\x00"))
	filesCount := uint32(len(fileEntries))
	binary.Write(dataFile, binary.LittleEndian, filesCount)

	// write file entries
	for _, entry := range fileEntries {
		filename := strings.ReplaceAll(entry.FilePath, "/", "\\")
		filename = strings.Replace(filename, "\\", ":\\", 1)
		filenameLength := uint8(len(filename))

		binary.Write(dataFile, binary.LittleEndian, filenameLength)
		dataFile.Write([]byte(filename))
		dataFile.Write([]byte{0x0})

		binary.Write(dataFile, binary.LittleEndian, uint32(entry.ContentOffset))

		// written twice
		binary.Write(dataFile, binary.LittleEndian, uint32(entry.FileSize))
		binary.Write(dataFile, binary.LittleEndian, uint32(entry.FileSize))

		log.Printf("Wrote file entry `%s`: offset=%d; size=%d", filename, entry.ContentOffset, entry.FileSize)
	}

	// write input file contents
	for i, entry := range fileEntries {
		contentOffset, _ := dataFile.Seek(0, io.SeekCurrent)
		fileContent, _ := os.ReadFile(filepath.Join(inputFolder, entry.FilePath))

		dataFile.Write(fileContent)
		log.Printf("Wrote `%s`", entry.FilePath)

		fileEntries[i].ContentOffset = contentOffset
		fileEntries[i].FileSize = int64(len(fileContent))
	}

	log.Printf("Pack complete: %d files packed.", filesCount)
	return nil
}

// Unpacks MHK 3 data files.
func unpackMhk3(dataFileLocation string, outputDirectory string) error {
	dataFile, err := os.Open(dataFileLocation)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	// header
	dataFileType := make([]byte, 0x8)
	dataFile.Read(dataFileType)
	log.Printf("Data file type: %s;", dataFileType)

	// read file count
	var fileCount uint32
	binary.Read(dataFile, binary.LittleEndian, &fileCount)
	log.Printf("File count: %d", fileCount)

	// read file entries
	unpackedFilesCount := 0
	fileEntries := make([]FileEntry, fileCount)
	log.Printf("Reading file entries...")
	for i := 0; i < int(fileCount); i++ {
		// read filename
		var filenameLength uint8
		binary.Read(dataFile, binary.LittleEndian, &filenameLength)
		filenameBytes := make([]byte, filenameLength)
		dataFile.Read(filenameBytes)

		dataFile.Seek(0x1, io.SeekCurrent) // skip null byte after filename

		// get file offset and length
		var fileOffset, fileLength uint32
		binary.Read(dataFile, binary.LittleEndian, &fileOffset)
		binary.Read(dataFile, binary.LittleEndian, &fileLength)

		// skip 4 bytes, file length is written twice
		dataFile.Seek(0x4, io.SeekCurrent)

		fileEntries[i] = FileEntry{
			FilePath:     string(filenameBytes),
			ContentOffset: int64(fileOffset),
			FileSize:     int64(fileLength),
		}
		log.Printf("Filename: %s; Offset: %d; Size: %d", fileEntries[i].FilePath, fileEntries[i].ContentOffset, fileEntries[i].FileSize)
	}

	dataStart, _ := dataFile.Seek(0, io.SeekCurrent)
	for _, entry := range fileEntries {
		// make parent dir for file
		filename := strings.ReplaceAll(entry.FilePath, ":", "")
		filename = strings.ReplaceAll(filename, "\\", "/")
		outputDirectory := filepath.Join(outputDirectory, filepath.Dir(filename))
		if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
			return err
		}
		outputPath := filepath.Join(outputDirectory, filepath.Base(filename))

		// offset of the next file entry
		nextOffset, _ := dataFile.Seek(0x4, io.SeekCurrent)

		// write file
		dataFile.Seek(dataStart + entry.ContentOffset, io.SeekStart)
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		if _, err := io.CopyN(outputFile, dataFile, entry.FileSize); err != nil {
			return err
		}
		outputFile.Close()
		unpackedFilesCount++

		// go to next file
		dataFile.Seek(nextOffset, io.SeekStart)
		log.Printf("Filename: %s; Offset: %d; Size: %d", filename, entry.ContentOffset, entry.FileSize)
	}

	log.Printf("Unpack complete: %d/%d.", unpackedFilesCount, fileCount)
	return nil
}

func transformMhk3(action string, dataFileLocation string, rootFolder string) error {
	switch action {
	case "pack":
		return packMhk3(dataFileLocation, rootFolder)
	case "unpack":
		return unpackMhk3(dataFileLocation, rootFolder)
	default:
		return errors.New("invalid action")
	}
}
