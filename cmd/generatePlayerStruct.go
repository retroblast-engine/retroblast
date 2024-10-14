package cmd

import (
	"fmt"
	"strings"
)

func generatePlayerGo(importPath, entityName, defaultState string) string {
	// Convert the entity name to lowercase for the package name and comments
	lowerEntityName := strings.ToLower(entityName)
	// Convert the entity name to title case for the struct name
	titleEntityName := strings.Title(entityName)
	// Get the first letter of the entity name in lowercase for the receiver
	receiver := strings.ToLower(string(entityName[0]))

	return fmt.Sprintf(`package %s

import (
    "fmt"

    "%s/settings"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/retroblast-engine/asevre"
    "github.com/solarlune/resolv"
)

// %s represents the %s entity.
type %s struct {
    x, y      float64
    file      asevre.ASEFile
    sprite    *ebiten.Image
    state     State
    animation *Animation
    direction Direction
    Object         *resolv.Object
    Speed          resolv.Vector
    OnGround       *resolv.Object
    SlidingOnWall  *resolv.Object
    FacingRight    bool
    IgnorePlatform *resolv.Object
}

// New%s creates a new %s instance.
func New%s(x, y float64, file string) (*%s, error) {
    aseprite, err := asevre.ParseAseprite(settings.Assets, file)
    if err != nil {
        return nil, fmt.Errorf("failed to parse aseprite file: %%w", err)
    }

    %s := &%s{
        x:         x,
        y:         y,
        file:      aseprite,
        state:     %s,
        direction: Right,
    }

    // Initialize the animation for %s's state
    %s.initAnimation(%s.state)

    // Create a new resolv object for %s
    width := float64(%s.sprite.Bounds().Dx())
    height := float64(%s.sprite.Bounds().Dy())
    %s.Object = resolv.NewObject(x, y, width, height, "%s")
    %s.Object.SetShape(resolv.NewRectangle(x, y, %s.Object.Size.X, %s.Object.Size.Y))

    return %s, nil
}

// GetObject returns the underlying resolv.Object.
func (%s *%s) GetObject() *resolv.Object {
    return %s.Object
}
`, lowerEntityName, importPath, titleEntityName, lowerEntityName, titleEntityName, titleEntityName, titleEntityName, titleEntityName, titleEntityName, receiver, titleEntityName, defaultState, lowerEntityName, receiver, receiver, receiver, receiver, receiver, receiver, lowerEntityName, receiver, receiver, receiver, receiver, receiver, titleEntityName, receiver)
}
