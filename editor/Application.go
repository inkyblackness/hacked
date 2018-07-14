package editor

import (
	"fmt"

	"github.com/inkyblackness/hacked/editor/about"
	"github.com/inkyblackness/hacked/editor/archives"
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/event"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/levels"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/editor/project"
	"github.com/inkyblackness/hacked/editor/texts"
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/hacked/ui/input"
	"github.com/inkyblackness/hacked/ui/opengl"
	"github.com/inkyblackness/imgui-go"
)

// Application is the root object of the graphical editor.
// It is set up by the main method.
type Application struct {
	// Version identifies the build of the application.
	Version string

	window    opengl.Window
	clipboard clipboardAdapter
	gl        opengl.OpenGL

	// FontFile specifies the font to use.
	FontFile string
	// FontSize specifies the font size to use.
	FontSize float32
	// GuiScale is applied when the window is initialized.
	GuiScale   float32
	guiContext *gui.Context

	lastModifier input.Modifier
	lastMouseX   float32
	lastMouseY   float32

	eventQueue      event.Queue
	eventDispatcher *event.Dispatcher

	cmdStack      *cmd.Stack
	mod           *model.Mod
	cp            text.Codepage
	textLineCache *text.Cache
	textPageCache *text.Cache
	paletteCache  *graphics.PaletteCache
	textureCache  *graphics.TextureCache

	mapDisplay *levels.MapDisplay

	levels [archive.MaxLevels]*level.Level

	projectView      *project.View
	archiveView      *archives.View
	levelControlView *levels.ControlView
	levelTilesView   *levels.TilesView
	levelObjectsView *levels.ObjectsView
	textsView        *texts.View
	aboutView        *about.View

	failureMessage string
	failurePending bool
}

// InitializeWindow takes the given window and attaches the callbacks.
func (app *Application) InitializeWindow(window opengl.Window) (err error) {
	app.window = window
	app.clipboard.window = window
	app.gl = window.OpenGL()

	app.initSignalling()
	app.initWindowCallbacks()
	app.initOpenGL()
	err = app.initGui()
	if err != nil {
		return
	}

	app.initModel()
	app.initView()

	app.onWindowResize(app.window.Size())

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
	app.window.OnFileDropCallback(app.onFilesDropped)

	app.window.OnResize(app.onWindowResize)

	app.window.OnKey(app.onKey)
	app.window.OnModifier(app.onModifier)

	app.window.OnMouseMove(app.onMouseMove)
	app.window.OnMouseScroll(app.onMouseScroll)
	app.window.OnMouseButtonDown(app.onMouseButtonDown)
	app.window.OnMouseButtonUp(app.onMouseButtonUp)

	app.window.OnRender(app.render)
}

func (app *Application) render() {
	app.dispatchEvents()
	app.guiContext.NewFrame()

	app.gl.Clear(opengl.COLOR_BUFFER_BIT)

	app.renderMainMenu()

	app.projectView.Render()
	app.archiveView.Render()
	activeLevel := app.levels[app.levelControlView.SelectedLevel()]
	app.levelControlView.Render(activeLevel)
	app.levelTilesView.Render(activeLevel)
	app.levelObjectsView.Render(activeLevel)
	app.textsView.Render()

	paletteTexture, _ := app.paletteCache.Palette(0)
	app.mapDisplay.Render(app.mod.ObjectProperties(), activeLevel,
		paletteTexture, app.textureCache.Texture,
		app.levelTilesView.TextureDisplay(), app.levelTilesView.ColorDisplay(activeLevel))

	// imgui.ShowDemoWindow(nil)

	app.handleFailure()
	app.aboutView.Render()

	app.guiContext.Render(app.bitmapTextureForUI)
}

