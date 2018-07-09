package project

import (
	"time"

	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/imgui-go"
)

type loadModWaitingState struct {
	view        *View
	failureTime time.Time
}

func (state *loadModWaitingState) Render() {
	if imgui.BeginPopupModalV("Load mod", nil,
		imgui.WindowFlagsNoResize|imgui.WindowFlagsNoMove|imgui.WindowFlagsNoSavedSettings|imgui.WindowFlagsAlwaysAutoResize) {

		imgui.Text("Waiting for folder.")
		if !state.failureTime.IsZero() {
			imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
			imgui.Text("Previous attempt failed, no usable data detected.\nPlease check and try again.")
			imgui.PopStyleColor()
			if time.Since(state.failureTime).Seconds() > 5 {
				state.failureTime = time.Time{}
			}
		}
		imgui.Text(`From your file browser drag'n'drop the folder
of the mod you want to work on into the editor window.
If you want to modify the main game files,
use the main "data" directory of the game.
`)
		imgui.Text("This action will clear the undo/redo buffer\nand you will lose any unsaved changes.")
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

func (state *loadModWaitingState) HandleFiles(names []string) {
	staging := fileStaging{
		resources: make(map[string]resource.Provider),
		savegames: make(map[string]resource.Provider),
	}

	for _, name := range names {
		staging.stage(name, true)
	}
	if len(staging.resources) > 0 {
		res := model.NewLocalizedResources()
		blacklisted := ids.LowResVideos()

		for filename, provider := range staging.resources {
			if !blacklisted.Matches(filename) {
				lang := ids.LocalizeFilename(filename)
				res[lang].Add(model.MutableResourcesFromProvider(filename, provider))
			}
		}

		state.view.fileState = &idlePopupState{}
		state.view.requestLoadMod(names[0], res, staging.objectProperties)
	} else {
		state.failureTime = time.Now()
	}
}
