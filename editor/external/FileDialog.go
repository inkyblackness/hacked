package external

import (
	"github.com/inkyblackness/imgui-go/v3"
	"github.com/sqweek/dialog"

	"github.com/inkyblackness/hacked/ui/gui"
)

// SingleFileHandler is a callback for a single selected file.
type SingleFileHandler func(string) error

type singleFileRunner func(*dialog.FileBuilder) (string, error)

// SaveFile starts a dialog series to save a file.
func SaveFile(machine gui.ModalStateMachine, types []TypeInfo, callback SingleFileHandler) {
	machine.SetState(&fileStartState{
		machine:  machine,
		title:    "Save",
		typeInfo: types,
		runner:   func(builder *dialog.FileBuilder) (string, error) { return builder.Save() },
		callback: callback,
	})
}

// LoadFile starts a dialog series to load a file.
func LoadFile(machine gui.ModalStateMachine, types []TypeInfo, callback SingleFileHandler) {
	machine.SetState(&fileStartState{
		machine:  machine,
		title:    "Load",
		typeInfo: types,
		runner:   func(builder *dialog.FileBuilder) (string, error) { return builder.Load() },
		callback: callback,
	})
}

type fileStartState struct {
	machine  gui.ModalStateMachine
	title    string
	typeInfo []TypeInfo
	runner   singleFileRunner
	callback SingleFileHandler
}

func (state fileStartState) HandleFiles(names []string) {
}

func (state fileStartState) Render() {
	imgui.OpenPopup(state.title)
	nextState := &fileWaitingState{
		machine:  state.machine,
		title:    state.title,
		typeInfo: state.typeInfo,
		runner:   state.runner,
		callback: state.callback,
	}
	state.machine.SetState(nextState)
}

type fileWaitingState struct {
	machine  gui.ModalStateMachine
	title    string
	typeInfo []TypeInfo
	runner   singleFileRunner
	callback SingleFileHandler

	renderCount      int
	shouldOpenDialog bool
	problem          string
}

func (state *fileWaitingState) HandleFiles(names []string) {
}

func (state *fileWaitingState) Render() {
	if imgui.BeginPopupModalV(state.title, nil,
		imgui.WindowFlagsNoResize|imgui.WindowFlagsNoMove|imgui.WindowFlagsNoSavedSettings|imgui.WindowFlagsAlwaysAutoResize) {
		if len(state.problem) == 0 {
			imgui.Text("Please use popup dialog...")
		} else {
			imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
			imgui.Text(state.problem)
			imgui.PopStyleColor()
		}
		state.renderControls()
		imgui.EndPopup()

		state.renderCount++
		if state.renderCount == 3 {
			state.shouldOpenDialog = true
		}

		if state.shouldOpenDialog {
			state.shouldOpenDialog = false
			state.openDialog()
		}
	} else {
		state.machine.SetState(nil)
	}
}

func (state *fileWaitingState) renderControls() {
	imgui.Separator()
	if imgui.Button("Browse...") {
		state.shouldOpenDialog = true
	}
	imgui.SameLine()
	if imgui.Button("Cancel") {
		state.machine.SetState(nil)
		imgui.CloseCurrentPopup()
	}
}

func (state *fileWaitingState) openDialog() {
	dlgBuilder := dialog.File()
	for _, info := range state.typeInfo {
		dlgBuilder = dlgBuilder.Filter(info.Title, info.Extensions...)
	}
	dlgBuilder = dlgBuilder.Filter("All files (*.*)", "*")

	filename, err := state.runner(dlgBuilder)
	state.machine.SetState(nil)
	if err != nil {
		return
	}
	state.reportFile(filename)
}

func (state *fileWaitingState) reportFile(filename string) {
	state.problem = ""

	err := state.callback(filename)
	if err != nil {
		state.problem = "Previous attempt failed.\nPlease check and try again.\nReason:\n" + err.Error()
		state.machine.SetState(state)
	}
}
