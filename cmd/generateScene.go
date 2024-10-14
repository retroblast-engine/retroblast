package cmd

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func generateSceneImports(removableEntities []string, importPath string) string {
	var imports string

	if len(removableEntities) > 0 {
		if len(removableEntities) == 1 {
			imports += fmt.Sprintf("\"%s/entities/%s\"", importPath, removableEntities[0])
		} else {
			for _, entity := range removableEntities {
				imports += fmt.Sprintf("\"%s/entities/%s\"\n", importPath, entity)
			}
		}
	}

	return imports
}

// cleanAndSortImports takes a string of import statements, removes empty lines, and sorts them alphabetically
func cleanAndSortImports(imports string) string {
	lines := strings.Split(imports, "\n")

	// Remove empty lines
	var nonEmptyLines []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			nonEmptyLines = append(nonEmptyLines, trimmedLine)
		}
	}

	// Sort the lines alphabetically, ignoring the first and last lines (import block)
	if len(nonEmptyLines) > 2 {
		sort.Strings(nonEmptyLines[1 : len(nonEmptyLines)-1])
	}

	// Add indentation to each line except the first and last
	for i := 1; i < len(nonEmptyLines)-1; i++ {
		nonEmptyLines[i] = "\t" + nonEmptyLines[i]
	}

	// Join the sorted lines back into a single string
	return strings.Join(nonEmptyLines, "\n")
}

func generateEntitiesPathOrNot(importPath string, removableEntities []string) string {
	fmt.Println("IMPORT PATH", importPath)
	if len(removableEntities) > 0 {
		// if the importPath has a slash at the end, remove it
		if importPath[len(importPath)-1] == '/' {
			importPath = importPath[:len(importPath)-1]
		}
		return fmt.Sprintf("\"%s/entities\"", importPath)
	}

	return ""
}

func generateScene(importPath, structName, settingsScene, player string, removableEntities []string) string {
	// Generate the import statements
	imports := fmt.Sprintf(`import (
    %s
    %s
    "%s/entities/%s"
    "%s/settings"
    "github.com/hajimehoshi/ebiten/v2"
)`, generateEntitiesPathOrNot(importPath, removableEntities), generateSceneImports(removableEntities, importPath), importPath, player, importPath)

	imports = cleanAndSortImports(imports)

	// Generate the removable entities code
	titleCaser := cases.Title(language.Und)
	var removableEntitiesCode strings.Builder
	for _, entity := range removableEntities {

		removableEntitiesCode.WriteString(fmt.Sprintf(`
        if e, ok := entity.(*%s.%s); ok && e.IsDestroyed {
            s.tiledMap.Space.Remove(e.Object)
            s.entities = entities.RemoveEntity(s.entities, e)
        }`, entity, titleCaser.String(entity)))
	}

	// Generate the player-specific code
	playerCode := fmt.Sprintf(`
        if p, ok := entity.(*%s.%s); ok {
            s.camera.FollowTargetX(p.GetX() + settings.MaxCellSize/2)

            // Calculate the width of the tile map in pixels
            tileMapWidthPixels := s.tiledMap.TiledMap.Width * s.tiledMap.TiledMap.TileWidth

            // Calculate the height of the tile map in pixels
            tileMapHeightPixels := s.tiledMap.TiledMap.Height * s.tiledMap.TiledMap.TileHeight

            // Constrain the camera to the tile map dimensions
            s.camera.Constrain(tileMapWidthPixels, tileMapHeightPixels)
        }`, player, titleCaser.String(player))

	return fmt.Sprintf(`package scenes

%s

// %s represents the %s scene
type %s struct {
    baseScene
}

// Init initializes the %s scene
func (s *%s) Init() {
    s.baseScene.Init(settings.Assets, settings.%s, settings.AssetsTiledPath, settings.AssetsAsepritePath)
    s.sceneID = Scene%s
}

// Update updates the %s scene
func (s *%s) Update() error {
    for _, entity := range s.entities {
        entity.Update()
        %s
        %s
    }

    return nil
}

// Draw draws the %s scene
func (s *%s) Draw(screen *ebiten.Image) {
    s.baseScene.Draw(screen)
}
`, imports, structName, structName, structName, structName, structName, settingsScene, structName, structName, structName, removableEntitiesCode.String(), playerCode, structName, structName)
}
