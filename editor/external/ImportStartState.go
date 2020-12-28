package external

import (
	"time"

	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/ui/gui"
)

type importStartState struct {
	machine   gui.ModalStateMachine
	callback  func(string)
	info      string
	typeInfo  []TypeInfo
	withError bool
}

func (state importStartState) Render() {
	imgui.OpenPopup("Import")
	nextState := &importWaitingState{
		machine:  state.machine,
		callback: state.callback,
		info:     state.info,
		typeInfo: state.typeInfo,
	}
	if state.withError {
		nextState.failureTime = time.Now()
	}
	state.machine.SetState(nextState)
}

func (state importStartState) HandleFiles(names []string) {
}
