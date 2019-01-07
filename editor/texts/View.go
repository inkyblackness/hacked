package texts

import (
	"fmt"

	"github.com/inkyblackness/imgui-go"

	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/edit/undoable"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
)

// View provides edit controls for texts.
type View struct {
	textService undoable.AugmentedTextService

	modalStateMachine gui.ModalStateMachine
	clipboard         external.Clipboard
	guiScale          float32

	model viewModel
}

// NewTextsView returns a new instance.
func NewTextsView(textService undoable.AugmentedTextService,
	modalStateMachine gui.ModalStateMachine, clipboard external.Clipboard,
	guiScale float32) *View {
	view := &View{
		textService: textService,

		modalStateMachine: modalStateMachine,
		clipboard:         clipboard,
		guiScale:          guiScale,

		model: freshViewModel(),
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
	knownTexts := edit.KnownTexts()
	imgui.PushItemWidth(-100 * view.guiScale)
	if imgui.BeginCombo("Text Type", knownTexts.Title(view.model.currentKey.ID)) {
		for _, info := range knownTexts {
			if imgui.SelectableV(info.Title, info.ID == view.model.currentKey.ID, 0, imgui.Vec2{}) {
				view.model.currentKey.ID = info.ID
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

	if view.textHasSound() {
		imgui.Separator()
		sound := view.currentSound()
		if !sound.Empty() {
			imgui.LabelText("Audio", fmt.Sprintf("%.2f sec", sound.Duration()))
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
	return view.textService.Text(view.model.currentKey)
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

	view.textService.RequestSetText(view.model.currentKey, value, view.restoreFunc())
}

func (view *View) clearText() {
	view.textService.RequestClear(view.model.currentKey, view.restoreFunc())
}

func (view *View) removeText() {
	view.textService.RequestRemove(view.model.currentKey, view.restoreFunc())
}

func (view View) textHasSound() bool {
	return view.textService.IsTrapMessage(view.model.currentKey)
}

func (view *View) currentSound() (sound audio.L8) {
	return view.textService.Sound(view.model.currentKey)
}

func (view *View) requestExportAudio(sound audio.L8) {
	if !view.textHasSound() {
		return
	}
	key := edit.TrapMessageSoundKeyFor(view.model.currentKey)
	filename := fmt.Sprintf("%05d_%s.wav", key.ID.Value(), view.model.currentKey.Lang.String())

	external.ExportAudio(view.modalStateMachine, filename, sound)
}

func (view *View) requestImportAudio() {
	external.ImportAudio(view.modalStateMachine, func(sound audio.L8) {
		view.textService.RequestSetSound(view.model.currentKey, sound, view.restoreFunc())
	})
}

func (view *View) restoreFunc() func() {
	oldKey := view.model.currentKey

	return func() {
		view.model.restoreFocus = true
		view.model.currentKey = oldKey
	}
}
