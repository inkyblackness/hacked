package sounds

import (
	"fmt"

	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/edit/undoable"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
)

// View provides edit controls for sound effects.
type View struct {
	soundEffectService undoable.SoundEffectService

	modalStateMachine gui.ModalStateMachine
	guiScale          float32

	model viewModel
}

// NewSoundEffectsView returns a new instance.
func NewSoundEffectsView(soundEffectService undoable.SoundEffectService,
	modalStateMachine gui.ModalStateMachine,
	guiScale float32) *View {
	view := &View{
		soundEffectService: soundEffectService,

		modalStateMachine: modalStateMachine,
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
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 300 * view.guiScale}, imgui.ConditionFirstUseEver)
		if imgui.BeginV("Sound Effects", view.WindowOpen(), imgui.WindowFlagsNoCollapse) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *View) renderContent() {
	info, _ := ids.Info(view.model.currentKey.ID)

	imgui.BeginChildV("SoundEffects", imgui.Vec2{X: -1, Y: -60 * view.guiScale}, true, 0)
	for i := 0; i < info.MaxCount; i++ {
		effects := ids.SoundEffectsForAudio(i)
		text := fmt.Sprintf("%3d", i)
		if len(effects) > 0 {
			text += " - "
			for effectIndex, effect := range effects {
				if effectIndex > 0 {
					text += ", "
				}
				text += effect.Name
			}
		} else {
			text += " (unidentified)"
		}
		if imgui.SelectableV(text, view.model.currentKey.Index == i, 0, imgui.Vec2{}) {
			view.model.currentKey.Index = i
		}
	}
	imgui.EndChild()

	imgui.PushItemWidth(-100 * view.guiScale)
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
	imgui.PopItemWidth()
	imgui.SameLine()
	if imgui.Button("Clear") {
		view.clearAudio()
	}
	if view.soundEffectService.Modified(view.model.currentKey) {
		imgui.SameLine()
		if imgui.Button("Remove") {
			view.removeAudio()
		}
	}
}

func (view *View) clearAudio() {
	view.soundEffectService.RequestClear(view.model.currentKey, view.restoreFunc())
}

func (view *View) removeAudio() {
	view.soundEffectService.RequestRemove(view.model.currentKey, view.restoreFunc())
}

func (view *View) currentSound() (sound audio.L8) {
	return view.soundEffectService.Audio(view.model.currentKey)
}

func (view *View) requestExportAudio(sound audio.L8) {
	filename := fmt.Sprintf("sfx_%03d.wav", view.model.currentKey.Index)

	external.ExportAudio(view.modalStateMachine, filename, sound)
}

func (view *View) requestImportAudio() {
	external.ImportAudio(view.modalStateMachine, func(sound audio.L8) {
		view.soundEffectService.RequestSetAudio(view.model.currentKey, sound, view.restoreFunc())
	})
}

func (view *View) restoreFunc() func() {
	oldKey := view.model.currentKey

	return func() {
		view.model.restoreFocus = true
		view.model.currentKey = oldKey
	}
}
