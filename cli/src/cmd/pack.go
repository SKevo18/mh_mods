package cmd

import (
	"github.com/spf13/cobra"
)

var packCmd = &cobra.Command{
	Use:   "pack <game ID> <data file location> <input folder path>",
	Short: "Pack all files into a data file",
	Long: `Packs all files into a data file for a specific game.
	The structure of the input folder must be the same as the original data file.
	All files that should be packed must be present in the input folder.`,
	Args: cobra.MinimumNArgs(3),
	Run:  genericTransform("pack"),
}
