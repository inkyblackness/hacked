package texts

import (
	"bytes"

	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/imgui-go"
)

// TextLinesView provides edit controls for simple text lines.
type TextLinesView struct {
	mod       *model.Mod
	adapter   *TextLinesAdapter
	clipboard external.Clipboard
	guiScale  float32
	commander cmd.Commander

	model viewModel
}

// NewTextLinesView returns a new instance.
func NewTextLinesView(mod *model.Mod, adapter *TextLinesAdapter, clipboard external.Clipboard, guiScale float32, commander cmd.Commander) *TextLinesView {
	view := &TextLinesView{
		mod:       mod,
		adapter:   adapter,
		clipboard: clipboard,
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
	imgui.PushItemWidth(-100 * view.guiScale)
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
	imgui.PopItemWidth()

	text := view.adapter.Line(view.model.currentKey)
	imgui.BeginChildV("Text", imgui.Vec2{X: -100 * view.guiScale, Y: 0}, true, 0)
	imgui.PushTextWrapPos()
	if len(text) == 0 {
		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1.0, Y: 1.0, Z: 1.0, W: 0.5})
		imgui.Text("(empty)")
		imgui.PopStyleColor()
	} else {
		imgui.Text(text)
	}
	imgui.PopTextWrapPos()
	imgui.EndChild()

	imgui.SameLine()

	imgui.BeginGroup()
	if imgui.ButtonV("-> Clip", imgui.Vec2{X: -1, Y: 0}) {
		view.copyTextToClipboard(text)
	}
	if imgui.ButtonV("<- Clip", imgui.Vec2{X: -1, Y: 0}) {
		view.setTextFromClipboard()
	}
	if imgui.ButtonV("Clear", imgui.Vec2{X: -1, Y: 0}) {
		view.clearTextLine()
	}
	if imgui.ButtonV("Remove", imgui.Vec2{X: -1, Y: 0}) {
		view.removeTextLine()
	}
	imgui.EndGroup()
}

func (view TextLinesView) copyTextToClipboard(text string) {
	if len(text) > 0 {
		view.clipboard.SetString(text)
	}
}

func (view *TextLinesView) setTextFromClipboard() {
	text, err := view.clipboard.String()
	if err != nil {
		return
	}

	old := view.mod.ModifiedBlock(view.model.currentKey)
	new := view.adapter.Codepage().Encode(text)
	view.requestSetTextLine(old, new)
}

func (view *TextLinesView) clearTextLine() {
	raw := view.mod.ModifiedBlock(view.model.currentKey)
	if !bytes.Equal(raw, []byte{0x00}) {
		view.requestSetTextLine(raw, []byte{0x00})
	}
}

func (view *TextLinesView) removeTextLine() {
	raw := view.mod.ModifiedBlock(view.model.currentKey)
	if len(raw) > 0 {
		view.requestSetTextLine(raw, nil)
	}
}

func (view *TextLinesView) requestSetTextLine(old []byte, new []byte) {
	command := setTextLineCommand{
		key:   view.model.currentKey,
		model: &view.model,
		old:   old,
		new:   new,
	}
	view.commander.Queue(command)
}
