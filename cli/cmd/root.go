/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), asciiArt+"\n")
			cmd.Help()
		},
	}

	cmd.AddCommand(VersionCommand())

	return cmd
}

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
