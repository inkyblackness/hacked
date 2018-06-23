package project

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/resource/lgres"
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

type fileStaging struct {
	failedFiles int
	savegames   map[string]resource.Provider
	resources   map[string]resource.Provider
}

func (staging *fileStaging) stage(name string, enterDir bool) {
	fileInfo, err := os.Stat(name)
	if err != nil {
		staging.failedFiles++
		return
	}
	file, err := os.Open(name)
	if err != nil {
		return
	}
	defer file.Close() // nolint: errcheck

	if fileInfo.IsDir() {
		if enterDir {
			subNames, _ := file.Readdirnames(0)
			for _, subName := range subNames {
				staging.stage(filepath.Join(name, subName), false)
			}
		}
	} else {
		fileData, err := ioutil.ReadAll(file)
		if err != nil {
			staging.failedFiles++
		}

		reader, err := lgres.ReaderFrom(bytes.NewReader(fileData))
		if err == nil {
			filename := filepath.Base(name)
			if resource.IsSavegame(reader) {
				staging.savegames[filename] = reader
			} else {
				staging.resources[filename] = reader
			}
		}

		if err != nil {
			staging.failedFiles++
		}
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

		state.view.requestAddManifestEntry(entry)
		state.view.fileState = &idlePopupState{}
	} else {
		state.failureTime = time.Now()
	}
}
