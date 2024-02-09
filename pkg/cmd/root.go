package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mhmods <pack/unpack/packmod> <game ID> <data file location> <input/output folder path(s)>",
	Short: "Moorhuhn modding tool for packing and unpacking game data files",
	Long: `A Fast and Flexible Moorhuhn modding tool built with love in Go.
	Complete documentation (including game IDs) is available at http://github.com/SKevo18/mhmods`,
	Args: cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Use)
	},
}

func EntryPoint() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	rootCmd.AddCommand(PackCmd())
	rootCmd.AddCommand(UnpackCmd())
}
