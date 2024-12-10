package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// getRepoURL retrieves the repository URL from the Git configuration
func getRepoURL() (string, error) {

	var cmd *exec.Cmd
	var err error
	var output []byte

	// Get the GOPATH from the environment variable
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		return "", fmt.Errorf("GOPATH is not set")
	}

	// Check the operating system
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf("$pwd = Get-Location; $pwd.Path -split \"%s\" | Select-Object -Index 1", "/"))
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("sh", "-c", fmt.Sprintf("pwd | awk -F \"%s\" '{print $2}'", "/"))
	} else {
		fmt.Println("Unsupported OS")
		return "", fmt.Errorf("unsupported OS")
	}

	output, err = cmd.Output()

	fmt.Println(output)

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
