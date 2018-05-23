package editor

import (
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/hacked/ui/opengl"
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

func (app *Application) windowClosing() {
	if app.guiContext != nil {
		app.guiContext.Destroy()
		app.guiContext = nil
	}
}

func (app *Application) initWindowCallbacks() {
	app.window.OnClosing(app.windowClosing)
	app.window.OnRender(app.render)
}

func (app *Application) render() {
	app.guiContext.NewFrame()

	app.gl.Clear(opengl.COLOR_BUFFER_BIT)

	app.guiContext.Render()
}

func (app *Application) initOpenGl() {
	app.gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

func (app *Application) initGui() (err error) {
	app.guiContext, err = gui.NewContext(app.window)
	return
}
