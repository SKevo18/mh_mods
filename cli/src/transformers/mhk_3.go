// "MHK 3 is almost the same format [as MHK 4], just the index is stored before the data.
// And the file offset is relative to the beginning of the data blob, not the sar file offset."
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
	dataFile.Write([]byte("SARC"))
	dataFile.Write([]byte{0x2, 0x0, 0x0, 0x0})
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
	dataFileType := make([]byte, 0x6)
	dataFile.Read(dataFileType)
	fileInfo, err := dataFile.Stat()
	if err != nil {
		return err
	}
	dataSize := fileInfo.Size()
	log.Printf("Data file type: %s; Data size: %d", dataFileType, dataSize)

	// read file count
	var fileCount uint32
	dataFile.Seek(0x8, io.SeekStart)
	binary.Read(dataFile, binary.LittleEndian, &fileCount)
	log.Printf("File count: %d", fileCount)

	// read file entries
	unpackedFilesCount := 0
	log.Printf("Reading file entries...")
	for i := 0; i < int(fileCount); i++ {
		// read filename
		var filenameLength uint8
		binary.Read(dataFile, binary.LittleEndian, &filenameLength)
		filenameBytes := make([]byte, filenameLength)
		dataFile.Read(filenameBytes)

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

		// offset of the next file entry
		nextOffset, _ := dataFile.Seek(0x4, io.SeekCurrent)

		// write file
		dataFile.Seek(int64(fileOffset), io.SeekStart)
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		io.CopyN(outputFile, dataFile, int64(fileLength))
		outputFile.Close()
		unpackedFilesCount++

		// go to next file
		dataFile.Seek(nextOffset, io.SeekStart)
		log.Printf("Filename: %s; Offset: %d; Size: %d", filename, fileOffset, fileLength)
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
