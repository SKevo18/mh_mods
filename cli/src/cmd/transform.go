package cmd

import (
	"github.com/spf13/cobra"
)

var unpackCmd = &cobra.Command{
	Use:   "unpack <game ID> <data file location> <output folder path>",
	Short: "Unpack all files into output folder",
	Long:  `Unpacks all files of a specific game into the output folder.`,
	Args:  cobra.MinimumNArgs(3),
	Run:   genericTransform("unpack"),
}
