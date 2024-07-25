/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/spf13/cobra"
)

type BuildInfo struct {
	Commit     string
	CommitTime string
	GoVersion  string
	Goos       string
	Goarch     string
	Modified   bool
}

var Version = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version and exit",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Print("command not found\n\n")
			cmd.Help()
		} else {

			// Print Version
			// fmt.Printf("retroblast version: \"%s\", commit: \"%s\", go version: \"%s\", GOOS: \"%s\", GOARCH: \"%s\" \n",
			//Version, Commit, runtime.Version(), runtime.GOOS, runtime.GOARCH)
			buildInfo := GetBuildInfo()

			fmt.Printf("%s\n", buildInfo)
		}
	},
}

func GetBuildInfo() BuildInfo {

	buildInfo := BuildInfo{}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return buildInfo
	}

	buildInfo.GoVersion = info.GoVersion
	buildInfo.Goos = runtime.GOOS
	buildInfo.Goarch = runtime.GOARCH

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			buildInfo.Commit = setting.Value
		case "vcs.time":
			buildInfo.CommitTime = setting.Value
		case "vcs.modified":
			buildInfo.Modified = setting.Value == "true"
		}
	}

	return buildInfo
}

func (i BuildInfo) String() string {

	return fmt.Sprintf("retroblast version: \"%s\", commit: \"%s\", go version: \"%s\", GOOS: \"%s\", GOARCH: \"%s\" \n",
		Version, i.Commit, i.GoVersion, i.Goos, i.Goarch)

}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
