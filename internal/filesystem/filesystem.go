// Package filesystem provides basic file operations required by the engine
package filesystem

import (
	"fmt"
	"os"
)

func fileOrDirExists(path string) (bool, error) {

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("fileOrDirExists: failed to check if file or directory exists at %s: %w", path, err)
	}

	return true, nil

}

// CreateFile creates a new file at the given `path` argument
func CreateFile(path string) (err error) {

	exists, err := fileOrDirExists(path)
	if err != nil {
		return fmt.Errorf("CreateFile: failed to check if file or directory exists at %s: %w", path, err)
	}
	if exists {
		return fmt.Errorf("CreateFile: file or directory already exists at path %s", path)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("CreateFile: failed to create file at %s: %w", path, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("CreateFile: failed to close file at %s: %w", path, cerr)
		}
	}()

	return nil

}

// DeleteFile deletes the file presented by the given `path` argument, provided it exists
func DeleteFile(filePath string) error {
	err := os.Remove(filePath)

	if err != nil {
		// Handle the error appropriately
		return fmt.Errorf("DeleteFile: failed to delete file at %s: %w", filePath, err)
	}

	return nil
}
