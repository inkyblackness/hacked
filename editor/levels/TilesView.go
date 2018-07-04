package levels

import (
	"fmt"

	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/event"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
)

// TilesView is for tile properties.
type TilesView struct {
	mod *model.Mod

	guiScale      float32
	commander     cmd.Commander
	eventListener event.Listener

	model tilesViewModel
}

// NewTilesView returns a new instance.
func NewTilesView(mod *model.Mod, guiScale float32, commander cmd.Commander, eventListener event.Listener, eventRegistry event.Registry) *TilesView {
	view := &TilesView{
		mod:           mod,
		guiScale:      guiScale,
		commander:     commander,
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

// Render renders the view.
func (view *TilesView) Render(lvl *level.Level) {
	if view.model.restoreFocus {
		imgui.SetNextWindowFocus()
		view.model.restoreFocus = false
		view.model.windowOpen = true
	}
	if view.model.windowOpen {
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 500 * view.guiScale}, imgui.ConditionOnce)
		title := fmt.Sprintf("Level Tiles, %d selected", len(view.model.selectedTiles.list))
		readOnly := !view.editingAllowed(lvl.ID())
		if readOnly {
			title += " (read-only)"
		}
		if imgui.BeginV(title+"###Level Tiles", view.WindowOpen(), 0) {
			view.renderContent(lvl, readOnly)
		}
		imgui.End()
	}
}

func (view *TilesView) renderContent(lvl *level.Level, readOnly bool) {
	imgui.PushItemWidth(-250 * view.guiScale)
	tileTypeUnifier := values.NewUnifier()
	floorHeightUnifier := values.NewUnifier()
	ceilingHeightUnifier := values.NewUnifier()
	slopeHeightUnifier := values.NewUnifier()
	multiple := len(view.model.selectedTiles.list) > 0
	for _, pos := range view.model.selectedTiles.list {
		tile := lvl.Tile(int(pos.X.Tile()), int(pos.Y.Tile()))
		tileTypeUnifier.Add(tile.Type)
		floorHeightUnifier.Add(tile.Floor.AbsoluteHeight())
		ceilingHeightUnifier.Add(tile.Ceiling.AbsoluteHeight())
		slopeHeightUnifier.Add(tile.SlopeHeight)
	}

	var selectedTileTypeString string
	if tileTypeUnifier.IsUnique() {
		selectedTileTypeString = tileTypeUnifier.Unified().(level.TileType).String()
	} else if multiple {
		selectedTileTypeString = "(multiple)"
	}
	if readOnly {
		imgui.LabelText("Tile Type", selectedTileTypeString)
	} else {
		if imgui.BeginCombo("Tile Type", selectedTileTypeString) {
			for _, tileType := range level.TileTypes() {
				tileTypeString := tileType.String()

				if imgui.SelectableV(tileTypeString, tileTypeString == selectedTileTypeString, 0, imgui.Vec2{}) {
					view.requestSetTileType(lvl, view.model.selectedTiles.list, tileType)
				}
			}
			imgui.EndCombo()
		}
	}

	view.renderSliderInt(readOnly, multiple, "Floor Height", floorHeightUnifier,
		func(u values.Unifier) int { return int(u.Unified().(level.TileHeightUnit)) },
		func(value int) string { return "%d" },
		0, int(level.TileHeightUnitMax)-1,
		func(newValue int) {
			view.requestSetFloorHeight(lvl, view.model.selectedTiles.list, level.TileHeightUnit(newValue))
		})
	view.renderSliderInt(readOnly, multiple, "Ceiling Height (abs)", ceilingHeightUnifier,
		func(u values.Unifier) int { return int(u.Unified().(level.TileHeightUnit)) },
		func(value int) string { return "%d" },
		0, int(level.TileHeightUnitMax)-1,
		func(newValue int) {
			view.requestSetCeilingHeight(lvl, view.model.selectedTiles.list, level.TileHeightUnit(newValue))
		})
	view.renderSliderInt(readOnly, multiple, "Slope Height", slopeHeightUnifier,
		func(u values.Unifier) int { return int(u.Unified().(level.TileHeightUnit)) },
		func(value int) string { return "%d" },
		0, int(level.TileHeightUnitMax)-1,
		func(newValue int) {
			view.requestSetSlopeHeight(lvl, view.model.selectedTiles.list, level.TileHeightUnit(newValue))
		})

	imgui.PopItemWidth()
}

