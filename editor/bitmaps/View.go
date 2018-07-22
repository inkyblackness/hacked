package bitmaps

import (
	"fmt"

	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
)

type bitmapInfo struct {
	title            string
	languageSpecific bool
}

var knownBitmapTypes = map[resource.ID]bitmapInfo{
	ids.MfdDataBitmaps: {"MFD Data Images", true},
}

var knownBitmapTypesOrder = []resource.ID{ids.MfdDataBitmaps}

// View provides edit controls for bitmaps.
type View struct {
	mod        *model.Mod
	imageCache *graphics.TextureCache

	modalStateMachine gui.ModalStateMachine
	clipboard         external.Clipboard
	guiScale          float32
	commander         cmd.Commander

	model viewModel
}

// NewBitmapsView returns a new instance.
func NewBitmapsView(mod *model.Mod, imageCache *graphics.TextureCache,
	modalStateMachine gui.ModalStateMachine, clipboard external.Clipboard,
	guiScale float32, commander cmd.Commander) *View {
	view := &View{
		mod:        mod,
		imageCache: imageCache,

		modalStateMachine: modalStateMachine,
		clipboard:         clipboard,
		guiScale:          guiScale,
		commander:         commander,

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
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 800 * view.guiScale, Y: 300 * view.guiScale}, imgui.ConditionOnce)
		if imgui.BeginV("Bitmaps", view.WindowOpen(), imgui.WindowFlagsNoCollapse|imgui.WindowFlagsHorizontalScrollbar) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *View) renderContent() {
	if imgui.BeginChildV("Properties", imgui.Vec2{X: 350 * view.guiScale, Y: 0}, false, 0) {
		imgui.PushItemWidth(-150 * view.guiScale)
		if imgui.BeginCombo("Bitmap Type", knownBitmapTypes[view.model.currentKey.ID].title) {
			for _, id := range knownBitmapTypesOrder {
				if imgui.SelectableV(knownBitmapTypes[id].title, id == view.model.currentKey.ID, 0, imgui.Vec2{}) {
					view.model.currentKey.ID = id
					view.model.currentKey.Index = 0
				}
			}
			imgui.EndCombo()
		}
		selectedType := knownBitmapTypes[view.model.currentKey.ID]
		if selectedType.languageSpecific {
			if imgui.BeginCombo("Language", view.model.currentKey.Lang.String()) {
				languages := resource.Languages()
				for _, lang := range languages {
					if imgui.SelectableV(lang.String(), lang == view.model.currentKey.Lang, 0, imgui.Vec2{}) {
						view.model.currentKey.Lang = lang
					}
				}
				imgui.EndCombo()
			}
		} else {
			view.model.currentKey.Lang = resource.LangAny
		}

		info, _ := ids.Info(view.model.currentKey.ID)

		gui.StepSliderInt("Index", &view.model.currentKey.Index, 0, info.MaxCount-1)

		render.TextureSelector("###"+"IndexBitmap", -1, view.guiScale, info.MaxCount,
			view.model.currentKey.Index, view.imageCache,
			func(index int) resource.Key {
				return resource.KeyOf(view.model.currentKey.ID, view.model.currentKey.Lang, index)
			},
			func(index int) string { return fmt.Sprintf("%d", index) },
			func(newValue int) {
				view.model.currentKey.Index = newValue
			})

		if imgui.Button("Export") {

		}
		imgui.SameLine()
		if imgui.Button("Import") {

		}
		imgui.PopItemWidth()
	}
	imgui.EndChild()
	imgui.SameLine()
	render.TextureImage("Big texture", view.imageCache, view.model.currentKey, imgui.Vec2{X: 320 * view.guiScale, Y: 240 * view.guiScale})
}
