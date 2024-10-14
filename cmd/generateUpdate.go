package cmd

func generateUpdateGo() string {
	return `package main

// update is the place where you can write your game logic
// such as changing scenes, deleting entities, pause, gameover, etc.
func (g *Game) update() error {
    // Write your game logic here
    return nil
}
`
}
