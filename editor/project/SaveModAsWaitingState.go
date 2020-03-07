package project

import (
	"os"
	"time"

	"github.com/inkyblackness/imgui-go"
	"github.com/sqweek/dialog"

	"github.com/inkyblackness/hacked/ui/gui"
)

type saveModAsWaitingState struct {
	machine     gui.ModalStateMachine
	view        *View
	failureTime time.Time
	errorInfo   string
}

func (state *saveModAsWaitingState) Render() {
	if imgui.BeginPopupModalV("Save mod as", nil,
		imgui.WindowFlagsNoResize|imgui.WindowFlagsNoMove|imgui.WindowFlagsNoSavedSettings|imgui.WindowFlagsAlwaysAutoResize) {
		imgui.Text("Waiting for folder.")
		if !state.failureTime.IsZero() {
			imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
			imgui.Text("Previous attempt failed, not a proper folder.\nPlease check and try again.")
			if len(state.errorInfo) > 0 {
				imgui.Text(state.errorInfo)
			}
			imgui.PopStyleColor()
			if time.Since(state.failureTime).Seconds() > 5 {
				state.failureTime = time.Time{}
			}
		}
		imgui.Text(`From your file browser drag'n'drop the folder
where to save the data into the editor window.
`)
		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
		imgui.Text("It is recommended to use an empty folder.\nSaving will potentially overwrite existing files.")
		imgui.PopStyleColor()
		imgui.Separator()
		if imgui.Button("Browse...") {
			dlgBuilder := dialog.Directory()
			filename, err := dlgBuilder.Browse()
			if err == nil {
				state.HandleFiles([]string{filename})
			}
		}
		imgui.SameLine()
		if imgui.Button("Cancel") {
			state.machine.SetState(nil)
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	} else {
		state.machine.SetState(nil)
	}
}

func (state *saveModAsWaitingState) HandleFiles(names []string) {
	modPath, ok := state.verifyDir(names)
	if ok {
		state.machine.SetState(nil)
		state.view.requestSaveMod(modPath)
	} else {
		state.failureTime = time.Now()
	}
}

func (state saveModAsWaitingState) verifyDir(names []string) (string, bool) {
	if len(names) != 1 {
		return "", false
	}
	fileInfo, err := os.Stat(names[0])
	if err != nil {
		return "", false
	}
	if !fileInfo.IsDir() {
		return "", false
	}
	return names[0], true
}
