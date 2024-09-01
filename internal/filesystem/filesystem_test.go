package filesystem

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var (
	fileToCreate                 = "testdata/example_createfile.txt"
	fileToDelete                 = "testdata/example_deletefile.txt"
	fileToCreateAndDelete        = "testdata/testcreateanddeletescenario.txt"
	fileThatAlreadyExists        = "testdata/test_file_do_not_delete.txt"
	fileThatDoesNotExist         = "dummydummy123456789.txt"
	filePathForDeletionThatIsDir = "testdata/test_delete_directory"
	directoryThatAlreadyExists   = "testdata"
)

// Helper function to create a test file.
func createTestFile(t *testing.T, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test file %q: %v", filePath, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Errorf("Failed to close file %q: %v", filePath, err)
		}
	}()
}

// Helper function to remove a test file.
func removeTestFile(t *testing.T, filePath string) {
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to remove test file %q: %v", filePath, err)
	}
}

// ExampleCreateFile demonstrates how to use the CreateFile function.
func ExampleCreateFile() {
	filePath := fileToCreate

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
	filePath := fileToDelete

	// Create a file to delete
	file, _ := os.Create(filePath)

	err := file.Close()

	if err != nil {
		fmt.Println("Failed to close file at ExampleDeleteFile invocation")
	}

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
	filePath := fileToCreateAndDelete
	defer removeTestFile(t, filePath)

	// Create the file
	err := CreateFile(filePath)
	if err != nil {
		t.Fatalf(`CreateFile(%q) returned error: %v`, filePath, err)
	}

	// Delete the file
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
	filePath := fileThatAlreadyExists
	createTestFile(t, filePath)
	defer removeTestFile(t, filePath)

	err := CreateFile(filePath)
	if err == nil {
		t.Fatalf(`CreateFile(%q) returned no error, expected an error`, filePath)
	} else {
		if strings.Contains(err.Error(), "file already exists") {
			t.Logf(`CreateFile(%q) returned expected error: %v`, filePath, err)
		} else {
			t.Fatalf(`CreateFile(%q) returned unexpected error: %v`, filePath, err)
		}
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
	filePath := fileThatDoesNotExist
	err := DeleteFile(filePath)
	if err == nil {
		t.Fatalf(`DeleteFile(%q) returned no error, expected an error`, filePath)
	} else {
		if strings.Contains(err.Error(), "file does not exist at") {
			t.Logf(`DeleteFile(%q) returned expected error: %v`, filePath, err)
		} else {
			t.Fatalf(`DeleteFile(%q) returned unexpected error: %v`, filePath, err)
		}
	}
}

// TestCreateFileInvalidPath tests CreateFile with an invalid path.
func TestCreateFileWithPathThatIsDirectory(t *testing.T) {
	filePath := directoryThatAlreadyExists

	err := CreateFile(filePath)
	if err == nil {
		t.Fatalf(`CreateFile(%q) returned no error, expected an error`, filePath)
	} else {
		if strings.Contains(err.Error(), "a directory already exists at path") {
			t.Logf(`CreateFile(%q) returned expected error: %v`, filePath, err)
		} else {
			t.Fatalf(`CreateFile(%q) returned unexpected error: %v`, filePath, err)
		}
	}
}

// TestDeleteFileDirectory tests DeleteFile with a directory path.
func TestDeleteFileDirectory(t *testing.T) {
	dirPath := filePathForDeletionThatIsDir
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory %q: %v", dirPath, err)
	}
	defer removeTestFile(t, dirPath)

	err = DeleteFile(dirPath)
	if err == nil {
		t.Fatalf(`DeleteFile(%q) returned no error, expected an error`, dirPath)
	} else {
		if strings.Contains(err.Error(), "is a directory") {
			t.Logf(`DeleteFile(%q) returned expected error: %v`, dirPath, err)
		} else {
			t.Fatalf(`DeleteFile(%q) returned unexpected error: %v`, dirPath, err)
		}
	}
}
