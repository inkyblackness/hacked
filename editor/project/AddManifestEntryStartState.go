package project

import (
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
)

type addManifestEntryStartState struct {
	machine gui.ModalStateMachine
	view    *View
}

func (state addManifestEntryStartState) Render() {
	imgui.OpenPopup("Add static world data")
	state.machine.SetState(&addManifestEntryWaitingState{
		machine: state.machine,
		view:    state.view,
	})
}

func (state addManifestEntryStartState) HandleFiles(names []string) {
}
