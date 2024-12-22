package cmd

import (
	"log"
	"os"
	"path/filepath"

	"idlemod/src/transformers"
	"idlemod/src/util"

	"github.com/spf13/cobra"
)

var packmodCmd = &cobra.Command{
	Use:   "packmod <game ID> <original data file> <output modded data file> <mod paths>... [-n/--no-merge]",
	Short: "Pack all mod paths into a single data file",
	Long:  `Packs all mod paths into a single data file for a specific game.`,
	Args:  cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := args[0]
		originalDataFile := args[1]
		outputDataFile := args[2]
		modPaths := args[3:]

		// create temp dirs
		tempDirUnpacked, err := os.MkdirTemp("", "idlemod")
		if err != nil {
			log.Fatalf("Fatal error: %s", err)
		}
		defer os.RemoveAll(tempDirUnpacked)
		tempDirPatched, err := os.MkdirTemp("", "idlemod")
		if err != nil {
			log.Fatalf("Fatal error: %s", err)
		}
		defer os.RemoveAll(tempDirPatched)

		// unpack
		if err := transformers.Transform(gameID, "unpack", []string{originalDataFile, tempDirUnpacked}); err != nil {
			log.Fatalf("Fatal error during unpacking: %s", err)
		}

		// copy mod files into unpacked dir, collect existing patch files
		log.Print("Copying mod files...")
		if err := util.CopyModFiles(modPaths, tempDirUnpacked); err != nil {
			log.Fatalf("Fatal error while copying mods: %s", err)
		}

		// patch
		patchFilePaths, err := filepath.Glob(filepath.Join(tempDirUnpacked, "_patches", "*.gopatch"))
		if err != nil {
			log.Fatalf("Fatal error while globbing patch files: %s", err)
		}
		if len(patchFilePaths) > 0 {
			log.Print("Patching mod files...")
			if err := util.PatchModFiles(tempDirUnpacked, tempDirPatched, patchFilePaths); err != nil {
				log.Fatalf("Fatal error while patching mods: %s", err)
			}
		} else {
			log.Print("No patch files found, skipping patching")
			tempDirPatched = tempDirUnpacked
		}

		// repack
		log.Print("Repacking...")
		if err := transformers.Transform(gameID, "pack", []string{outputDataFile, tempDirPatched}); err != nil {
			log.Fatalf("Fatal error during repacking: %s", err)
		}

		log.Printf("Packed modded data file: %s (paths: %v)", outputDataFile, modPaths)
	},
}
