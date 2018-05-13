package ui

import "github.com/inkyblackness/hacked/ui/opengl"

// Application is the root object of the graphical editor.
// It is set up by the main method.
type Application struct {
	window opengl.Window
	gl     opengl.OpenGl
}

// InitializeWindow takes the given window and attaches the callbacks.
func (app *Application) InitializeWindow(window opengl.Window) error {
	app.window = window
	app.gl = window.OpenGl()

	app.initWindowCallbacks()
	app.initOpenGl()

	return nil
}

func (app *Application) initWindowCallbacks() {
	app.window.OnRender(app.render)
}

func (app *Application) render() {
	app.gl.Clear(opengl.COLOR_BUFFER_BIT)
}

func (app *Application) initOpenGl() {
	app.gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}
