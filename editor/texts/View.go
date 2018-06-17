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

type textInfo struct {
	id    resource.ID
	title string
}

var knownTextTypes = []textInfo{
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

// View provides edit controls for texts.
type View struct {
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

// NewTextsView returns a new instance.
func NewTextsView(mod *model.Mod, lineCache *text.Cache, pageCache *text.Cache, cp text.Codepage,
	clipboard external.Clipboard, guiScale float32, commander cmd.Commander) *View {
	view := &View{
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
	for _, info := range knownTextTypes {
		view.textTypeTitleByID[info.id] = info.title
	}
	return view
}

// WindowOpen returns the flag address, to be used with the main menu.
func (view *View) WindowOpen() *bool {
	return &view.model.windowOpen
}

// Render renders the view.
func (view *View) Render() {
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

func (view *View) renderContent() {
	imgui.PushItemWidth(-100 * view.guiScale)
	if imgui.BeginCombo("Text Type", view.textTypeTitleByID[view.model.currentKey.ID]) {
		for _, info := range knownTextTypes {
			if imgui.SelectableV(info.title, info.id == view.model.currentKey.ID, 0, imgui.Vec2{}) {
				view.model.currentKey.ID = info.id
				view.model.currentKey.Index = 0
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
	info, _ := ids.Info(view.model.currentKey.ID)
	if imgui.Button("-") && view.model.currentKey.Index > 0 {
		view.model.currentKey.Index--
	}
	imgui.SameLine()
	if imgui.Button("+") && view.model.currentKey.Index < (info.MaxCount-1) {
		view.model.currentKey.Index++
	}
	imgui.SameLine()
	index := int32(view.model.currentKey.Index)
	if imgui.SliderInt("Index", &index, 0, int32(info.MaxCount-1)) {
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
		view.clearText()
	}
	if imgui.ButtonV("Remove", imgui.Vec2{X: -1, Y: 0}) {
		view.removeText()
	}
	imgui.EndGroup()
}

func (view View) currentText() string {
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

func (view View) currentModification() (data [][]byte, isList bool) {
	info, _ := ids.Info(view.model.currentKey.ID)
	if info.List {
		data = [][]byte{view.mod.ModifiedBlock(view.model.currentKey.Lang, view.model.currentKey.ID, view.model.currentKey.Index)}
	} else {
		data = view.mod.ModifiedBlocks(view.model.currentKey.Lang, view.model.currentKey.ID.Plus(view.model.currentKey.Index))
	}
	return data, info.List
}

func (view View) copyTextToClipboard(text string) {
	if len(text) > 0 {
		view.clipboard.SetString(text)
	}
}

func (view *View) setTextFromClipboard() {
	value, err := view.clipboard.String()
	if err != nil {
		return
	}

	oldData, isList := view.currentModification()
	newData := view.cp.Encode(value)
	if isList {
		view.requestSetTextLine(oldData[0], newData)
	} else {
		view.requestSetTextPage(oldData, [][]byte{newData})
	}
}

func (view *View) clearText() {
	emptyLine := view.cp.Encode("")
	currentModification, isList := view.currentModification()
	if isList && !bytes.Equal(currentModification[0], emptyLine) {
		view.requestSetTextLine(currentModification[0], emptyLine)
	} else if !isList && ((len(currentModification) != 1) || !bytes.Equal(currentModification[0], []byte{0x00})) {
		view.requestSetTextPage(currentModification, [][]byte{emptyLine})
	}
}

func (view *View) removeText() {
	currentModification, isList := view.currentModification()
	if isList && (len(currentModification[0]) > 0) {
		view.requestSetTextLine(currentModification[0], nil)
	} else if !isList && (len(currentModification) > 0) {
		view.requestSetTextPage(currentModification, nil)
	}
}

func (view *View) requestSetTextLine(oldData []byte, newData []byte) {
	command := setTextLineCommand{
		key:     view.model.currentKey,
		model:   &view.model,
		oldData: oldData,
		newData: newData,
	}
	view.commander.Queue(command)
}

func (view *View) requestSetTextPage(oldData [][]byte, newData [][]byte) {
	command := setTextPageCommand{
		key:     view.model.currentKey,
		model:   &view.model,
		oldData: oldData,
		newData: newData,
	}
	view.commander.Queue(command)
}
