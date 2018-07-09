package project

import (
	"time"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/imgui-go"
)

type addManifestEntryWaitingState struct {
	view        *View
	failureTime time.Time
}

func (state *addManifestEntryWaitingState) Render() {
	if imgui.BeginPopupModalV("Add static world data", nil,
		imgui.WindowFlagsNoResize|imgui.WindowFlagsNoMove|imgui.WindowFlagsNoSavedSettings|imgui.WindowFlagsAlwaysAutoResize) {

		imgui.Text("Waiting for folders/files.")
		if !state.failureTime.IsZero() {
			imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
			imgui.Text("Previous attempt failed, no usable data detected.\nPlease check and try again.")
			imgui.PopStyleColor()
			if time.Since(state.failureTime).Seconds() > 5 {
				state.failureTime = time.Time{}
			}
		}
		imgui.Text(`From your file browser drag'n'drop the folder (or files)
of the static data you want to reference into the editor window.
Typically, you would use the main "data" directory of the game
(where all the .res files are).
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

func (state *addManifestEntryWaitingState) HandleFiles(names []string) {
	staging := fileStaging{
		resources: make(map[string]resource.Provider),
		savegames: make(map[string]resource.Provider),
	}

	for _, name := range names {
		staging.stage(name, true)
	}
	if len(staging.resources) > 0 {
		entry := &world.ManifestEntry{
			ID: names[0],
		}
		blacklisted := ids.LowResVideos()

		for filename, provider := range staging.resources {
			if !blacklisted.Matches(filename) {
				localized := resource.LocalizedResources{
					ID:       filename,
					Language: ids.LocalizeFilename(filename),
					Provider: provider,
				}
				entry.Resources = append(entry.Resources, localized)
			}
		}
		if len(staging.objectProperties) > 0 {
			entry.ObjectProperties = staging.objectProperties
		}

		state.view.requestAddManifestEntry(entry)
		state.view.fileState = &idlePopupState{}
	} else {
		state.failureTime = time.Now()
	}
}
