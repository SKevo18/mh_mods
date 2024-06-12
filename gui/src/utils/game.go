package utils

import (
	"os/exec"
	"path/filepath"
)

type Game struct {
	Name string
	Id string
	ExecutablePath string
	EnabledMods []string
}

func (g *Game) FromConfig(gameId string) {
	g.Name = "Moorhuhn Kart 1"
	g.Id = gameId
	g.ExecutablePath = "/usr/games/mhk1"
	g.EnabledMods = []string{}
}

func (g *Game) ModFolder() string {
	return filepath.Join(DataDir, "mods", g.Id)
}

func (g *Game) Launch() error {
	cmd := exec.Command(g.ExecutablePath)
	return cmd.Start()
}

