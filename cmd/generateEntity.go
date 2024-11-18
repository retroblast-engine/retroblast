package cmd

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func generateEntityGo(importPath string, entities []string) string {
	// Generate the import statements for entities
	var entityImports strings.Builder
	for _, entity := range entities {
		entityImports.WriteString(fmt.Sprintf("\t\"%s/entities/%s\"\n", importPath, entity))
	}

	// Generate the EntityConstructors map entries
	var entityConstructors strings.Builder
	for _, entity := range entities {
		c := cases.Title(language.Und)
		entityConstructors.WriteString(fmt.Sprintf(
			"\t\"%s\": func(x, y float64, assetPath string) (Entity, error) {\n\t\treturn %s.New%s(x, y, assetPath)\n\t},\n",
			entity, entity, c.String(entity),
		))
	}

	return fmt.Sprintf(`package entities

import (
%s
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/solarlune/resolv"
)

// Entity represents a game entity with update and draw methods.
type Entity interface {
    Update()
    Draw(screen *ebiten.Image, camX, camY float64)
    GetObject() *resolv.Object
}

// RemoveEntity removes a specific entity from the slice of entities.
func RemoveEntity(entities []Entity, target Entity) []Entity {
    for i, entity := range entities {
        if entity == target {
            // Remove the entity by slicing around it
            return append(entities[:i], entities[i+1:]...)
        }
    }
    return entities
}

type EntityConstructor func(x, y float64, assetPath string) (Entity, error)

var EntityConstructors = map[string]EntityConstructor{
%s
}
`, entityImports.String(), entityConstructors.String())
}
