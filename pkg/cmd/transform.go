package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"mhmods/pkg/transformers"
)

func genericTransform(action string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		gameID := args[0]
		dataFileLocation := args[1]
		folderPath := args[2]

		err := transformers.Transform(action, gameID, dataFileLocation, folderPath)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

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

func UnpackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unpack <game ID> <data file location> <output folder path>",
		Short: "Unpack all files into output folder",
		Long:  `Unpacks all files of a specific game into the output folder.`,
		Args:  cobra.MinimumNArgs(3),
		Run:   genericTransform("unpack"),
	}
}
