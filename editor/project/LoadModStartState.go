package project

import "github.com/inkyblackness/imgui-go"

type loadModStartState struct {
	view *View
}

func (state loadModStartState) Render() {
	imgui.OpenPopup("Load mod")
	state.view.fileState = &loadModWaitingState{
		view: state.view,
	}
	state.view.fileState.Render()
}

func (state loadModStartState) HandleFiles(names []string) {
}
