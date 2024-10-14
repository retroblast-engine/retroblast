package cmd

// generateHelper generates the helper.go file.
func generateHelper() string {
	return `package main

import (
    "image/color"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/vector"
    "github.com/solarlune/resolv"
)

// debugDraw draws the debug information of the game.
func (g *Game) debugDraw(screen *ebiten.Image, space *resolv.Space) {

    for y := 0; y < space.Height(); y++ {

        for x := 0; x < space.Width(); x++ {

            cell := space.Cell(x, y)

            cw := float32(space.CellWidth)
            ch := float32(space.CellHeight)
            cam := g.Scene.GetCamera()
            cx := (float32(cell.X)*cw + float32(cam.X))
            cy := (float32(cell.Y) * ch)

            drawColor := color.RGBA{20, 20, 20, 255}

            if cell.Occupied() {
                drawColor = color.RGBA{255, 255, 0, 255}
            }

            // Draw the right and bottom lines of the current cell
            vector.StrokeLine(screen, cx+cw, cy, cx+cw, cy+ch, 1, drawColor, false) // draw right line
            vector.StrokeLine(screen, cx, cy+ch, cx+cw, cy+ch, 1, drawColor, false) // draw bottom line

            // Draw the top line if this is the first row or the cell above is not occupied
            if y == 0 || !space.Cell(x, y-1).Occupied() {
                vector.StrokeLine(screen, cx, cy, cx+cw, cy, 1, drawColor, false) // draw top line
            }

            // Draw the left line if this is the first column or the cell to the left is not occupied
            if x == 0 || !space.Cell(x-1, y).Occupied() {
                vector.StrokeLine(screen, cx, cy, cx, cy+ch, 1, drawColor, false) // draw left line
            }
        }
    }
}
`
}
