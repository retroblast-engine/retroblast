package cmd

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func generateAnimationGo(entityName string) string {
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

    "github.com/hajimehoshi/ebiten/v2"
)

// Animation holds the animation details for a state.
type Animation struct {
    index       int
    totalFrames int
    lastChange  time.Time
    duration    []time.Duration
}

// initAnimation initializes the animation for the given state.
func (%s *%s) initAnimation(state State) {
    %s.animation = &Animation{
        index:       0,
        totalFrames: len(%s.file.State[state].Frames),
        lastChange:  time.Now(),
        duration:    %s.file.State[state].FrameDuration[state],
    }
    %s.sprite = %s.file.State[state].Frames[0]
}

// spriteFrame returns the frame image for the given state and frame index.
func (%s *%s) spriteFrame(state State, frameIdx int) *ebiten.Image {
    return %s.file.State[state].Frames[frameIdx]
}

// nextFrame returns the next frame image for the current state.
func (%s *%s) nextFrame() *ebiten.Image {
    if time.Since(%s.animation.lastChange) >= %s.animation.duration[%s.animation.index] {
        %s.animation.index++
        if %s.animation.index >= %s.animation.totalFrames {
            %s.animation.index = 0
        }
        %s.animation.lastChange = time.Now()
    }

    return %s.file.State[%s.state].Frames[%s.animation.index]
}

// currentFrame returns the current frame image for the %s's state.
func (%s *%s) currentFrame() *ebiten.Image {
    if %s.file.State[%s.state].HasAnimations {
        return %s.nextFrame()
    }
    return %s.spriteFrame(%s.state, 0)
}

// SetState updates the %s's state and reinitializes the animation if the state has changed.
func (%s *%s) SetState(state State) {
    if %s.state != state {
        %s.state = state
        %s.initAnimation(state)
    }
}
`, lowerEntityName, receiver, titleEntityName, receiver, receiver, receiver, receiver, receiver, receiver, titleEntityName, receiver, receiver, titleEntityName, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, receiver, lowerEntityName, receiver, titleEntityName, receiver, receiver, receiver, receiver, receiver, lowerEntityName, receiver, titleEntityName, receiver, receiver, receiver)
}
