package cmd

import (
	"fmt"
	"log"

	"idlemod/src/transformers"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "Idlemod tool for packing and unpacking game data files",
	Long: `
A fast and flexible old game modding tool built with the Go programming language.
More information is available at http://github.com/SKevo18/idlemod`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
		Idlemod Tool (https://github.com/SKevo18/idlemod)

		Welcome to the Idlemod game modding tool!
		Use the 'interactive' subcommand to run the tool in interactive mode (to download and manage existing mods, recommended for normal players).
		Otherwise, please use the '--help' flag for more help about creating mods.
		`)
	},
}

func EntryPoint() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func genericTransform(action string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		log.Printf("Beginning the %sing of files for game `%s`... (args: %v)\n", action, args[0], args[1:])

		if err := transformers.Transform(args[0], action, args[1:]); err != nil {
			log.Fatal("Fatal error: ", err.Error())
		}

		log.Printf("Successfully(?) performed action `%s` for game `%s`!\n", action, args[0])
	}
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
	rootCmd.AddCommand(packCmd)
	rootCmd.AddCommand(unpackCmd)
	rootCmd.AddCommand(packmodCmd)
}
