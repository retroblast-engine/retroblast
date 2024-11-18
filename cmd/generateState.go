package cmd

import (
	"fmt"
	"strings"
)

func generateStateGo(entityName, statesString string) string {
	// Convert the entity name to lowercase for the package name and comments
	lowerEntityName := strings.ToLower(entityName)

	return fmt.Sprintf(`package %s

// State represents the state of the %s.
type State int

%s
`, lowerEntityName, lowerEntityName, statesString)
}
