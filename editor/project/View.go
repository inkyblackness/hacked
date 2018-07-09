package project

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/world"
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

// WindowOpen returns the flag address, to be used with the main menu.
func (view *View) WindowOpen() *bool {
	return &view.model.windowOpen
}

// Render requests to render the view.
func (view *View) Render() {
	if view.model.restoreFocus {
		imgui.SetNextWindowFocus()
		view.model.restoreFocus = false
		view.model.windowOpen = true
	}
	if view.model.windowOpen {
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 300 * view.guiScale}, imgui.ConditionOnce)
		if imgui.BeginV("Project", view.WindowOpen(), 0) {
			imgui.Text("Mod Location")
			imgui.PushStyleVarVec2(imgui.StyleVarWindowPadding, imgui.Vec2{X: 1, Y: 0})
			imgui.BeginChildV("ModLocation", imgui.Vec2{X: -200*view.guiScale - 10*view.guiScale, Y: imgui.TextLineHeight() * 1.5}, true,
				imgui.WindowFlagsNoScrollbar|imgui.WindowFlagsNoScrollWithMouse)
			modPath := view.mod.Path()
			if len(modPath) > 0 {
				imgui.Text(modPath)
			} else {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1.0, Y: 1.0, Z: 1.0, W: 0.5})
				imgui.Text("(new mod)")
				imgui.PopStyleColor()
			}
			imgui.EndChild()
			imgui.PopStyleVar()
			imgui.BeginGroup()
			imgui.SameLine()
			if imgui.ButtonV("Save", imgui.Vec2{X: 100 * view.guiScale, Y: 0}) {
				view.startSavingMod()
			}
			imgui.SameLine()
			if imgui.ButtonV("Load...", imgui.Vec2{X: 100 * view.guiScale, Y: 0}) {
				view.startLoadingMod()
			}
			imgui.EndGroup()

			imgui.Text("Static World Data")
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
	}

	view.fileState.Render()
}

// HandleFiles is called when the user wants to add the given files to the library.
func (view *View) HandleFiles(names []string) {
	view.fileState.HandleFiles(names)
}

func (view *View) startLoadingMod() {
	view.fileState = &loadModStartState{
		view: view,
	}
}

func (view *View) startSavingMod() {
	modPath := view.mod.Path()
	if len(modPath) > 0 {
		view.requestSaveMod(modPath)
	} else {
		view.fileState = &saveModAsStartState{
			view: view,
		}
	}
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

func (view *View) requestAddManifestEntry(entry *world.ManifestEntry) {
	at := view.model.selectedManifestEntry + 1
	command := listManifestEntryCommand{
		keeper: view.mod.World(),
		model:  &view.model,

		at:    at,
		entry: entry,
		adder: true,
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
		adder: false,
	}
	view.commander.Queue(command)
}

func (view *View) requestLoadMod(modPath string, resources model.LocalizedResources, objectProperties object.PropertiesTable) {
	view.mod.SetPath(modPath)
	view.mod.Reset(resources, objectProperties)
}

func (view *View) requestSaveMod(modPath string) {
	err := saveModResourcesTo(view.mod.ModifiedResources(), modPath)
	if err != nil {
		view.fileState = &saveModFailedState{
			view: view,
		}
	} else {
		view.mod.SetPath(modPath)
	}
}
