package project

import (
	"github.com/inkyblackness/imgui-go/v2"

	"github.com/inkyblackness/hacked/ui/gui"
)

type loadModStartState struct {
	machine gui.ModalStateMachine
	view    *View
}

func (state loadModStartState) Render() {
	imgui.OpenPopup("Load mod")
	state.machine.SetState(&loadModWaitingState{
		machine: state.machine,
		view:    state.view,
	})
}

func (state loadModStartState) HandleFiles(names []string) {
}
