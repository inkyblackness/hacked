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
	"github.com/inkyblackness/imgui-go"
)

// TilesView is for tile properties.
type TilesView struct {
	mod *model.Mod

	guiScale  float32
	commander cmd.Commander

	model tilesViewModel
}

// NewTilesView returns a new instance.
func NewTilesView(mod *model.Mod, guiScale float32, commander cmd.Commander, eventRegistry event.Registry) *TilesView {
	view := &TilesView{
		mod:       mod,
		guiScale:  guiScale,
		commander: commander,
		model:     freshTilesViewModel(),
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
		title := "Level Tiles"
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
	imgui.LabelText("Selected Tiles", fmt.Sprintf("%d", len(view.model.selectedTiles.list)))

	tileTypeUnifier := values.NewUnifier()
	multiple := len(view.model.selectedTiles.list) > 0
	for _, pos := range view.model.selectedTiles.list {
		tile := lvl.Tile(int(pos.X.Tile()), int(pos.Y.Tile()))
		tileTypeUnifier.Add(tile.Type)
	}

	var unifiedTileTypeString string
	if tileTypeUnifier.IsUnique() {
		unifiedTileTypeString = tileTypeUnifier.Unified().(level.TileType).String()
	} else if multiple {
		unifiedTileTypeString = "(multiple)"
	}
	if readOnly {
		imgui.LabelText("Tile Type", unifiedTileTypeString)
	} else {
		if imgui.BeginCombo("Tile Type", unifiedTileTypeString) {
			for _, tileType := range level.TileTypes() {
				tileTypeString := tileType.String()

				imgui.SelectableV(tileTypeString, tileTypeString == unifiedTileTypeString, 0, imgui.Vec2{})
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
