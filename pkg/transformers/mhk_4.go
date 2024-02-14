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

// Packs MHK 4 (Thunder) data files.
// Unpacks MHK 4 (Thunder) data files.
// MHK 4 has 2 data files in the installation directory: `data.sar` (main one),
// and `data.s01` (whose purpose is unknown, but it might be some demo data?).
func packMhk4(dataFileLocation string, inputFolder string) error {
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

	// reserve space for header + file entries offset
	dataFile.Seek(0x10, io.SeekStart)

	// write input file contents
	for i, entry := range fileEntries {
		contentOffset, _ := dataFile.Seek(0, io.SeekCurrent)
		fileContent, _ := os.ReadFile(filepath.Join(inputFolder, entry.FilePath))

		dataFile.Write(fileContent)
		log.Printf("Wrote `%s`", entry.FilePath)

		fileEntries[i].ContentOffset = contentOffset
		fileEntries[i].FileSize = int64(len(fileContent))
	}

	// write file entries
	fileEntriesBegin, _ := dataFile.Seek(0, io.SeekCurrent) // original data file: 874881073
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

	// write header
	dataFile.Seek(0x0, io.SeekStart)
	dataFile.Write([]byte("SARCFV"))
	dataFile.Write([]byte{0x1, 0x1})

	// file count
	filesCount := uint32(len(fileEntries) + 5) // duplicates
	binary.Write(dataFile, binary.LittleEndian, filesCount)
	// write file entries count
	binary.Write(dataFile, binary.LittleEndian, uint32(fileEntriesBegin))

	log.Printf("Pack complete: %d files packed. File entries begin at `%d`.", filesCount, fileEntriesBegin)
	return nil
}

// Unpacks MHK 4 (Thunder) data files.
// MHK 4 has 2 data files in the installation directory: `data.sar` (main one),
// and `data.s01` (whose purpose is unknown, but it might be some demo data?).
/*
Note: For some reason, the following files are present twice in the original datafile index:
	"unpacked/mhk_4/D/Mhkart4/3dobjects_extras/hammer_addon.lwo"
	"unpacked/mhk_4/D/Mhkart4/3dobjects_extras/health_addon.lwo"
	"unpacked/mhk_4/D/Mhkart4/3dobjects_extras/quaddamage_addon.lwo"
	"unpacked/mhk_4/D/Mhkart4/textures_cars/kart_blau_0001.bmp.dds"
	"unpacked/mhk_4/D/Mhkart4/textures_cars/kartlogo_01.png.dds"
*/
func unpackMhk4(dataFileLocation string, outputDirectory string) error {
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

	// read file entries offset
	var fileEntriesBegin uint32
	dataFile.Seek(0xC, io.SeekStart)
	binary.Read(dataFile, binary.LittleEndian, &fileEntriesBegin)

	// read file entries
	// note: these are listed at end of the data file
	unpackedFilesCount := 0
	log.Printf("Reading file entries from offset %d...", fileEntriesBegin)
	dataFile.Seek(int64(fileEntriesBegin), io.SeekStart)
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
