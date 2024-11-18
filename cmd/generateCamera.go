package cmd

import "fmt"

func generateCameraGo(importPath string) string {
	return fmt.Sprintf(`package camera

import (
    "math"

    "%s/settings"
)

// Camera is the main struct that holds the camera state
type Camera struct {
    X, Y float64
}

// New creates a new camera
func New(x, y float64) *Camera {
    return &Camera{
        X: x,
        Y: y,
    }
}

// FollowTargetX moves the camera to follow the target on the X axis
func (c *Camera) FollowTargetX(targetX float64) {
    c.X = -targetX + float64(settings.ScaledWidth)/2
}

// FollowTargetY moves the camera to follow the target on the Y axis
func (c *Camera) FollowTargetY(targetY float64) {
    c.Y = -targetY + float64(settings.ScaledHeight)/2
}

// Constrain constrains the camera to the tile map width and height so it doesn't go out of bounds
func (c *Camera) Constrain(tileMapWidthPixels, tileMapHeightPixels int) {
    // Constrain to the top left corner
    c.X = math.Min(c.X, 0)
    c.Y = math.Min(c.Y, 0)

    c.X = math.Max(c.X, float64(settings.ScaledWidth-tileMapWidthPixels))
    c.Y = math.Max(c.Y, float64(settings.ScaledHeight-tileMapHeightPixels))
}
`, importPath)
}
