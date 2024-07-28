package filesystem

import (
	"strings"
	"testing"
)

// TestCreateAndDeleteSuccessScenario tests successful creation and deletion of a file.
func TestCreateAndDeleteSuccessScenario(t *testing.T) {
	filePath := "testdata/testcreateanddeletescenario.txt"
	err := CreateFile(filePath)
	if err != nil {
		t.Fatalf(`CreateFile(%q) returned error: %v`, filePath, err)
	}

	err = DeleteFile(filePath)
	if err != nil {
		t.Fatalf(`DeleteFile(%q) returned error: %v`, filePath, err)
	}
}

// TestCreateFileEmpty tests CreateFile with an empty string.
func TestCreateFileEmpty(t *testing.T) {
	err := CreateFile("")
	if err == nil {
		t.Fatalf(`CreateFile("") returned no error, expected error`)
	} else {
		t.Logf(`CreateFile("") returned expected error: %v`, err)
	}
}

// TestCreateFileAlreadyExists tests CreateFile with a file that already exists.
func TestCreateFileAlreadyExists(t *testing.T) {
	filePath := "testdata/test_file_do_not_delete.txt"
	err := CreateFile(filePath)
	if err == nil {
		t.Fatalf(`CreateFile(%q) returned no error, expected error`, filePath)
	} else {
		t.Logf(`CreateFile(%q) returned expected error: %v`, filePath, err)
	}
}

// TestDeleteFileEmpty tests DeleteFile with an empty string.
func TestDeleteFileEmpty(t *testing.T) {
	err := DeleteFile("")
	if err == nil {
		t.Fatalf(`DeleteFile("") returned no error, expected error`)
	} else {
		t.Logf(`DeleteFile("") returned expected error: %v`, err)
	}
}

// TestDeleteFileNotExisting tests DeleteFile with a non-existing file.
func TestDeleteFileNotExisting(t *testing.T) {
	err := DeleteFile("dummydummy123456789.txt")
	if err == nil {
		t.Fatalf(`DeleteFile("dummydummy123456789.txt") returned no error, expected error`)
	} else if strings.Contains(err.Error(), "no such file or directory") || strings.Contains(err.Error(), "The system cannot find the file specified") {
		// Expected error
		t.Logf(`DeleteFile("dummydummy123456789.txt") returned expected error: %v`, err)
	} else {
		t.Fatalf(`DeleteFile("dummydummy123456789.txt") returned unexpected error: %v`, err)
	}
}
