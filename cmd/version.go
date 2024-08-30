// Package cmd provides the command-line interface for the application,
// including the root command and its associated subcommands.
package cmd

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// BuildInfo holds information about the build, including the version,
// commit hash, and the time when the build was created.
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
		Long:  "Displays detailed information about the current build of the application.",
		Args:  cobra.NoArgs, // No arguments are allowed
		Run: func(cmd *cobra.Command, args []string) {
			// Print build information
			cmd.Println(CreateBuildInfo().String())
		},
	}

	return cmd
}

// CreateBuildInfo retrieves build information, such as Version, CommitHash, goVersion, Goos, Goarch
func CreateBuildInfo() BuildInfo {

	buildInfo = BuildInfo{}

	// The version is currently labeled "dev" because the engine is still under development
	buildInfo.Version = "dev"

	info, ok := debug.ReadBuildInfo()
	if !ok {
		log.Printf("Failed to read build information")
	}

	buildInfo.GoVersion = info.GoVersion
	buildInfo.Goos = runtime.GOOS
	buildInfo.Goarch = runtime.GOARCH

	// Iterate over the settings in the info object and update the buildInfo
	// struct based on specific keys.
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
		i.Version, i.Commit, i.GoVersion, i.Goos, i.Goarch)

}
