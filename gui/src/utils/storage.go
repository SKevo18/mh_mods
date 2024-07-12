package utils

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"mhmods_gui/src"
)

// The path to the data directory.
var DataDir = GetDataDir()

// Returns the path to the data directory.
func GetDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/mhmods"
	}

	return filepath.Join(homeDir, src.AppID)
}

// Ensures that the data directories exist, prior to app launch.
func EnsureDataDirs() error {
	for _, gameId := range src.SupportedGames {
		// create "game dir / mods dir"
		if err := os.MkdirAll(filepath.Join(DataDir, gameId, "mods"), 0o755); err != nil {
			return err
		}

		// write default game config, if it doesn't exist
		configPath := filepath.Join(DataDir, gameId, "config.json")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			data := []byte(`{"id":"` + gameId + `","name":"` + gameId + `","executable":"","enabledMods":[]}`)
			if err := os.WriteFile(configPath, data, 0o644); err != nil {
				return err
			}
		}
	}

	return nil
}

// Returns the names of all folders in the given parent folder.
func GetFolderNames(parentFolder string) ([]string, error) {
	paths, err := os.ReadDir(parentFolder)
	if err != nil {
		return nil, err
	}

	folders := make([]string, 0, len(paths))
	for _, path := range paths {
		if path.IsDir() {
			folders = append(folders, path.Name())
		}
	}

	return folders, nil
}

// A generic function to load a JSON file.
func LoadJsonFile(filePath string, target any) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(fileData, target)
}

// A generic function to write a JSON file.
func WriteJsonFile(filePath string, data any) error {
	fileData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, fileData, 0o644)
}

// Unzips the given mod ZIP at its location
func UnzipMod(modZipPath string) error {
	// reader
	zipReader, err := zip.OpenReader(modZipPath)
	if err != nil {
		return fmt.Errorf("error opening ZIP file: %s", err)
	}
	defer zipReader.Close()

	// extract
	for _, file := range zipReader.File {
		// open zipped file
		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("error opening file: %s", err)
		}
		defer fileReader.Close()

		// create file
		filePath := filepath.Join(filepath.Dir(modZipPath), file.Name)
		if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
			return fmt.Errorf("error creating file dirs: %s", err)
		}
		newFile, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("error creating file: %s", err)
		}
		defer newFile.Close()

		// copy
		if _, err := io.Copy(newFile, fileReader); err != nil {
			return fmt.Errorf("error copying file: %s", err)
		}
	}

	return nil
}
