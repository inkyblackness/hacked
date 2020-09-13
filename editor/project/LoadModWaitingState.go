package project

import (
	"time"

	"github.com/inkyblackness/imgui-go/v2"
	"github.com/sqweek/dialog"

	"github.com/inkyblackness/hacked/ui/gui"
)

type loadModWaitingState struct {
	machine     gui.ModalStateMachine
	view        *View
	failureTime time.Time
}

func (state *loadModWaitingState) Render() {
	if imgui.BeginPopupModalV("Load mod", nil,
		imgui.WindowFlagsNoResize|imgui.WindowFlagsNoMove|imgui.WindowFlagsNoSavedSettings|imgui.WindowFlagsAlwaysAutoResize) {
		imgui.Text("Waiting for folder.")
		if !state.failureTime.IsZero() {
			imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
			imgui.Text("Previous attempt failed, no usable data detected.\nPlease check and try again.")
			imgui.PopStyleColor()
			if time.Since(state.failureTime).Seconds() > 5 {
				state.failureTime = time.Time{}
			}
		}
		imgui.Text(`From your file browser drag'n'drop the folder
of the mod you want to work on into the editor window.
If you want to modify the main game files,
use the main "data" directory of the game.
`)
		imgui.Text("This action will clear the undo/redo buffer\nand you will lose any unsaved changes.")
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

func (state *loadModWaitingState) HandleFiles(names []string) {
	err := state.view.tryLoadModFrom(names)
	if err != nil {
		state.failureTime = time.Now()
	} else {
		state.machine.SetState(nil)
	}
}
