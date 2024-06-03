package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mhmods <pack/unpack/packmod> ...",
	Short: "Moorhuhn modding tool for packing and unpacking game data files",
	Long: `
A fast and flexible Moorhuhn modding tool built with the Go programming language.
More information is available at http://github.com/SKevo18/mh_mods`,
	Args: cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Use)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of mhmods",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mhmods v1.0.0 (2024-05-10)")
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

	rootCmd.AddCommand(versionCmd)
}
