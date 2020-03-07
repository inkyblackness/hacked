package project

import (
	"time"

	"github.com/inkyblackness/imgui-go"
	"github.com/sqweek/dialog"

	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
)

type loadModWaitingState struct {
	machine     gui.ModalStateMachine
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
		if imgui.Button("Browse...") {
			dlgBuilder := dialog.Directory()
			filename, err := dlgBuilder.Browse()
			if err == nil {
				state.HandleFiles([]string{filename})
			}
		}
		imgui.SameLine()
		if imgui.Button("Cancel") {
			state.machine.SetState(nil)
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	} else {
		state.machine.SetState(nil)
	}
}

func (state *loadModWaitingState) HandleFiles(names []string) {
	staging := newFileStaging()

	staging.stageAll(names)

	if len(staging.resources) > 0 {
		var locs []*world.LocalizedResources

		for filename, viewer := range staging.resources {
			lang := ids.LocalizeFilename(filename)
			loc := &world.LocalizedResources{
				Filename: filename,
				Language: lang,
			}
			for _, id := range viewer.IDs() {
				view, err := viewer.View(id)
				if err == nil {
					_ = loc.Store.Put(id, view)
				}
				// TODO: handle error?
			}
			locs = append(locs, loc)
		}

		state.machine.SetState(nil)
		state.view.requestLoadMod(names[0], locs, staging.objectProperties, staging.textureProperties)
	} else {
		state.failureTime = time.Now()
	}
}
