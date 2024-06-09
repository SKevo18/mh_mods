package api_test

import (
	"fmt"
	"testing"

	"mhmods_gui/src/api"
)

func TestGetMods(t *testing.T) {
	allMods, err := api.GetAllMods()
	if err != nil {
		t.Fatalf("GetAllMods() returned error: %v", err)
	}

	fmt.Println("All mods:")
	for id, mods := range allMods {
		fmt.Printf("\t%s:\n", id)
		for modId, name := range mods {
			fmt.Printf("\t\t%s: %s\n", modId, name)
		}
	}
}

func TestGetSpecificMods(t *testing.T) {
	game := "mhk_3"
	mods, err := api.GetMods(game)
	if err != nil {
		t.Fatalf("GetMods(%q) returned error: %v", game, err)
	}

	fmt.Printf("Mods for %q:\n", game)
	for id, name := range mods {
		fmt.Printf("\t%s: %s\n", id, name)
	}
}
