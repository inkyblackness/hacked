package project

import "github.com/inkyblackness/imgui-go"

type saveModAsStartState struct {
	view *View
}

func (state saveModAsStartState) Render() {
	imgui.OpenPopup("Save mod as")
	state.view.fileState = &saveModAsWaitingState{
		view: state.view,
	}
	state.view.fileState.Render()
}

func (state saveModAsStartState) HandleFiles(names []string) {
}
