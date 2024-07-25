/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "retroblast",
	Short: "A brief description of your application",
	Long:  asciiArt + "\n" + `A retro 2D game engine using Go`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.retroblast.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
