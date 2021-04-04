package levels

import (
	"fmt"

	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// TilesView is for tile properties.
type TilesView struct {
	editor       *edit.LevelEditorService
	textCache    *text.Cache
	textureCache *graphics.TextureCache

	guiScale float32
	registry cmd.Registry
	model    tilesViewModel
}

// NewTilesView returns a new instance.
func NewTilesView(editor *edit.LevelEditorService,
	guiScale float32, textCache *text.Cache, textureCache *graphics.TextureCache, registry cmd.Registry) *TilesView {
	view := &TilesView{
		editor:       editor,
		textCache:    textCache,
		textureCache: textureCache,

		guiScale: guiScale,
		model:    freshTilesViewModel(),
		registry: registry,
	}
	return view
}

// WindowOpen returns the flag address, to be used with the main menu.
func (view *TilesView) WindowOpen() *bool {
	return &view.model.windowOpen
}

// TextureDisplay returns the current setting which textures should be displayed.
func (view TilesView) TextureDisplay() TextureDisplay {
	return view.model.textureDisplay
}

// ColorDisplay returns the current setting which colors should be displayed.
func (view TilesView) ColorDisplay() ColorDisplay {
	lvl := view.editor.Level()
	if lvl.IsCyberspace() {
		return view.model.cyberColorDisplay
	}
	return view.model.shadowDisplay
}

// Render renders the view.
func (view *TilesView) Render() {
	if view.model.restoreFocus {
		imgui.SetNextWindowFocus()
		view.model.restoreFocus = false
		view.model.windowOpen = true
	}
	if view.model.windowOpen {
		lvl := view.editor.Level()
		tiles := view.editor.Tiles()
		readOnly := view.editor.IsReadOnly()
		title := fmt.Sprintf("Level Tiles, %d selected", len(tiles))
		if readOnly {
			title += hintReadOnly
		}
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 500 * view.guiScale}, imgui.ConditionFirstUseEver)
		if imgui.BeginV(title+"###Level Tiles", view.WindowOpen(), 0) {
			view.renderContent(lvl, tiles, readOnly)
		}
		imgui.End()
	}
}

