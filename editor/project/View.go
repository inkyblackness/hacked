package project

import (
	"fmt"
	"math"
	"time"

	"github.com/inkyblackness/imgui-go/v2"

	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
)

// View handles the project display.
type View struct {
	mod *world.Mod

	modalStateMachine gui.ModalStateMachine
	guiScale          float32
	commander         cmd.Commander

	model viewModel
}

// NewView creates a new instance for the project display.
func NewView(mod *world.Mod, modalStateMachine gui.ModalStateMachine,
	guiScale float32, commander cmd.Commander) *View {
	return &View{
		mod: mod,

		modalStateMachine: modalStateMachine,
		guiScale:          guiScale,
		commander:         commander,

		model: freshViewModel(),
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
	title := "Project"
	changedFiles := len(view.mod.ModifiedFilenames())
	if changedFiles > 0 {
		title += fmt.Sprintf(" - %d file(s) pending save", changedFiles)
		lastChangeTime := view.mod.LastChangeTime()

		if (len(view.mod.Path()) > 0) && !lastChangeTime.IsZero() {
			saveAt := lastChangeTime.Add(time.Duration(view.model.autosaveTimeoutSec) * time.Second)
			autoSaveIn := time.Until(saveAt)
			if autoSaveIn.Seconds() < 4 {
				title += fmt.Sprintf(" - auto-save in %d", int(math.Max(autoSaveIn.Seconds(), 0.0)))
			}
			if autoSaveIn.Seconds() <= 0 {
				view.mod.ResetLastChangeTime()
				view.StartSavingMod()
			}
		}
	}
	if view.model.windowOpen {
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 300 * view.guiScale}, imgui.ConditionOnce)
		if imgui.BeginV(title+"###Project", view.WindowOpen(), 0) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *View) renderContent() {
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
		view.StartSavingMod()
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

func (view *View) startLoadingMod() {
	view.modalStateMachine.SetState(&loadModStartState{
		machine: view.modalStateMachine,
		view:    view,
	})
}

// StartSavingMod initiates to save the mod.
// It either opens the save-as dialog, or simply saves under the current folder.
func (view *View) StartSavingMod() {
	modPath := view.mod.Path()
	if len(modPath) > 0 {
		view.requestSaveMod(modPath)
	} else {
		view.modalStateMachine.SetState(&saveModAsStartState{
			machine: view.modalStateMachine,
			view:    view,
		})
	}
}

func (view *View) startAddingManifestEntry() {
	view.modalStateMachine.SetState(&addManifestEntryStartState{
		machine: view.modalStateMachine,
		view:    view,
	})
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

func (view *View) tryLoadModFrom(names []string) error {
	staging := newFileStaging(false)

	staging.stageAll(names)

	resourcesToTake := staging.resources
	isSavegame := false
	if (len(resourcesToTake) == 0) && (len(staging.savegames) == 1) {
		resourcesToTake = staging.savegames
		isSavegame = true
	}
	if len(resourcesToTake) == 0 {
		return fmt.Errorf("no resources found")
	}
	var locs []*world.LocalizedResources
	modPath := ""

	for location := range resourcesToTake {
		if (len(modPath) == 0) || (len(location.DirPath) < len(modPath)) {
			modPath = location.DirPath
		}
	}

	for location, viewer := range resourcesToTake {
		lang := ids.LocalizeFilename(location.Name)
		template := location.Name
		if isSavegame {
			template = string(ids.Archive)
		}
		loc := &world.LocalizedResources{
			File:     location,
			Template: template,
			Language: lang,
		}
		for _, id := range viewer.IDs() {
			view, err := viewer.View(id)
			if err == nil {
				_ = loc.Store.Put(id, view)
			}
			// TODO: handle error?
		}
		locs = append(locs, loc)
	}

	view.setActiveMod(modPath, locs, staging.objectProperties, staging.textureProperties)
	return nil
}

func (view *View) setActiveMod(modPath string, resources []*world.LocalizedResources,
	objectProperties object.PropertiesTable, textureProperties texture.PropertiesList) {
	view.mod.SetPath(modPath)
	view.mod.Reset(resources, objectProperties, textureProperties)
	// fix list resources for any "old" mod.
	view.mod.FixListResources()
}

func (view *View) requestSaveMod(modPath string) {
	view.mod.FixListResources()
	err := saveModResourcesTo(view.mod, modPath)
	if err != nil {
		view.modalStateMachine.SetState(&saveModFailedState{
			machine:   view.modalStateMachine,
			view:      view,
			errorInfo: err.Error(),
		})
	} else {
		view.mod.SetPath(modPath)
		view.mod.MarkSave()
	}
}
