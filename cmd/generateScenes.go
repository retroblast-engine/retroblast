package cmd

import (
	"fmt"
	"strings"
)

func generateSceneGo(baseImportPath string, sceneNames []string, worldNames []string) string {
	// Generate the scene map entries
	var sceneMapEntries strings.Builder
	for i, sceneName := range sceneNames {
		sceneMapEntries.WriteString(fmt.Sprintf("%s: &%s{},\n", sceneName, worldNames[i]))
	}

	return fmt.Sprintf(`package scenes

import (
    "%s/camera"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/solarlune/resolv"
)

// Scene defines the methods that a scene should have.
type Scene interface {
    Init()
    Update() error
    Draw(screen *ebiten.Image)
    Destroy()
    GetCamera() *camera.Camera
    GetSpace() *resolv.Space
    GetSceneID() SceneID
}

// SceneID represents the identifier of a scene
type SceneID int

const (
    %s
)

// SceneMap maps scene identifiers to scene instances
var SceneMap = map[SceneID]Scene{
    %s
}
`, baseImportPath, generateSceneIDConstants(sceneNames), sceneMapEntries.String())
}

func generateSceneIDConstants(sceneNames []string) string {
	var constants strings.Builder
	for i, sceneName := range sceneNames {
		constants.WriteString(fmt.Sprintf("%s = iota + %d\n", sceneName, i+1))
	}
	return constants.String()
}
