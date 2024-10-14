package cmd

import "fmt"

func generateGameGo(importPath string) string {
	return fmt.Sprintf(`package main

import (
    "%s/scenes"
    "%s/settings"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game is the main struct that holds the game state
type Game struct {
    Debug bool
    Scene scenes.Scene
}

// Update is the main function that updates the game
func (g *Game) Update() error {
    // Store the current scene ID
    currentSceneID := g.Scene.GetSceneID()

    // Toggle debug mode when F1 is pressed
    if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
        g.Debug = !g.Debug
    }

    // Injecting your game logic here
    if err := g.update(); err != nil {
        return err
    }

    // Update the current scene, if it exists
    if g.Scene != nil {
        if err := g.Scene.Update(); err != nil {
            return err
        }
    }

    // Check if the scene has changed
    newSceneID := g.Scene.GetSceneID()
    if newSceneID != currentSceneID {
        g.changeScene(newSceneID)
    }

    return nil
}

// Draw is the main function that draws the game
func (g *Game) Draw(screen *ebiten.Image) {
    // Fill the screen with black
    screen.Clear()

    // Draw the current scene, if it exists
    if g.Scene != nil {
        g.Scene.Draw(screen)
    }

    g.draw(screen) // Write your drawing logic in this method

    // Draw debug information if debug mode is enabled
    if g.Debug {
        space := g.Scene.GetSpace()
        g.debugDraw(screen, space)
    }
}

// Layout is the main function that sets the window size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return settings.ScaledWidth, settings.ScaledHeight
}

// changeScene changes the current scene to a new one and runs the Init() function
func (g *Game) changeScene(newSceneID scenes.SceneID) {
    newScene := scenes.SceneMap[newSceneID]

    // Destroy the current scene, if it exists
    if g.Scene != nil {
        g.Scene.Destroy()
    }

    // Set and Initialize the new scene
    g.Scene = newScene
    g.Scene.Init()
}
`, importPath, importPath)
}
