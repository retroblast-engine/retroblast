package filesystem

import (
	"testing"
)

func TestCreateAndDeleteSuccessScenario(t *testing.T) {
	err := CreateFile("testdata/testcreateanddeletescenario.txt")

	if err != nil {
		t.Fatalf(`CreateFile("testdata/testcreateanddeletescenario.txt") returned %v , error`, err)
	}

	err = DeleteFile("testdata/testcreateanddeletescenario.txt")

	if err != nil {
		t.Fatalf(`DeleteFile("testdata/testcreateanddeletescenario.txt") returned %v , error`, err)
	}

}

// TestCreateFileEmpty calls filesystem.CreateFile with an empty string,
// checking for an error.
func TestCreateFileEmpty(t *testing.T) {
	err := CreateFile("")
	if err == nil {
		t.Fatalf(`CreateFile("") returned %v , error`, err)
	}
}

// TestCreateFileEmpty calls filesystem.CreateFile with an a filepath to a file that already exists
// checking for an error
func TestCreateFileAlreadyExists(t *testing.T) {
	err := CreateFile("testdata/test_file_do_not_delete.txt")
	if err == nil {
		t.Fatalf(`CreateFile("testdata/test_file_do_not_delete.txt") returned %v , error`, err)
	}
}

// TestDeleteFileEmpty calls filesystem.DeleteFile with an empty string,
// checking for an error.
func TestDeleteFileEmpty(t *testing.T) {
	err := DeleteFile("")
	if err == nil {
		t.Fatalf(`DeleteFile("") returned %v , error`, err)
	}
}

// TestDeleteFileEmpty calls filesystem.DeleteFile with an empty string,
// checking for an error.
func TestDeleteFileNotExisting(t *testing.T) {
	err := DeleteFile("dummydummy123456789.txt")

	if err == nil {
		t.Fatalf(`DeleteFile("dummydummy123456789.txt") returned %v , error`, err)
	}
}
