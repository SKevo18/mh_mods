package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Game struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	ExecutablePath string   `json:"executable"`
	EnabledMods    []string `json:"enabledMods"`
}

func (g *Game) ConfigPath() string {
	return filepath.Join(DataDir, g.Id, "config.json")
}

func (g *Game) WriteConfig() error {
	return WriteJsonFile(g.ConfigPath(), g)
}

func (g *Game) ModFolder() string {
	return filepath.Join(DataDir, g.Id, "mods")
}

func (g *Game) Launch() error {
	return exec.Command(g.ExecutablePath).Start()
}

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
