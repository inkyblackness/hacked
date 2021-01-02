package levels

import (
	"fmt"
	"sort"

	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/editor/event"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/input"
	"github.com/inkyblackness/hacked/ui/opengl"
)

type hoverItem interface {
	Pos() MapPosition
	Size() float32
}

type tileHoverItem struct {
	pos MapPosition
}

func (item tileHoverItem) Pos() MapPosition {
	return item.pos
}

func (item tileHoverItem) Size() float32 {
	return fineCoordinatesPerTileSide
}

type objectHoverItem struct {
	id  level.ObjectID
	pos MapPosition
}

func (item objectHoverItem) Pos() MapPosition {
	return item.pos
}

func (item objectHoverItem) Size() float32 {
	return fineCoordinatesPerTileSide / 4
}

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
	icons       *MapIcons

	moveCapture func(pixelX, pixelY float32)
	mouseMoved  bool

	positionPopupPos imgui.Vec2
	positionValid    bool
	position         MapPosition

	selectedTiles   tileCoordinates
	selectedObjects objectIDs

	activeLevel         *level.Level
	availableHoverItems []hoverItem
	activeHoverIndex    int
	activeHoverItem     hoverItem
}

// NewMapDisplay returns a new instance.
func NewMapDisplay(gl opengl.OpenGL, guiScale float32,
	textureQuery TextureQuery,
	eventListener event.Listener, eventRegistry event.Registry) *MapDisplay {
	tilesPerMapSide := float32(64)

	tileBaseLength := float32(fineCoordinatesPerTileSide)
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
	display.icons = NewMapIcons(&display.context)

	centerX, centerY := (tilesPerMapSide*tileBaseLength)/-2.0, (tilesPerMapSide*tileBaseLength)/-2.0
	display.camera.ZoomAt(-3+zoomShift, centerX, centerY)
	display.camera.MoveTo(centerX, centerY)

	display.selectedTiles.registerAt(eventRegistry)
	display.selectedObjects.registerAt(eventRegistry)

	return display
}

