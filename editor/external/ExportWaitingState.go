package external

import (
	"os"
	"time"

	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
)

type exportWaitingState struct {
	machine  gui.ModalStateMachine
	callback func(string)
	info     string

	failureTime time.Time
}

func (state *exportWaitingState) Render() {
	if imgui.BeginPopupModalV("Export", nil,
		imgui.WindowFlagsNoResize|imgui.WindowFlagsNoMove|imgui.WindowFlagsNoSavedSettings|imgui.WindowFlagsAlwaysAutoResize) {

		imgui.Text("Waiting for folder.")
		if !state.failureTime.IsZero() {
			imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
			imgui.Text("Previous attempt failed, could not save file.\nPlease check and try again.")
			imgui.PopStyleColor()
			if time.Since(state.failureTime).Seconds() > 5 {
				state.failureTime = time.Time{}
			}
		}
		imgui.Text(`From your file browser drag'n'drop the folder
where to save the data into the editor window.
`)
		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
		imgui.Text("Saving will potentially overwrite existing files in the folder.")
		imgui.PopStyleColor()
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

func (state *exportWaitingState) HandleFiles(names []string) {
	dirPath, ok := state.verifyDir(names)
	if ok {
		state.machine.SetState(nil)
		state.callback(dirPath)
	} else {
		state.failureTime = time.Now()
	}
}

func (state exportWaitingState) verifyDir(names []string) (string, bool) {
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
