package levels

import (
	"fmt"

	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/editor/event"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// TilesView is for tile properties.
type TilesView struct {
	levels         *edit.EditableLevels
	levelSelection *edit.LevelSelectionService
	mod            *world.Mod
	textCache      *text.Cache
	textureCache   *graphics.TextureCache

	guiScale      float32
	registry      cmd.Registry
	eventListener event.Listener

	model tilesViewModel
}

// NewTilesView returns a new instance.
func NewTilesView(levels *edit.EditableLevels, levelSelection *edit.LevelSelectionService, mod *world.Mod,
	guiScale float32, textCache *text.Cache, textureCache *graphics.TextureCache,
	registry cmd.Registry, eventListener event.Listener, eventRegistry event.Registry) *TilesView {
	view := &TilesView{
		levels:         levels,
		levelSelection: levelSelection,
		mod:            mod,
		textCache:      textCache,
		textureCache:   textureCache,

		guiScale:      guiScale,
		registry:      registry,
		eventListener: eventListener,
		model:         freshTilesViewModel(),
	}
	view.model.selectedTiles.registerAt(eventRegistry)
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
func (view TilesView) ColorDisplay(lvl *level.Level) ColorDisplay {
	if lvl.IsCyberspace() {
		return view.model.cyberColorDisplay
	}
	return view.model.shadowDisplay
}

// Render renders the view.
func (view *TilesView) Render(lvl *level.Level) {
	if view.model.restoreFocus {
		imgui.SetNextWindowFocus()
		view.model.restoreFocus = false
		view.model.windowOpen = true
	}
	if view.model.windowOpen {
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 500 * view.guiScale}, imgui.ConditionFirstUseEver)
		title := fmt.Sprintf("Level Tiles, %d selected", len(view.model.selectedTiles.list))
		readOnly := !view.editingAllowed(lvl.ID())
		if readOnly {
			title += hintReadOnly
		}
		if imgui.BeginV(title+"###Level Tiles", view.WindowOpen(), 0) {
			view.renderContent(lvl, readOnly)
		}
		imgui.End()
	}
}

