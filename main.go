// main package is the entry point of the program
package main

import (
	"fmt"
  
  "github.com/retroblast-engine/retroblast/cmd"
)

func main() {
	fmt.Println("Hello!")
  cmd.Execute()
}
