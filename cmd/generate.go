package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/retroblast-engine/asevre"
	maploader "github.com/retroblast-engine/tilesre"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TODO: Read the tileset and get the maxCellSize from there, instead of hardcoding it
var maxCellSize = 16 // HARDCODED VALUE !!!
var minCellSize = 2  // HARDCODED VALUE !!!

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

	// Get Git import paths using the Git configuration
	// ------------------------------------------------------------------------------------------------------------ //
	importPath, err := getRepoURL()
	if err != nil {
		log.Fatalf("Failed to get repository URL: %v", err)
	}

	log.Printf("Repository import path: %s", importPath)
	// ------------------------------------------------------------------------------------------------------------ //

	// Generate helper.go
	// ------------------------------------------------------------------------------------------------------------ //
	helperGoFile := "helper.go"
	helperGoContent := generateHelper()
	if err := os.WriteFile(helperGoFile, []byte(helperGoContent), 0644); err != nil {
		log.Fatal(err)
	}
	// ------------------------------------------------------------------------------------------------------------ //

	// Generate embed.go
	// ------------------------------------------------------------------------------------------------------------ //
	embedGoFile := "embed.go"
	embedGoContent := generateEmbedGo()
	if err := os.WriteFile(embedGoFile, []byte(embedGoContent), 0644); err != nil {
		log.Fatal(err)
	}
	// ------------------------------------------------------------------------------------------------------------ //

	// Generate game.go
	// ------------------------------------------------------------------------------------------------------------ //
	gameGoFile := "game.go"
	gameGoContent := generateGameGo(importPath)
	if err := os.WriteFile(gameGoFile, []byte(gameGoContent), 0644); err != nil {
		log.Fatal(err)
	}
	// ------------------------------------------------------------------------------------------------------------ //

	// Generate update.go, only if it doesn't already exist
	// ------------------------------------------------------------------------------------------------------------ //
	updateGoFile := "update.go"
	if _, err := os.Stat(updateGoFile); os.IsNotExist(err) {
		updateGoContent := generateUpdateGo()
		if err := os.WriteFile(updateGoFile, []byte(updateGoContent), 0644); err != nil {
			log.Fatal(err)
		}
	}
	// ------------------------------------------------------------------------------------------------------------ //

	// Generate draw.go (only if it doesn't already exist)
	// ------------------------------------------------------------------------------------------------------------ //
	drawGoFile := "draw.go"
	if _, err := os.Stat(drawGoFile); os.IsNotExist(err) {
		drawGoContent := generateDrawGo()
		if err := os.WriteFile(drawGoFile, []byte(drawGoContent), 0644); err != nil {
			log.Fatal(err)
		}
	}
	// ------------------------------------------------------------------------------------------------------------ //

	// Generate settings/settings.go file
	// TODO: Read the title from the input
	gameTitle := "My Game" // HARDCODED VALUE !!!
	// ------------------------------------------------------------------------------------------------------------ //
	// Create directory settings, or skip if it already exists
	settingsDir := "settings"
	if _, err := os.Stat(settingsDir); os.IsNotExist(err) {
		if err := os.Mkdir(settingsDir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	// Create settings.go file
	settingsGoFile := filepath.Join(settingsDir, "settings.go")

	// Search in the disk if assets/tiled directory exists, and exit if it doesn not
	tiledDir := "assets/tiled"
	asepriteDir := "assets/aseprite"
	if _, err := os.Stat(tiledDir); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist. Please create the directory or check the path.", tiledDir)
	}

	// Find all the *.tmx files in the tiledDir directory, and store their filenames in the slice tmxFiles
	var tmxFiles []string
	pattern := filepath.Join(tiledDir, "*.tmx")
	tmxFiles, err = filepath.Glob(pattern)

	if err != nil {
		log.Fatalf("Failed to find .tmx files in %s: %v", tiledDir, err)
	}

	if len(tmxFiles) == 0 {
		log.Fatalf("No .tmx files found in %s", tiledDir)
	}

	settingsGoContent := generateSettingsGo(minCellSize, maxCellSize, tmxFiles, gameTitle)
	if err := os.WriteFile(settingsGoFile, []byte(settingsGoContent), 0644); err != nil {
		log.Fatal(err)
	}
	// ------------------------------------------------------------------------------------------------------------ //

	// Generate main.go
	// ------------------------------------------------------------------------------------------------------------ //
	mainGoFile := "main.go"

	// Pick the first element from the tmxFiles slice, remove the .tmx extension
	// and nd ensure the string starts with a capital letter
	initialScene := ""
	if len(tmxFiles) > 0 {
		initialScene = strings.TrimSuffix(filepath.Base(tmxFiles[0]), ".tmx")
		initialScene = cases.Title(language.Und).String(initialScene)
		initialScene = "Scene" + initialScene
	}

	mainGoContent := generateMainGo(initialScene, importPath)
	if err := os.WriteFile(mainGoFile, []byte(mainGoContent), 0644); err != nil {
		log.Fatal(err)
	}
	// ------------------------------------------------------------------------------------------------------------ //

	// Generate camera/camera.go
	// ------------------------------------------------------------------------------------------------------------ //
	cameraDir := "camera"
	if _, err := os.Stat(cameraDir); os.IsNotExist(err) {
		if err := os.Mkdir(cameraDir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	cameraGoFile := filepath.Join(cameraDir, "camera.go")
	cameraGoContent := generateCameraGo(importPath)
	if err := os.WriteFile(cameraGoFile, []byte(cameraGoContent), 0644); err != nil {
		log.Fatal(err)
	}
	// ------------------------------------------------------------------------------------------------------------ //

	// TODO: Check if there are tiled assets and generate the scenes
	// For the time being assume there are already tiled assets

	// Create directory scenes, or skip if it already exists
	scenesDir := "scenes"
	if _, err := os.Stat(scenesDir); os.IsNotExist(err) {
		if err := os.Mkdir(scenesDir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	// Generate scenes/base.go file
	// ------------------------------------------------------------------------------------------------------------ //
	baseSceneFile := filepath.Join(scenesDir, "base.go")
	baseSceneContent := generateBaseScene(importPath)
	if err := os.WriteFile(baseSceneFile, []byte(baseSceneContent), 0644); err != nil {
		log.Fatal(err)
	}
	// ------------------------------------------------------------------------------------------------------------ //

	// Generate scenes/scene.go file
	// ------------------------------------------------------------------------------------------------------------ //
	scenesGoFile := filepath.Join(scenesDir, "scene.go")

	// Generate sceneIDNames and worldNames based on tmxFiles
	var sceneIDNames []string
	var sceneNames []string
	titleCaser := cases.Title(language.Und)

	for _, file := range tmxFiles {
		baseName := strings.TrimSuffix(titleCaser.String(filepath.Base(file)), ".tmx") // Capitalize the first letter, remove the extension
		pattern := "Scene" + titleCaser.String(baseName)                               // add Scene prefix
		sceneIDNames = append(sceneIDNames, pattern)
		sceneNames = append(sceneNames, baseName)
	}

	scenesGoContent := generateSceneGo(importPath, sceneIDNames, sceneNames)
	if err := os.WriteFile(scenesGoFile, []byte(scenesGoContent), 0644); err != nil {
		log.Fatal(err)
	}
	// ------------------------------------------------------------------------------------------------------------ //

	// Generate individual scene files, for every tmx file
	// ------------------------------------------------------------------------------------------------------------ //
	// This is temporary, until we have a way to pass some data in the aseprite files
	// signaling that this entity is destroyable.
	// Current approach, assumes all entities are destroyable, except the player!
	// TODO: Fix this properly
	// -- start of temporary code
	asepriteFiles, err := filepath.Glob("assets/aseprite/*.aseprite")
	if err != nil {
		log.Fatal(err)
	}

	var destroyableEntities []string
	for _, file := range asepriteFiles {
		// skip player.aseprite
		if strings.Contains(file, "player") {
			continue
		}
		destroyableEntities = append(destroyableEntities, strings.ToLower(strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))))
	}
	// -- end of temporary code

	allEntities, err := processTMXFiles(tmxFiles, tiledDir, asepriteDir)
	if err != nil {
		log.Fatalf("Error processing TMX files: %v", err)
	}
	// Print the entities for each tmxFile
	for i := range allEntities {
		fullPath := tmxFiles[i]
		justTheFile := filepath.Base(fullPath)
		justTheFileLower := strings.ToLower(justTheFile)
		makeTheFirstLetterCapital := titleCaser.String(justTheFileLower)
		WithoutSuffix := strings.TrimSuffix(makeTheFirstLetterCapital, ".tmx")
		tmxFiles[i] = WithoutSuffix
	}

	var characters []string
	for _, sceneName := range sceneNames {
		fmt.Printf("Entities for Scene %s:\n", sceneName)
		for i, entities := range allEntities {
			if sceneName == tmxFiles[i] {
				// Create this Scene file

				// The destroyableEntities has all the entities that are destroyable
				// but not all of those can be necessarily part of this particular scene
				// so we need to filter them out

				// Create list of type Entity for this scene
				var entitiesOFThisScene []Entity

				// These are all the entities in the scene
				for _, entity := range entities {
					fmt.Printf("Name: %s, X: %f, Y: %f\n", entity.Name, entity.X, entity.Y)

					// For every entity.Name there has to be an assets/aseprite/entity.Name.aseprite file
					// If there is no such file, then we skip this entity and print a warning message to the user
					asepriteFile := filepath.Join(asepriteDir, strings.ToLower(entity.Name)+".aseprite")
					if _, err := os.Stat(asepriteFile); os.IsNotExist(err) {
						log.Printf("Warning: No aseprite file found for entity %s. Skipping this entity.", entity.Name)
						continue
					}

					entitiesOFThisScene = append(entitiesOFThisScene, entity)
				}

				// From 'entitiesOFThisScene' and 'destroyableEntities', filter out the entities that are destroyable
				// and create a new []string with the names of the destroyable entities which has to be unique names
				// and lowercase
				var destroyableEntitiesOFThisScene []string
				for _, entity := range entitiesOFThisScene {
					for _, destroyableEntity := range destroyableEntities {
						if entity.Name == destroyableEntity {
							destroyableEntitiesOFThisScene = append(destroyableEntitiesOFThisScene, destroyableEntity)
						}
					}
				}

				// If destroyableEntitiesOFThisScene has multiple instances of the same entity (meaning, same Name)
				// then we need to remove the duplicates
				destroyableEntitiesOFThisScene = unique(destroyableEntitiesOFThisScene)
				characters = append(characters, destroyableEntitiesOFThisScene...)

				fmt.Printf("Destroyable entities for Scene %s:\n", sceneName)
				fmt.Println(destroyableEntitiesOFThisScene)

				sceneFile := filepath.Join(scenesDir, strings.ToLower(sceneName)+".go")
				settingsScene := fmt.Sprintf("Scene%d", i+1)

				sceneContent := generateScene(importPath, sceneName, settingsScene, "player", destroyableEntitiesOFThisScene)
				if err :=
					os.WriteFile(sceneFile, []byte(sceneContent), 0644); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	// ------------------------------------------------------------------------------------------------------------ //
	// Generate eneities/entity.go
	// ------------------------------------------------------------------------------------------------------------ //
	entitiesDir := "entities"
	if _, err := os.Stat(entitiesDir); os.IsNotExist(err) {
		if err := os.Mkdir(entitiesDir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	entitiesGoFile := filepath.Join(entitiesDir, "entity.go")

	// Add the player to the characters list
	characters = append(characters, "player")
	uniqueCharacters := unique(characters)

	entitiesGoContent := generateEntityGo(importPath, uniqueCharacters)
	if err := os.WriteFile(entitiesGoFile, []byte(entitiesGoContent), 0644); err != nil {
		log.Fatal(err)
	}

	// ------------------------------------------------------------------------------------------------------------ //
	// Generate individual entity packages inside the entity directory, for every entity, basically for every aseprite file
	// ------------------------------------------------------------------------------------------------------------ //

	// if err := generateAseprite(); err != nil {
	// 	return err
	// }

	// Loop through all the uniqueCharacters
	for _, character := range uniqueCharacters {
		// Create directory for the aseprite file if it doesn't already exist, otherwise skip
		characterDir := filepath.Join(entitiesDir, character)
		if _, err := os.Stat(characterDir); os.IsNotExist(err) {
			if err := os.Mkdir(characterDir, 0755); err != nil {
				log.Fatal(err)
			}
		}

		// Create the characterDir/animation.go file
		animationGoFile := filepath.Join(characterDir, "animation.go")
		animationGoContent := generateAnimationGo(character)
		if err := os.WriteFile(animationGoFile, []byte(animationGoContent), 0644); err != nil {
			log.Fatal(err)
		}

		// Create the characterDir/direction.go file
		directionGoFile := filepath.Join(characterDir, "direction.go")
		directionGoContent := generateDirectionGo(character)
		if err := os.WriteFile(directionGoFile, []byte(directionGoContent), 0644); err != nil {
			log.Fatal(err)
		}

		// Create the characterDir/draw.go file
		drawGoFile := filepath.Join(characterDir, "draw.go")
		drawGoContent := generateEntityDrawGo(character)
		if err := os.WriteFile(drawGoFile, []byte(drawGoContent), 0644); err != nil {
			log.Fatal(err)
		}

		// Create the characterDir/state.go file
		stateGoFile := filepath.Join(characterDir, "state.go")

		// Loop through all the aseprite files
		var defautState asevre.ASETag

		// Construct the file path for the aseprite file
		file := filepath.Join("assets/aseprite", character+".aseprite")

		// Make sure it exists
		if _, err := os.Stat(file); os.IsNotExist(err) {
			// Try if the file is Tiled file (so first letter is capital)
			c := cases.Title(language.Und)
			file = filepath.Join("assets/aseprite", c.String(character)+".aseprite")
			if _, err := os.Stat(file); os.IsNotExist(err) {
				// Try if the file is all lowercase
				file = filepath.Join("assets/aseprite", strings.ToLower(character)+".aseprite")
				if _, err := os.Stat(file); os.IsNotExist(err) {
					log.Fatalf("No aseprite file found for entity %s", character)
				}
			}
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

		c := cases.Title(language.Und)
		statesString += "const (\n"
		for i, state := range states {
			if i == 0 {
				statesString += fmt.Sprintf("\t%s State = iota\n", c.String(state.Name))
			} else {
				statesString += fmt.Sprintf("\t%s\n", c.String(state.Name))
			}
		}
		statesString += ")\n"
		defautState = states[0]

		stateGoContent := generateStateGo(character, statesString)
		if err := os.WriteFile(stateGoFile, []byte(stateGoContent), 0644); err != nil {
			log.Fatal(err)
		}

		// Create the characterDir/player.go file
		playerGoFile := filepath.Join(characterDir, character+".go")
		var entityGoContent string
		if character == "player" {
			entityGoContent = generatePlayerGo(importPath, character, c.String(defautState.Name))
		} else {
			entityGoContent = generateCharacterGo(importPath, character, c.String(defautState.Name))
		}
		if err := os.WriteFile(playerGoFile, []byte(entityGoContent), 0644); err != nil {
			log.Fatal(err)
		}

		// Create the characterDir/update.go file
		updateGoFile := filepath.Join(characterDir, "update.go")
		var updateGoContent string
		if character == "player" {
			updateGoContent = generatePlayerUpdateGo(importPath, character)
		} else {
			updateGoContent = generateCharacterUpdateGo(importPath, character)
		}
		if err := os.WriteFile(updateGoFile, []byte(updateGoContent), 0644); err != nil {
			log.Fatal(err)
		}
	}

	// ------------------------------------------------------------------------------------------------------------ //

	if err := initializeGoModule(); err != nil {
		return err
	}

	if err := tidyGoModule(); err != nil {
		return err
	}

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

func unique(list []string) []string {
	keys := make(map[string]bool)
	var listUnique []string
	for _, entry := range list {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			listUnique = append(listUnique, entry)
		}
	}

	return listUnique
}

var importPaths []string

// generateAseprite generates the code for the game based on the Aseprite files
func generateAseprite() error {
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

	return nil
}

// Entity represents an entity with a name and coordinates
type Entity struct {
	Name string
	X    float64
	Y    float64
}

// readDirAsFS reads a directory and returns an fs.FS and a cleanup function
func readDirAsFS(dirPath string) (fs.FS, func(), error) {
	tempFS, err := fs.Sub(os.DirFS(dirPath), ".")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create fs.FS: %w", err)
	}

	// Cleanup function (no-op in this case)
	cleanup := func() {}

	return tempFS, cleanup, nil
}

// processTMXFiles processes all the tmxFiles and extracts entities
func processTMXFiles(tmxFiles []string, tiledPath, asepritePath string) ([][]Entity, error) {
	var allEntities [][]Entity

	// Create a tmpFolder to copy the assets folder
	tmpFolder := "tmp"

	// Ensure the tmp folder is recreated
	if _, err := os.Stat(tmpFolder); !os.IsNotExist(err) {
		if err := os.RemoveAll(tmpFolder); err != nil {
			log.Fatalf("Error removing existing directory: %v", err)
		}
	}

	if err := os.Mkdir(tmpFolder, 0755); err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	// Create the assets folder if it doesn't exist, inside the tmp folder
	if _, err := os.Stat(filepath.Join(tmpFolder, "assets")); os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Join(tmpFolder, "assets"), 0755); err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
	}

	// Copy the assets folder to the tmp folder
	// So the expected file structure would be: tmp/assets/tiled/*.tmx, tmp/assets/aseprite/*.aseprite
	// this is what tilesre.Load() expects as input for 'assets fs.FS'
	if err := filepath.Walk("assets", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the assets folder itself
		if path == "assets" {
			return nil
		}

		// Get the relative path

		relPath, err := filepath.Rel("assets", path)
		if err != nil {
			return err
		}

		// Create the destination path
		dstPath := filepath.Join(tmpFolder, "assets", relPath)

		// If it's a directory, create it

		if info.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err

			}
			return nil
		}

		// If it's a file, copy it
		if err := copyFile(path, dstPath, info); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("Error copying directory: %v", err)
	}

	fs, cleanup, err := readDirAsFS(tmpFolder)
	if err != nil {
		log.Fatalf("Error reading directory as fs.FS: %v", err)
	}
	defer cleanup()

	for _, tmxFile := range tmxFiles {
		var entities []Entity

		// Load the tiled map using the temporary fs.FS
		tiledMap, err := maploader.Load(fs, tmxFile, tiledPath, asepritePath, minCellSize, minCellSize)
		if err != nil {
			log.Fatalf("Failed to load map: %v", err)
		}

		// Process object layers to extract entities
		for _, objectLayer := range tiledMap.TiledMap.ObjectGroups {
			for _, object := range objectLayer.Objects {
				entity := Entity{
					Name: object.Name,
					X:    object.X,
					Y:    object.Y,
				}
				entities = append(entities, entity)
			}
		}

		// Add the entities for this tmxFile to the allEntities slice
		allEntities = append(allEntities, entities)
	}

	// delete the tmp folder
	if err := os.RemoveAll(tmpFolder); err != nil {
		log.Fatalf("Error removing directory: %v", err)
	}

	return allEntities, nil
}

// copyFile copies a file from src to dst, preserving its permissions.
func copyFile(src string, dst string, info os.FileInfo) error {
	// Open the source file for reading.
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %v", src, err)
	}
	defer srcFile.Close()

	// Create the destination file for writing.
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %v", dst, err)
	}
	defer dstFile.Close()

	// Copy the contents of the source file to the destination file.
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents from %s to %s: %v", src, dst, err)
	}

	// Preserve the file mode (permissions).
	err = os.Chmod(dst, info.Mode())
	if err != nil {
		return fmt.Errorf("failed to set file mode for %s: %v", dst, err)
	}

	return nil
}
