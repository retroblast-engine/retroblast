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
	Version    string
	Commit     string
	CommitTime string
	GoVersion  string
	Goos       string
	Goarch     string
	Modified   bool
}

var buildInfo BuildInfo

// VersionCommand represents the version command
func VersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version and build information",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				fmt.Print("command not found\n\n")
				cmd.Help()
			} else {

				buildInfo := CreateBuildInfo()

				fmt.Fprintf(cmd.OutOrStdout(), "%s\n", buildInfo)

			}
		},
	}

	return cmd
}

// Retrieve build information, such as Version, CommitHash, goVersion, Goos, Goarch
func CreateBuildInfo() BuildInfo {

	buildInfo = BuildInfo{}

	// The version is currently labeled "dev" because the engine is still under development
	buildInfo.Version = "dev"

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

func GetBuildInfo() BuildInfo {
	return buildInfo
}

func (i BuildInfo) String() string {

	return fmt.Sprintf("retroblast version: \"%s\", commit: \"%s\", go version: \"%s\", GOOS: \"%s\", GOARCH: \"%s\" \n",
		i.Version, i.Commit, i.GoVersion, i.Goos, i.Goarch)

}

func init() {

	// Flags and Configuration Settings should be defined here

}