// Render renders the whole map display.
func (display *MapDisplay) Render(properties object.PropertiesTable, lvl *level.Level,
	paletteTexture *graphics.PaletteTexture, textureRetriever func(resource.Key) (*graphics.BitmapTexture, error),
	textureDisplay TextureDisplay, colorDisplay ColorDisplay) {
	columns, rows, _ := lvl.Size()

	display.selectedObjects.filterInvalid(lvl)

	display.setActiveLevel(lvl)
	display.background.Render(columns, rows)
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
			display.textures.Render(columns, rows, func(x, y int) (level.TileType, level.TextureIndex, int) {
				tile := lvl.Tile(x, y)
				if tile == nil {
					return level.TileTypeSolid, 0, 0
				}
				atlasIndex, textureRotations := textureDisplay.Func()(tile)
				atlas := lvl.TextureAtlas()
				textureIndex := level.TextureIndex(-1)
				if (int(atlasIndex) >= 0) && (int(atlasIndex) < len(atlas)) {
					textureIndex = atlas[atlasIndex]
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
	display.mapGrid.Render(columns, rows, lvl)
	if display.positionValid {
		if len(display.availableHoverItems) == 0 {
			display.availableHoverItems = display.nearestHoverItems(lvl, display.position)
			display.activeHoverIndex = 0
			display.activeHoverItem = display.availableHoverItems[0]
		}
	}
	display.highlighter.Render(display.selectedTiles.list, fineCoordinatesPerTileSide, [4]float32{0.0, 0.8, 0.2, 0.5})
	{
		var objects []MapPosition
		lvl.ForEachObject(func(id level.ObjectID, entry level.ObjectMainEntry) {
			objects = append(objects, MapPosition{X: entry.X, Y: entry.Y})
		})
		display.highlighter.Render(objects, fineCoordinatesPerTileSide/4, [4]float32{1.0, 1.0, 1.0, 0.3})
	}
	if paletteTexture != nil {
		tripleOffsets := make(map[object.Triple]int)

		{
			offset := 0
			properties.Iterate(func(triple object.Triple, prop *object.Properties) bool {
				numExtra := int(prop.Common.Bitmap3D.FrameNumber())

				if triple.Class != object.ClassTrap {
					tripleOffsets[triple] = offset + 2
				} else {
					tripleOffsets[triple] = offset
				}
				offset += 3 + numExtra
				return true
			})
		}
		var icons []iconData
		var highlightIcon iconData
		var highlightID level.ObjectID

		if display.positionValid {
			objectItem, isObject := display.activeHoverItem.(objectHoverItem)
			if isObject {
				highlightID = objectItem.id
			}
		}
		lvl.ForEachObject(func(id level.ObjectID, entry level.ObjectMainEntry) {
			triple := entry.Triple()
			index, cached := tripleOffsets[triple]
			if cached {
				key := resource.KeyOf(ids.ObjectBitmaps, resource.LangAny, index+1)
				texture, err := textureRetriever(key)
				if err == nil {
					icon := iconData{pos: MapPosition{X: entry.X, Y: entry.Y}, texture: texture}
					if highlightID == id {
						highlightIcon = icon
					} else {
						icons = append(icons, icon)
					}
				}
			}
		})
		if (highlightID != 0) && (highlightIcon.texture != nil) {
			icons = append(icons, highlightIcon)
		}
		display.icons.Render(paletteTexture, fineCoordinatesPerTileSide/4, icons)
	}
	{
		selectedObjectHighlights := make([]MapPosition, 0, len(display.selectedObjects.list))
		for _, entry := range display.selectedObjects.list {
			obj := lvl.Object(entry)
			if obj != nil {
				objPos := MapPosition{X: obj.X, Y: obj.Y}
				selectedObjectHighlights = append(selectedObjectHighlights, objPos)
			}
		}
		display.highlighter.Render(selectedObjectHighlights, fineCoordinatesPerTileSide/4, [4]float32{0.0, 0.8, 0.2, 0.5})
	}
	if display.activeHoverItem != nil {
		display.highlighter.Render([]MapPosition{display.activeHoverItem.Pos()}, display.activeHoverItem.Size(), [4]float32{0.0, 0.2, 0.8, 0.3})
	}

	display.renderPositionOverlay(lvl)
}

func (display *MapDisplay) nearestHoverItems(lvl *level.Level, ref MapPosition) []hoverItem {
	var items []hoverItem
	var distances []float32

	refVec := mgl.Vec2{float32(ref.X), float32(ref.Y)}

	lvl.ForEachObject(func(id level.ObjectID, entry level.ObjectMainEntry) {
		entryVec := mgl.Vec2{float32(entry.X), float32(entry.Y)}
		distance := refVec.Sub(entryVec).Len()
		if distance < fineCoordinatesPerTileSide/4 {
			items = append(items, objectHoverItem{id: id, pos: MapPosition{X: entry.X, Y: entry.Y}})
			distances = append(distances, distance)
		}
	})
	items = append(items, tileHoverItem{pos: MapPosition{
		X: level.CoordinateAt(ref.X.Tile(), 128),
		Y: level.CoordinateAt(ref.Y.Tile(), 128),
	}})
	distances = append(distances, fineCoordinatesPerTileSide)

	sort.Slice(items, func(a, b int) bool { return distances[a] < distances[b] })

	return items
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
	imgui.SetNextWindowSize(imgui.Vec2{X: 140 * display.guiScale, Y: 0})
	imgui.SetNextWindowBgAlpha(0.3)
	if imgui.BeginV("Position", nil, imgui.WindowFlagsNoMove|imgui.WindowFlagsNoTitleBar|imgui.WindowFlagsNoResize|imgui.WindowFlagsAlwaysAutoResize|
		imgui.WindowFlagsNoSavedSettings|imgui.WindowFlagsNoFocusOnAppearing|imgui.WindowFlagsNoNav) {
		typeString := "---"
		hasPos := false
		var pos MapPosition
		hasFloor := false
		var floorRaw int
		floorString := hintUnknown
		hasCeiling := false
		var ceilingRaw int
		ceilingString := hintUnknown

		if display.activeHoverItem != nil {
			pos = display.activeHoverItem.Pos()
			hasPos = true

			if _, isTileItem := display.activeHoverItem.(tileHoverItem); isTileItem {
				pos = display.position // use raw cursor position for this display
				typeString = "Tile"
				tile := lvl.Tile(int(pos.X.Tile()), int(pos.Y.Tile()))
				if (tile != nil) && (tile.Type != level.TileTypeSolid) {
					_, _, heightShift := lvl.Size()
					floorHeight := tile.Floor.AbsoluteHeight()
					floorHeightInTiles, err := heightShift.ValueFromTileHeight(floorHeight)
					if err == nil {
						floorString = fmt.Sprintf("%2.3f", floorHeightInTiles)
					}
					floorRaw = int(floorHeight)
					hasFloor = true

					ceilingHeight := tile.Ceiling.AbsoluteHeight()
					ceilingHeightInTiles, err := heightShift.ValueFromTileHeight(ceilingHeight)
					if err == nil {
						ceilingString = fmt.Sprintf("%2.3f", ceilingHeightInTiles)
					}
					ceilingRaw = int(ceilingHeight)
					hasCeiling = true
				}
			} else if objectItem, isObjectItem := display.activeHoverItem.(objectHoverItem); isObjectItem {
				_, _, heightShift := lvl.Size()
				obj := lvl.Object(objectItem.id)
				if obj != nil {
					typeString = fmt.Sprintf("%3d = %v", objectItem.id, obj.Triple())
					heightInTiles, err := heightShift.ValueFromObjectHeight(obj.Z)
					if err == nil {
						floorString = fmt.Sprintf("%2.3f", heightInTiles)
					}
					floorRaw = int(obj.Z)
					hasFloor = true
				}
			}
		}
		imgui.Text("T: " + typeString)
		if hasPos {
			imgui.Text(fmt.Sprintf("X: T %2d F %3d", pos.X.Tile(), pos.X.Fine()))
			imgui.Text(fmt.Sprintf("Y: T %2d F %3d", pos.Y.Tile(), pos.Y.Fine()))
		} else {
			imgui.Text("X: T -- F ---")
			imgui.Text("Y: T -- F ---")
		}
		if hasFloor {
			imgui.Text(fmt.Sprintf("F: %3d = %s", floorRaw, floorString))
		} else {
			imgui.Text("F: -- = --.---")
		}
		if hasCeiling {
			imgui.Text(fmt.Sprintf("C: %3d = %s", ceilingRaw, ceilingString))
		} else {
			imgui.Text("C: -- = --.---")
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
			switch {
			case modifier.Has(input.ModControl):
				display.toggleSelectionAtActiveHoverItem()
			case modifier.Has(input.ModShift) && (len(display.selectedTiles.list) > 0):
				firstPos := display.selectedTiles.list[0]

				fromX := int(firstPos.X.Tile())
				fromY := int(firstPos.Y.Tile())
				toX := int(display.position.X.Tile())
				toY := int(display.position.Y.Tile())
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
				display.eventListener.Event(ObjectSelectionSetEvent{objects: display.objectsInTiles(newList)})
			default:
				display.setSelectionByActiveHoverItem()
			}
		}
	} else if button == input.MouseSecondary {
		if display.positionValid {
			evt := ObjectRequestCreateEvent{Pos: display.position}
			if modifier.Has(input.ModShift) {
				toGrid := func(value byte) byte {
					switch {
					case value < 0x40:
						return 0x00
					case value >= 0xC0:
						return 0xFF
					default:
						return 0x80
					}
				}
				evt.Pos.X = level.CoordinateAt(evt.Pos.X.Tile(), toGrid(evt.Pos.X.Fine()))
				evt.Pos.Y = level.CoordinateAt(evt.Pos.Y.Tile(), toGrid(evt.Pos.Y.Fine()))
			}
			display.eventListener.Event(evt)
		}
	}
}

func (display *MapDisplay) setSelectionByActiveHoverItem() {
	var tiles []MapPosition
	var objects []level.ObjectID
	if display.activeHoverItem != nil {
		if tileItem, isTile := display.activeHoverItem.(tileHoverItem); isTile {
			tiles = append(tiles, tileItem.pos)
		} else if objectItem, isObject := display.activeHoverItem.(objectHoverItem); isObject {
			objects = append(objects, objectItem.id)
		}
	}
	display.eventListener.Event(TileSelectionSetEvent{tiles: tiles})
	if len(tiles) > 0 {
		display.eventListener.Event(ObjectSelectionSetEvent{objects: display.objectsInTiles(tiles)})
	} else {
		display.eventListener.Event(ObjectSelectionSetEvent{objects: objects})
	}
}

func (display *MapDisplay) toggleSelectionAtActiveHoverItem() {
	if display.activeHoverItem != nil {
		if tileItem, isTile := display.activeHoverItem.(tileHoverItem); isTile {
			wasSelected := display.selectedTiles.contains(tileItem.pos)
			tiles := []MapPosition{tileItem.pos}
			if wasSelected {
				display.eventListener.Event(TileSelectionRemoveEvent{tiles: tiles})
				display.eventListener.Event(ObjectSelectionRemoveEvent{objects: display.objectsInTiles(tiles)})
			} else {
				display.eventListener.Event(TileSelectionAddEvent{tiles: tiles})
				display.eventListener.Event(ObjectSelectionAddEvent{objects: display.objectsInTiles(tiles)})
			}
		} else if objectItem, isObject := display.activeHoverItem.(objectHoverItem); isObject {
			wasSelected := display.selectedObjects.contains(objectItem.id)
			if wasSelected {
				display.eventListener.Event(ObjectSelectionRemoveEvent{objects: []level.ObjectID{objectItem.id}})
			} else {
				display.eventListener.Event(ObjectSelectionAddEvent{objects: []level.ObjectID{objectItem.id}})
			}
		}
	}
}

func (display *MapDisplay) objectsInTiles(tiles []MapPosition) []level.ObjectID {
	tilesContain := func(pos MapPosition) bool {
		for _, entry := range tiles {
			if entry == pos {
				return true
			}
		}
		return false
	}

	var objects []level.ObjectID
	if display.activeLevel != nil {
		display.activeLevel.ForEachObject(func(id level.ObjectID, entry level.ObjectMainEntry) {
			tilePos := MapPosition{X: level.CoordinateAt(entry.X.Tile(), 128), Y: level.CoordinateAt(entry.Y.Tile(), 128)}
			if tilesContain(tilePos) {
				objects = append(objects, id)
			}
		})
	}
	return objects
}

// MouseMoved must be called for a mouse move.
func (display *MapDisplay) MouseMoved(mouseX, mouseY float32) {
	display.updateMouseWorldPosition(mouseX, mouseY)
	display.resetHoverItems()
	display.moveCapture(mouseX, mouseY)
}

// MouseScrolled must be called for a mouse scroll.
func (display *MapDisplay) MouseScrolled(mouseX, mouseY float32, deltaX, deltaY float32, modifier input.Modifier) {
	if modifier.Has(input.ModControl) {
		hoverItems := len(display.availableHoverItems)
		if hoverItems > 0 {
			diff := 1
			if deltaY < 0 {
				diff = -1
			}
			display.activeHoverIndex = (hoverItems + (display.activeHoverIndex + diff)) % hoverItems
			display.activeHoverItem = display.availableHoverItems[display.activeHoverIndex]
		}
	} else {
		worldX, worldY := display.unprojectPixel(mouseX, mouseY)

		if deltaY < 0 {
			display.camera.ZoomAt(-0.5, worldX, worldY)
		}
		if deltaY > 0 {
			display.camera.ZoomAt(0.5, worldX, worldY)
		}
	}
}

func (display *MapDisplay) updateMouseWorldPosition(mouseX, mouseY float32) {
	worldX, worldY := display.unprojectPixel(mouseX, mouseY)
	var worldWidth float32
	var worldHeight float32
	if display.activeLevel != nil {
		width, height, _ := display.activeLevel.Size()
		worldWidth = float32(width) * fineCoordinatesPerTileSide
		worldHeight = float32(height) * fineCoordinatesPerTileSide
	}
	display.positionValid = (worldX >= 0.0) && (worldX < worldWidth) && (worldY >= 0.0) && (worldY < worldHeight)
	if display.positionValid {
		display.position = MapPosition{X: level.Coordinate(worldX + 0.5), Y: level.Coordinate(worldY + 0.5)}
	} else {
		display.position = MapPosition{}
	}
}

func (display *MapDisplay) setActiveLevel(lvl *level.Level) {
	oldIsNil := display.activeLevel == nil
	newIsNil := lvl == nil
	isChanged := (oldIsNil != newIsNil) || (!oldIsNil && !newIsNil && display.activeLevel.ID() != lvl.ID())
	if isChanged {
		display.activeLevel = lvl
		display.resetHoverItems()
	}
}

func (display *MapDisplay) resetHoverItems() {
	display.availableHoverItems = nil
	display.activeHoverIndex = 0
	display.activeHoverItem = nil
}
