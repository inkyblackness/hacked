package project

import (
	"time"

	"github.com/inkyblackness/imgui-go"
)

type saveModFailedState struct {
	view *View
}

func (state saveModFailedState) Render() {
	imgui.OpenPopup("Save mod as")
	state.view.fileState = &saveModAsWaitingState{
		view:        state.view,
		failureTime: time.Now(),
	}
	state.view.fileState.Render()
}

func (state saveModFailedState) HandleFiles(names []string) {
}
