package editor

import (
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/hacked/ui/input"
	"github.com/inkyblackness/hacked/ui/opengl"
	"github.com/inkyblackness/imgui-go"
)

// Application is the root object of the graphical editor.
// It is set up by the main method.
type Application struct {
	window opengl.Window
	gl     opengl.OpenGL

	// GuiScale is applied when the window is initialized.
	GuiScale   float32
	guiContext *gui.Context
}

// InitializeWindow takes the given window and attaches the callbacks.
func (app *Application) InitializeWindow(window opengl.Window) (err error) {
	app.window = window
	app.gl = window.OpenGL()

	app.initWindowCallbacks()
	app.initOpenGL()
	err = app.initGui()
	if err != nil {
		return
	}

	return
}

func (app *Application) onWindowClosed() {
	if app.guiContext != nil {
		app.guiContext.Destroy()
		app.guiContext = nil
	}
}

func (app *Application) initWindowCallbacks() {
	app.window.OnClosing(app.onWindowClosing)
	app.window.OnClosed(app.onWindowClosed)

	app.window.OnMouseMove(app.onMouseMove)
	app.window.OnMouseScroll(app.onMouseScroll)
	app.window.OnMouseButtonDown(app.onMouseButtonDown)
	app.window.OnMouseButtonUp(app.onMouseButtonUp)

	app.window.OnRender(app.render)
}

func (app *Application) render() {
	app.guiContext.NewFrame()

	app.gl.Clear(opengl.COLOR_BUFFER_BIT)

	if imgui.BeginMainMenuBar() {
		if imgui.BeginMenu("File") {
			if imgui.MenuItem("Exit") {
				app.window.SetCloseRequest(true)
			}
			imgui.EndMenu()
		}
		imgui.EndMainMenuBar()
	}

	imgui.SetNextWindowSizeV(imgui.Vec2{400 * app.GuiScale, 300 * app.GuiScale}, imgui.ConditionOnce)
	if imgui.Begin("Project") {
		imgui.TextUnformatted("Static World Data")
		imgui.BeginChildV("ManifestEntries", imgui.Vec2{-100 * app.GuiScale, 0}, true, 0)
		for _, text := range []string{"a", "b", "c"} {
			imgui.Selectable(text)
		}
		imgui.EndChild()
		imgui.SameLine()
		imgui.BeginGroup()
		imgui.ButtonV("Add...", imgui.Vec2{-1, 0})
		imgui.ButtonV("Up", imgui.Vec2{-1, 0})
		imgui.ButtonV("Down", imgui.Vec2{-1, 0})
		imgui.ButtonV("Remove...", imgui.Vec2{-1, 0})
		imgui.EndGroup()
	}
	imgui.End()

	imgui.ShowDemoWindow(nil)

	app.guiContext.Render()
}

func (app *Application) initOpenGL() {
	app.gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

func (app *Application) initGui() (err error) {
	app.guiContext, err = gui.NewContext(app.window)
	if err != nil {
		return
	}
	app.initGuiScaling()
	app.initGuiStyle()

	return
}

func (app *Application) onWindowClosing() {

}

func (app *Application) onMouseMove(x, y float32) {
	app.guiContext.SetMousePosition(x, y)
}

func (app *Application) onMouseScroll(dx, dy float32) {
	app.guiContext.MouseScroll(dx, -dy)
}

func (app *Application) onMouseButtonDown(buttonMask uint32, modifier input.Modifier) {
	app.reportButtonChange(buttonMask, true)
}

func (app *Application) onMouseButtonUp(buttonMask uint32, modifier input.Modifier) {
	app.reportButtonChange(buttonMask, false)
}

func (app *Application) reportButtonChange(buttonMask uint32, down bool) {
	for buttonIndex := 0; buttonIndex < 32; buttonIndex++ {
		if (buttonMask & (1 << uint32(buttonIndex))) != 0 {
			app.guiContext.MouseButtonChanged(buttonIndex, down)
		}
	}
}

func (app *Application) initGuiScaling() {
	if app.GuiScale < 0.5 {
		app.GuiScale = 1.0
	} else if app.GuiScale > 10.0 {
		app.GuiScale = 10.0
	}

	imgui.CurrentStyle().ScaleAllSizes(app.GuiScale)
	imgui.CurrentIO().SetFontGlobalScale(app.GuiScale)
}

func (app *Application) initGuiStyle() {
	color := func(r, g, b byte, alpha float32) imgui.Vec4 {
		return imgui.Vec4{X: float32(r) / 255.0, Y: float32(g) / 255.0, Z: float32(b) / 255.0, W: alpha}
	}
	colorDoubleFull := func(alpha float32) imgui.Vec4 { return color(0xC4, 0x38, 0x9F, alpha) }
	colorDoubleDark := func(alpha float32) imgui.Vec4 { return color(0x31, 0x01, 0x38, alpha) }

	colorTripleFull := func(alpha float32) imgui.Vec4 { return color(0x21, 0xFF, 0x43, alpha) }
	colorTripleDark := func(alpha float32) imgui.Vec4 { return color(0x06, 0xCC, 0x94, alpha) }
	colorTripleLight := func(alpha float32) imgui.Vec4 { return color(0x51, 0x99, 0x58, alpha) }

	style := imgui.CurrentStyle()
	style.SetColor(imgui.StyleColorText, colorTripleFull(1.0))
	style.SetColor(imgui.StyleColorTextDisabled, colorTripleDark(1.0))

	style.SetColor(imgui.StyleColorWindowBg, colorDoubleDark(0.80))
	style.SetColor(imgui.StyleColorPopupBg, colorDoubleDark(0.75))

	style.SetColor(imgui.StyleColorTitleBgActive, colorTripleLight(1.0))
	style.SetColor(imgui.StyleColorFrameBg, colorTripleLight(0.54))

	style.SetColor(imgui.StyleColorFrameBgHovered, colorTripleDark(0.4))
	style.SetColor(imgui.StyleColorFrameBgActive, colorTripleDark(0.67))
	style.SetColor(imgui.StyleColorCheckMark, colorTripleDark(1.0))
	style.SetColor(imgui.StyleColorSliderGrabActive, colorTripleDark(1.0))
	style.SetColor(imgui.StyleColorButton, colorTripleDark(0.4))
	style.SetColor(imgui.StyleColorButtonHovered, colorTripleDark(1.0))
	style.SetColor(imgui.StyleColorHeader, colorTripleLight(0.70))
	style.SetColor(imgui.StyleColorHeaderHovered, colorTripleDark(0.8))
	style.SetColor(imgui.StyleColorHeaderActive, colorTripleDark(1.0))
	style.SetColor(imgui.StyleColorResizeGrip, colorTripleDark(0.25))
	style.SetColor(imgui.StyleColorResizeGripHovered, colorTripleDark(0.67))
	style.SetColor(imgui.StyleColorResizeGripActive, colorTripleDark(0.95))
	style.SetColor(imgui.StyleColorTextSelectedBg, colorTripleDark(0.35))

	style.SetColor(imgui.StyleColorSliderGrab, colorDoubleFull(1.0))
	style.SetColor(imgui.StyleColorButtonActive, colorDoubleFull(1.0))
	style.SetColor(imgui.StyleColorSeparatorHovered, colorDoubleFull(0.78))
	style.SetColor(imgui.StyleColorSeparatorActive, colorTripleLight(1.0))
}
