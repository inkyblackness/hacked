package project

import (
	"time"

	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
)

type saveModFailedState struct {
	machine   gui.ModalStateMachine
	view      *View
	errorInfo string
}

func (state saveModFailedState) Render() {
	imgui.OpenPopup("Save mod as")
	state.machine.SetState(&saveModAsWaitingState{
		machine:     state.machine,
		view:        state.view,
		failureTime: time.Now(),
		errorInfo:   state.errorInfo,
	})
}

func (state saveModFailedState) HandleFiles(names []string) {
}
