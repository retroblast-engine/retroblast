package cmd

func generateDrawGo() string {
	return `package main

import "github.com/hajimehoshi/ebiten/v2"

// draw is the place where you can write your drawing logic
func (g *Game) draw(screen *ebiten.Image) {
    // Write your drawing logic here
}
`
}
