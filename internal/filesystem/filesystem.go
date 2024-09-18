// Package filesystem provides basic file operations required by the engine
package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// fileOrDirExists checks if a file or directory exists at the given path and returns whether it's a directory.
func fileOrDirExists(path string) (exists bool, isDir bool, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false, nil
		}
		return false, false, fmt.Errorf("failed to check if file or directory exists at %s: %w", path, err)
	}
	return true, fileInfo.IsDir(), nil
}

// sanitizePath validates and sanitizes the input path.
func sanitizePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path to sanitize cannot be empty")
	}

	// Clean the path to remove any redundant or potentially harmful elements.
	cleanPath := filepath.Clean(path)

	// Ensure the path is absolute.
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %s: %w", path, err)
	}

	// Check for potentially dangerous characters or sequences.
	if strings.Contains(absPath, "..") {
		return "", fmt.Errorf("path contains invalid sequences: %s", absPath)
	}

	return absPath, nil
}

// isDiskError checks if the error is related to disk issues, such as no space left.
func isDiskError(err error) bool {
	return strings.Contains(err.Error(), "no space left") ||
		strings.Contains(err.Error(), "disk quota exceeded")
}

// CreateFile creates a new file at the given `path` argument.
func CreateFile(path string) error {
	// Sanitize the input path.
	sanitizedPath, err := sanitizePath(path)
	if err != nil {
		return fmt.Errorf("invalid path provided: %s - %w", path, err)
	}

	// Check if the file or directory already exists.
	exists, isDir, err := fileOrDirExists(sanitizedPath)
	if err != nil {
		return fmt.Errorf("failed to check existence of file at %s: %w", sanitizedPath, err)
	}
	if exists {
		if isDir {
			return fmt.Errorf("a directory already exists at path: %s", sanitizedPath)
		}
		return fmt.Errorf("file already exists at %s", sanitizedPath)
	}

	// Create the file.
	file, err := os.Create(sanitizedPath)
	if err != nil {
		// Handle permission errors.
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied when trying to create file at %s: %w", sanitizedPath, err)
		}

		// Handle path not found errors.
		if os.IsNotExist(err) {
			dir := filepath.Dir(sanitizedPath)
			return fmt.Errorf("directory %s does not exist. Please create the directory or check the path: %w", dir, err)
		}

		// Handle disk-related errors (e.g., no space left).
		if isDiskError(err) {
			return fmt.Errorf("disk error (e.g., no space left) when creating file at %s: %w", sanitizedPath, err)
		}

		// Generic fallback error handler.
		return fmt.Errorf("unexpected error while creating file at %s: %w", sanitizedPath, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close file at %s: %w", sanitizedPath, cerr)
		}
	}()

	return nil
}

// isFileInUseError checks if the error is related to the file being in use or locked.
func isFileInUseError(err error) bool {
	// On Windows, "The process cannot access the file because it is being used by another process" or "resource busy"
	return strings.Contains(err.Error(), "file is in use") ||
		strings.Contains(err.Error(), "resource busy")
}

// DeleteFile deletes the file at the given `path` argument, provided it exists.
func DeleteFile(filePath string) error {
	// Sanitize the input path.
	sanitizedPath, err := sanitizePath(filePath)
	if err != nil {
		return fmt.Errorf("invalid path provided for removal: %s - %w", filePath, err)
	}

	// Check if the file or directory exists.
	exists, isDir, err := fileOrDirExists(sanitizedPath)
	if err != nil {
		return fmt.Errorf("failed to check existence at %s: %w", sanitizedPath, err)
	}
	// Handle cases where the filepath is a directory.
	if isDir {
		return fmt.Errorf("path: %s is a directory, not a file. Use a function that deletes directories instead", sanitizedPath)
	}
	if !exists {
		return fmt.Errorf("file does not exist at %s", sanitizedPath)
	}

	// Attempt to delete the file.
	err = os.Remove(sanitizedPath)
	if err != nil {
		// Handle permission errors.
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied when trying to delete %s: %w", sanitizedPath, err)
		}

		// Handle cases where the file does not exist (race condition).
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist at %s: %w", sanitizedPath, err)
		}

		// Handle cases where the file is in use or locked by another process.
		if isFileInUseError(err) {
			return fmt.Errorf("file at %s is in use or locked by another process: %w", sanitizedPath, err)
		}

		// Generic fallback error handler.
		return fmt.Errorf("unexpected error while deleting file at %s: %w", sanitizedPath, err)
	}

	return nil
}