func (view *TilesView) renderContent(lvl *level.Level, tiles []*level.TileMapEntry, readOnly bool) {
	isCyberspace := lvl.IsCyberspace()
	tileTypeUnifier := values.NewUnifier()
	floorHeightUnifier := values.NewUnifier()
	ceilingHeightUnifier := values.NewUnifier()
	slopeHeightUnifier := values.NewUnifier()
	slopeControlUnifier := values.NewUnifier()
	musicIndexUnifier := values.NewUnifier()

	floorPaletteIndexUnifier := values.NewUnifier()
	ceilingPaletteIndexUnifier := values.NewUnifier()
	flightPullTypeUnifier := values.NewUnifier()
	gameOfLightStateUnifier := values.NewUnifier()

	floorTextureIndexUnifier := values.NewUnifier()
	floorTextureRotationsUnifier := values.NewUnifier()
	ceilingTextureIndexUnifier := values.NewUnifier()
	ceilingTextureRotationsUnifier := values.NewUnifier()
	wallTextureIndexUnifier := values.NewUnifier()
	wallTextureOffsetUnifier := values.NewUnifier()
	useAdjacentWallTextureUnifier := values.NewUnifier()
	wallTexturePatternUnifier := values.NewUnifier()
	floorLightUnifier := values.NewUnifier()
	ceilingLightUnifier := values.NewUnifier()
	deconstructedUnifier := values.NewUnifier()
	floorHazardUnifier := values.NewUnifier()
	ceilingHazardUnifier := values.NewUnifier()

	for _, tile := range tiles {
		tileTypeUnifier.Add(tile.Type)
		floorHeightUnifier.Add(tile.Floor.AbsoluteHeight())
		ceilingHeightUnifier.Add(tile.Ceiling.AbsoluteHeight())
		slopeHeightUnifier.Add(tile.SlopeHeight)
		slopeControlUnifier.Add(tile.Flags.SlopeControl())
		musicIndexUnifier.Add(tile.Flags.MusicIndex())
		if isCyberspace {
			floorPaletteIndexUnifier.Add(tile.TextureInfo.FloorPaletteIndex())
			ceilingPaletteIndexUnifier.Add(tile.TextureInfo.CeilingPaletteIndex())
			flightPullTypeUnifier.Add(tile.Flags.ForCyberspace().FlightPull())
			gameOfLightStateUnifier.Add(tile.Flags.ForCyberspace().GameOfLifeState())
		} else {
			flags := tile.Flags.ForRealWorld()
			floorTextureIndexUnifier.Add(int(tile.TextureInfo.FloorTextureIndex()))
			floorTextureRotationsUnifier.Add(tile.Floor.TextureRotations())
			ceilingTextureIndexUnifier.Add(int(tile.TextureInfo.CeilingTextureIndex()))
			ceilingTextureRotationsUnifier.Add(tile.Ceiling.TextureRotations())
			wallTextureIndexUnifier.Add(int(tile.TextureInfo.WallTextureIndex()))
			wallTextureOffsetUnifier.Add(flags.WallTextureOffset())
			useAdjacentWallTextureUnifier.Add(flags.UseAdjacentWallTexture())
			wallTexturePatternUnifier.Add(flags.WallTexturePattern())
			floorLightUnifier.Add(level.GradesOfShadow - 1 - flags.FloorShadow())
			ceilingLightUnifier.Add(level.GradesOfShadow - 1 - flags.CeilingShadow())
			deconstructedUnifier.Add(flags.Deconstructed())
			floorHazardUnifier.Add(tile.Floor.HasHazard())
			ceilingHazardUnifier.Add(tile.Ceiling.HasHazard())
		}
	}

	imgui.PushItemWidth(-250 * view.guiScale)

	_, _, levelHeight := lvl.Size()
	tileHeightFormatter := tileHeightFormatterFor(levelHeight)

	tileTypes := level.TileTypes()
	values.RenderUnifiedCombo(readOnly, "Tile Type", tileTypeUnifier,
		func(u values.Unifier) int { return int(u.Unified().(level.TileType)) },
		func(value int) string { return tileTypes[value].String() },
		len(tileTypes),
		func(newValue int) { view.changeTiles(setTileTypeTo(level.TileType(newValue))) })
	values.RenderUnifiedSliderInt(readOnly, "Floor Height", floorHeightUnifier,
		func(u values.Unifier) int { return int(u.Unified().(level.TileHeightUnit)) },
		tileHeightFormatter,
		0, int(level.TileHeightUnitMax)-1,
		func(newValue int) { view.changeTiles(setFloorHeightTo(level.TileHeightUnit(newValue))) })
	values.RenderUnifiedSliderInt(readOnly, "Ceiling Height (abs)", ceilingHeightUnifier,
		func(u values.Unifier) int { return int(u.Unified().(level.TileHeightUnit)) },
		tileHeightFormatter,
		1, int(level.TileHeightUnitMax),
		func(newValue int) { view.changeTiles(setCeilingHeightTo(level.TileHeightUnit(newValue))) })
	values.RenderUnifiedSliderInt(readOnly, "Slope Height", slopeHeightUnifier,
		func(u values.Unifier) int { return int(u.Unified().(level.TileHeightUnit)) },
		tileHeightFormatter,
		0, int(level.TileHeightUnitMax)-1,
		func(newValue int) { view.changeTiles(setSlopeHeightTo(level.TileHeightUnit(newValue))) })
	slopeControls := level.TileSlopeControls()
	values.RenderUnifiedCombo(readOnly, "Slope Control", slopeControlUnifier,
		func(u values.Unifier) int { return int(u.Unified().(level.TileSlopeControl)) },
		func(value int) string { return slopeControls[value].String() },
		len(slopeControls),
		func(newValue int) { view.changeTiles(setSlopeControlTo(slopeControls[newValue])) })
	values.RenderUnifiedSliderInt(readOnly, "Music Index", musicIndexUnifier,
		func(u values.Unifier) int { return u.Unified().(int) },
		func(value int) string { return "%d" },
		0, 15,
		func(newValue int) { view.changeTiles(setMusicIndexTo(newValue)) })

	imgui.Separator()

	if isCyberspace {
		if imgui.BeginCombo("Color View", view.model.cyberColorDisplay.String()) {
			displays := ColorDisplays()
			for _, display := range displays {
				displayString := display.String()

				if imgui.SelectableV(displayString, display == view.model.cyberColorDisplay, 0, imgui.Vec2{}) {
					view.model.cyberColorDisplay = display
				}
			}
			imgui.EndCombo()
		}

		values.RenderUnifiedSliderInt(readOnly, "Floor Color", floorPaletteIndexUnifier,
			func(u values.Unifier) int { return int(u.Unified().(byte)) },
			func(value int) string { return "%d" },
			0, bitmap.PaletteSize-1,
			func(newValue int) { view.changeTiles(setFloorPaletteIndexTo(newValue)) })
		values.RenderUnifiedSliderInt(readOnly, "Ceiling Color", ceilingPaletteIndexUnifier,
			func(u values.Unifier) int { return int(u.Unified().(byte)) },
			func(value int) string { return "%d" },
			0, bitmap.PaletteSize-1,
			func(newValue int) { view.changeTiles(setCeilingPaletteIndexTo(newValue)) })

		imgui.Separator()

		flightPulls := level.CyberspaceFlightPulls()
		values.RenderUnifiedCombo(readOnly, "Flight Pull Type", flightPullTypeUnifier,
			func(u values.Unifier) int { return int(u.Unified().(level.CyberspaceFlightPull)) },
			func(value int) string { return flightPulls[value].String() },
			len(flightPulls),
			func(newValue int) { view.changeTiles(setFlightPullTypeTo(flightPulls[newValue])) })
		values.RenderUnifiedSliderInt(readOnly, "Game Of Life State", gameOfLightStateUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string { return "%d" },
			0, 3,
			func(newValue int) { view.changeTiles(setGameOfLightStateTo(newValue)) })
	} else {
		atlas := lvl.TextureAtlas()

		if imgui.BeginCombo("Texture View", view.model.textureDisplay.String()) {
			displays := TextureDisplays()
			for _, display := range displays {
				displayString := display.String()

				if imgui.SelectableV(displayString, display == view.model.textureDisplay, 0, imgui.Vec2{}) {
					view.model.textureDisplay = display
				}
			}
			imgui.EndCombo()
		}

		values.RenderUnifiedSliderInt(readOnly, "Floor Texture (atlas index)", floorTextureIndexUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string { return "%d" },
			0, level.FloorCeilingTextureLimit-1,
			func(newValue int) { view.changeTiles(setFloorTextureIndexTo(level.AtlasIndex(newValue))) })
		view.renderTextureSelector(readOnly, "Floor Texture", floorTextureIndexUnifier, atlas, 0, level.FloorCeilingTextureLimit-1,
			func(newValue int) { view.changeTiles(setFloorTextureIndexTo(level.AtlasIndex(newValue))) })
		values.RenderUnifiedSliderInt(readOnly, "Floor Texture Rotations", floorTextureRotationsUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string { return "%d" },
			0, 3,
			func(newValue int) { view.changeTiles(setFloorTextureRotationsTo(newValue)) })

		values.RenderUnifiedSliderInt(readOnly, "Ceiling Texture (atlas index)", ceilingTextureIndexUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string { return "%d" },
			0, level.FloorCeilingTextureLimit-1,
			func(newValue int) { view.changeTiles(setCeilingTextureIndexTo(level.AtlasIndex(newValue))) })
		view.renderTextureSelector(readOnly, "Ceiling Texture", ceilingTextureIndexUnifier, atlas, 0, level.FloorCeilingTextureLimit-1,
			func(newValue int) { view.changeTiles(setCeilingTextureIndexTo(level.AtlasIndex(newValue))) })
		values.RenderUnifiedSliderInt(readOnly, "Ceiling Texture Rotations", ceilingTextureRotationsUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string { return "%d" },
			0, 3,
			func(newValue int) { view.changeTiles(setCeilingTextureRotationsTo(newValue)) })

		values.RenderUnifiedSliderInt(readOnly, "Wall Texture (atlas index)", wallTextureIndexUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string { return "%d" },
			0, len(atlas)-1,
			func(newValue int) { view.changeTiles(setWallTextureIndexTo(level.AtlasIndex(newValue))) })
		view.renderTextureSelector(readOnly, "Wall Texture", wallTextureIndexUnifier, atlas, 0, len(atlas)-1,
			func(newValue int) { view.changeTiles(setWallTextureIndexTo(level.AtlasIndex(newValue))) })
		values.RenderUnifiedSliderInt(readOnly, "Wall Texture Offset", wallTextureOffsetUnifier,
			func(u values.Unifier) int { return int(u.Unified().(level.TileHeightUnit)) },
			tileHeightFormatter,
			0, int(level.TileHeightUnitMax)-1,
			func(newValue int) { view.changeTiles(setWallTextureOffsetTo(level.TileHeightUnit(newValue))) })

		values.RenderUnifiedCheckboxCombo(readOnly, "Use Adjacent Wall Texture", useAdjacentWallTextureUnifier,
			func(newValue bool) { view.changeTiles(setUseAdjacentWallTextureTo(newValue)) })
		wallTexturePatterns := level.WallTexturePatterns()
		values.RenderUnifiedCombo(readOnly, "Wall Texture Pattern", wallTexturePatternUnifier,
			func(u values.Unifier) int { return int(u.Unified().(level.WallTexturePattern)) },
			func(value int) string { return wallTexturePatterns[value].String() },
			len(wallTexturePatterns),
			func(newValue int) { view.changeTiles(setWallTexturePatternTo(wallTexturePatterns[newValue])) })

		imgui.Separator()

		if imgui.BeginCombo("Shadow View", view.model.shadowDisplay.String()) {
			displays := ColorDisplays()
			for _, display := range displays {
				displayString := display.String()

				if imgui.SelectableV(displayString, display == view.model.shadowDisplay, 0, imgui.Vec2{}) {
					view.model.shadowDisplay = display
				}
			}
			imgui.EndCombo()
		}

		values.RenderUnifiedSliderInt(readOnly, "Floor Light", floorLightUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string { return "%d" },
			0, level.GradesOfShadow-1,
			func(newValue int) { view.changeTiles(setFloorLightTo(newValue)) })
		values.RenderUnifiedSliderInt(readOnly, "Ceiling Light", ceilingLightUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string { return "%d" },
			0, level.GradesOfShadow-1,
			func(newValue int) { view.changeTiles(setCeilingLightTo(newValue)) })

		imgui.Separator()

		values.RenderUnifiedCheckboxCombo(readOnly, "Deconstructed", deconstructedUnifier,
			func(newValue bool) { view.changeTiles(setDeconstructedTo(newValue)) })
		values.RenderUnifiedCheckboxCombo(readOnly, "Floor Hazard", floorHazardUnifier,
			func(newValue bool) { view.changeTiles(setFloorHazardTo(newValue)) })
		values.RenderUnifiedCheckboxCombo(readOnly, "Ceiling Hazard", ceilingHazardUnifier,
			func(newValue bool) { view.changeTiles(setCeilingHazardTo(newValue)) })
	}

	imgui.PopItemWidth()
}

