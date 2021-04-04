package levels

import (
	"fmt"
	"sort"

	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/input"
	"github.com/inkyblackness/hacked/ui/opengl"
)

const contextMenuName = "MapContextMenu"

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
	return level.FineCoordinatesPerTileSide
}

type objectHoverItem struct {
	id  level.ObjectID
	pos MapPosition
}

func (item objectHoverItem) Pos() MapPosition {
	return item.pos
}

func (item objectHoverItem) Size() float32 {
	return level.FineCoordinatesPerTileSide / 4
}

// MapDisplay renders a level map.
type MapDisplay struct {
	levelSelection *edit.LevelSelectionService
	editor         *edit.LevelEditorService

	context  render.Context
	camera   *LimitedCamera
	guiScale float32

	background  *BackgroundGrid
	textures    *MapTextures
	colors      *MapColors
	mapGrid     *MapGrid
	highlighter *Highlighter
	icons       *MapIcons

	moveCapture func(pixelX, pixelY float32)
	mouseMoved  bool

	positionPopupPos     imgui.Vec2
	positionValid        bool
	position             MapPosition
	contextMenuRequested bool

	selectionReference *level.TilePosition

	activeLevel         *level.Level
	availableHoverItems []hoverItem
	activeHoverIndex    int
	activeHoverItem     hoverItem
}

