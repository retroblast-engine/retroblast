package cmd_test

import (
	"regexp"
	"testing"

	"github.com/retroblast-engine/retroblast/cli/cmd"
)

func TestVersionCommand(t *testing.T) {

	buildInfo := cmd.CreateBuildInfo()

	// Regex Patterns to check syntax
	testVersion := "dev"
	testCommit := "[a-zA-Z0-9]*"
	testGoVersion := "go[0-9]+.[0-9]*.[0-9]*"
	testGoos := "[a-zA-Z0-9]+"
	testGoarch := "[a-zA-Z0-9]+"

	// Test Version
	expectedVersion, err := regexp.Compile(testVersion)
	if err != nil {
		t.Fatalf("Error compiling regex: %v", err)
	}

	if !expectedVersion.MatchString(buildInfo.Version) {
		t.Errorf("\nExpected output: %q\n\nActual output: %q", expectedVersion, buildInfo.Version)
	} else {
		t.Log("Version OK")
	}

	// Test Commit
	expectedCommit, err := regexp.Compile(testCommit)
	if err != nil {
		t.Fatalf("Error compiling regex: %v", err)
	}

	if !expectedCommit.MatchString(buildInfo.Commit) {
		t.Errorf("\nExpected output: %q\n\nActual output: %q", expectedCommit, buildInfo.Commit)
	} else {
		t.Log("Commit OK")
	}

	// Test Go Version
	expectedGoVersion, err := regexp.Compile(testGoVersion)
	if err != nil {
		t.Fatalf("Error compiling regex: %v", err)
	}

	if !expectedGoVersion.MatchString(buildInfo.GoVersion) {
		t.Errorf("\nExpected output: %q\n\nActual output: %q", expectedGoVersion, buildInfo.GoVersion)
	} else {
		t.Log("GoVersion OK")
	}

	// Test Goos
	expectedGoos, err := regexp.Compile(testGoos)
	if err != nil {
		t.Fatalf("Error compiling regex: %v", err)
	}

	if !expectedGoos.MatchString(buildInfo.Goos) {
		t.Errorf("\nExpected output: %q\n\nActual output: %q", expectedGoos, buildInfo.Goos)
	} else {
		t.Log("Goos OK")
	}

	// Test Goarch
	expectedGoarch, err := regexp.Compile(testGoarch)
	if err != nil {
		t.Fatalf("Error compiling regex: %v", err)
	}

	if !expectedGoarch.MatchString(buildInfo.Goarch) {
		t.Errorf("\nExpected output: %q\n\nActual output: %q", expectedGoarch, buildInfo.Goarch)
	} else {
		t.Log("Goarch OK")
	}
}
