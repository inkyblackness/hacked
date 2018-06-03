package project

import "github.com/inkyblackness/imgui-go"

type addManifestEntryStartState struct {
	view *View
}

func (state addManifestEntryStartState) Render() {
	imgui.OpenPopup("Add static world data")
	state.view.fileState = &addManifestEntryWaitingState{
		view: state.view,
	}
	state.view.fileState.Render()
}

func (state addManifestEntryStartState) HandleFiles(names []string) {
}
