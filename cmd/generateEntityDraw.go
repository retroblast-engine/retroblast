package cmd

import (
	"fmt"
	"strings"
)

func generateEntityDrawGo(entityName string) string {
	// Convert the entity name to lowercase for the package name and comments
	lowerEntityName := strings.ToLower(entityName)
	// Convert the entity name to title case for the struct name
	titleEntityName := strings.Title(entityName)
	// Get the first letter of the entity name in lowercase for the receiver
	receiver := strings.ToLower(string(entityName[0]))

	return fmt.Sprintf(`package %s

import (
    "github.com/hajimehoshi/ebiten/v2"
)

// draw draws the %s's sprite on the screen.
func (%s *%s) draw(screen *ebiten.Image, camX, camY float64) {
    opts := &ebiten.DrawImageOptions{}

    // Flip the sprite, if left is pressed
    if %s.direction == Left {
        %s.hFlip(opts)
    }

    opts.GeoM.Translate(%s.x, %s.y)
    opts.GeoM.Translate(camX, 0)
    // fmt.Printf("%s x: %%f, y: %%f\n", %s.x, %s.y)

    // Cache the current frame in the sprite field
    %s.sprite = %s.currentFrame()
    screen.DrawImage(%s.sprite, opts)
}

// hFlip modifies the DrawImageOptions to flip the sprite horizontally.
func (%s *%s) hFlip(opts *ebiten.DrawImageOptions) {
    opts.GeoM.Scale(-1, 1)
    opts.GeoM.Translate(float64(%s.sprite.Bounds().Dx()), 0)
}

// vFlip modifies the DrawImageOptions to flip the sprite vertically.
func (%s *%s) vFlip(opts *ebiten.DrawImageOptions) { // nolint: unused
    opts.GeoM.Scale(1, -1)
    opts.GeoM.Translate(0, float64(%s.sprite.Bounds().Dy()))
}

// hvFlip modifies the DrawImageOptions to flip the sprite horizontally and vertically.
func (%s *%s) hvFlip(opts *ebiten.DrawImageOptions) { // nolint: unused
    opts.GeoM.Scale(-1, -1)
    opts.GeoM.Translate(float64(%s.sprite.Bounds().Dx()), float64(%s.sprite.Bounds().Dy()))
}

// Draw draws the %s on the screen.
func (%s *%s) Draw(screen *ebiten.Image, camX, camY float64) {
    %s.draw(screen, camX, camY)
}
`, lowerEntityName, lowerEntityName, receiver, titleEntityName, receiver, receiver, receiver, receiver, lowerEntityName, receiver, receiver, receiver, receiver, receiver, receiver, titleEntityName, receiver, receiver, titleEntityName, receiver, receiver, titleEntityName, receiver, receiver, lowerEntityName, receiver, titleEntityName, receiver)
}
