package project

import (
	"github.com/inkyblackness/imgui-go/v2"

	"github.com/inkyblackness/hacked/ui/gui"
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
