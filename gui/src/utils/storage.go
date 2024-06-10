package utils

import (
	"os"
	"path/filepath"

	"mhmods_gui/src"
)

func EnsureDataDirs() (string, error) {
	dataDirRoot, err := GetDataDir()
	if err != nil {
		return "", err
	}

	for gameId, _ := range src.SupportedGames {
		if err := os.MkdirAll(filepath.Join(dataDirRoot, "mods", gameId), 0755); err != nil {
			return "", err
		}
	}

	return dataDirRoot, nil
}

func GetDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, src.AppID), nil
}

func GetInstalledMods(parentFolder string) ([]string, error) {
	files, err := os.ReadDir(parentFolder)
	if err != nil {
		return nil, err
	}

	mods := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			mods = append(mods, file.Name())
		}
	}

	return mods, nil
}