func (view *TilesView) renderTextureSelector(readOnly bool, label string, unifier values.Unifier,
	atlas level.TextureAtlas, minIndex, maxIndex int, changeHandler func(int)) {
	selectedIndex := -1
	if unifier.IsUnique() {
		selectedIndex = unifier.Unified().(int)
	}

	count := maxIndex - minIndex + 1
	if count > len(atlas) {
		count = len(atlas)
	}
	render.TextureSelector(label, -1, view.guiScale, count, selectedIndex-minIndex,
		view.textureCache,
		func(index int) resource.Key {
			return resource.KeyOf(ids.LargeTextures.Plus(int(atlas[minIndex+index])), resource.LangAny, 0)
		},
		func(index int) string { return view.textureName(int(atlas[minIndex+index])) },
		func(index int) {
			if !readOnly {
				changeHandler(index)
			}
		})
}

func (view *TilesView) textureName(index int) string {
	key := resource.KeyOf(ids.TextureNames, resource.LangDefault, index)
	name, err := view.textCache.Text(key)
	suffix := ""
	if err == nil {
		suffix = ": " + name
	}
	return fmt.Sprintf("%3d", index) + suffix
}

func setTileTypeTo(tileType level.TileType) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Type = tileType
	}
}

func setFloorHeightTo(height level.TileHeightUnit) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Floor = tile.Floor.WithAbsoluteHeight(height)
	}
}

