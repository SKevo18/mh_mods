package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var ApiServer string = "https://mh-mods.svit.ac"

// Fetches a JSON object from the API server.
func fetchJson(path string, target interface{}) (error) {
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

// GetAllMods returns a list of all mods.
func GetAllMods() (map[string]map[string]string, error) {
	var mods map[string]map[string]string
	if err := fetchJson("/mods", &mods); err != nil {
		return nil, err
	}
	return mods, nil
}

// GetMods returns a list of mods for the given game.
func GetMods(gameId string) (map[string]string, error) {
	var mods map[string]string
	if err := fetchJson(fmt.Sprintf("/mods/%s", gameId), &mods); err != nil {
		return nil, err
	}
	return mods, nil
}
