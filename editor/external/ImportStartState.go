package external

import (
	"time"

	"github.com/inkyblackness/imgui-go"

	"github.com/inkyblackness/hacked/ui/gui"
)

type importStartState struct {
	machine   gui.ModalStateMachine
	callback  func(string)
	info      string
	withError bool
}

func (state importStartState) Render() {
	imgui.OpenPopup("Import")
	nextState := &importWaitingState{
		machine:  state.machine,
		callback: state.callback,
		info:     state.info,
	}
	if state.withError {
		nextState.failureTime = time.Now()
	}
	state.machine.SetState(nextState)
}

func (state importStartState) HandleFiles(names []string) {
}
