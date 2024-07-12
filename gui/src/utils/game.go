package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"mhmods_gui/src/api"
)

type Game struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	ExecutablePath string   `json:"executable"`
	EnabledMods    []string `json:"enabledMods"`
}

// Returns the path to the game config file
func (g *Game) ConfigPath() string {
	return filepath.Join(DataDir, g.Id, "config.json")
}

// Writes the game attributes to the game config file
func (g *Game) WriteConfig() error {
	return WriteJsonFile(g.ConfigPath(), g)
}

// Returns the path to the root mods folder for the game
func (g *Game) ModsFolder() string {
	return filepath.Join(DataDir, g.Id, "mods")
}

// Launches the game using the configured executable path
func (g *Game) Launch() error {
	return exec.Command(g.ExecutablePath).Start()
}

// Downloads a mod with a specific ID from the API and unzips it to the game's mods folder
// The downloaded mod folder already contains its current commit hash
// If the mod is already downloaded, it will update (overwrite) it
func (g *Game) DownloadMod(modId string) error {
	modPath := filepath.Join(g.ModsFolder(), modId+".zip")

	if err := api.DownloadMod(g.Id, modId, modPath); err != nil {
		return err
	}
	if err := UnzipMod(modPath); err != nil {
		return err
	}
	defer os.Remove(modPath)

	return nil
}

// Loads a game from the game config file and returns its struct
func LoadGame(gameId string) (*Game, error) {
	g := &Game{Id: gameId}
	configPath := g.ConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("game config file not found: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, g)
	return g, err
}
