package archives

import (
	"fmt"
	"strings"

	"github.com/inkyblackness/imgui-go/v2"

	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// View provides edit controls for the archive.
type View struct {
	mod       *world.Mod
	textCache *text.Cache
	cp        text.Codepage

	guiScale  float32
	commander cmd.Commander

	model viewModel
}

// NewArchiveView returns a new instance.
func NewArchiveView(mod *world.Mod, textCache *text.Cache, cp text.Codepage, guiScale float32, commander cmd.Commander) *View {
	view := &View{
		mod:       mod,
		textCache: textCache,
		cp:        cp,

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
	imgui.BeginTabBar("archive-tab")

	if imgui.BeginTabItem("Levels") {
		view.renderLevelsContent()
		imgui.EndTabItem()
	}
	if imgui.BeginTabItem("Game State") {
		view.renderGameStateContent()
		imgui.EndTabItem()
	}

	imgui.EndTabBar()
}

func (view *View) renderLevelsContent() {
	imgui.BeginChildV("Levels", imgui.Vec2{X: -100 * view.guiScale, Y: 0}, true, 0)
	for id := 0; id < archive.MaxLevels; id++ {
		inMod := view.hasLevelInMod(id)
		info := fmt.Sprintf("%d", id)
		if !inMod {
			if view.hasGameStateInMod() && (id == view.effectiveGameState().CurrentLevel()) {
				// current level should really be in the archive.
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
	if imgui.ButtonV("Remove", imgui.Vec2{X: -1, Y: 0}) {
		view.requestRemoveLevel(view.model.selectedLevel)
	}
	imgui.EndGroup()
}

func (view *View) renderGameStateContent() {
	rawGameState := view.gameStateData()
	var state *archive.GameState
	if len(rawGameState) > 0 {
		state = archive.NewGameState(rawGameState)
	}

	imgui.PushItemWidth(-260 * view.guiScale)
	if (state != nil) && !state.IsSavegame() {
		resetText := "Override"
		if !state.IsDefaulting() {
			if imgui.Button("Remove") {
				view.requestSetGameState(archive.ZeroGameStateData())
			}
			imgui.SameLine()
			resetText = "Reset"
		}
		if imgui.Button(resetText) {
			view.requestSetGameState(archive.DefaultGameStateData())
		}
		if imgui.IsItemHovered() {
			imgui.BeginTooltip()
			imgui.SetTooltip("Values of new game archives are only considered by special engines.")
			imgui.EndTooltip()
		}
	}

	readOnly := false
	if (state == nil) || state.IsDefaulting() {
		state = archive.NewGameState(archive.DefaultGameStateData())
		readOnly = true
	}
	view.createPropertyControls(readOnly, state.Instance, func(key string, modifier func(uint32) uint32) {
		view.setInterpreterValueKeyed(state.Instance, key, modifier)
		view.requestSetGameState(state.Raw())
	})

	imgui.PopItemWidth()
}

func (view *View) gameStateData() []byte {
	raw := view.mod.ModifiedBlock(resource.LangAny, ids.GameState, 0)
	if len(raw) == 0 {
		return nil
	}
	copied := make([]byte, len(raw))
	copy(copied, raw)
	return copied
}

func (view *View) effectiveGameState() *archive.GameState {
	raw := view.gameStateData()
	if raw == nil {
		return archive.NewGameState(archive.DefaultGameStateData())
	}
	state := archive.NewGameState(raw)
	if state.IsDefaulting() {
		return archive.NewGameState(archive.DefaultGameStateData())
	}
	return state
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
			command.newData[ids.GameState] = archive.ZeroGameStateData()
		}

		param := level.EmptyLevelParameters{
			Cyberspace:  world.IsConsideredCyberspaceByDefault(id),
			MapModifier: func(level.TileMap) {},
		}
		effectiveGameState := view.effectiveGameState()
		if id == effectiveGameState.CurrentLevel() {
			posX, posY := effectiveGameState.HackerMapPosition()
			param.MapModifier = func(m level.TileMap) {
				m.Tile(int(posX.Tile()), int(posY.Tile())).Type = level.TileTypeOpen
			}
		}

		levelData := level.EmptyLevelData(param)
		levelIDBegin := ids.LevelResourcesStart.Plus(lvlids.PerLevel * id)
		for offset, newData := range &levelData {
			if (offset < lvlids.FirstUsed) || (offset >= lvlids.FirstUnused) {
				continue
			}
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

func (view *View) requestSetGameState(newData []byte) {
	command := setArchiveDataCommand{
		model:         &view.model,
		selectedLevel: view.model.selectedLevel,
		newData:       make(map[resource.ID][]byte),
		oldData:       make(map[resource.ID][]byte),
	}

	command.oldData[ids.GameState] = view.gameStateData()
	command.newData[ids.GameState] = newData

	view.commander.Queue(command)
}

func (view *View) createPropertyControls(readOnly bool, rootInterpreter *interpreters.Instance,
	updater func(string, func(uint32) uint32)) {
	objTypeRenderer := values.ObjectTypeControlRenderer{
		Meta:      view.mod.ObjectProperties(),
		TextCache: view.textCache,
	}

	var processInterpreter func(string, *interpreters.Instance)
	processInterpreter = func(path string, interpreter *interpreters.Instance) {
		for _, key := range interpreter.Keys() {
			fullKey := path + key
			unifier := values.NewUnifier()
			unifier.Add(int32(interpreter.Get(key)))
			simplifier := values.StandardSimplifier(readOnly, false, fullKey, unifier,
				func(modifier func(uint32) uint32) {
					updater(fullKey, modifier)
				}, objTypeRenderer)

			interpreter.Describe(key, simplifier)
		}

		for _, key := range interpreter.ActiveRefinements() {
			fullKey := path + key
			if len(fullKey) > 0 {
				imgui.Separator()
				imgui.Text(fullKey + ":")
			}
			processInterpreter(fullKey+".", interpreter.Refined(key))
		}
	}
	processInterpreter("", rootInterpreter)
}

func (view *View) setInterpreterValueKeyed(instance *interpreters.Instance, key string, modifier func(uint32) uint32) {
	resolvedInterpreter := instance
	keys := strings.Split(key, ".")
	keyCount := len(keys)
	if len(keys) > 1 {
		for _, subKey := range keys[:keyCount-1] {
			resolvedInterpreter = resolvedInterpreter.Refined(subKey)
		}
	}
	resolvedInterpreter.Set(keys[keyCount-1], modifier(resolvedInterpreter.Get(keys[keyCount-1])))
}
