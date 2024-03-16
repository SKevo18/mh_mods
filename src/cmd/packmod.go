package cmd

import (
	"log"
	"os"

	"mhmods/src/transformers"
	"mhmods/src/util"

	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

func PackmodCmd() *cobra.Command {
	var noMerge bool

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
			tempDirUnpacked, err := os.MkdirTemp("", "mhmods_temp_unpacked")
			if err != nil {
				log.Fatalf("Fatal error: %s", err)
			}
			defer os.RemoveAll(tempDirUnpacked)
			tempDirMerged, err := os.MkdirTemp("", "mhmods_temp_merged")
			if err != nil {
				log.Fatalf("Fatal error: %s", err)
			}
			defer os.RemoveAll(tempDirMerged)

			// unpack
			if err := transformers.Transform("unpack", gameID, originalDataFile, tempDirUnpacked); err != nil {
				log.Fatalf("Fatal error during unpacking: %s", err)
			}

			// merge
			if noMerge {
				log.Print("Copying mod files without merging...")

				if err := cp.Copy(tempDirUnpacked, tempDirMerged); err != nil {
					log.Fatalf("Fatal error while copying original files: %s", err)
				}

				for _, modPath := range modPaths {
					if err := cp.Copy(modPath, tempDirMerged); err != nil {
						log.Fatalf("Fatal error while copying mods: %s", err)
					}
				}
			} else {
				log.Print("Merging mod files...")
				if err := util.MergeModFilesRecursively(tempDirUnpacked, modPaths, tempDirMerged); err != nil {
					log.Fatalf("Fatal error during merging: %s", err)
				}
			}

			// repack
			if err := transformers.Transform("pack", gameID, outputDataFile, tempDirMerged); err != nil {
				log.Fatalf("Fatal error during repacking: %s", err)
			}

			log.Printf("Packed modded data file: %s (paths: %v)", outputDataFile, modPaths)
		},
	}

	packmodCmd.Flags().BoolVarP(&noMerge, "no-merge", "n", false, "Do not merge mod files by comparing unique lines, just overwrite original files in the data file.")

	return packmodCmd
}
