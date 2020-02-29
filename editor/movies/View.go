package movies

import (
	"github.com/inkyblackness/imgui-go"

	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
)

type movieInfo struct {
	title     string
	multilang bool
}

var knownMovies = map[resource.ID]movieInfo{
	ids.MovieIntro: {title: "Intro", multilang: true},
	ids.MovieDeath: {title: "Death", multilang: false},
	ids.MovieEnd:   {title: "End", multilang: false},
}

var knownMoviesOrder = []resource.ID{ids.MovieIntro, ids.MovieDeath, ids.MovieEnd}

// View provides edit controls for animations.
type View struct {
	mod *world.Mod

	modalStateMachine gui.ModalStateMachine
	guiScale          float32
	commander         cmd.Commander

	model viewModel
}

// NewMoviesView returns a new instance.
func NewMoviesView(mod *world.Mod,
	modalStateMachine gui.ModalStateMachine, guiScale float32, commander cmd.Commander) *View {
	view := &View{
		mod: mod,

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
		if imgui.BeginV("Movies", view.WindowOpen(), imgui.WindowFlagsNoCollapse|imgui.WindowFlagsHorizontalScrollbar) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *View) renderContent() {
	if imgui.BeginChildV("Properties", imgui.Vec2{X: 350 * view.guiScale, Y: 0}, false, 0) {
		imgui.PushItemWidth(-150 * view.guiScale)
		if imgui.BeginCombo("Movie", knownMovies[view.model.currentKey.ID].title) {
			for _, id := range knownMoviesOrder {
				if imgui.SelectableV(knownMovies[id].title, id == view.model.currentKey.ID, 0, imgui.Vec2{}) {
					view.model.currentKey.ID = id
					view.model.currentKey.Index = 0
				}
			}
			imgui.EndCombo()
		}

		if knownMovies[view.model.currentKey.ID].multilang {
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
			// TODO: skip a line
		}

		//		info, _ := ids.Info(view.model.currentKey.ID)
		//		if gui.StepSliderInt("Index", &view.model.currentKey.Index, 0, info.MaxCount-1) {
		//			view.model.currentFrame = 0
		//		}

		imgui.Separator()

		view.renderProperties()

		imgui.PopItemWidth()
	}
	imgui.EndChild()
	imgui.SameLine()

	if imgui.BeginChildV("Frames", imgui.Vec2{X: -1, Y: 0}, false, 0) {

	}
	imgui.EndChild()
}

func (view *View) renderProperties() {
	// gui.StepSliderInt("Frame Index", &view.model.currentFrame, 0, lastFrame)
}
