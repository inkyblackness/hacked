package project

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/imgui-go"
)

// View handles the project display.
type View struct {
	mod       *model.Mod
	guiScale  float32
	commander cmd.Commander

	model viewModel

	fileState popupState
}

// NewView creates a new instance for the project display.
func NewView(mod *model.Mod, guiScale float32, commander cmd.Commander) *View {
	return &View{
		mod:       mod,
		guiScale:  guiScale,
		commander: commander,

		model:     freshViewModel(),
		fileState: &idlePopupState{},
	}
}

// Render requests to render the view.
func (view *View) Render() {
	if view.model.restoreFocus {
		imgui.SetNextWindowFocus()
		view.model.restoreFocus = false
	}
	imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 300 * view.guiScale}, imgui.ConditionOnce)
	if imgui.Begin("Project") {
		imgui.TextUnformatted("Mod Location")
		imgui.BeginChildV("ModLocation", imgui.Vec2{X: -200*view.guiScale - 10*view.guiScale, Y: imgui.TextLineHeight() * view.guiScale}, true,
			imgui.WindowFlagsNoScrollbar|imgui.WindowFlagsNoScrollWithMouse)
		imgui.TextUnformatted("some/long/path/that/should/be/cut/at/a/point")
		imgui.EndChild()
		imgui.BeginGroup()
		imgui.SameLine()
		imgui.ButtonV("Save", imgui.Vec2{X: 100 * view.guiScale, Y: 0})
		imgui.SameLine()
		imgui.ButtonV("Load...", imgui.Vec2{X: 100 * view.guiScale, Y: 0})
		imgui.EndGroup()

		imgui.TextUnformatted("Static World Data")
		imgui.BeginChildV("ManifestEntries", imgui.Vec2{X: -100 * view.guiScale, Y: 0}, true, 0)
		manifest := view.mod.World()
		entries := manifest.EntryCount()
		for i := entries - 1; i >= 0; i-- {
			entry, _ := manifest.Entry(i)
			if imgui.SelectableV(entry.ID, view.model.selectedManifestEntry == i, 0, imgui.Vec2{}) {
				view.model.selectedManifestEntry = i
			}
		}
		imgui.EndChild()
		imgui.SameLine()
		imgui.BeginGroup()
		if imgui.ButtonV("Add...", imgui.Vec2{X: -1, Y: 0}) {
			view.startAddingManifestEntry()
		}
		if imgui.ButtonV("Up", imgui.Vec2{X: -1, Y: 0}) {
			view.requestMoveManifestEntryUp()
		}
		if imgui.ButtonV("Down", imgui.Vec2{X: -1, Y: 0}) {
			view.requestMoveManifestEntryDown()
		}
		if imgui.ButtonV("Remove", imgui.Vec2{X: -1, Y: 0}) {
			view.requestRemoveManifestEntry()
		}
		imgui.EndGroup()
	}
	imgui.End()

	view.fileState.Render()
}

func (view *View) startAddingManifestEntry() {
	view.fileState = &addManifestEntryStartState{
		view: view,
	}
}

func (view *View) requestMoveManifestEntryUp() {
	manifest := view.mod.World()
	entries := manifest.EntryCount()
	if (view.model.selectedManifestEntry >= 0) && (view.model.selectedManifestEntry < (entries - 1)) {
		view.requestMoveManifestEntry(view.model.selectedManifestEntry+1, view.model.selectedManifestEntry)
	}
}

func (view *View) requestMoveManifestEntryDown() {
	if view.model.selectedManifestEntry > 0 {
		view.requestMoveManifestEntry(view.model.selectedManifestEntry-1, view.model.selectedManifestEntry)
	}
}

func (view *View) requestMoveManifestEntry(to, from int) {
	command := moveManifestEntryCommand{
		mover: view.mod.World(),
		model: &view.model,
		to:    to,
		from:  from,
	}
	view.commander.Queue(command)
}
func (view *View) requestRemoveManifestEntry() {
	manifest := view.mod.World()
	at := view.model.selectedManifestEntry
	if (at < 0) || (at >= manifest.EntryCount()) {
		return
	}
	entry, err := manifest.Entry(at)
	if err != nil {
		return
	}
	command := listManifestEntryCommand{
		keeper: view.mod.World(),
		model:  &view.model,

		at:    view.model.selectedManifestEntry,
		entry: entry,
	}
	view.commander.Queue(command)
}
