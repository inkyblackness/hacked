package project

import (
	"time"

	"github.com/inkyblackness/imgui-go/v2"
	"github.com/sqweek/dialog"

	"github.com/inkyblackness/hacked/ui/gui"
)

type addManifestEntryWaitingState struct {
	machine     gui.ModalStateMachine
	view        *View
	failureTime time.Time
}

func (state *addManifestEntryWaitingState) Render() {
	if imgui.BeginPopupModalV("Add static world data", nil,
		imgui.WindowFlagsNoResize|imgui.WindowFlagsNoMove|imgui.WindowFlagsNoSavedSettings|imgui.WindowFlagsAlwaysAutoResize) {
		imgui.Text("Waiting for folders/files.")
		if !state.failureTime.IsZero() {
			imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
			imgui.Text("Previous attempt failed, no usable data detected.\nPlease check and try again.")
			imgui.PopStyleColor()
			if time.Since(state.failureTime).Seconds() > 5 {
				state.failureTime = time.Time{}
			}
		}
		imgui.Text(`From your file browser drag'n'drop the folder (or files)
of the static data you want to reference into the editor window.
Typically, you would use the main "data" directory of the game
(where all the .res files are).
`)
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

func (state *addManifestEntryWaitingState) HandleFiles(names []string) {
	err := state.view.tryAddManifestEntryFrom(names)
	if err != nil {
		state.failureTime = time.Now()
	} else {
		state.machine.SetState(nil)
	}
}
