package external

import (
	"time"

	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/ui/gui"
)

type exportStartState struct {
	machine   gui.ModalStateMachine
	callback  func(string)
	info      string
	withError bool
}

func (state exportStartState) Render() {
	imgui.OpenPopup("Export")
	nextState := &exportWaitingState{
		machine:  state.machine,
		callback: state.callback,
		info:     state.info,
	}
	if state.withError {
		nextState.failureTime = time.Now()
	}
	state.machine.SetState(nextState)
}

func (state exportStartState) HandleFiles(names []string) {
}