func (view *TilesView) renderContent(lvl *level.Level, readOnly bool) {
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

	for _, pos := range view.model.selectedTiles.list {
		tile := lvl.Tile(pos.Tile())
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
			floorLightUnifier.Add(15 - flags.FloorShadow())
			ceilingLightUnifier.Add(15 - flags.CeilingShadow())
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
		func(newValue int) {
			view.requestSetTileType(lvl, level.TileType(newValue))
		})
	values.RenderUnifiedSliderInt(readOnly, "Floor Height", floorHeightUnifier,
		func(u values.Unifier) int { return int(u.Unified().(level.TileHeightUnit)) },
		tileHeightFormatter,
		0, int(level.TileHeightUnitMax)-1,
		func(newValue int) {
			view.requestSetFloorHeight(lvl, level.TileHeightUnit(newValue))
		})
	values.RenderUnifiedSliderInt(readOnly, "Ceiling Height (abs)", ceilingHeightUnifier,
		func(u values.Unifier) int { return int(u.Unified().(level.TileHeightUnit)) },
		tileHeightFormatter,
		1, int(level.TileHeightUnitMax),
		func(newValue int) {
			view.requestSetCeilingHeight(lvl, level.TileHeightUnit(newValue))
		})
	values.RenderUnifiedSliderInt(readOnly, "Slope Height", slopeHeightUnifier,
		func(u values.Unifier) int { return int(u.Unified().(level.TileHeightUnit)) },
		tileHeightFormatter,
		0, int(level.TileHeightUnitMax)-1,
		func(newValue int) {
			view.requestSetSlopeHeight(lvl, level.TileHeightUnit(newValue))
		})
	slopeControls := level.TileSlopeControls()
	values.RenderUnifiedCombo(readOnly, "Slope Control", slopeControlUnifier,
		func(u values.Unifier) int { return int(u.Unified().(level.TileSlopeControl)) },
		func(value int) string { return slopeControls[value].String() },
		len(slopeControls),
		func(newValue int) {
			view.requestSetSlopeControl(lvl, slopeControls[newValue])
		})
	values.RenderUnifiedSliderInt(readOnly, "Music Index", musicIndexUnifier,
		func(u values.Unifier) int { return u.Unified().(int) },
		func(value int) string { return "%d" },
		0, 15,
		func(newValue int) {
			view.requestMusicIndex(lvl, newValue)
		})

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
			0, 255,
			func(newValue int) {
				view.requestFloorPaletteIndex(lvl, newValue)
			})
		values.RenderUnifiedSliderInt(readOnly, "Ceiling Color", ceilingPaletteIndexUnifier,
			func(u values.Unifier) int { return int(u.Unified().(byte)) },
			func(value int) string { return "%d" },
			0, 255,
			func(newValue int) {
				view.requestCeilingPaletteIndex(lvl, newValue)
			})

		imgui.Separator()

		flightPulls := level.CyberspaceFlightPulls()
		values.RenderUnifiedCombo(readOnly, "Flight Pull Type", flightPullTypeUnifier,
			func(u values.Unifier) int { return int(u.Unified().(level.CyberspaceFlightPull)) },
			func(value int) string { return flightPulls[value].String() },
			len(flightPulls),
			func(newValue int) {
				view.requestFlightPullType(lvl, flightPulls[newValue])
			})
		values.RenderUnifiedSliderInt(readOnly, "Game Of Life State", gameOfLightStateUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string { return "%d" },
			0, 3,
			func(newValue int) {
				view.requestGameOfLightState(lvl, newValue)
			})
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
			func(newValue int) {
				view.requestFloorTextureIndex(lvl, newValue)
			})
		view.renderTextureSelector(readOnly, "Floor Texture", floorTextureIndexUnifier, atlas, 0, level.FloorCeilingTextureLimit-1,
			func(newValue int) {
				view.requestFloorTextureIndex(lvl, newValue)
			})
		values.RenderUnifiedSliderInt(readOnly, "Floor Texture Rotations", floorTextureRotationsUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string { return "%d" },
			0, 3,
			func(newValue int) {
				view.requestFloorTextureRotations(lvl, newValue)
			})

		values.RenderUnifiedSliderInt(readOnly, "Ceiling Texture (atlas index)", ceilingTextureIndexUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string { return "%d" },
			0, level.FloorCeilingTextureLimit-1,
			func(newValue int) {
				view.requestCeilingTextureIndex(lvl, newValue)
			})
		view.renderTextureSelector(readOnly, "Ceiling Texture", ceilingTextureIndexUnifier, atlas, 0, level.FloorCeilingTextureLimit-1,
			func(newValue int) {
				view.requestCeilingTextureIndex(lvl, newValue)
			})
		values.RenderUnifiedSliderInt(readOnly, "Ceiling Texture Rotations", ceilingTextureRotationsUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string { return "%d" },
			0, 3,
			func(newValue int) {
				view.requestCeilingTextureRotations(lvl, newValue)
			})

		values.RenderUnifiedSliderInt(readOnly, "Wall Texture (atlas index)", wallTextureIndexUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string { return "%d" },
			0, len(atlas)-1,
			func(newValue int) {
				view.requestWallTextureIndex(lvl, newValue)
			})
		view.renderTextureSelector(readOnly, "Wall Texture", wallTextureIndexUnifier, atlas, 0, len(atlas)-1,
			func(newValue int) {
				view.requestWallTextureIndex(lvl, newValue)
			})
		values.RenderUnifiedSliderInt(readOnly, "Wall Texture Offset", wallTextureOffsetUnifier,
			func(u values.Unifier) int { return int(u.Unified().(level.TileHeightUnit)) },
			tileHeightFormatter,
			0, int(level.TileHeightUnitMax)-1,
			func(newValue int) {
				view.requestWallTextureOffset(lvl, level.TileHeightUnit(newValue))
			})

		values.RenderUnifiedCheckboxCombo(readOnly, "Use Adjacent Wall Texture", useAdjacentWallTextureUnifier,
			func(newValue bool) {
				view.requestUseAdjacentWallTexture(lvl, newValue)
			})
		wallTexturePatterns := level.WallTexturePatterns()
		values.RenderUnifiedCombo(readOnly, "Wall Texture Pattern", wallTexturePatternUnifier,
			func(u values.Unifier) int { return int(u.Unified().(level.WallTexturePattern)) },
			func(value int) string { return wallTexturePatterns[value].String() },
			len(wallTexturePatterns),
			func(newValue int) {
				view.requestWallTexturePattern(lvl, wallTexturePatterns[newValue])
			})

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
			0, 15,
			func(newValue int) {
				view.requestFloorLight(lvl, newValue)
			})
		values.RenderUnifiedSliderInt(readOnly, "Ceiling Light", ceilingLightUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string { return "%d" },
			0, 15,
			func(newValue int) {
				view.requestCeilingLight(lvl, newValue)
			})

		imgui.Separator()

		values.RenderUnifiedCheckboxCombo(readOnly, "Deconstructed", deconstructedUnifier,
			func(newValue bool) {
				view.requestDeconstructed(lvl, newValue)
			})
		values.RenderUnifiedCheckboxCombo(readOnly, "Floor Hazard", floorHazardUnifier,
			func(newValue bool) {
				view.requestFloorHazard(lvl, newValue)
			})
		values.RenderUnifiedCheckboxCombo(readOnly, "Ceiling Hazard", ceilingHazardUnifier,
			func(newValue bool) {
				view.requestCeilingHazard(lvl, newValue)
			})
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

func (view *TilesView) editingAllowed(id int) bool {
	moddedLevel := len(view.mod.ModifiedBlocks(resource.LangAny, ids.LevelResourcesStart.Plus(lvlids.PerLevel*id+lvlids.FirstUsed))) > 0
	return moddedLevel
}

func (view *TilesView) requestSetTileType(lvl *level.Level, tileType level.TileType) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Type = tileType
	})
}

func (view *TilesView) requestSetFloorHeight(lvl *level.Level, height level.TileHeightUnit) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Floor = tile.Floor.WithAbsoluteHeight(height)
	})
}

