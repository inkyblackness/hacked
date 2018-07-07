package levels

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/inkyblackness/hacked/editor/event"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ui/input"
	"github.com/inkyblackness/hacked/ui/opengl"
	"github.com/inkyblackness/imgui-go"
)

// MapDisplay renders a level map.
type MapDisplay struct {
	context  render.Context
	camera   *LimitedCamera
	guiScale float32

	eventListener event.Listener

	background  *BackgroundGrid
	textures    *MapTextures
	colors      *MapColors
	mapGrid     *MapGrid
	highlighter *Highlighter

	moveCapture func(pixelX, pixelY float32)
	mouseMoved  bool

	positionPopupPos imgui.Vec2
	positionValid    bool
	positionX        level.Coordinate
	positionY        level.Coordinate

	selectedTiles tileCoordinates
}

// NewMapDisplay returns a new instance.
func NewMapDisplay(gl opengl.OpenGL, guiScale float32,
	textureQuery TextureQuery,
	eventListener event.Listener, eventRegistry event.Registry) *MapDisplay {
	tilesPerMapSide := float32(64)

	tileBaseLength := fineCoordinatesPerTileSide
	tileBaseHalf := tileBaseLength / 2.0
	camLimit := tilesPerMapSide*tileBaseLength - tileBaseHalf
	zoomShift := guiScale - 1.0
	zoomLevelMin := float32(-5) + zoomShift
	zoomLevelMax := float32(1) + zoomShift

	display := &MapDisplay{
		context: render.Context{
			OpenGL:           gl,
			ProjectionMatrix: mgl.Ident4(),
		},
		camera:        NewLimitedCamera(zoomLevelMin, zoomLevelMax, -tileBaseHalf, camLimit),
		guiScale:      guiScale,
		eventListener: eventListener,
		moveCapture:   func(float32, float32) {},
	}
	display.context.ViewMatrix = display.camera.ViewMatrix()
	display.background = NewBackgroundGrid(&display.context)
	display.textures = NewMapTextures(&display.context, textureQuery)
	display.colors = NewMapColors(&display.context)
	display.mapGrid = NewMapGrid(&display.context)
	display.highlighter = NewHighlighter(&display.context)

	centerX, centerY := (tilesPerMapSide*tileBaseLength)/-2.0, (tilesPerMapSide*tileBaseLength)/-2.0
	display.camera.ZoomAt(-3+zoomShift, centerX, centerY)
	display.camera.MoveTo(centerX, centerY)

	display.selectedTiles.registerAt(eventRegistry)

	return display
}

// Render renders the whole map display.
func (display *MapDisplay) Render(lvl *level.Level, paletteTexture *graphics.PaletteTexture,
	textureDisplay TextureDisplay, colorDisplay ColorDisplay) {
	columns, rows, _ := lvl.Size()

	display.background.Render()
	if lvl.IsCyberspace() {
		if paletteTexture != nil {
			var colorQuery ColorQuery
			palette := paletteTexture.Palette()
			if colorDisplay == ColorDisplayFloor {
				colorQuery = display.colorQueryFor(lvl, func(tile *level.TileMapEntry) [4]float32 {
					rgb := palette[tile.TextureInfo.FloorPaletteIndex()]
					return [4]float32{float32(rgb.Red) / 255, float32(rgb.Green) / 255, float32(rgb.Blue) / 255, 0.8}
				})
			} else if colorDisplay == ColorDisplayCeiling {
				colorQuery = display.colorQueryFor(lvl, func(tile *level.TileMapEntry) [4]float32 {
					rgb := palette[tile.TextureInfo.CeilingPaletteIndex()]
					return [4]float32{float32(rgb.Red) / 255, float32(rgb.Green) / 255, float32(rgb.Blue) / 255, 0.8}
				})
			}
			if colorQuery != nil {
				display.colors.Render(columns, rows, colorQuery)
			}
		}
	} else {
		if paletteTexture != nil {
			display.textures.Render(columns, rows, func(x, y int) (level.TileType, int, int) {
				tile := lvl.Tile(x, y)
				if tile == nil {
					return level.TileTypeSolid, 0, 0
				}
				atlasIndex, textureRotations := textureDisplay.Func()(tile)
				atlas := lvl.TextureAtlas()
				textureIndex := -1
				if (atlasIndex >= 0) && (atlasIndex < len(atlas)) {
					textureIndex = int(atlas[atlasIndex])
				}
				return tile.Type, textureIndex, textureRotations
			}, paletteTexture)
		}

		var colorQuery ColorQuery
		if colorDisplay == ColorDisplayFloor {
			colorQuery = display.colorQueryFor(lvl, func(tile *level.TileMapEntry) [4]float32 {
				return [4]float32{0.0, 0.0, 0.0, float32(tile.Flags.ForRealWorld().FloorShadow()) / 15.0}
			})
		} else if colorDisplay == ColorDisplayCeiling {
			colorQuery = display.colorQueryFor(lvl, func(tile *level.TileMapEntry) [4]float32 {
				return [4]float32{0.0, 0.0, 0.0, float32(tile.Flags.ForRealWorld().CeilingShadow()) / 15.0}
			})
		}
		if colorQuery != nil {
			display.colors.Render(columns, rows, colorQuery)
		}
	}
	display.mapGrid.Render(lvl)
	display.highlighter.Render(display.selectedTiles.list, fineCoordinatesPerTileSide, [4]float32{0.0, 0.8, 0.2, 0.5})
	if display.positionValid {
		tilePos := MapPosition{
			X: level.CoordinateAt(display.positionX.Tile(), 128),
			Y: level.CoordinateAt(display.positionY.Tile(), 128),
		}
		display.highlighter.Render([]MapPosition{tilePos}, fineCoordinatesPerTileSide, [4]float32{0.0, 0.2, 0.8, 0.3})
	}

	display.renderPositionOverlay(lvl)
}