func (app *Application) initOpenGL() {
	app.gl.Disable(opengl.DEPTH_TEST)
	app.gl.Enable(opengl.BLEND)
	app.gl.BlendFunc(opengl.SRC_ALPHA, opengl.ONE_MINUS_SRC_ALPHA)
	app.gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

func (app *Application) initGui() (err error) {
	app.initGuiSizes()
	param := gui.ContextParameters{
		FontFile: app.FontFile,
		FontSize: app.FontSize,
	}
	app.guiContext, err = gui.NewContext(app.window, param)
	if err != nil {
		return
	}
	app.initGuiStyle()

	app.mapDisplay = levels.NewMapDisplay(app.gl, app.GuiScale,
		app.gameTexture,
		&app.eventQueue, app.eventDispatcher)

	return
}

func (app *Application) gameTexture(index int) (*graphics.BitmapTexture, error) {
	key := resource.KeyOf(ids.LargeTextures.Plus(index), resource.LangAny, 0)
	return app.textureCache.Texture(key)
}

func (app *Application) bitmapTextureForUI(textureID imgui.TextureID) (palette uint32, texture uint32) {
	paletteTexture, _ := app.paletteCache.Palette(0)
	if paletteTexture == nil {
		return 0, 0
	}

	lang := resource.Language((textureID >> 32) & 0xFF)
	resourceID := resource.ID((textureID >> 16) & 0xFFFF)
	blockIndex := int(textureID & 0xFFFF)
	key := resource.KeyOf(resourceID, lang, blockIndex)
	tex, err := app.textureCache.Texture(key)
	if err != nil {
		return 0, 0
	}

	return paletteTexture.Handle(), tex.Handle()
}

func (app *Application) onWindowClosing() {

}

func (app *Application) onWindowResize(width int, height int) {
	app.mapDisplay.WindowResized(width, height)
	app.gl.Viewport(0, 0, int32(width), int32(height))
}

func (app *Application) onFilesDropped(names []string) {
	app.projectView.HandleFiles(names)
}

func (app *Application) onKey(key input.Key, modifier input.Modifier) {
	app.lastModifier = modifier
	if key == input.KeyUndo {
		app.tryUndo()
	} else if key == input.KeyRedo {
		app.tryRedo()
	} else if key == input.KeyF1 {
		*app.projectView.WindowOpen() = !*app.projectView.WindowOpen()
	} else if key == input.KeyF2 {
		*app.levelControlView.WindowOpen() = !*app.levelControlView.WindowOpen()
	} else if key == input.KeyF3 {
		*app.levelTilesView.WindowOpen() = !*app.levelTilesView.WindowOpen()
	} else if key == input.KeyF4 {
		*app.levelObjectsView.WindowOpen() = !*app.levelObjectsView.WindowOpen()
	}
}

func (app *Application) onModifier(modifier input.Modifier) {
	app.lastModifier = modifier
}

func (app *Application) tryUndo() {
	if !app.cmdStack.CanUndo() {
		return
	}
	err := app.modifyModByCommand(app.cmdStack.Undo)
	if err != nil {
		app.onFailure("Undo", "", err)
	}
}

func (app *Application) tryRedo() {
	if !app.cmdStack.CanRedo() {
		return
	}
	err := app.modifyModByCommand(app.cmdStack.Redo)
	if err != nil {
		app.onFailure("Redo", "", err)
	}
}

func (app *Application) modifyModByCommand(modifier func(cmd.Transaction) error) (err error) {
	app.mod.Modify(func(trans *model.ModTransaction) {
		err = modifier(trans)
	})
	return
}

func (app *Application) onMouseMove(x, y float32) {
	app.lastMouseX = x
	app.lastMouseY = y
	app.guiContext.SetMousePosition(x, y)
	if !app.guiContext.IsUsingMouse() {
		app.mapDisplay.MouseMoved(x, y)
	}
}

func (app *Application) onMouseScroll(dx, dy float32) {
	if !app.guiContext.IsUsingMouse() {
		app.mapDisplay.MouseScrolled(app.lastMouseX, app.lastMouseY, dx, dy, app.lastModifier)
	}
	app.guiContext.MouseScroll(dx, dy)
}

func (app *Application) onMouseButtonDown(buttonMask uint32, modifier input.Modifier) {
	if !app.guiContext.IsUsingMouse() {
		app.mapDisplay.MouseButtonDown(app.lastMouseX, app.lastMouseY, buttonMask)
	}
	app.reportButtonChange(buttonMask, true)
}

func (app *Application) onMouseButtonUp(buttonMask uint32, modifier input.Modifier) {
	if !app.guiContext.IsUsingMouse() {
		app.mapDisplay.MouseButtonUp(app.lastMouseX, app.lastMouseY, buttonMask, modifier)
	}
	app.reportButtonChange(buttonMask, false)
}

func (app *Application) reportButtonChange(buttonMask uint32, down bool) {
	for buttonIndex := 0; buttonIndex < 32; buttonIndex++ {
		if (buttonMask & (1 << uint32(buttonIndex))) != 0 {
			app.guiContext.MouseButtonChanged(buttonIndex, down)
		}
	}
}

func (app *Application) initGuiSizes() {
	if app.GuiScale < 0.5 {
		app.GuiScale = 1.0
	} else if app.GuiScale > 10.0 {
		app.GuiScale = 10.0
	}

	if app.FontSize <= 0 {
		app.FontSize = 16.0
	}

	app.FontSize *= app.GuiScale
}

func (app *Application) initGuiStyle() {
	if len(app.FontFile) == 0 {
		imgui.CurrentIO().SetFontGlobalScale(app.GuiScale)
	}
	imgui.CurrentStyle().ScaleAllSizes(app.GuiScale)

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

func (app *Application) initSignalling() {
	app.eventDispatcher = event.NewDispatcher()
	app.cmdStack = new(cmd.Stack)
}

func (app *Application) initModel() {
	app.mod = model.NewMod(app.resourcesChanged, app.modReset)

	app.cp = text.DefaultCodepage()
	app.textLineCache = text.NewLineCache(app.cp, app.mod)
	app.textPageCache = text.NewPageCache(app.cp, app.mod)

	for i := 0; i < archive.MaxLevels; i++ {
		app.levels[i] = level.NewLevel(ids.LevelResourcesStart, i, app.mod)
	}

	app.paletteCache = graphics.NewPaletteCache(app.gl, app.mod)
	app.textureCache = graphics.NewTextureCache(app.gl, app.mod)
}

func (app *Application) resourcesChanged(modifiedIDs []resource.ID, failedIDs []resource.ID) {
	app.textLineCache.InvalidateResources(modifiedIDs)
	app.textPageCache.InvalidateResources(modifiedIDs)
	for _, lvl := range app.levels {
		lvl.InvalidateResources(modifiedIDs)
	}
	app.paletteCache.InvalidateResources(modifiedIDs)
	app.textureCache.InvalidateResources(modifiedIDs)
}

func (app *Application) modReset() {
	app.cmdStack = new(cmd.Stack)
}

func (app *Application) initView() {
	app.projectView = project.NewView(app.mod, app.GuiScale, app)
	app.archiveView = archives.NewArchiveView(app.mod, app.GuiScale, app)
	app.levelControlView = levels.NewControlView(app.mod, app.GuiScale, app.textLineCache, app, &app.eventQueue, app.eventDispatcher)
	app.levelTilesView = levels.NewTilesView(app.mod, app.GuiScale, app, &app.eventQueue, app.eventDispatcher)
	app.levelObjectsView = levels.NewObjectsView(app.mod, app.GuiScale, app.textLineCache, app, &app.eventQueue, app.eventDispatcher)
	app.textsView = texts.NewTextsView(app.mod, app.textLineCache, app.textPageCache, app.cp, app.clipboard, app.GuiScale, app)
	app.aboutView = about.NewView(app.clipboard, app.GuiScale, app.Version)

	app.eventDispatcher.RegisterHandler(app.onLevelObjectRequestCreateEvent)
}

// Queue requests to perform the given command.
func (app *Application) Queue(command cmd.Command) {
	err := app.modifyModByCommand(func(trans cmd.Transaction) error {
		return app.cmdStack.Perform(command, trans)
	})
	if err != nil {
		app.onFailure("Command", "", err)
	}
}

func (app *Application) dispatchEvents() {
	for iteration := 0; (iteration < 100) && !app.eventQueue.IsEmpty(); iteration++ {
		app.eventQueue.ForwardTo(app.eventDispatcher)
	}
}

func (app *Application) onFailure(source string, details string, err error) {
	app.failurePending = true
	app.failureMessage = fmt.Sprintf("Source: %v\nDetails: %v\nError: %v", source, details, err)
}

func (app *Application) onLevelObjectRequestCreateEvent(evt levels.ObjectRequestCreateEvent) {
	lvl := app.levels[app.levelControlView.SelectedLevel()]
	app.levelObjectsView.RequestCreateObject(lvl, evt.Pos)
}

func (app *Application) renderMainMenu() {
	windowEntry := func(name string, shortcut string, isOpen *bool) {
		if imgui.MenuItemV(name, shortcut, *isOpen, true) {
			*isOpen = !*isOpen
		}
	}

	if imgui.BeginMainMenuBar() {
		if imgui.BeginMenu("File") {
			windowEntry("Project", "F1", app.projectView.WindowOpen())
			imgui.Separator()
			if imgui.MenuItem("Exit") {
				app.window.SetCloseRequest(true)
			}
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Edit") {
			if imgui.MenuItemV("Undo", "Ctrl+Z", false, app.cmdStack.CanUndo()) {
				app.tryUndo()
			}
			if imgui.MenuItemV("Redo", "Ctrl+Y / Ctrl+Shift+Z", false, app.cmdStack.CanRedo()) {
				app.tryRedo()
			}
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Window") {
			windowEntry("Archive", "", app.archiveView.WindowOpen())
			windowEntry("Level Control", "F2", app.levelControlView.WindowOpen())
			windowEntry("Level Tiles", "F3", app.levelTilesView.WindowOpen())
			windowEntry("Level Objects", "F4", app.levelObjectsView.WindowOpen())
			windowEntry("Texts", "", app.textsView.WindowOpen())
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Help") {
			if imgui.MenuItem("About...") {
				app.aboutView.Show()
			}
			imgui.EndMenu()
		}
		imgui.EndMainMenuBar()
	}
}

func (app *Application) handleFailure() {
	if app.failurePending {
		imgui.OpenPopup("Failure Message")
		app.failurePending = false
	}
	if imgui.BeginPopupModal("Failure Message") {
		imgui.Text(`Something went wrong. This is bad and I am sorry.

You have the option to "Ignore" this and hope for the best.
This action also clears the undo/redo buffer.

Or you can simply "Exit" the application and then restart it.
This action loses any pending changes.

Perhaps you can make something with the details of the error below.
If you can reproduce this, please make a screenshot and
report it with details on how to reproduce it on the
http://www.systemshock.org forums. Thank you!
`)
		imgui.Text("Version: " + app.Version)
		imgui.Separator()
		imgui.Text(app.failureMessage)
		imgui.Separator()
		if imgui.Button("Ignore") {
			app.failureMessage = ""
			app.cmdStack = new(cmd.Stack)
			imgui.CloseCurrentPopup()
		}
		imgui.SameLine()
		if imgui.Button("Exit") {
			app.window.SetCloseRequest(true)
		}
		imgui.EndPopup()
	}
}
