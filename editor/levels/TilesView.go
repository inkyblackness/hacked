package levels

import (
	"bytes"
	"fmt"

	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/event"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial/rle"
	"github.com/inkyblackness/hacked/ss1/world/ids"
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
	tileTypeUnifier := values.NewUnifier()
	multiple := len(view.model.selectedTiles.list) > 0
	for _, pos := range view.model.selectedTiles.list {
		tile := lvl.Tile(int(pos.X.Tile()), int(pos.Y.Tile()))
		tileTypeUnifier.Add(tile.Type)
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
			oldData := view.mod.ModifiedBlock(resource.LangAny, resourceID, 0)
			if !bytes.Equal(oldData, newData) {
				forwardPatch := bytes.NewBuffer(nil)
				rle.Compress(forwardPatch, newData, oldData)
				command.forwardPatches = append(command.forwardPatches, patchEntry{resourceID, forwardPatch.Bytes()})

				reversePatch := bytes.NewBuffer(nil)
				rle.Compress(reversePatch, oldData, newData)
				command.reversePatches = append(command.reversePatches, patchEntry{resourceID, reversePatch.Bytes()})
			}
		}
	}

	view.commander.Queue(command)
}

func (view *TilesView) setSelectedTiles(positions []MapPosition) {
	view.eventListener.Event(TileSelectionSetEvent{tiles: positions})
}