func setCeilingHeightTo(height level.TileHeightUnit) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Ceiling = tile.Ceiling.WithAbsoluteHeight(height)
	}
}

func setSlopeHeightTo(height level.TileHeightUnit) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.SlopeHeight = height
	}
}

func setSlopeControlTo(value level.TileSlopeControl) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.WithSlopeControl(value)
	}
}

func setMusicIndexTo(value int) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.WithMusicIndex(value)
	}
}

func setFloorTextureIndexTo(value level.AtlasIndex) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.TextureInfo = tile.TextureInfo.WithFloorTextureIndex(value)
	}
}

func setFloorTextureRotationsTo(value int) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Floor = tile.Floor.WithTextureRotations(value)
	}
}

func setCeilingTextureIndexTo(value level.AtlasIndex) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.TextureInfo = tile.TextureInfo.WithCeilingTextureIndex(value)
	}
}

func setCeilingTextureRotationsTo(value int) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Ceiling = tile.Ceiling.WithTextureRotations(value)
	}
}

func setWallTextureIndexTo(value level.AtlasIndex) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.TextureInfo = tile.TextureInfo.WithWallTextureIndex(value)
	}
}

func setWallTextureOffsetTo(value level.TileHeightUnit) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForRealWorld().WithWallTextureOffset(value).AsTileFlag()
	}
}

