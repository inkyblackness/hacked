package archives

import (
	"fmt"

	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/imgui-go"
)

// View provides edit controls for the archive.
type View struct {
	mod *model.Mod

	guiScale  float32
	commander cmd.Commander

	model viewModel
}

// NewArchiveView returns a new instance.
func NewArchiveView(mod *model.Mod, guiScale float32, commander cmd.Commander) *View {
	view := &View{
		mod: mod,

		guiScale:  guiScale,
		commander: commander,

		model: freshViewModel(),
	}
	return view
}

// WindowOpen returns the flag address, to be used with the main menu.
func (view *View) WindowOpen() *bool {
	return &view.model.windowOpen
}

// Render renders the view.
func (view *View) Render() {
	if view.model.restoreFocus {
		imgui.SetNextWindowFocus()
		view.model.restoreFocus = false
		view.model.windowOpen = true
	}
	if view.model.windowOpen {
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 350 * view.guiScale, Y: 400 * view.guiScale}, imgui.ConditionOnce)
		if imgui.BeginV("Archive", view.WindowOpen(), imgui.WindowFlagsNoCollapse) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *View) renderContent() {
	imgui.Text("Levels")
	imgui.BeginChildV("Levels", imgui.Vec2{X: -100 * view.guiScale, Y: 0}, true, 0)
	for id := 0; id < archive.MaxLevels; id++ {
		inMod := view.hasLevelInMod(id)
		info := fmt.Sprintf("%d", id)
		if !inMod {
			if (id == world.StartingLevel) && view.hasGameStateInMod() {
				// Starting level should really be in a new archive.
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1.0, Y: 0.0, Z: 0.0, W: 1.0})
			} else {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1.0, Y: 1.0, Z: 1.0, W: 0.8})
			}
			info += " (not in mod, read-only)"
		}
		if imgui.SelectableV(info, id == view.model.selectedLevel, 0, imgui.Vec2{}) {
			view.model.selectedLevel = id
		}
		if !inMod {
			imgui.PopStyleColor()
		}
	}
	imgui.EndChild()
	imgui.SameLine()
	imgui.BeginGroup()
	if imgui.ButtonV("Clear", imgui.Vec2{X: -1, Y: 0}) {
		view.requestClearLevel(view.model.selectedLevel)
	}
	if view.model.selectedLevel != world.StartingLevel {
		if imgui.ButtonV("Remove", imgui.Vec2{X: -1, Y: 0}) {
			view.requestRemoveLevel(view.model.selectedLevel)
		}
	}
	imgui.EndGroup()
}

func (view *View) hasGameStateInMod() bool {
	return len(view.mod.ModifiedBlocks(resource.LangAny, ids.GameState)) > 0
}

func (view *View) hasLevelInMod(id int) bool {
	return len(view.mod.ModifiedBlocks(resource.LangAny, ids.LevelResourcesStart.Plus(lvlids.PerLevel*id+lvlids.FirstUsed))) > 0
}

func (view *View) requestClearLevel(id int) {
	if (id >= 0) && (id < archive.MaxLevels) {
		command := setArchiveDataCommand{
			model:         &view.model,
			selectedLevel: id,
			newData:       make(map[resource.ID][]byte),
			oldData:       make(map[resource.ID][]byte),
		}

		if !view.hasGameStateInMod() {
			command.newData[ids.ArchiveName] = text.DefaultCodepage().Encode("Starting Game | by InkyBlackness HackEd")
			command.newData[ids.GameState] = make([]byte, archive.GameStateSize)
		}

		param := level.EmptyLevelParameters{
			Cyberspace:  world.IsConsideredCyberspaceByDefault(id),
			MapModifier: func(level.TileMap) {},
		}
		if id == world.StartingLevel {
			param.MapModifier = func(m level.TileMap) {
				m.Tile(world.StartingTileX, world.StartingTileY).Type = level.TileTypeOpen
			}
		}

		levelData := level.EmptyLevelData(param)
		levelIDBegin := ids.LevelResourcesStart.Plus(lvlids.PerLevel * id)
		for offset, newData := range &levelData {
			resourceID := levelIDBegin.Plus(offset)
			oldData := view.mod.ModifiedBlock(resource.LangAny, resourceID, 0)
			if len(oldData) > 0 {
				command.oldData[resourceID] = oldData
			}
			if len(newData) > 0 {
				command.newData[resourceID] = newData
			}
		}

		view.commander.Queue(command)
	}
}

func (view *View) requestRemoveLevel(id int) {
	if (id >= 0) && (id < archive.MaxLevels) && view.hasLevelInMod(id) {
		command := setArchiveDataCommand{
			model:         &view.model,
			selectedLevel: id,
			newData:       make(map[resource.ID][]byte),
			oldData:       make(map[resource.ID][]byte),
		}

		levelIDBegin := ids.LevelResourcesStart.Plus(lvlids.PerLevel * id)
		for offset := lvlids.FirstUsed; offset < lvlids.PerLevel; offset++ {
			resourceID := levelIDBegin.Plus(offset)
			oldData := view.mod.ModifiedBlock(resource.LangAny, resourceID, 0)
			if len(oldData) > 0 {
				command.oldData[resourceID] = oldData
			}
		}

		view.commander.Queue(command)
	}
}
