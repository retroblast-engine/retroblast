package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

// getRepoURL retrieves the repository URL from the Git configuration
func getRepoURL() (string, error) {
	cmd := exec.Command("sh", "-c", "git config --list | grep 'github.com'")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "github.com") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				path := parts[1]
				// Remove the /.git/.gitcookies part
				path = strings.TrimSuffix(path, "/.git/.gitcookies")
				// Split the path and keep only the relevant parts
				urlParts := strings.Split(path, "/")
				if len(urlParts) >= 3 {
					return strings.Join(urlParts[len(urlParts)-3:], "/"), nil
				}
			}
		}
	}

	return "", fmt.Errorf("no valid repository URL found")
}