func (display *MapDisplay) colorQueryFor(lvl *level.Level, tileToColor func(*level.TileMapEntry) [4]float32) func(int, int) [4]float32 {
	return func(x, y int) [4]float32 {
		tile := lvl.Tile(x, y)
		if tile == nil {
			return [4]float32{}
		}
		return tileToColor(tile)
	}
}

func (display *MapDisplay) renderPositionOverlay(lvl *level.Level) {
	imgui.SetNextWindowPosV(display.positionPopupPos, imgui.ConditionAlways, imgui.Vec2{X: 1.0, Y: 1.0})
	imgui.SetNextWindowSize(imgui.Vec2{X: 120 * display.guiScale, Y: 0})
	imgui.SetNextWindowBgAlpha(0.3)
	if imgui.BeginV("Position", nil, imgui.WindowFlagsNoMove|imgui.WindowFlagsNoTitleBar|imgui.WindowFlagsNoResize|imgui.WindowFlagsAlwaysAutoResize|
		imgui.WindowFlagsNoSavedSettings|imgui.WindowFlagsNoFocusOnAppearing|imgui.WindowFlagsNoNav) {

		if display.positionValid {
			tile := lvl.Tile(int(display.positionX.Tile()), int(display.positionY.Tile()))

			imgui.Text(fmt.Sprintf("X: T %2d F %3d", display.positionX.Tile(), display.positionX.Fine()))
			imgui.Text(fmt.Sprintf("Y: T %2d F %3d", display.positionY.Tile(), display.positionY.Fine()))
			if (tile != nil) && (tile.Type != level.TileTypeSolid) {
				_, _, heightShift := lvl.Size()
				height := tile.Floor.AbsoluteHeight()
				imgui.Text(fmt.Sprintf("Z: %2d = %2.3f", height, heightShift.ValueFromTileHeight(height)))
			} else {
				imgui.Text("Z: -- = --.---")
			}
		} else {
			imgui.Text("X: T -- F ---")
			imgui.Text("Y: T -- F ---")
			imgui.Text("Z: -- = --.---")
		}
		imgui.End()
	}

}

// WindowResized must be called to notify of a change in window geometry.
func (display *MapDisplay) WindowResized(width int, height int) {
	display.context.ProjectionMatrix = mgl.Ortho2D(0.0, float32(width), float32(height), 0.0)
	display.camera.SetViewportSize(float32(width), float32(height))
	display.positionPopupPos.X = float32(width) - 10.0
	display.positionPopupPos.Y = float32(height) - 10.0
}

