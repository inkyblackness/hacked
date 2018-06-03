package project

import "github.com/inkyblackness/imgui-go"

type addManifestEntryWaitingState struct {
	view *View
}

func (state addManifestEntryWaitingState) Render() {
	if imgui.BeginPopupModalV("Add static world data", nil,
		imgui.WindowFlagsNoResize|imgui.WindowFlagsNoMove|imgui.WindowFlagsNoSavedSettings) {

		imgui.TextUnformatted(`Waiting for folders/files.

From your file browser drag'n'drop the folder (or files)
of the static data you want to reference into the editor window.
`)
		imgui.Separator()
		if imgui.Button("Cancel") {
			state.view.fileState = &idlePopupState{}
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	} else {
		state.view.fileState = &idlePopupState{}
	}
}

func (state addManifestEntryWaitingState) HandleFiles(names []string) {

}