func (view *TilesView) renderSliderInt(readOnly, multiple bool, label string, unifier values.Unifier,
	intConverter func(values.Unifier) int, formatter func(int) string, min, max int, changeHandler func(int)) {

	var selectedString string
	selectedValue := -1
	if unifier.IsUnique() {
		selectedValue = intConverter(unifier)
		selectedString = formatter(selectedValue)
	} else if multiple {
		selectedString = "(multiple)"
	}
	if readOnly {
		imgui.LabelText(label, selectedString)
	} else {
		if gui.StepSliderIntV(label, &selectedValue, min, max, selectedString) {
			changeHandler(selectedValue)
		}
	}
}

func (view *TilesView) editingAllowed(id int) bool {
	gameStateData := view.mod.ModifiedBlocks(resource.LangAny, ids.GameState)
	isSavegame := (len(gameStateData) == 1) && (len(gameStateData[0]) == archive.GameStateSize) && (gameStateData[0][0x009C] > 0)
	moddedLevel := len(view.mod.ModifiedBlocks(resource.LangAny, ids.LevelResourcesStart.Plus(lvlids.PerLevel*id+lvlids.FirstUsed))) > 0

	return moddedLevel && !isSavegame
}

func (view *TilesView) requestSetTileType(lvl *level.Level, positions []MapPosition, tileType level.TileType) {
	view.changeTiles(lvl, positions, func(tile *level.TileMapEntry) {
		tile.Type = tileType
	})
}

func (view *TilesView) requestSetFloorHeight(lvl *level.Level, positions []MapPosition, height level.TileHeightUnit) {
	view.changeTiles(lvl, positions, func(tile *level.TileMapEntry) {
		tile.Floor = tile.Floor.WithAbsoluteHeight(height)
	})
}

func (view *TilesView) requestSetCeilingHeight(lvl *level.Level, positions []MapPosition, height level.TileHeightUnit) {
	view.changeTiles(lvl, positions, func(tile *level.TileMapEntry) {
		tile.Ceiling = tile.Ceiling.WithAbsoluteHeight(height)
	})
}

func (view *TilesView) requestSetSlopeHeight(lvl *level.Level, positions []MapPosition, height level.TileHeightUnit) {
	view.changeTiles(lvl, positions, func(tile *level.TileMapEntry) {
		tile.SlopeHeight = height
	})
}
func (view *TilesView) changeTiles(lvl *level.Level, positions []MapPosition, modifier func(*level.TileMapEntry)) {
	for _, pos := range positions {
		tile := lvl.Tile(int(pos.X.Tile()), int(pos.Y.Tile()))
		modifier(tile)
	}

	command := patchLevelDataCommand{
		restoreState: func() {
			view.model.restoreFocus = true
			view.setSelectedTiles(positions)
		},
	}

	newDataSet := lvl.EncodeState()
	for id, newData := range newDataSet {
		if len(newData) > 0 {
			resourceID := ids.LevelResourcesStart.Plus(lvlids.PerLevel*lvl.ID() + id)
			patch, changed, err := view.mod.CreateBlockPatch(resource.LangAny, resourceID, 0, newData)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				// TODO how to handle this? We're not expecting this, so crash and burn?
			} else if changed {
				command.patches = append(command.patches, patch)
			}
		}
	}

	view.commander.Queue(command)
}

func (view *TilesView) setSelectedTiles(positions []MapPosition) {
	view.eventListener.Event(TileSelectionSetEvent{tiles: positions})
}
