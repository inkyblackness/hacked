package project

import (
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
)

type saveModAsStartState struct {
	machine gui.ModalStateMachine
	view    *View
}

func (state saveModAsStartState) Render() {
	imgui.OpenPopup("Save mod as")
	state.machine.SetState(&saveModAsWaitingState{
		machine: state.machine,
		view:    state.view,
	})
}

func (state saveModAsStartState) HandleFiles(names []string) {
}
