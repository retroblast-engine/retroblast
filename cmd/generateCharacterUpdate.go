package cmd

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func generateCharacterUpdateGo(importPath, entityName string) string {
	// Convert the entity name to lowercase for the package name and comments
	lowerEntityName := strings.ToLower(entityName)
	// Convert the entity name to title case for the struct name
	c := cases.Title(language.Und)
	titleEntityName := c.String(entityName)
	// Get the first letter of the entity name in lowercase for the receiver
	receiver := strings.ToLower(string(entityName[0]))

	return fmt.Sprintf(`package %s

import (
    "time"

    "%s/settings"
)

// handleInput handles the %s's input.
func (%s *%s) handleInput() { // nolint: unused
}

// Update updates the %s's state based on input.
func (%s *%s) Update() {
    if %s.state == Dead {
        // if 2 seconds have passed since the %s died, destroy it.
        if time.Since(%s.deadTime) >= 2*time.Second {
            %s.IsDestroyed = true
        }
        return
    }

    friction := 0.25
    accel := 0.25 + friction
    maxSpeed := 0.5

    // Move left or right.
    if %s.direction == Right {
        %s.Speed.X += accel
    } else {
        %s.Speed.X -= accel
    }

    // Apply friction and horizontal speed limiting.
    if %s.Speed.X > friction {
        %s.Speed.X -= friction
    } else if %s.Speed.X < -friction {
        %s.Speed.X += friction
    } else {
        %s.Speed.X = 0
    }

    // Prevent %s from moving too fast.
    if %s.Speed.X > maxSpeed {
        %s.Speed.X = maxSpeed
    } else if %s.Speed.X < -maxSpeed {
        %s.Speed.X = -maxSpeed
    }

    // Handle horizontal movement.
    dx := %s.Speed.X

    if check := %s.Object.Check(%s.Speed.X, 0, "solid", "%s"); check != nil {
        // If we come into contact with a solid object, we move as close as possible to contact with it, and stop horizontal movement.
        dx = check.ContactWithCell(check.Cells[0]).X
        %s.Speed.X = 0
        if %s.direction == Right {
            %s.direction = Left
        } else {
            %s.direction = Right
        }
        dx = -dx
    }

    %s.Object.Position.X += dx

    if check := %s.Object.Check(%s.Speed.X, 0, "player"); check != nil {
        // get object from tag
        player := check.Objects[0]
        if player.Position.Y+settings.MaxCellSize < %s.Object.Position.Y {
            %s.Speed.X = 0
            %s.Object.Position.X -= dx
            %s.SetState(Dead)
            %s.deadTime = time.Now()
        }
    }

    %s.Object.Update()
    %s.x = %s.Object.Position.X
    %s.y = %s.Object.Position.Y
}
`, lowerEntityName, importPath, lowerEntityName, receiver, titleEntityName, lowerEntityName, receiver, titleEntityName, receiver, lowerEntityName, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, lowerEntityName, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver)
}
