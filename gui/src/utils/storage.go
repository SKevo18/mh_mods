package utils

import (
	"os"
	"path/filepath"

	"mhmods_gui/src"
)

var DataDir string = GetDataDir()

func GetDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/mhmods"
	}

	return filepath.Join(homeDir, src.AppID)
}

func EnsureDataDirs() error {
	for _, gameId := range src.SupportedGames {
		if err := os.MkdirAll(filepath.Join(DataDir, "mods", gameId), 0o755); err != nil {
			return err
		}
	}

	return nil
}

func GetModFolders(parentFolder string) ([]string, error) {
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
