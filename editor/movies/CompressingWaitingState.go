package movies

import (
	"github.com/inkyblackness/imgui-go"

	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ui/gui"
)

type compressingWaitingState struct {
	machine gui.ModalStateMachine
	view    *View

	input    movie.Scene
	listener compressionListenerFunc
	task     *compressionTask
}

func (state *compressingWaitingState) Render() {
	if imgui.BeginPopupModalV("Compressing...", nil,
		imgui.WindowFlagsNoResize|imgui.WindowFlagsNoMove|imgui.WindowFlagsNoSavedSettings|imgui.WindowFlagsAlwaysAutoResize) {

		imgui.Text("This may take a while.\n" +
			"Yes, your PC is probably capable of recoding HD movies way quicker;\n" +
			"Sadly, this codec of '94 is quite tricky.")

		if imgui.Button("Cancel") {
			state.task.cancel()
		}

		result := state.task.update()
		if result != nil {
			state.machine.SetState(nil)
			imgui.CloseCurrentPopup()

			state.listener(result)
		}

		imgui.EndPopup()
	} else {
		state.machine.SetState(nil)
	}
}

func (state *compressingWaitingState) HandleFiles(names []string) {
}