func setUseAdjacentWallTextureTo(value bool) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForRealWorld().WithUseAdjacentWallTexture(value).AsTileFlag()
	}
}

func setWallTexturePatternTo(value level.WallTexturePattern) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForRealWorld().WithWallTexturePattern(value).AsTileFlag()
	}
}

func setFloorLightTo(value int) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForRealWorld().WithFloorShadow(level.GradesOfShadow - 1 - value).AsTileFlag()
	}
}

func setCeilingLightTo(value int) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForRealWorld().WithCeilingShadow(level.GradesOfShadow - 1 - value).AsTileFlag()
	}
}

func setDeconstructedTo(value bool) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForRealWorld().WithDeconstructed(value).AsTileFlag()
	}
}

func setFloorHazardTo(value bool) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Floor = tile.Floor.WithHazard(value)
	}
}

func setCeilingHazardTo(value bool) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Ceiling = tile.Ceiling.WithHazard(value)
	}
}

func setFloorPaletteIndexTo(value int) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.TextureInfo = tile.TextureInfo.WithFloorPaletteIndex(byte(value))
	}
}

func setCeilingPaletteIndexTo(value int) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.TextureInfo = tile.TextureInfo.WithCeilingPaletteIndex(byte(value))
	}
}

func setFlightPullTypeTo(value level.CyberspaceFlightPull) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForCyberspace().WithFlightPull(value).AsTileFlag()
	}
}

func setGameOfLightStateTo(value int) tileMapEntryModifier {
	return func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForCyberspace().WithGameOfLifeState(value).AsTileFlag()
	}
}

type tileMapEntryModifier func(*level.TileMapEntry)

func (view *TilesView) changeTiles(modifier tileMapEntryModifier) {
	err := view.registry.Register(cmd.Named("ChangeTiles"),
		cmd.Forward(view.restoreFocusTask()),
		cmd.Nested(func() error { return view.editor.ChangeTiles(modifier) }),
		cmd.Reverse(view.restoreFocusTask()))
	if err != nil {
		panic(err)
	}
}

func (view *TilesView) restoreFocusTask() cmd.Task {
	return func(modder world.Modder) error {
		view.model.restoreFocus = true
		return nil
	}
}
