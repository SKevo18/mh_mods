package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"mhmods/pkg/transformers"
)

// Returns a function that can be used to pack or unpack files.
func genericTransform(action string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		gameID := args[0]
		dataFileLocation := args[1]
		folderPath := args[2]

		log.Printf("Beginning the %sing of files from data file `%s` into output folder `%s` for game `%s`...\n", action, dataFileLocation, folderPath, gameID)
		err := transformers.Transform(action, gameID, dataFileLocation, folderPath)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("Successfully %sed files for game `%s`!\n", action, gameID)
	}
}

// Returns the root command to pack files.
func PackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pack <game ID> <data file location> <input folder path>",
		Short: "Pack all files into a data file",
		Long: `Packs all files into a data file for a specific game.
		The structure of the input folder must be the same as the original data file.
		All files that should be packed must be present in the input folder.`,
		Args: cobra.MinimumNArgs(3),
		Run:  genericTransform("pack"),
	}
}

// Returns the command used to unpack files.
func UnpackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unpack <game ID> <data file location> <output folder path>",
		Short: "Unpack all files into output folder",
		Long:  `Unpacks all files of a specific game into the output folder.`,
		Args:  cobra.MinimumNArgs(3),
		Run:   genericTransform("unpack"),
	}
}
