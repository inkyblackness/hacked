package texts

import (
	"bytes"

	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/imgui-go"
)

type textLineInfo struct {
	id    resource.ID
	title string
}

var knownTextLineTypes = []textLineInfo{
	{ids.TrapMessageTexts, "Trap Messages"},
	{ids.WordTexts, "Words"},
	{ids.LogCategoryTexts, "Log Categories"},
	{ids.VariousMessageTexts, "Various Messages"},
	{ids.ScreenMessageTexts, "Screen Messages"},
	{ids.InfoNodeMessageTexts, "Info Node Message Texts (8/5/6)"},
	{ids.AccessCardNameTexts, "Access Card Names"},
	{ids.DataletMessageTexts, "Datalet Messages (8/5/8)"},
	{ids.PaperTextsStart, "Papers"},
	{ids.PanelNameTexts, "Panel Names"},
}

// TextLinesView provides edit controls for simple text lines.
type TextLinesView struct {
	mod       *model.Mod
	lineCache *text.Cache
	pageCache *text.Cache
	cp        text.Codepage

	clipboard external.Clipboard
	guiScale  float32
	commander cmd.Commander

	model viewModel

	textTypeTitleByID map[resource.ID]string
}

// NewTextLinesView returns a new instance.
func NewTextLinesView(mod *model.Mod, lineCache *text.Cache, pageCache *text.Cache, cp text.Codepage,
	clipboard external.Clipboard, guiScale float32, commander cmd.Commander) *TextLinesView {
	view := &TextLinesView{
		mod:       mod,
		lineCache: lineCache,
		pageCache: pageCache,
		cp:        cp,

		clipboard: clipboard,
		guiScale:  guiScale,
		commander: commander,

		model: freshViewModel(),

		textTypeTitleByID: make(map[resource.ID]string),
	}
	for _, info := range knownTextLineTypes {
		view.textTypeTitleByID[info.id] = info.title
	}
	return view
}

// WindowOpen returns the flag address, to be used with the main menu.
func (view *TextLinesView) WindowOpen() *bool {
	return &view.model.windowOpen
}

// Render renders the view.
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
	if imgui.BeginCombo("Text Type", view.textTypeTitleByID[view.model.currentKey.ID]) {
		for _, info := range knownTextLineTypes {
			if imgui.SelectableV(info.title, info.id == view.model.currentKey.ID, 0, imgui.Vec2{}) {
				view.model.currentKey.ID = info.id
			}
		}
		imgui.EndCombo()
	}
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

	currentText := view.currentText()
	imgui.BeginChildV("Text", imgui.Vec2{X: -100 * view.guiScale, Y: 0}, true, 0)
	imgui.PushTextWrapPos()
	if len(currentText) == 0 {
		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1.0, Y: 1.0, Z: 1.0, W: 0.5})
		imgui.Text("(empty)")
		imgui.PopStyleColor()
	} else {
		imgui.Text(currentText)
	}
	imgui.PopTextWrapPos()
	imgui.EndChild()

	imgui.SameLine()

	imgui.BeginGroup()
	if imgui.ButtonV("-> Clip", imgui.Vec2{X: -1, Y: 0}) {
		view.copyTextToClipboard(currentText)
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

func (view TextLinesView) currentText() string {
	var cache *text.Cache
	resourceInfo, existing := ids.Info(view.model.currentKey.ID)
	if !existing || resourceInfo.List {
		cache = view.lineCache
	} else {
		cache = view.pageCache
	}
	currentValue, cacheErr := cache.Text(view.model.currentKey)
	if cacheErr != nil {
		currentValue = ""
	}
	return currentValue
}

func (view TextLinesView) copyTextToClipboard(text string) {
	if len(text) > 0 {
		view.clipboard.SetString(text)
	}
}

func (view *TextLinesView) setTextFromClipboard() {
	value, err := view.clipboard.String()
	if err != nil {
		return
	}

	oldData := view.mod.ModifiedBlock(view.model.currentKey)
	newData := view.cp.Encode(value)
	view.requestSetTextLine(oldData, newData)
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

func (view *TextLinesView) requestSetTextLine(oldData []byte, newData []byte) {
	command := setTextLineCommand{
		key:     view.model.currentKey,
		model:   &view.model,
		oldData: oldData,
		newData: newData,
	}
	view.commander.Queue(command)
}
