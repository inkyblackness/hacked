package levels

import "github.com/inkyblackness/imgui-go"

// TilesView is for tile properties.
type TilesView struct {
	guiScale float32
	model    tilesViewModel
}

// NewTilesView returns a new instance.
func NewTilesView(guiScale float32) *TilesView {
	view := &TilesView{
		guiScale: guiScale,
		model:    freshTilesViewModel(),
	}
	return view
}

// WindowOpen returns the flag address, to be used with the main menu.
func (view *TilesView) WindowOpen() *bool {
	return &view.model.windowOpen
}

// Render renders the view.
func (view *TilesView) Render() {
	if view.model.restoreFocus {
		imgui.SetNextWindowFocus()
		view.model.restoreFocus = false
		view.model.windowOpen = true
	}
	if view.model.windowOpen {
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 500 * view.guiScale}, imgui.ConditionOnce)
		if imgui.BeginV("Level Tiles", view.WindowOpen(), 0) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *TilesView) renderContent() {

}