func (view *TilesView) requestSetCeilingHeight(lvl *level.Level, height level.TileHeightUnit) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Ceiling = tile.Ceiling.WithAbsoluteHeight(height)
	})
}

func (view *TilesView) requestSetSlopeHeight(lvl *level.Level, height level.TileHeightUnit) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.SlopeHeight = height
	})
}

func (view *TilesView) requestSetSlopeControl(lvl *level.Level, value level.TileSlopeControl) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.WithSlopeControl(value)
	})
}

func (view *TilesView) requestMusicIndex(lvl *level.Level, value int) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.WithMusicIndex(value)
	})
}

func (view *TilesView) requestFloorTextureIndex(lvl *level.Level, value int) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.TextureInfo = tile.TextureInfo.WithFloorTextureIndex(level.AtlasIndex(value))
	})
}

func (view *TilesView) requestFloorTextureRotations(lvl *level.Level, value int) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Floor = tile.Floor.WithTextureRotations(value)
	})
}

func (view *TilesView) requestCeilingTextureIndex(lvl *level.Level, value int) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.TextureInfo = tile.TextureInfo.WithCeilingTextureIndex(level.AtlasIndex(value))
	})
}

func (view *TilesView) requestCeilingTextureRotations(lvl *level.Level, value int) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Ceiling = tile.Ceiling.WithTextureRotations(value)
	})
}

func (view *TilesView) requestWallTextureIndex(lvl *level.Level, value int) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.TextureInfo = tile.TextureInfo.WithWallTextureIndex(level.AtlasIndex(value))
	})
}

func (view *TilesView) requestWallTextureOffset(lvl *level.Level, value level.TileHeightUnit) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForRealWorld().WithWallTextureOffset(value).AsTileFlag()
	})
}

func (view *TilesView) requestUseAdjacentWallTexture(lvl *level.Level, value bool) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForRealWorld().WithUseAdjacentWallTexture(value).AsTileFlag()
	})
}

func (view *TilesView) requestWallTexturePattern(lvl *level.Level, value level.WallTexturePattern) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForRealWorld().WithWallTexturePattern(value).AsTileFlag()
	})
}

func (view *TilesView) requestFloorLight(lvl *level.Level, value int) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForRealWorld().WithFloorShadow(15 - value).AsTileFlag()
	})
}

func (view *TilesView) requestCeilingLight(lvl *level.Level, value int) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForRealWorld().WithCeilingShadow(15 - value).AsTileFlag()
	})
}

func (view *TilesView) requestDeconstructed(lvl *level.Level, value bool) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForRealWorld().WithDeconstructed(value).AsTileFlag()
	})
}

func (view *TilesView) requestFloorHazard(lvl *level.Level, value bool) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Floor = tile.Floor.WithHazard(value)
	})
}

func (view *TilesView) requestCeilingHazard(lvl *level.Level, value bool) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Ceiling = tile.Ceiling.WithHazard(value)
	})
}

func (view *TilesView) requestFloorPaletteIndex(lvl *level.Level, value int) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.TextureInfo = tile.TextureInfo.WithFloorPaletteIndex(byte(value))
	})
}

func (view *TilesView) requestCeilingPaletteIndex(lvl *level.Level, value int) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.TextureInfo = tile.TextureInfo.WithCeilingPaletteIndex(byte(value))
	})
}

func (view *TilesView) requestFlightPullType(lvl *level.Level, value level.CyberspaceFlightPull) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForCyberspace().WithFlightPull(value).AsTileFlag()
	})
}

func (view *TilesView) requestGameOfLightState(lvl *level.Level, value int) {
	view.changeTiles(lvl, func(tile *level.TileMapEntry) {
		tile.Flags = tile.Flags.ForCyberspace().WithGameOfLifeState(value).AsTileFlag()
	})
}

func (view *TilesView) changeTiles(lvl *level.Level, modifier func(*level.TileMapEntry)) {
	positions := view.model.selectedTiles.list
	for _, pos := range positions {
		tile := lvl.Tile(pos.Tile())
		modifier(tile)
	}

	oldLevelID := view.levelSelection.CurrentLevelID()

	err := view.registry.Register(cmd.Named("PatchLevel"),
		cmd.Forward(view.restoreFocusTask()),
		cmd.Forward(view.levelSelection.SetCurrentLevelIDTask(lvl.ID())),
		cmd.Forward(view.setSelectedTilesTask(positions)),
		cmd.Nested(func() error { return view.levels.CommitLevelChanges(lvl.ID()) }),
		cmd.Reverse(view.setSelectedTilesTask(positions)),
		cmd.Reverse(view.levelSelection.SetCurrentLevelIDTask(oldLevelID)),
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

func (view *TilesView) setSelectedTilesTask(positions []MapPosition) cmd.Task {
	return func(modder world.Modder) error {
		view.setSelectedTiles(positions)
		// TODO replace with level selection
		// view.levelSelection.SetCurrentSelectedTiles(positions)
		return nil
	}
}

func (view *TilesView) setSelectedTiles(positions []MapPosition) {
	view.eventListener.Event(TileSelectionSetEvent{tiles: positions})
}
