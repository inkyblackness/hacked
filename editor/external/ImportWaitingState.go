package external

import (
	"os"
	"time"

	"github.com/inkyblackness/imgui-go"

	"github.com/inkyblackness/hacked/ui/gui"
)

type importWaitingState struct {
	machine  gui.ModalStateMachine
	callback func(string)
	info     string

	failureTime time.Time
}

func (state *importWaitingState) Render() {
	if imgui.BeginPopupModalV("Import", nil,
		imgui.WindowFlagsNoResize|imgui.WindowFlagsNoMove|imgui.WindowFlagsNoSavedSettings|imgui.WindowFlagsAlwaysAutoResize) {

		imgui.Text("Waiting for file.")
		if !state.failureTime.IsZero() {
			imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
			imgui.Text("Previous attempt failed, could not load file.\nPlease check and try again.")
			imgui.PopStyleColor()
			if time.Since(state.failureTime).Seconds() > 5 {
				state.failureTime = time.Time{}
			}
		}
		imgui.Text(`From your file browser drag'n'drop the file
that shall be loaded into the editor window.
`)
		imgui.Text(state.info)
		imgui.Separator()
		if imgui.Button("Cancel") {
			state.machine.SetState(nil)
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	} else {
		state.machine.SetState(nil)
	}
}

func (state *importWaitingState) HandleFiles(names []string) {
	filename, ok := state.verifyFile(names)
	if ok {
		state.machine.SetState(nil)
		state.callback(filename)
	} else {
		state.failureTime = time.Now()
	}
}

func (state importWaitingState) verifyFile(names []string) (string, bool) {
	if len(names) != 1 {
		return "", false
	}
	fileInfo, err := os.Stat(names[0])
	if err != nil {
		return "", false
	}
	if fileInfo.IsDir() {
		return "", false
	}
	return names[0], true
}
