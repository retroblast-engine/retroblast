package filesystem

import (
	"errors"
	"fmt"
	"os"
)

func fileOrDirExists(path string) (bool, error) {

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("failed to check if file or directory exists: %w", err)
		}
	} else {
		return true, nil
	}

}

// isDirectory determines if a file represented
// by `path` is a directory or not
func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false, fmt.Errorf("failed to stat path %s: %w", path, err)
	}

	return fileInfo.IsDir(), err
}

func CreateFile(path string) (err error) {

	exists, err := fileOrDirExists(path)
	if err != nil {
		return fmt.Errorf("failed to check if file or directory exists: %w", err)
	}
	if exists {
		return errors.New("file or directory already exists at path: " + path)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	return nil

}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)

	if err != nil {
		// Handle the error appropriately
		return fmt.Errorf("Failed to delete file %s: %w", filePath, err)
	}

	return nil
}
