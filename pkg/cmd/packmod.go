package cmd

import (
	"log"
	"os"

	"mhmods/pkg/transformers"
    "mhmods/pkg/util"

	"github.com/spf13/cobra"
)

func PackmodCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "packmod <game ID> <original data file> <output modded data file> <mod paths>...",
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

			// merge and repack
			modPaths = append(modPaths, tempDirUnpacked)
			if err := util.MergeRecursively(tempDirMerged, modPaths); err != nil {
				log.Fatalf("Fatal error during merging: %s", err)
			}
			if err := transformers.Transform("pack", gameID, outputDataFile, tempDirMerged); err != nil {
				log.Fatalf("Fatal error during repacking: %s", err)
			}

			log.Printf("Packed modded data file: %s (paths: %v)", outputDataFile, modPaths)
		},
	}
}
