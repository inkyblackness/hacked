package texts

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/imgui-go"
)

// TextLinesView provides edit controls for simple text lines.
type TextLinesView struct {
	adapter   *TextLinesAdapter
	guiScale  float32
	commander cmd.Commander

	model viewModel
}

// NewTextLinesView returns a new instance.
func NewTextLinesView(adapter *TextLinesAdapter, guiScale float32, commander cmd.Commander) *TextLinesView {
	view := &TextLinesView{
		adapter:   adapter,
		guiScale:  guiScale,
		commander: commander,

		model: freshViewModel(),
	}
	return view
}

// WindowOpen returns the flag address, to be used with the main menu.
func (view *TextLinesView) WindowOpen() *bool {
	return &view.model.windowOpen
}

func (view *TextLinesView) Render() {
	if view.model.restoreFocus {
		imgui.SetNextWindowFocus()
		view.model.restoreFocus = false
		view.model.windowOpen = true
	}
	if view.model.windowOpen {
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 300 * view.guiScale}, imgui.ConditionOnce)
		if imgui.BeginV("Texts", view.WindowOpen(), imgui.WindowFlagsNoCollapse) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *TextLinesView) renderContent() {
	if imgui.BeginCombo("Language", view.model.currentKey.Lang.String()) {
		languages := resource.Languages()
		for _, lang := range languages {
			if imgui.SelectableV(lang.String(), lang == view.model.currentKey.Lang, 0, imgui.Vec2{}) {
				view.model.currentKey.Lang = lang
			}
		}
		imgui.EndCombo()
	}
	if imgui.Button("-") && view.model.currentKey.Index > 0 {
		view.model.currentKey.Index--
	}
	imgui.SameLine()
	if imgui.Button("+") && view.model.currentKey.Index < 255 {
		view.model.currentKey.Index++
	}
	imgui.SameLine()
	index := int32(view.model.currentKey.Index)
	if imgui.SliderInt("Index", &index, 0, 255) {
		view.model.currentKey.Index = int(index)
	}

	text := view.adapter.Line(view.model.currentKey)
	imgui.TextUnformatted(text)
}
