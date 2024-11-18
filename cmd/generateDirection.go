package cmd

import (
	"fmt"
	"strings"
)

func generateDirectionGo(entityName string) string {
	// Convert the entity name to lowercase for the package name and comments
	lowerEntityName := strings.ToLower(entityName)

	return fmt.Sprintf(`package %s

// Direction represents %s's direction.
type Direction int

const (
    Right Direction = iota
    Left
)
`, lowerEntityName, lowerEntityName)
}
