package cmd_test

import (
	"regexp"
	"testing"

	"github.com/retroblast-engine/retroblast/cmd"
)

func TestVersionCommand(t *testing.T) {

	buildInfo := cmd.CreateBuildInfo()

	// Regex Patterns to check syntax
	tests := []struct {
		name    string
		value   string
		pattern string
	}{
		{"Version", buildInfo.Version, "dev"},
		{"Commit", buildInfo.Commit, "[a-zA-Z0-9]*"},
		{"GoVersion", buildInfo.GoVersion, "go[0-9]+.[0-9]*.[0-9]*"},
		{"Goos", buildInfo.Goos, "[a-zA-Z0-9]+"},
		{"Goarch", buildInfo.Goarch, "[a-zA-Z0-9]+"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected, err := regexp.Compile(tt.pattern)
			if err != nil {
				t.Fatalf("Error compiling regex for %s: %v", tt.name, err)
			}

			if !expected.MatchString(tt.value) {
				t.Errorf("\nExpected %s to match: %q\n\nActual: %q", tt.name, expected, tt.value)
			} else {
				t.Logf("%s OK", tt.name)
			}

		})
	}
}
