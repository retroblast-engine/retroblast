package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/retroblast-engine/asevre"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// GenerateCommand represents the generate command
func GenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate code for the game",
		RunE:  runGenerateCommand,
	}

	// Add any flags here if needed
	// cmd.Flags().StringVarP(&flagVar, "flag", "f", "", "Description of the flag")

	return cmd
}

// runGenerateCommand is the function that runs when the generate command is called
func runGenerateCommand(_ *cobra.Command, _ []string) error {

	if err := generateAseprite(); err != nil {
		return err
	}

	if err := initializeGoModule(); err != nil {
		return err
	}

	if err := tidyGoModule(); err != nil {
		return err
	}

	return nil
}

// generateAseprite handles the generation of Aseprite assets
func generateAseprite() error {
	// Add logic for generating Aseprite assets here
	log.Println("Generating Aseprite assets...")
	return nil
}

// initializeGoModule initializes the Go module if it doesn't already exist
func initializeGoModule() error {
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		cmd := exec.Command("go", "mod", "init")
		if err := cmd.Run(); err != nil {
			return err
		}
		log.Println("Go module initialized.")
	} else {
		log.Println("Go module already initialized. Skipping 'go mod init'.")
	}
	return nil
}

// tidyGoModule tidies up the Go module
func tidyGoModule() error {
	cmd := exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Println("Go module tidied.")
	return nil
}

// genDeclarations generates the declaration string for a player
func genDeclarations(base string, x, y int) string {
	// make sure base is allo lowercase
	base = strings.ToLower(base)
	// Make first letter of the base name uppercase
	Base := cases.Title(language.Und).String(base) // Capitalize the first letter
	return fmt.Sprintf("var my%s, _ = %s.New%s(%d, %d, \"assets/aseprite/%s.aseprite\")\n", base, base, Base, x, y, base)
}

func genImportPaths(p []string) string {
	var imports string
	for _, path := range p {
		// If it has suffix "*.aseprite", remove it
		str := strings.TrimSuffix(path, ".aseprite")
		imports += fmt.Sprintf("\"%s\"\n", str)
	}
	return imports
}

var importPaths []string

