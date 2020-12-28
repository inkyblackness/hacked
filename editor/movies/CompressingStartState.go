package movies

import (
	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ui/gui"
)

type compressionListenerFunc func(compressionResult)

type compressingStartState struct {
	machine gui.ModalStateMachine
	view    *View

	width    int
	height   int
	input    movie.Scene
	listener compressionListenerFunc
}

func (state compressingStartState) Render() {
	imgui.OpenPopup("Compressing...")
	task := newCompressionTask(state.input, state.width, state.height)
	state.machine.SetState(&compressingWaitingState{
		machine:  state.machine,
		view:     state.view,
		input:    state.input,
		listener: state.listener,
		task:     task,
	})
	go task.run()
}

func (state compressingStartState) HandleFiles(names []string) {
}
