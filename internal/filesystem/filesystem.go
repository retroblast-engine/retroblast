package filesystem

import (
	"fmt"
	"os"
)

func fileOrDirExists(path string) (bool, error) {

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("fileOrDirExists: failed to check if file or directory exists at %s: %w", path, err)
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
		return false, fmt.Errorf("isDirectory: failed to stat path %s: %w", path, err)
	}

	return fileInfo.IsDir(), err
}

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

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)

	if err != nil {
		// Handle the error appropriately
		return fmt.Errorf("DeleteFile: failed to delete file at %s: %w", filePath, err)
	}

	return nil
}
