package cmd_test

import (
	"bytes"
	"testing"

	"github.com/retroblast-engine/retroblast/cli/cmd"
)

func TestRootCommand(t *testing.T) {

	cmd := cmd.RootCommand()

	// Redirect the stdout to a buffer
	var buffer bytes.Buffer
	cmd.SetOut(&buffer)

	// Execute the root command with no arguments
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check output
	expectedOutput := `		
  _____      _             ____  _           _   
 |  __ \    | |           |  _ \| |         | |  
 | |__) |___| |_ _ __ ___ | |_) | | __ _ ___| |_ 
 |  _  // _ \ __| '__/ _ \|  _ <| |/ _` + "`" + ` / __| __|
 | | \ \  __/ |_| | | (_) | |_) | | (_| \__ \ |_ 
 |_|  \_\___|\__|_|  \___/|____/|_|\__,_|___/\__|
 ` + "\n" + `A retro 2D game engine using Go`

	if buffer.String() != expectedOutput {
		t.Errorf("Expected output: %q/nActual output: %q", expectedOutput, buffer)
	}
}
