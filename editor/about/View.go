package about

import (
	"github.com/inkyblackness/imgui-go/v2"

	"github.com/inkyblackness/hacked/editor/external"
)

// View handles the about display.
type View struct {
	version   string
	guiScale  float32
	clipboard external.Clipboard

	model viewModel
}

// NewView creates a new instance for the about display.
func NewView(clipboard external.Clipboard, guiScale float32, version string) *View {
	return &View{
		version:   version,
		guiScale:  guiScale,
		clipboard: clipboard,

		model: freshViewModel(),
	}
}

// Show requests to show the about view.
func (view *View) Show() {
	view.model.windowOpen = true
}

// Render requests to render the view.
func (view *View) Render() {
	if view.model.windowOpen {
		imgui.OpenPopup("About")
		view.model.windowOpen = false
	}
	if imgui.BeginPopupModalV("About", nil, imgui.WindowFlagsAlwaysAutoResize|imgui.WindowFlagsNoMove|imgui.WindowFlagsHorizontalScrollbar) {
		view.renderContent()
		imgui.EndPopup()
	}
}

func (view *View) renderContent() {
	projectURL := "https://inkyblackness.github.io"
	communityURL := "https://www.systemshock.org"
	userguideURL := "https://github.com/inkyblackness/hacked/wiki"
	urlLine := func(title, url string) {
		imgui.Text(title + ": " + url)
		imgui.SameLine()
		if imgui.Button("-> Clip##" + title) {
			view.clipboard.SetString(url)
		}
	}

	imgui.Text("InkyBlackness - HackEd - " + view.version)
	imgui.Separator()
	urlLine("User guide", userguideURL)
	urlLine("Community", communityURL)
	urlLine("Project", projectURL)
	imgui.Separator()
	imgui.Text("Thanks to: 3RDPlayer: testing, ToxicFrog: name idea.")
	imgui.Text("Written in Go. Because I can.")
	imgui.Separator()
	if imgui.Button("OK") {
		imgui.CloseCurrentPopup()
	}
}
