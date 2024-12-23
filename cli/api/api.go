// WIP
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// The API server to use
var ApiServer string = "https://idlemod.svit.ac"

type (
	Mods  map[string]string
	Games map[string]Mods
)

// Fetches a JSON object from the API server
func fetchJson(path string, target interface{}) error {
	resp, err := http.Get(fmt.Sprintf("%s%s", ApiServer, path))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		return err
	}

	return nil
}

// Returns a list of all games and their mods
func GetAllMods() (Games, error) {
	var mods Games
	if err := fetchJson("/mods", &mods); err != nil {
		return nil, err
	}
	return mods, nil
}

// Returns a list of mods for the given game
func GetMods(gameId string) (Mods, error) {
	var mods Mods
	if err := fetchJson(fmt.Sprintf("/mods/%s", gameId), &mods); err != nil {
		return nil, err
	}
	return mods, nil
}

// Returns the current (up-to-date) commit hash of the given mod from the API
func GetCurrentModHash(gameId string, modId string) (string, error) {
	allMods, err := GetMods(gameId)
	if err != nil {
		return "", err
	}
	for id, hash := range allMods {
		if id == modId {
			return hash, nil
		}
	}

	return "", fmt.Errorf("mod not found: %s", modId)
}

// Downloads a given mod of a game as a ZIP file into destination folder
// Unzipping the file is handled in the `utils/storage` package
func DownloadMod(gameId string, modId string, destFile string) error {
	// prepare file
	file, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// GET
	resp, err := http.Get(fmt.Sprintf("%s/mods/%s/%s", ApiServer, gameId, modId))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// write
	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}