// GenerateAseprite generates the code for the game based on the Aseprite files
func GenerateAseprite() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Extract the relevant part of the path for the import statement
	// Assuming the structure is always /Users/username/gocode/src/github.com/drpaneas/tut1
	parts := strings.Split(cwd, "/src/")
	if len(parts) < 2 {
		log.Fatal("Unexpected directory structure")
	}

	// Search for *.aseprite files in the assets/aseprite/ directory
	files, err := filepath.Glob("assets/aseprite/*.aseprite")
	if err != nil {
		log.Fatal(err)
	}

	asepriteImportPath := "github.com/retroblast-engine/asevre"

	// Generate the variable declarations and function calls for each file
	var declarations, updates, draws string

	for _, file := range files {
		importPath := filepath.Join(parts[1], filepath.Base(file))
		importPaths = append(importPaths, importPath)

		base := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		Base := cases.Title(language.Und).String(base) // Capitalize the first letter
		declarations += genDeclarations(base, 64, 64)
		updates += fmt.Sprintf("my%s.Update()\n", base)
		draws += fmt.Sprintf("my%s.Draw(screen)\n", base)

		// Create directory for the aseprite file
		err := os.MkdirAll(base, 0755)
		if err != nil {
			log.Fatal(err)
		}

		aseprite, err := asevre.ParseAseprite(file)
		if err != nil {
			log.Fatal(err)
		}

		states := aseprite.State
		statesString := ""
		if len(states) < 1 {
			log.Fatal("No states found in the aseprite file")
		}

		statesString += "const (\n"
		for i := range states {
			state := &states[i]
			statesString += fmt.Sprintf("\t%s State = iota\n", state.Name)
		}
		statesString += ")\n"
		defautState := states[0]
		newFactoryString := fmt.Sprintf("New%s", Base)

		// Generate the player.go content
		playerGoContent := fmt.Sprintf(`package %s

import (
    "fmt"

    "%s"
    "github.com/hajimehoshi/ebiten/v2"
)

// %s represents the %s entity.
type %s struct {
	x, y      float64
	file      asevre.ASEFile
	sprite    *ebiten.Image
	state     State
	animation *Animation
	direction Direction
}

// %s creates a new %s instance.
func %s(x, y float64, file string) (*%s, error) {
	aseprite, err := asevre.ParseAseprite(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse aseprite file: %%w", err)
	}
		
	%s := &%s{
		x:     x,
		y:     y,
		file:  aseprite,
		state: %s,
		direction: Right,
	}
		
	%s.initAnimation(%s.state)
	
	return %s, nil
}
`, base, asepriteImportPath, Base, base, Base, newFactoryString, base, newFactoryString, Base, base, Base, defautState.Name, base, base, base)

		// Write the player.go file
		playerGoPath := filepath.Join(base, base+".go")
		err = os.WriteFile(playerGoPath, []byte(playerGoContent), 0644)
		if err != nil {
			log.Fatal(err)
		}

		// ================================= //
		// Generate the animation.go file    //
		// ================================= //

		l := strings.ToLower(base[:1]) // first later of the base name, used for pointer receiver
		animationGoContent := fmt.Sprintf(`package %s

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
	`, base, l, Base, l, l, l, l, l, l, Base, l, l, Base, l, l, l, l, l, l, l, l, l, l, l, base, l, Base, l, l, l, l, l, base, l, Base, l, l, l)

		// Write the player.go file
		playerAnimationGoPath := filepath.Join(base, "animation.go")
		err = os.WriteFile(playerAnimationGoPath, []byte(animationGoContent), 0644)
		if err != nil {
			log.Fatal(err)
		}

		// Generate the direction.go file
		directionGoContent := fmt.Sprintf(`package %s
// Direction represents %s's direction.
type Direction int

const (
	Right Direction = iota
	Left
	Up
	Down
)
`, base, base)

		// Write the direction.go file
		playerDirectionGoPath := filepath.Join(base, "direction.go")
		err = os.WriteFile(playerDirectionGoPath, []byte(directionGoContent), 0644)
		if err != nil {
			log.Fatal(err)
		}

		// ================================= //
		// Generate the draw.go file         //
		// ================================= //

		drawGoContent := fmt.Sprintf(`package %s

import "github.com/hajimehoshi/ebiten/v2"

// draw draws the %s's sprite on the screen.
func (%s *%s) draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	// Flip the sprite, if left is pressed
	if %s.direction == Left {
		%s.hFlip(opts)
	}

	opts.GeoM.Translate(%s.x, %s.y)

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
func (%s *%s) vFlip(opts *ebiten.DrawImageOptions) {
	opts.GeoM.Scale(1, -1)
	opts.GeoM.Translate(0, float64(%s.sprite.Bounds().Dy()))
}

// hvFlip modifies the DrawImageOptions to flip the sprite horizontally and vertically.
func (%s *%s) hvFlip(opts *ebiten.DrawImageOptions) {
	opts.GeoM.Scale(-1, -1)
	opts.GeoM.Translate(float64(%s.sprite.Bounds().Dx()), float64(%s.sprite.Bounds().Dy()))
}

// Draw draws the %s on the screen.
func (%s *%s) Draw(screen *ebiten.Image) {
	%s.draw(screen)
}

`, base, base, l, Base, l, l, l, l, l, l, l, l, Base, l, l, Base, l, l, Base, l, l, base, l, Base, l)

		// Write the draw.go file
		playerDrawGoPath := filepath.Join(base, "draw.go")
		err = os.WriteFile(playerDrawGoPath, []byte(drawGoContent), 0644)
		if err != nil {
			log.Fatal(err)
		}

		// ================================= //
		// Generate the state.go file        //
		// ================================= //

		stateGoContent := fmt.Sprintf(`package %s

// State represents the state of the %s.
type State int

%s`, base, base, statesString)

		// Write the state.go file
		playerStateGoPath := filepath.Join(base, "state.go")
		err = os.WriteFile(playerStateGoPath, []byte(stateGoContent), 0644)
		if err != nil {
			log.Fatal(err)
		}

		// ================================= //
		// Generate the update.go file       //
		// ================================= //

		updateGoContent := fmt.Sprintf(`package %s

// import (
// 	"github.com/hajimehoshi/ebiten/v2"
// )

// handleInput handles the %s's input.
func (%s *%s) handleInput() {
	// Modify this :)
}

// Update updates the %s's state based on input.
func (%s *%s) Update() {
	%s.handleInput()
}

`, base, base, l, Base, base, l, Base, l)

		// Write the update.go file
		playerUpdateGoPath := filepath.Join(base, "update.go")
		err = os.WriteFile(playerUpdateGoPath, []byte(updateGoContent), 0644)
		if err != nil {
			log.Fatal(err)
		}

	}

	// Generate the main.go content
	code := fmt.Sprintf(`package main

import (
    "image/color"
    "log"

    %s
    "github.com/hajimehoshi/ebiten/v2"
)

%s

type Game struct {
}

func (g *Game) Update() error {
    %s
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(color.NRGBA{0, 0, 0, 255})
    %s
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return 128, 128
}

func main() {

    ebiten.SetWindowSize(128*4, 128*4)
    ebiten.SetWindowTitle("Tutorial 1")

    if err := ebiten.RunGame(&Game{}); err != nil {
        log.Fatal(err)
    }

}
`, genImportPaths(importPaths), declarations, updates, draws)

	// Write the generated code to main.go
	err = os.WriteFile("main.go", []byte(code), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
