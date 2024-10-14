package cmd

import (
	"fmt"
	"strings"
)

func generateSettingsGo(MinCellSize, maxCellSize int, sceneFiles []string, gameTitle string) string {
	// Generate the scene variables
	var sceneVars strings.Builder
	for i, filepath := range sceneFiles {
		// from the filepath, keep only the last file part
		// e.g. "assets/tiled/scene1.json" -> "scene1.json"
		file := filepath[strings.LastIndex(filepath, "/")+1:]
		sceneVars.WriteString(fmt.Sprintf("var Scene%d string = AssetsTiledPath + \"%s\"\n", i+1, file))
	}

	return fmt.Sprintf(`package settings

import (
    "embed"

    "github.com/hajimehoshi/ebiten/v2"
)

const (
    ScreenWidth  = 1280
    ScreenHeight = 800
    Scale        = 3
    ScaledWidth  = ScreenWidth / Scale
    ScaledHeight = ScreenHeight / Scale
    WindowTitle  = "%s"

    // Sometimes we might want precise collision detection, so we can use a smaller cell size as well
    // however, this will increase the number of cells and the number of checks
    // so we need to find a balance between speed and precision
    MaxCellSize = %d // should be your sprite size in Aseprite
    MinCellSize = %d  // (must be a power of 2) if unsure, use the same as MaxCellSize
)

var (
    BackgroundImageOffsetX, BackgroundImageOffsetY float64
    BackgroundImage                                *ebiten.Image
    Assets                                         embed.FS
)

// SetEmbeddedFiles sets the embedded files
func SetEmbeddedFiles(fs embed.FS) {
    Assets = fs
}

var AssetsTiledPath string = "assets/tiled/"
var AssetsAsepritePath string = "assets/aseprite/"
`, gameTitle, maxCellSize, MinCellSize) + sceneVars.String()
}