func (display *MapDisplay) unprojectPixel(pixelX, pixelY float32) (x, y float32) {
	pixelVec := mgl.Vec4{pixelX, pixelY, 0.0, 1.0}
	invertedView := display.camera.ViewMatrix().Inv()
	result := invertedView.Mul4x1(pixelVec)

	return result[0], result[1]
}

// MouseButtonDown must be called when a button was pressed.
func (display *MapDisplay) MouseButtonDown(mouseX, mouseY float32, button uint32) {
	display.updateMouseWorldPosition(mouseX, mouseY)
	if button == input.MousePrimary {
		lastPixelX, lastPixelY := mouseX, mouseY

		display.mouseMoved = false
		display.moveCapture = func(pixelX, pixelY float32) {
			lastWorldX, lastWorldY := display.unprojectPixel(lastPixelX, lastPixelY)
			worldX, worldY := display.unprojectPixel(pixelX, pixelY)

			display.camera.MoveBy(worldX-lastWorldX, worldY-lastWorldY)
			lastPixelX, lastPixelY = pixelX, pixelY
			display.mouseMoved = true
		}
	}
}

// MouseButtonUp must be called when a button was released.
func (display *MapDisplay) MouseButtonUp(mouseX, mouseY float32, button uint32, modifier input.Modifier) {
	display.updateMouseWorldPosition(mouseX, mouseY)
	if button == input.MousePrimary {
		display.moveCapture = func(float32, float32) {}
		if !display.mouseMoved && display.positionValid {
			tilePos := MapPosition{
				X: level.CoordinateAt(display.positionX.Tile(), 128),
				Y: level.CoordinateAt(display.positionY.Tile(), 128),
			}
			if modifier.Has(input.ModControl) {
				wasSelected := display.selectedTiles.contains(tilePos)
				if wasSelected {
					display.eventListener.Event(TileSelectionRemoveEvent{tiles: []MapPosition{tilePos}})
				} else {
					display.eventListener.Event(TileSelectionAddEvent{tiles: []MapPosition{tilePos}})
				}
			} else if modifier.Has(input.ModShift) && (len(display.selectedTiles.list) > 0) {
				firstPos := display.selectedTiles.list[0]

				fromX := int(firstPos.X.Tile())
				fromY := int(firstPos.Y.Tile())
				toX := int(tilePos.X.Tile())
				toY := int(tilePos.Y.Tile())
				xIncrement := 1
				yIncrement := 1
				if fromX > toX {
					xIncrement = -1
				}
				if fromY > toY {
					yIncrement = -1
				}
				toX += xIncrement
				toY += yIncrement
				var newList []MapPosition
				for y := fromY; y != toY; y += yIncrement {
					for x := fromX; x != toX; x += xIncrement {
						newList = append(newList, MapPosition{X: level.CoordinateAt(byte(x), 128), Y: level.CoordinateAt(byte(y), 128)})
					}
				}
				display.eventListener.Event(TileSelectionSetEvent{tiles: newList})
			} else {
				display.eventListener.Event(TileSelectionSetEvent{tiles: []MapPosition{tilePos}})
			}
		}
	}
}

// MouseMoved must be called for a mouse move.
func (display *MapDisplay) MouseMoved(mouseX, mouseY float32) {
	display.updateMouseWorldPosition(mouseX, mouseY)
	display.moveCapture(mouseX, mouseY)
}

// MouseScrolled must be called for a mouse scroll
func (display *MapDisplay) MouseScrolled(mouseX, mouseY float32, deltaX, deltaY float32) {
	worldX, worldY := display.unprojectPixel(mouseX, mouseY)

	if deltaY < 0 {
		display.camera.ZoomAt(-0.5, worldX, worldY)
	}
	if deltaY > 0 {
		display.camera.ZoomAt(0.5, worldX, worldY)
	}
}

func (display *MapDisplay) updateMouseWorldPosition(mouseX, mouseY float32) {
	worldX, worldY := display.unprojectPixel(mouseX, mouseY)

	display.positionValid = (worldX >= 0.0) && (worldX < (64.0 * fineCoordinatesPerTileSide)) &&
		(worldY >= 0.0) && (worldY < (64.0 * fineCoordinatesPerTileSide))
	if display.positionValid {
		display.positionX = level.Coordinate(worldX + 0.5)
		display.positionY = level.Coordinate(worldY + 0.5)
	}
}
