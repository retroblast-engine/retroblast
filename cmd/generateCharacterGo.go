package cmd

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func generateCharacterGo(importPath, entityName, defaultState string) string {
	// Convert the entity name to lowercase for the package name and comments
	lowerEntityName := strings.ToLower(entityName)
	// Convert the entity name to title case for the struct name
	c := cases.Title(language.Und)
	titleEntityName := c.String(entityName)
	// Get the first letter of the entity name in lowercase for the receiver
	receiver := strings.ToLower(string(entityName[0]))

	return fmt.Sprintf(`package %s

import (
    "fmt"
    "time"

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

    // New ones
    Object      *resolv.Object
    Speed       resolv.Vector
    deadTime    time.Time
    IsDestroyed bool
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

    %s.initAnimation(%s.state)

    // Create new resolv object for %s
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
`, lowerEntityName, importPath, titleEntityName, lowerEntityName, titleEntityName, titleEntityName, titleEntityName, titleEntityName, titleEntityName, receiver, titleEntityName, defaultState, receiver, receiver, receiver, receiver, receiver, receiver, lowerEntityName, receiver, receiver, receiver, receiver, receiver, titleEntityName, receiver)
}

//  titleEntityName, receiver
