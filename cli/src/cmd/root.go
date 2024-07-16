package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "Moorhuhn modding tool for packing and unpacking game data files",
	Long: `
A fast and flexible Moorhuhn modding tool built with the Go programming language.
More information is available at http://github.com/SKevo18/mh_mods`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: interactive mode (https://github.com/charmbracelet/huh)
		cmd.Help()
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
	rootCmd.AddCommand(PackmodCmd())
}
