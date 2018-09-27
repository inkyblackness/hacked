package project

import (
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
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
