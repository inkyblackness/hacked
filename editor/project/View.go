package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/inkyblackness/imgui-go/v2"

	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ui/gui"
)

// View handles the project display.
type View struct {
	service *edit.ProjectService

	modalStateMachine gui.ModalStateMachine
	guiScale          float32
	commander         cmd.Registry

	model viewModel
}

// NewView creates a new instance for the project display.
func NewView(service *edit.ProjectService, modalStateMachine gui.ModalStateMachine, guiScale float32, commander cmd.Registry) *View {
	return &View{
		service: service,

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
	saveStatus := view.service.SaveStatus()
	if saveStatus.FilesModified > 0 {
		title += fmt.Sprintf(" - %d file(s) pending save", saveStatus.FilesModified)
		if saveStatus.SavePending && (saveStatus.SaveIn.Seconds() < 4) {
			title += fmt.Sprintf(" - auto-save in %d", int(saveStatus.SaveIn.Seconds()))
		}
	}
	if view.model.windowOpen {
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 300 * view.guiScale}, imgui.ConditionFirstUseEver)
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
	modPath := view.service.ModPath()
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
	manifest := view.service.Mod().World()
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

// NewProject resets the current project and prepares for a new one.
func (view *View) NewProject() {
	view.service.ResetProject()
}

// LoadProject requests to load a stored project from storage.
func (view *View) LoadProject() {
	importProjectFile(view.modalStateMachine, func(filename string, settings edit.ProjectSettings) {
		view.service.RestoreProject(settings)
	})
}

const settingsFileExtension = "hacked-project"

func projectFileTypes() []external.TypeInfo {
	return []external.TypeInfo{
		{
			Title:      "HackEd Project File (*." + settingsFileExtension + ")",
			Extensions: []string{settingsFileExtension},
		},
	}
}

func importProjectFile(machine gui.ModalStateMachine, callback func(string, edit.ProjectSettings)) {
	fileHandler := func(filename string) error {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return errors.New("could not load file")
		}
		var settings edit.ProjectSettings
		err = json.Unmarshal(data, &settings)
		if err != nil {
			return errors.New("could not read file")
		}
		callback(filename, settings)
		return nil
	}

	external.LoadFile(machine, projectFileTypes(), fileHandler)
}

// SaveProject requests to save the project to storage.
func (view *View) SaveProject() {
	external.SaveFile(view.modalStateMachine, projectFileTypes(), func(filename string) error {
		completeFilename := filename
		dotExtension := "." + settingsFileExtension
		if !strings.HasSuffix(completeFilename, dotExtension) {
			completeFilename += dotExtension
		}
		return view.saveCurrentSettingsAs(completeFilename)
	})
}

func (view *View) saveCurrentSettingsAs(filename string) error {
	settings := view.service.CurrentSettings()
	data, err := json.Marshal(settings)
	if err != nil {
		return errors.New("could not encode settings")
	}
	err = ioutil.WriteFile(filename, data, 0640)
	if err != nil {
		return errors.New("could not write file")
	}
	return nil
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
	if view.service.ModHasStorageLocation() {
		err := view.service.SaveMod()
		if err != nil {
			view.handleSaveModFailure(err)
		}
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
	manifest := view.service.Mod().World()
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
	_ = view.commander.Register(
		cmd.Named("moveManifestEntry"),
		cmd.Reverse(view.taskToRestoreFocus()),
		cmd.Reverse(view.taskToSelectEntry(from)),
		cmd.Nested(func() error { return view.service.MoveManifestEntry(to, from) }),
		cmd.Forward(view.taskToSelectEntry(to)),
		cmd.Forward(view.taskToRestoreFocus()),
	)
}

func (view *View) tryAddManifestEntryFrom(names []string) error {
	entry, err := world.NewManifestEntryFrom(names)
	if err != nil {
		return err
	}

	view.requestAddManifestEntry(entry)
	return nil
}

func (view *View) requestAddManifestEntry(entry *world.ManifestEntry) {
	at := view.model.selectedManifestEntry + 1
	_ = view.commander.Register(
		cmd.Named("addManifestEntry"),
		cmd.Reverse(view.taskToRestoreFocus()),
		cmd.Reverse(view.taskToSelectEntry(view.model.selectedManifestEntry)),
		cmd.Nested(func() error { return view.service.AddManifestEntry(at, entry) }),
		cmd.Forward(view.taskToSelectEntry(at)),
		cmd.Forward(view.taskToRestoreFocus()),
	)
}

func (view *View) requestRemoveManifestEntry() {
	manifest := view.service.Mod().World()
	at := view.model.selectedManifestEntry
	if (at < 0) || (at >= manifest.EntryCount()) {
		return
	}

	_ = view.commander.Register(
		cmd.Named("removeManifestEntry"),
		cmd.Reverse(view.taskToRestoreFocus()),
		cmd.Reverse(view.taskToSelectEntry(view.model.selectedManifestEntry)),
		cmd.Nested(func() error { return view.service.RemoveManifestEntry(at) }),
		cmd.Forward(view.taskToSelectEntry(view.model.selectedManifestEntry-1)),
		cmd.Forward(view.taskToRestoreFocus()),
	)
}

func (view *View) tryLoadModFrom(names []string) error {
	return view.service.TryLoadModFrom(names)
}

func (view *View) requestSaveMod(modPath string) {
	err := view.service.SaveModUnder(modPath)
	if err != nil {
		view.handleSaveModFailure(err)
	}
}

func (view *View) handleSaveModFailure(err error) {
	view.modalStateMachine.SetState(&saveModFailedState{
		machine:   view.modalStateMachine,
		view:      view,
		errorInfo: err.Error(),
	})
}

func (view *View) taskToRestoreFocus() cmd.Task {
	return func(modder world.Modder) error {
		view.model.restoreFocus = true
		return nil
	}
}

func (view *View) taskToSelectEntry(index int) cmd.Task {
	return func(modder world.Modder) error {
		view.model.selectedManifestEntry = index
		return nil
	}
}
