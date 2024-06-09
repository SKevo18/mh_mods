package storage

import (
	"os"
	"path/filepath"

	"mhmods_gui/src"
)

func EnsureAppDir() (string, error) {
	appDirRoot, err := GetAppDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(appDirRoot, os.ModePerm); err != nil {
		return "", err
	}
	return appDirRoot, nil
}

func GetAppDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, src.AppID), nil
}
