package animations

import (
	"fmt"

	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
)

type animationInfo struct {
	title string
}

var knownAnimationTypes = map[resource.ID]animationInfo{
	ids.VideoMailAnimationsStart: {title: "Video Mail Parts"},
}

var knownAnimationTypesOrder = []resource.ID{
	ids.VideoMailAnimationsStart,
}

// View provides edit controls for animations.
type View struct {
	mod            *model.Mod
	imageCache     *graphics.TextureCache
	paletteCache   *graphics.PaletteCache
	animationCache *bitmap.AnimationCache

	modalStateMachine gui.ModalStateMachine
	guiScale          float32
	commander         cmd.Commander

	model viewModel
}

// NewAnimationsView returns a new instance.
func NewAnimationsView(mod *model.Mod, imageCache *graphics.TextureCache, paletteCache *graphics.PaletteCache,
	animationCache *bitmap.AnimationCache,
	modalStateMachine gui.ModalStateMachine, guiScale float32, commander cmd.Commander) *View {
	view := &View{
		mod:            mod,
		imageCache:     imageCache,
		paletteCache:   paletteCache,
		animationCache: animationCache,

		modalStateMachine: modalStateMachine,
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
		if imgui.BeginV("Animations", view.WindowOpen(), imgui.WindowFlagsNoCollapse|imgui.WindowFlagsHorizontalScrollbar) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *View) renderContent() {
	if imgui.BeginChildV("Properties", imgui.Vec2{X: 350 * view.guiScale, Y: 0}, false, 0) {
		imgui.PushItemWidth(-150 * view.guiScale)
		if imgui.BeginCombo("Animation Type", knownAnimationTypes[view.model.currentKey.ID].title) {
			for _, id := range knownAnimationTypesOrder {
				if imgui.SelectableV(knownAnimationTypes[id].title, id == view.model.currentKey.ID, 0, imgui.Vec2{}) {
					view.model.currentKey.ID = id
					view.model.currentKey.Index = 0
					view.model.currentFrame = 0
				}
			}
			imgui.EndCombo()
		}
		// selectedType := knownAnimationTypes[view.model.currentKey.ID]

		info, _ := ids.Info(view.model.currentKey.ID)

		if gui.StepSliderInt("Index", &view.model.currentKey.Index, 0, info.MaxCount-1) {
			view.model.currentFrame = 0
		}

		imgui.Separator()

		view.renderProperties()

		imgui.PopItemWidth()
	}
	imgui.EndChild()
	imgui.SameLine()

	anim, hasAnim, _ := view.currentAnimation()

	if hasAnim {
		frameKey := resource.KeyOf(anim.ResourceID, resource.LangAny, view.model.currentFrame)
		if view.cacheFrame(frameKey) {
			render.TextureImage("Frame", view.imageCache, frameKey,
				imgui.Vec2{X: float32(anim.Width) * view.guiScale, Y: float32(anim.Height) * view.guiScale})
		}
	}
}

func (view *View) cacheFrame(key resource.Key) bool {
	lastKey := resource.KeyOf(key.ID, key.Lang, 0)
	_, err := view.imageCache.Texture(lastKey)
	for index := 1; (index <= key.Index) && (err == nil); index++ {
		nextKey := resource.KeyOf(key.ID, key.Lang, index)
		_, err = view.imageCache.TextureReferenced(nextKey, lastKey)
		lastKey = nextKey
	}
	return err == nil
}

func (view *View) renderProperties() {
	anim, hasAnim, _ := view.currentAnimation()
	widthString := ""
	heightString := ""
	lastFrame := 0

	if hasAnim {
		widthString = fmt.Sprintf("%d", anim.Width)
		heightString = fmt.Sprintf("%d", anim.Height)
		for _, entry := range anim.Entries {
			lastFrame = int(entry.LastFrame)
		}
	}

	imgui.LabelText("Width", widthString)
	imgui.LabelText("Height", heightString)

	gui.StepSliderInt("Frame Index", &view.model.currentFrame, 0, lastFrame)
}

func (view *View) currentAnimation() (bitmap.Animation, bool, bool) {
	key := resource.KeyOf(view.model.currentKey.ID.Plus(view.model.currentKey.Index), resource.LangAny, 0)
	anim, err := view.animationCache.Animation(key)
	if err != nil {
		return anim, false, true
	}
	readOnly := len(view.mod.ModifiedBlocks(resource.LangAny, key.ID)) == 0
	return anim, true, readOnly
}
