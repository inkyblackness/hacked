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
	gl     opengl.OpenGl

	guiContext *gui.Context
}

// InitializeWindow takes the given window and attaches the callbacks.
func (app *Application) InitializeWindow(window opengl.Window) (err error) {
	app.window = window
	app.gl = window.OpenGl()

	app.initWindowCallbacks()
	app.initOpenGl()
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
	}
	imgui.EndMainMenuBar()

	app.guiContext.Render()
}

func (app *Application) initOpenGl() {
	app.gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

func (app *Application) initGui() (err error) {
	app.guiContext, err = gui.NewContext(app.window)
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
