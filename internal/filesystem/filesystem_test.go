package filesystem

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// ExampleCreateFile demonstrates how to use the CreateFile function.
func ExampleCreateFile() {
	filePath := "testdata/example_createfile.txt"

	// Create the file
	err := CreateFile(filePath)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("File created successfully")
	}

	// Clean up by deleting the file
	_ = os.Remove(filePath)

	// Output:
	// File created successfully
}

// ExampleDeleteFile demonstrates how to use the DeleteFile function.
func ExampleDeleteFile() {
	filePath := "testdata/example_deletefile.txt"

	// Create a file to delete
	file, err := os.Create(filePath)

	file.Close()

	// Delete the file
	err = DeleteFile(filePath)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("File deleted successfully")
	}

	// Output:
	// File deleted successfully
}

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
		t.Fatalf(`CreateFile("") returned no error, expected an error`)
	} else {
		t.Logf(`CreateFile("") returned expected error: %v`, err)
	}
}

// TestCreateFileAlreadyExists tests CreateFile with a file that already exists.
func TestCreateFileAlreadyExists(t *testing.T) {
	filePath := "testdata/test_file_do_not_delete.txt"
	err := CreateFile(filePath)
	if err == nil {
		t.Fatalf(`CreateFile(%q) returned no error, expected an error`, filePath)
	} else {
		t.Logf(`CreateFile(%q) returned expected error: %v`, filePath, err)
	}
}

// TestDeleteFileEmpty tests DeleteFile with an empty string.
func TestDeleteFileEmpty(t *testing.T) {
	err := DeleteFile("")
	if err == nil {
		t.Fatalf(`DeleteFile("") returned no error, expected an error`)
	} else {
		t.Logf(`DeleteFile("") returned expected error: %v`, err)
	}
}

// TestDeleteFileNotExisting tests DeleteFile with a non-existing file.
func TestDeleteFileNotExisting(t *testing.T) {
	err := DeleteFile("dummydummy123456789.txt")
	if err == nil {
		t.Fatalf(`DeleteFile("dummydummy123456789.txt") returned no error, expected an error`)
	} else if strings.Contains(err.Error(), "no such file or directory") || strings.Contains(err.Error(), "The system cannot find the file specified") {
		// Expected error
		t.Logf(`DeleteFile("dummydummy123456789.txt") returned expected error: %v`, err)
	} else {
		t.Fatalf(`DeleteFile("dummydummy123456789.txt") returned unexpected error: %v`, err)
	}
}
