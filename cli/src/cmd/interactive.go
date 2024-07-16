package cmd

import (
	"fmt"
	"os"

	"mhmods/src/util"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	action  string
	game_id string
	mod     string

	actionForm = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose an action to perform").
				Options(
					huh.NewOption("Manage downloaded mods", "manage"),
					huh.NewOption("Download/update mods", "download"),
					huh.NewOption("Open mod folder", "open"),
				).Value(&action),
			huh.NewSelect[string]().
				Title("Choose a game").
				Options(
					huh.NewOption("Moorhuhn Kart Extra", "mhk_1"),
					huh.NewOption("Moorhuhn Kart 2", "mhk_2"),
					huh.NewOption("Moorhuhn Kart 3", "mhk_3"),
					huh.NewOption("Moorhuhn Kart: Thunder", "mhk_4"),
				).Value(&game_id),
		),
	)
)

func manageAllModsGroup(game_id string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a mod to manage").
				Options(
					huh.NewOption("Mod 1", game_id+"_mod1"),
					huh.NewOption("Mod 2", "mod2"),
					huh.NewOption("Mod 3", "mod3"),
				).Value(&mod),
		),
	)
}

func runInteractive(cmd *cobra.Command, args []string) {
	if err := actionForm.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	switch action {
	case "open":
		// TODO: data folder to store mods, etc...
		if err := util.OpenFolder("../mods/" + game_id); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(`
		Opening the mod folder for ` + game_id + `...
		You can use this folder to place or configure your local mods for the game.
		`)
	case "manage":
		if err := manageAllModsGroup(game_id).Run(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "download":
		fmt.Println("Not implemented yet.")
	}
}

func InteractiveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "interactive",
		Short: "Run the tool in interactive mode",
		Long:  `Run the tool in interactive mode to download and manage existing mods (recommended for non-mod developers).`,
		Args:  cobra.NoArgs,
		Run:   runInteractive,
	}
}
