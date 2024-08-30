// Package cmd provides the command-line interface for the application,
// including the root command and its associated subcommands.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const asciiArt string = `		
  _____      _             ____  _           _   
 |  __ \    | |           |  _ \| |         | |  
 | |__) |___| |_ _ __ ___ | |_) | | __ _ ___| |_ 
 |  _  // _ \ __| '__/ _ \|  _ <| |/ _` + "`" + ` / __| __|
 | | \ \  __/ |_| | | (_) | |_) | | (_| \__ \ |_ 
 |_|  \_\___|\__|_|  \___/|____/|_|\__,_|___/\__|
 `

// RootCommand represents the base command when called without any subcommands
func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retroblast",
		Short: "A retro 2D game engine using Go",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Println(asciiArt)

			cmd.Help()
		},
	}

	cmd.AddCommand(VersionCommand())

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := RootCommand().Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {

	// Flags and Configuration Settings should be defined here

	RootCommand().Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
