package cmd

import "fmt"

func generateMainGo(sceneName, importPath string) string {
	return fmt.Sprintf(`package main

import (
    "log"

    "%s/scenes"
    "%s/settings"
    "github.com/hajimehoshi/ebiten/v2"
)

func main() {
    // Load the embedded files
    settings.SetEmbeddedFiles(embeddedFiles)

    // Set the window size and title
    ebiten.SetWindowSize(settings.ScreenWidth, settings.ScreenHeight)
    ebiten.SetWindowTitle(settings.WindowTitle)

    // Initialize the game
    game := &Game{}

    // Set to the initial scene
    game.changeScene(scenes.%s)

    // Run the game
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}
`, importPath, importPath, sceneName)
}
