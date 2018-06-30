package levels

import (
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/imgui-go"
)

// ControlView is the core view for level editing.
type ControlView struct {
	guiScale float32

	model viewModel
}

// NewControlView returns a new instance.
func NewControlView(guiScale float32) *ControlView {
	view := &ControlView{
		guiScale: guiScale,
		model:    freshViewModel(),
	}
	return view
}

// WindowOpen returns the flag address, to be used with the main menu.
func (view *ControlView) WindowOpen() *bool {
	return &view.model.windowOpen
}

// SelectedLevel returns the currently selected level.
func (view *ControlView) SelectedLevel() int {
	return view.model.selectedLevel
}

// Render renders the view.
func (view *ControlView) Render() {
	if view.model.restoreFocus {
		imgui.SetNextWindowFocus()
		view.model.restoreFocus = false
		view.model.windowOpen = true
	}
	if view.model.windowOpen {
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 300 * view.guiScale}, imgui.ConditionOnce)
		if imgui.BeginV("Level Control", view.WindowOpen(), imgui.WindowFlagsNoCollapse) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *ControlView) renderContent() {
	imgui.PushItemWidth(-100 * view.guiScale)
	if imgui.Button("-") && view.model.selectedLevel > 0 {
		view.model.selectedLevel--
	}
	imgui.SameLine()
	if imgui.Button("+") && view.model.selectedLevel < (archive.MaxLevels-1) {
		view.model.selectedLevel++
	}
	imgui.SameLine()
	level := int32(view.model.selectedLevel)
	if imgui.SliderInt("Active Level", &level, 0, int32(archive.MaxLevels-1)) {
		view.model.selectedLevel = int(level)
	}
	imgui.PopItemWidth()

}
