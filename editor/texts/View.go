package texts

import (
	"fmt"

	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/modding"
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
)

type textInfo struct {
	id    resource.ID
	title string
}

var knownTextTypes = []textInfo{
	{id: ids.TrapMessageTexts, title: "Trap Messages"},
	{id: ids.WordTexts, title: "Words"},
	{id: ids.LogCategoryTexts, title: "Log Categories"},
	{id: ids.VariousMessageTexts, title: "Various Messages"},
	{id: ids.ScreenMessageTexts, title: "Screen Messages"},
	{id: ids.InfoNodeMessageTexts, title: "Info Node Message Texts (8/5/6)"},
	{id: ids.AccessCardNameTexts, title: "Access Card Names"},
	{id: ids.DataletMessageTexts, title: "Datalet Messages (8/5/8)"},
	{id: ids.PaperTextsStart, title: "Papers"},
	{id: ids.PanelNameTexts, title: "Panel Names"},
}

var textToAudio = map[resource.ID]resource.ID{
	ids.TrapMessageTexts: ids.TrapMessagesAudioStart,
}

// View provides edit controls for texts.
type View struct {
	getText  modding.GetTextService
	setText  modding.SetTextService
	getAudio modding.GetAudioService
	setAudio modding.SetAudioService

	modalStateMachine gui.ModalStateMachine
	clipboard         external.Clipboard
	guiScale          float32
	commander         cmd.Commander

	model viewModel

	textTypeTitleByID map[resource.ID]string
}

// NewTextsView returns a new instance.
func NewTextsView(
	getText modding.GetTextService, setText modding.SetTextService,
	getAudio modding.GetAudioService, setAudio modding.SetAudioService,
	modalStateMachine gui.ModalStateMachine, clipboard external.Clipboard,
	guiScale float32, commander cmd.Commander) *View {
	view := &View{
		getText:  getText,
		setText:  setText,
		getAudio: getAudio,
		setAudio: setAudio,

		modalStateMachine: modalStateMachine,
		clipboard:         clipboard,
		guiScale:          guiScale,
		commander:         commander,

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
	info, _ := ids.Info(view.model.currentKey.ID)
	gui.StepSliderInt("Index", &view.model.currentKey.Index, 0, info.MaxCount-1)

	if imgui.BeginCombo("Language", view.model.currentKey.Lang.String()) {
		languages := resource.Languages()
		for _, lang := range languages {
			if imgui.SelectableV(lang.String(), lang == view.model.currentKey.Lang, 0, imgui.Vec2{}) {
				view.model.currentKey.Lang = lang
			}
		}
		imgui.EndCombo()
	}

	if view.hasAudio() {
		imgui.Separator()
		sound := view.currentSound()
		hasSound := len(sound.Samples) > 0
		if hasSound {
			imgui.LabelText("Audio", fmt.Sprintf("%.2f sec", float32(len(sound.Samples))/sound.SampleRate))
			if imgui.Button("Export") {
				view.requestExportAudio(sound)
			}
			imgui.SameLine()
		} else {
			imgui.LabelText("Audio", "(no sound)")
		}
		if imgui.Button("Import") {
			view.requestImportAudio()
		}
	}
	imgui.Separator()

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
	return view.getText.Current(view.model.currentKey)
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

	key := view.model.currentKey
	oldText := view.currentText()
	isModified := view.getText.Modified(key)
	view.requestCommand(
		func(trans cmd.Transaction) {
			view.setText.Set(trans, key, value)
		},
		func(trans cmd.Transaction) {
			if isModified {
				view.setText.Set(trans, key, oldText)
			} else {
				view.setText.Remove(trans, key)
			}
		})
}

func (view *View) clearText() {
	textKey := view.model.currentKey
	oldText := view.currentText()
	isTextModified := view.getText.Modified(textKey)

	audioKey, textWithAudio := view.currentSoundKey()
	oldAudio := view.currentSound()
	isAudioModified := textWithAudio && view.getAudio.Modified(audioKey)

	view.requestCommand(
		func(trans cmd.Transaction) {
			view.setText.Clear(trans, textKey)
			if textWithAudio {
				view.setAudio.Clear(trans, audioKey)
			}
		},
		func(trans cmd.Transaction) {
			if isTextModified {
				view.setText.Set(trans, textKey, oldText)
			} else {
				view.setText.Remove(trans, textKey)
			}
			if isAudioModified {
				view.setAudio.Set(trans, audioKey, oldAudio)
			} else if textWithAudio {
				view.setAudio.Remove(trans, audioKey)
			}
		})
}

func (view *View) removeText() {
	textKey := view.model.currentKey
	oldText := view.currentText()
	isTextModified := view.getText.Modified(textKey)

	audioKey, textWithAudio := view.currentSoundKey()
	oldAudio := view.currentSound()
	isAudioModified := textWithAudio && view.getAudio.Modified(audioKey)

	view.requestCommand(
		func(trans cmd.Transaction) {
			view.setText.Remove(trans, textKey)
			if textWithAudio {
				view.setAudio.Remove(trans, audioKey)
			}
		},
		func(trans cmd.Transaction) {
			if isTextModified {
				view.setText.Set(trans, textKey, oldText)
			} else {
				view.setText.Remove(trans, textKey)
			}
			if isAudioModified {
				view.setAudio.Set(trans, audioKey, oldAudio)
			} else if textWithAudio {
				view.setAudio.Remove(trans, audioKey)
			}
		})
}

// --- audio

func (view View) hasAudio() bool {
	return view.model.currentKey.ID == ids.TrapMessageTexts
}

func (view *View) currentSoundKey() (key resource.Key, textWithAudio bool) {
	if !view.hasAudio() {
		return key, false
	}
	key = view.model.currentKey
	key.ID = textToAudio[key.ID].Plus(key.Index)
	key.Index = 0
	return key, true
}

func (view *View) currentSound() (sound audio.L8) {
	key, hasAudio := view.currentSoundKey()
	if !hasAudio {
		return
	}
	return view.getAudio.Get(key)
}

func (view *View) requestExportAudio(sound audio.L8) {
	key, hasAudio := view.currentSoundKey()
	if !hasAudio {
		return
	}
	filename := fmt.Sprintf("%05d_%s.wav", key.ID.Value(), view.model.currentKey.Lang.String())

	external.ExportAudio(view.modalStateMachine, filename, sound)
}

func (view *View) requestImportAudio() {
	external.ImportAudio(view.modalStateMachine, func(sound audio.L8) {
		view.requestSetAudio(sound)
	})
}

func (view *View) requestSetAudio(sound audio.L8) {
	audioKey, textWithAudio := view.currentSoundKey()
	oldAudio := view.currentSound()
	isAudioModified := textWithAudio && view.getAudio.Modified(audioKey)

	view.requestCommand(
		func(trans cmd.Transaction) {
			if textWithAudio {
				view.setAudio.Set(trans, audioKey, sound)
			}
		},
		func(trans cmd.Transaction) {
			if isAudioModified {
				view.setAudio.Set(trans, audioKey, oldAudio)
			} else if textWithAudio {
				view.setAudio.Remove(trans, audioKey)
			}
		})
}

func (view *View) requestCommand(forward func(trans cmd.Transaction), backward func(trans cmd.Transaction)) {
	c := command{
		key:      view.model.currentKey,
		model:    &view.model,
		forward:  forward,
		backward: backward,
	}
	view.commander.Queue(c)
}
