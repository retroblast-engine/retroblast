package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// getRepoURL retrieves the repository URL from the Git configuration
func getRepoURL() (string, error) {
	var err error

	// Get the GOPATH from the environment variable
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		return "", fmt.Errorf("GOPATH is not set")
	}

	// Retrieves the current working directory directly, ensuring compatibility across platforms.
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// Ensure the current directory is under GOPATH/src
	if !strings.HasPrefix(currentDir, filepath.Join(goPath, "src")) {
		return "", fmt.Errorf("current directory is not under GOPATH/src")
	}

	// Extract the path relative to GOPATH/src
	relativePath := strings.TrimPrefix(currentDir, filepath.Join(goPath, "src")+string(os.PathSeparator))

	// Find the "github.com" segment in the path
	index := strings.Index(relativePath, "github.com")
	if index == -1 {
		return "", fmt.Errorf("no valid repository URL found")
	}

	// Return the repository path starting from "github.com"
	return relativePath[index:], nil
}
