package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// getRepoURL retrieves the repository URL from the Git configuration
func getRepoURL() (string, error) {

	// Get the GOPATH from the environment variable
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		return "", fmt.Errorf("GOPATH is not set")
	}

	cmd := exec.Command("sh", "-c", "pwd | awk -F "+goPath+" '{print $2}'")
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	path := string(output)
	index := strings.Index(path, "github.com")
	if index == -1 {
		return "", fmt.Errorf("no valid repository URL found")
	}
	path = path[index:]

	return path, nil
}