// NewMapDisplay returns a new instance.
func NewMapDisplay(levelSelection *edit.LevelSelectionService, editor *edit.LevelEditorService,
	gl opengl.OpenGL, guiScale float32,
	textureQuery TextureQuery) *MapDisplay {
	tilesPerMapSide := float32(64)

	tileBaseLength := float32(level.FineCoordinatesPerTileSide)
	tileBaseHalf := tileBaseLength / 2.0
	camLimit := tilesPerMapSide*tileBaseLength - tileBaseHalf
	zoomShift := guiScale - 1.0
	zoomLevelMin := float32(-5) + zoomShift
	zoomLevelMax := float32(1) + zoomShift

	display := &MapDisplay{
		levelSelection: levelSelection,
		editor:         editor,
		context: render.Context{
			OpenGL:           gl,
			ProjectionMatrix: mgl.Ident4(),
		},
		camera:      NewLimitedCamera(zoomLevelMin, zoomLevelMax, -tileBaseHalf, camLimit),
		guiScale:    guiScale,
		moveCapture: func(float32, float32) {},
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

	return display
}

// Render renders the whole map display.
func (display *MapDisplay) Render(properties object.PropertiesTable,
	paletteTexture *graphics.PaletteTexture, textureRetriever func(resource.Key) (*graphics.BitmapTexture, error),
	textureDisplay TextureDisplay, colorDisplay ColorDisplay) {
	lvl := display.editor.Level()
	columns, rows, _ := lvl.Size()

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
			display.textures.Render(columns, rows, func(pos level.TilePosition) (level.TileType, level.TextureIndex, int) {
				tile := lvl.Tile(pos)
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
	{
		selectedTiles := display.levelSelection.CurrentSelectedTiles()
		tileMapPositions := make([]MapPosition, 0, len(selectedTiles))
		for _, pos := range selectedTiles {
			tileMapPositions = append(tileMapPositions, MapPosition{
				X: level.CoordinateAt(pos.X, level.FineCoordinatesPerTileSide/2),
				Y: level.CoordinateAt(pos.Y, level.FineCoordinatesPerTileSide/2),
			})
		}
		display.highlighter.Render(tileMapPositions, level.FineCoordinatesPerTileSide, [4]float32{0.0, 0.8, 0.2, 0.5})
	}
	{
		var objects []MapPosition
		lvl.ForEachObject(func(id level.ObjectID, entry level.ObjectMainEntry) {
			objects = append(objects, MapPosition{X: entry.X, Y: entry.Y})
		})
		display.highlighter.Render(objects, level.FineCoordinatesPerTileSide/4, [4]float32{1.0, 1.0, 1.0, 0.3})
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
		display.icons.Render(paletteTexture, level.FineCoordinatesPerTileSide/4, icons)
	}
	{
		selectedObjects := display.editor.Objects()
		selectedObjectHighlights := make([]MapPosition, 0, len(selectedObjects))
		for _, obj := range selectedObjects {
			objPos := MapPosition{X: obj.X, Y: obj.Y}
			selectedObjectHighlights = append(selectedObjectHighlights, objPos)
		}
		display.highlighter.Render(selectedObjectHighlights, level.FineCoordinatesPerTileSide/4, [4]float32{0.0, 0.8, 0.2, 0.5})
	}
	if display.activeHoverItem != nil {
		display.highlighter.Render([]MapPosition{display.activeHoverItem.Pos()}, display.activeHoverItem.Size(), [4]float32{0.0, 0.2, 0.8, 0.3})
	}

	display.renderPositionOverlay(lvl)
	if display.contextMenuRequested {
		imgui.OpenPopup(contextMenuName)
		display.contextMenuRequested = false
	}
	display.renderContextMenu()
}

func (display *MapDisplay) nearestHoverItems(lvl *level.Level, ref MapPosition) []hoverItem {
	var items []hoverItem
	var distances []float32

	refVec := mgl.Vec2{float32(ref.X), float32(ref.Y)}

	lvl.ForEachObject(func(id level.ObjectID, entry level.ObjectMainEntry) {
		entryVec := mgl.Vec2{float32(entry.X), float32(entry.Y)}
		distance := refVec.Sub(entryVec).Len()
		if distance < level.FineCoordinatesPerTileSide/4 {
			items = append(items, objectHoverItem{id: id, pos: MapPosition{X: entry.X, Y: entry.Y}})
			distances = append(distances, distance)
		}
	})
	items = append(items, tileHoverItem{pos: MapPosition{
		X: level.CoordinateAt(ref.X.Tile(), level.FineCoordinatesPerTileSide/2),
		Y: level.CoordinateAt(ref.Y.Tile(), level.FineCoordinatesPerTileSide/2),
	}})
	distances = append(distances, level.FineCoordinatesPerTileSide)

	sort.Slice(items, func(a, b int) bool { return distances[a] < distances[b] })

	return items
}

func (display *MapDisplay) colorQueryFor(lvl *level.Level, tileToColor func(*level.TileMapEntry) [4]float32) func(level.TilePosition) [4]float32 {
	return func(pos level.TilePosition) [4]float32 {
		tile := lvl.Tile(pos)
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
				tile := lvl.Tile(pos.Tile())
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

func (display *MapDisplay) renderContextMenu() {
	if imgui.BeginPopupV(contextMenuName, imgui.PopupFlagsMouseButtonRight) {
		if imgui.BeginMenu("New...") {
			implicitTriple := display.editor.NewObjectTriple()
			canCreateImplicitClass := display.canCreateObjectOf(implicitTriple.Class)
			if imgui.MenuItemV("Object", "2ndClick", false, canCreateImplicitClass) {
				display.requestCreateNewObject(false, implicitTriple)
			}
			if imgui.MenuItemV("Object (at grid)", "Shift+2ndClick", false, canCreateImplicitClass) {
				display.requestCreateNewObject(true, implicitTriple)
			}
			imgui.EndMenu()
		}
		imgui.Separator()
		if imgui.MenuItemV("Delete Objects", "", false, !display.editor.IsReadOnly()) {
			_ = display.editor.DeleteObjects()
		}
		imgui.EndPopup()
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
			case modifier.Has(input.ModShift) && (display.selectionReference != nil):
				firstPos := *display.selectionReference

				fromX := int(firstPos.X)
				fromY := int(firstPos.Y)
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
				var newList []level.TilePosition
				for y := fromY; y != toY; y += yIncrement {
					for x := fromX; x != toX; x += xIncrement {
						newList = append(newList, level.TilePosition{X: byte(x), Y: byte(y)})
					}
				}
				display.levelSelection.SetCurrentSelectedTiles(newList)
				display.levelSelection.SetCurrentSelectedObjects(display.objectsInTiles(newList))
			default:
				display.setSelectionByActiveHoverItem()
			}
		}
	} else if button == input.MouseSecondary {
		if modifier.Has(input.ModControl) {
			display.contextMenuRequested = true
		} else {
			display.requestCreateNewObject(modifier.Has(input.ModShift), display.editor.NewObjectTriple())
		}
	}
}

func (display *MapDisplay) canCreateObjectOf(class object.Class) bool {
	if display.editor.IsReadOnly() {
		return false
	}
	lvl := display.editor.Level()
	return lvl.HasRoomForObjectOf(class) && display.positionValid
}

func (display *MapDisplay) requestCreateNewObject(snapToGrid bool, triple object.Triple) {
	if !display.canCreateObjectOf(triple.Class) {
		return
	}
	pos := display.position
	if snapToGrid {
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
		pos.X = level.CoordinateAt(pos.X.Tile(), toGrid(pos.X.Fine()))
		pos.Y = level.CoordinateAt(pos.Y.Tile(), toGrid(pos.Y.Fine()))
	}
	err := display.editor.CreateNewObject(triple, func(obj *level.ObjectMainEntry) {
		obj.X = pos.X
		obj.Y = pos.Y
	})
	if err != nil {
		panic(err)
	}
}

func (display *MapDisplay) setSelectionByActiveHoverItem() {
	var tiles []level.TilePosition
	var objects []level.ObjectID
	if display.activeHoverItem != nil {
		if tileItem, isTile := display.activeHoverItem.(tileHoverItem); isTile {
			tiles = append(tiles, tileItem.pos.Tile())
		} else if objectItem, isObject := display.activeHoverItem.(objectHoverItem); isObject {
			objects = append(objects, objectItem.id)
		}
	}
	display.levelSelection.SetCurrentSelectedTiles(tiles)
	if len(tiles) > 0 {
		display.selectionReference = &tiles[0]
		display.levelSelection.SetCurrentSelectedObjects(display.objectsInTiles(tiles))
	} else {
		display.levelSelection.SetCurrentSelectedObjects(objects)
	}
}

func (display *MapDisplay) toggleSelectionAtActiveHoverItem() {
	if display.activeHoverItem != nil {
		if tileItem, isTile := display.activeHoverItem.(tileHoverItem); isTile {
			wasSelected := display.levelSelection.IsTileSelected(tileItem.pos.Tile())
			tiles := []level.TilePosition{tileItem.pos.Tile()}
			if wasSelected {
				display.levelSelection.RemoveCurrentSelectedTiles(tiles)
				display.levelSelection.RemoveCurrentSelectedObjects(display.objectsInTiles(tiles))
			} else {
				display.levelSelection.AddCurrentSelectedTiles(tiles)
				display.levelSelection.AddCurrentSelectedObjects(display.objectsInTiles(tiles))
			}
		} else if objectItem, isObject := display.activeHoverItem.(objectHoverItem); isObject {
			display.levelSelection.ToggleObjectSelection([]level.ObjectID{objectItem.id})
		}
	}
}

func (display *MapDisplay) objectsInTiles(tiles []level.TilePosition) []level.ObjectID {
	tilesContain := func(pos level.TilePosition) bool {
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
			if tilesContain(entry.TilePosition()) {
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
		worldWidth = float32(width) * level.FineCoordinatesPerTileSide
		worldHeight = float32(height) * level.FineCoordinatesPerTileSide
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
