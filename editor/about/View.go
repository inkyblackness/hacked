package about

import (
	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/imgui-go"
)

// View handles the about display.
type View struct {
	guiScale  float32
	clipboard external.Clipboard

	model viewModel
}

// NewView creates a new instance for the about display.
func NewView(clipboard external.Clipboard, guiScale float32) *View {
	return &View{
		guiScale:  guiScale,
		clipboard: clipboard,

		model: freshViewModel(),
	}
}

// Show requets to show the about view.
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
	projectUrl := "https://inkyblackness.github.io"
	communityUrl := "https://www.systemshock.org"
	userguideUrl := "https://github.com/inkyblackness/hacked/wiki"
	urlLine := func(title, url string) {
		imgui.Text(title + ": " + url)
		imgui.SameLine()
		if imgui.Button("-> Clip##" + title) {
			view.clipboard.SetString(url)
		}
	}

	imgui.Text("InkyBlackness - HackEd")
	imgui.Separator()
	urlLine("User guide", userguideUrl)
	urlLine("Community", communityUrl)
	urlLine("Project", projectUrl)

	if imgui.Button("OK") {
		imgui.CloseCurrentPopup()
	}
}
