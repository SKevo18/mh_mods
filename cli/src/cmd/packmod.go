package cmd

import (
	"log"
	"os"

	"mhmods/src/transformers"
	"mhmods/src/util"

	"github.com/spf13/cobra"
)

func PackmodCmd() *cobra.Command {
	packmodCmd := &cobra.Command{
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
			tempDirUnpacked, err := os.MkdirTemp("", "mhmods_temp_patched")
			if err != nil {
				log.Fatalf("Fatal error: %s", err)
			}
			defer os.RemoveAll(tempDirUnpacked)
			tempDirPatched, err := os.MkdirTemp("", "mhmods_temp_patched")
			if err != nil {
				log.Fatalf("Fatal error: %s", err)
			}
			defer os.RemoveAll(tempDirPatched)

			// unpack
			if err := transformers.Transform("unpack", gameID, originalDataFile, tempDirUnpacked); err != nil {
				log.Fatalf("Fatal error during unpacking: %s", err)
			}

			// copy mod files into unpacked dir, collect existing patch files
			log.Print("Copying mod files...")
			patchFilePaths, err := util.CopyModFiles(modPaths, tempDirUnpacked)
			if err != nil {
				log.Fatalf("Fatal error while copying mods: %s", err)
			}

			// patch
			log.Print("Patching mod files...")
			if err := util.PatchModFiles(tempDirUnpacked, tempDirPatched, patchFilePaths); err != nil {
				log.Fatalf("Fatal error while patching mods: %s", err)
			}

			// repack
			log.Print("Repacking...")
			if err := transformers.Transform("pack", gameID, outputDataFile, tempDirPatched); err != nil {
				log.Fatalf("Fatal error during repacking: %s", err)
			}

			log.Printf("Packed modded data file: %s (paths: %v)", outputDataFile, modPaths)
		},
	}

	return packmodCmd
}
