package cmd

import (
	"fmt"
)

func generateBaseScene(baseImportPath string) string {
	return fmt.Sprintf(`package scenes

import (
    "embed"
    "fmt"
    "image"
    "image/color"
    "log"

    "%s/camera"
    "%s/entities"
    "%s/settings"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/lafriks/go-tiled"
    maploader "github.com/retroblast-engine/tilesre"
    "github.com/solarlune/resolv"
)

// baseScene holds the common state and methods for all scenes
type baseScene struct {
    tiledMap *maploader.Map
    entities []entities.Entity
    camera   *camera.Camera
    sceneID  SceneID
}

// Init initializes the scene
func (s *baseScene) Init(assets embed.FS, mapFile, tiledPath, asepritePath string) {
    var err error

    s.camera = camera.New(0, 0)

    s.tiledMap, err = maploader.Load(assets, mapFile, tiledPath, asepritePath, settings.MinCellSize, settings.MinCellSize)
    if err != nil {
        log.Fatalf("Failed to load map: %%v", err)
    }

    s.processMapLayers()
    s.processObjectLayers(asepritePath)
    s.loadImageLayers()

    for _, obj := range s.tiledMap.Objects {
        s.tiledMap.Space.Add(obj.Physics)
    }

    // Make all interactive objects also solid ones
	for _, o := range s.tiledMap.Space.Objects() {
		if o.HasTags("interactive") {
			o.AddTags("solid")
		}
	}
}

// Update updates the scene
func (s *baseScene) Update() error {
    return nil
}

// Draw draws the scene
func (s *baseScene) Draw(screen *ebiten.Image) {
    screen.Fill(color.NRGBA{0, 0, 0, 255})
    s.tiledMap.MapDraw(screen, s.camera.X, s.camera.Y)
    for _, entity := range s.entities {
        entity.Draw(screen, s.camera.X, s.camera.Y)
    }
}

// Destroy cleans up resources
func (s *baseScene) Destroy() {
    s.tiledMap = nil
    s.entities = nil
}

// GetCamera returns the camera
func (s *baseScene) GetCamera() *camera.Camera {
    return s.camera
}

// GetSpace returns the physics space
func (s *baseScene) GetSpace() *resolv.Space {
    return s.tiledMap.Space
}

// processMapLayers processes the map layers
func (s *baseScene) processMapLayers() {
    for i, layer := range s.tiledMap.TiledMap.Layers {
        fmt.Println("Processing layer", layer.Name)
        if err := s.tiledMap.ProcessLayer(i, layer); err != nil {
            log.Fatalf("Failed to process layer: %%v", err)
        }
    }
}

// processObjectLayers processes the object layers
func (s *baseScene) processObjectLayers(asepritePath string) {
    for _, objectLayer := range s.tiledMap.TiledMap.ObjectGroups {
        fmt.Println("Processing object layer", objectLayer.Name)
        s.processObjectLayer(objectLayer, asepritePath, s.tiledMap)
    }
}

// loadImageLayers loads the image layers
func (s *baseScene) loadImageLayers() {
    for _, layer := range s.tiledMap.TiledMap.ImageLayers {
        fmt.Println("Processing layer", layer.Name)
        s.tiledMap.BackgroundImageOffsetX = float64(layer.OffsetX)
        s.tiledMap.BackgroundImageOffsetY = float64(layer.OffsetY)
        img, err := createEbitenImageFromSource(settings.Assets, settings.AssetsTiledPath+layer.Image.Source)
        if err != nil {
            log.Fatalf("Error creating Ebiten image: %%v", err)
        }
        s.tiledMap.BackgroundImage = img
    }
}

// createEbitenImageFromSource creates an Ebiten image from a source file
func createEbitenImageFromSource(assets embed.FS, source string) (*ebiten.Image, error) {
    file, err := assets.Open(source)
    if err != nil {
        return nil, fmt.Errorf("failed to open image file: %%w", err)
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        return nil, fmt.Errorf("failed to decode image file: %%w", err)
    }

    return ebiten.NewImageFromImage(img), nil
}

// processObjectLayer processes an object layer
func (s *baseScene) processObjectLayer(objectLayer *tiled.ObjectGroup, asepritePath string, m *maploader.Map) {
    for _, object := range objectLayer.Objects {
        x, y, name := object.X, object.Y, object.Name

        constructor, exists := entities.EntityConstructors[name]
        if !exists {
            log.Printf("No constructor found for entity: %%s\n", name)
            continue
        }

        entity, err := constructor(x, y, asepritePath+name+".aseprite")
        if err != nil {
            log.Fatalf("Failed to create entity: %%v", err)
        }
        s.entities = append(s.entities, entity)
        m.Space.Add(entity.GetObject())
    }
}

// GetSceneID returns the scene ID
func (s *baseScene) GetSceneID() SceneID {
    return s.sceneID
}
`, baseImportPath, baseImportPath, baseImportPath)
}
