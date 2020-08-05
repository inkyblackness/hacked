package about

import (
	"github.com/inkyblackness/imgui-go/v2"
)

// LicensesView handles the about display.
type LicensesView struct {
	guiScale float32

	model viewModel
}

// NewLicensesView creates a new instance for the about display.
func NewLicensesView(guiScale float32) *LicensesView {
	return &LicensesView{
		guiScale: guiScale,

		model: freshViewModel(),
	}
}

// Show requests to show the licenses view.
func (view *LicensesView) Show() {
	view.model.windowOpen = true
}

// Render requests to render the view.
func (view *LicensesView) Render() {
	if view.model.windowOpen {
		imgui.OpenPopup("Licenses")
		view.model.windowOpen = false
		imgui.SetNextWindowSize(imgui.Vec2{X: 640 * view.guiScale, Y: 480 * view.guiScale})
	}
	if imgui.BeginPopupModalV("Licenses", nil, imgui.WindowFlagsHorizontalScrollbar|imgui.WindowFlagsNoSavedSettings) {
		view.renderContent()
		imgui.EndPopup()
	}
}

func (view *LicensesView) renderContent() {
	if imgui.TreeNodeV("System Shock", imgui.TreeNodeFlagsDefaultOpen) {
		imgui.Text("System Shock by Night Dive Studios, LLC.\nOriginally by Looking Glass Technologies")
		imgui.TreePop()
	}
	add := func(title, license string) {
		if imgui.TreeNode(title) {
			imgui.Text(license)
			imgui.TreePop()
		}
	}

	add("InkyBlackness - HackEd", inkyblacknessHackedLicense)
	add("InkyBlackness - imgui", inkyblacknessImGuiLicense)
	add("go-gl/glfw - Go bindings for GLFW 3", goglGlfwLicense)
	add("go-gl/gl", goglOpenGLLicense)
	add("go-gl/mathgl", goglMathGLLicense)
	add("Dear ImGui", dearImGuiLicense)
	add("GLFW", glfwLicense)
	add("MinGW-w64", mingwLicense)

	imgui.Separator()
	if imgui.Button("OK") {
		imgui.CloseCurrentPopup()
	}
}
