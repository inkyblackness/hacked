package editor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/inkyblackness/imgui-go/v2"

	"github.com/inkyblackness/hacked/editor/about"
	"github.com/inkyblackness/hacked/editor/animations"
	"github.com/inkyblackness/hacked/editor/archives"
	"github.com/inkyblackness/hacked/editor/bitmaps"
	"github.com/inkyblackness/hacked/editor/event"
	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/levels"
	"github.com/inkyblackness/hacked/editor/messages"
	"github.com/inkyblackness/hacked/editor/movies"
	"github.com/inkyblackness/hacked/editor/objects"
	"github.com/inkyblackness/hacked/editor/project"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/editor/sounds"
	"github.com/inkyblackness/hacked/editor/texts"
	"github.com/inkyblackness/hacked/editor/textures"
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ss1/content/sound"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/edit/media"
	"github.com/inkyblackness/hacked/ss1/edit/undoable"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/hacked/ui/input"
	"github.com/inkyblackness/hacked/ui/opengl"
)

type projectState struct {
	Settings         *edit.ProjectSettings `json:",omitempty"`
	OpenWindows      []string              `json:",omitempty"`
	ActiveLevelIndex *int                  `json:",omitempty"`
}

func projectStateFromFile(filename string) (projectState, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return projectState{}, errors.New("could not open file")
	}
	var state projectState
	err = json.Unmarshal(data, &state)
	if err != nil {
		return projectState{}, errors.New("could not read file")
	}
	return state, nil
}

// SaveTo stores the state in a file with given filename.
func (state projectState) SaveTo(filename string) error {
	bytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, 0640)
}

type lastProjectState struct {
	projectState
	LastProject string `json:",omitempty"`
}

func lastProjectStateFromFile(filename string) lastProjectState {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return lastProjectState{}
	}
	var state lastProjectState
	err = json.Unmarshal(data, &state)
	if err != nil {
		return lastProjectState{}
	}
	return state
}

// SaveTo stores the state in a file with given filename.
func (state lastProjectState) SaveTo(filename string) error {
	bytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, 0640)
}

// Application is the root object of the graphical editor.
// It is set up by the main method.
type Application struct {
	// Version identifies the build of the application.
	Version string

	// ConfigDir is the base path to store configuration in.
	ConfigDir string

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

	txnBuilder       cmd.TransactionBuilder
	cmdStack         *cmd.Stack
	mod              *world.Mod
	cp               text.Codepage
	textLineCache    *text.Cache
	textPageCache    *text.Cache
	messagesCache    *text.ElectronicMessageCache
	paletteCache     *graphics.PaletteCache
	textureCache     *graphics.TextureCache
	frameCache       *graphics.FrameCache
	animationCache   *bitmap.AnimationCache
	movieCache       *movie.Cache
	soundEffectCache *sound.SoundEffectCache

	mapDisplay *levels.MapDisplay

	levels [archive.MaxLevels]*level.Level

	projectService   *edit.ProjectService
	gameStateService *edit.GameStateService

	projectView      *project.View
	archiveView      *archives.View
	levelControlView *levels.ControlView
	levelTilesView   *levels.TilesView
	levelObjectsView *levels.ObjectsView
	messagesView     *messages.View
	textsView        *texts.View
	bitmapsView      *bitmaps.View
	texturesView     *textures.View
	animationsView   *animations.View
	moviesView       *movies.View
	soundEffectsView *sounds.View
	objectsView      *objects.View
	aboutView        *about.View
	licensesView     *about.LicensesView

	modalState gui.ModalStateWrapper

	failureMessage string
	failurePending bool
}

// InitializeWindow takes the given window and attaches the callbacks.
func (app *Application) InitializeWindow(window opengl.Window) (err error) {
	app.window = window
	app.clipboard.window = window
	app.gl = window.OpenGL()

	app.window.RestoreState(opengl.WindowStateFromFile(app.windowStateConfigFilename(), app.window.StateSnapshot()))

	app.initSignalling()
	app.initWindowCallbacks()
	app.initOpenGL()
	err = app.initGui()
	if err != nil {
		return
	}

	app.initModel()
	app.initView()
	app.restoreWorkspace()

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

	app.window.OnKeyPress(app.onKeyPress)
	app.window.OnKeyRelease(app.onKeyRelease)
	app.window.OnCharCallback(app.onChar)
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
	app.updateAutoSave()

	app.archiveView.Render()
	activeLevel := app.levels[app.levelControlView.SelectedLevel()]
	app.levelControlView.Render(activeLevel)
	app.levelTilesView.Render(activeLevel)
	app.levelObjectsView.Render(activeLevel)
	app.messagesView.Render()
	app.textsView.Render()
	app.bitmapsView.Render()
	app.texturesView.Render()
	app.animationsView.Render()
	app.moviesView.Render()
	app.soundEffectsView.Render()
	app.objectsView.Render()

	paletteTexture, _ := app.paletteCache.Palette(0)
	app.mapDisplay.Render(app.mod.ObjectProperties(), activeLevel,
		paletteTexture, app.textureCache.Texture,
		app.levelTilesView.TextureDisplay(), app.levelTilesView.ColorDisplay(activeLevel))

	app.handleFailure()
	app.aboutView.Render()
	app.licensesView.Render()

	app.modalState.Render()

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
		ConfigDir: app.ConfigDir,
		FontFile:  app.FontFile,
		FontSize:  app.FontSize,
	}
	app.guiContext, err = gui.NewContext(app.window, param)
	if err != nil {
		return
	}

	imguiIO := imgui.CurrentIO()
	imguiIO.KeyMap(imgui.KeyTab, int(input.KeyTab))
	imguiIO.KeyMap(imgui.KeyLeftArrow, int(input.KeyLeft))
	imguiIO.KeyMap(imgui.KeyRightArrow, int(input.KeyRight))
	imguiIO.KeyMap(imgui.KeyUpArrow, int(input.KeyUp))
	imguiIO.KeyMap(imgui.KeyDownArrow, int(input.KeyDown))
	imguiIO.KeyMap(imgui.KeyPageUp, int(input.KeyPageUp))
	imguiIO.KeyMap(imgui.KeyPageDown, int(input.KeyPageDown))
	imguiIO.KeyMap(imgui.KeyHome, int(input.KeyHome))
	imguiIO.KeyMap(imgui.KeyEnd, int(input.KeyEnd))
	imguiIO.KeyMap(imgui.KeyInsert, int(input.KeyInsert))
	imguiIO.KeyMap(imgui.KeyDelete, int(input.KeyDelete))
	imguiIO.KeyMap(imgui.KeyBackspace, int(input.KeyBackspace))
	imguiIO.KeyMap(imgui.KeyEnter, int(input.KeyEnter))
	imguiIO.KeyMap(imgui.KeyEscape, int(input.KeyEscape))
	imguiIO.KeyMap(imgui.KeyA, int(input.KeyA))
	imguiIO.KeyMap(imgui.KeyC, int(input.KeyC))
	imguiIO.KeyMap(imgui.KeyV, int(input.KeyV))
	imguiIO.KeyMap(imgui.KeyX, int(input.KeyX))

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
	idType := byte(textureID >> 48)
	switch idType {
	case render.BitmapTextureTypeResource:
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
	case render.BitmapTextureTypeFrame:
		key := graphics.FrameCacheKey(textureID & 0xFFFF)
		return app.frameCache.HandlesForKey(key)
	}
	return 0, 0
}

func (app *Application) onWindowClosing() {
	windowState := app.window.StateSnapshot()
	_ = windowState.SaveTo(app.windowStateConfigFilename())

	app.saveWorkspace()
}

func (app *Application) onWindowResize(width int, height int) {
	app.mapDisplay.WindowResized(width, height)
	app.gl.Viewport(0, 0, int32(width), int32(height))
}

func (app *Application) onFilesDropped(names []string) {
	app.modalState.HandleFiles(names)
}

func (app *Application) onKeyPress(key input.Key, modifier input.Modifier) {
	app.lastModifier = modifier

	usedGlobalShortcut := true
	switch {
	case key == input.KeyEscape:
		app.modalState.SetState(nil)
	case key == input.KeyUndo:
		app.tryUndo()
	case key == input.KeyRedo:
		app.tryRedo()
	case key == input.KeyNew:
		app.newProject()
	case key == input.KeySave:
		app.projectView.StartSavingMod()
	case key == input.KeyF1 && modifier.IsClear():
		*app.projectView.WindowOpen() = !*app.projectView.WindowOpen()
	case key == input.KeyF2 && modifier.IsClear():
		*app.levelControlView.WindowOpen() = !*app.levelControlView.WindowOpen()
	case key == input.KeyF3 && modifier.IsClear():
		*app.levelTilesView.WindowOpen() = !*app.levelTilesView.WindowOpen()
	case key == input.KeyF4 && modifier.IsClear():
		*app.levelObjectsView.WindowOpen() = !*app.levelObjectsView.WindowOpen()
	case key == input.KeyF5 && modifier.IsClear():
		*app.messagesView.WindowOpen() = !*app.messagesView.WindowOpen()
	default:
		usedGlobalShortcut = false
	}
	if usedGlobalShortcut {
		return
	}

	if app.shouldReportKeyToImgui(key) {
		imgui.CurrentIO().KeyPress(int(key))
		app.updateKeyModifier()
	}
}

func (app *Application) onKeyRelease(key input.Key, modifier input.Modifier) {
	app.lastModifier = modifier

	if app.shouldReportKeyToImgui(key) {
		imgui.CurrentIO().KeyRelease(int(key))
		app.updateKeyModifier()
	}
}

func (app *Application) shouldReportKeyToImgui(key input.Key) bool {
	keyIsModifier := key.AsModifier() != input.ModNone
	return keyIsModifier || app.guiContext.IsUsingKeyboard()
}

func (app *Application) onChar(char rune) {
	if app.guiContext.IsUsingKeyboard() {
		imgui.CurrentIO().AddInputCharacters(string(char))
		return
	}

	activeLevel := app.levels[app.levelControlView.SelectedLevel()]
	switch char {
	case 'v':
		app.levelObjectsView.PlaceSelectedObjectsOnFloor(activeLevel)
	case 'f':
		app.levelObjectsView.PlaceSelectedObjectsOnEyeLevel(activeLevel)
	case 'r':
		app.levelObjectsView.PlaceSelectedObjectsOnCeiling(activeLevel)
	}
}

func (app *Application) updateKeyModifier() {
	imguiIO := imgui.CurrentIO()
	imguiIO.KeyCtrl(int(input.KeyControl), int(input.KeyControl))
	imguiIO.KeyShift(int(input.KeyShift), int(input.KeyShift))
	imguiIO.KeyAlt(int(input.KeyAlt), int(input.KeyAlt))
	imguiIO.KeySuper(int(input.KeySuper), int(input.KeySuper))
}

func (app *Application) onModifier(modifier input.Modifier) {
	app.lastModifier = modifier
}

func (app *Application) modalActive() bool {
	return (app.modalState.State != nil) || (len(app.failureMessage) > 0)
}

func (app *Application) tryUndo() {
	if !app.cmdStack.CanUndo() || app.modalActive() {
		return
	}
	imgui.ClearActiveID()
	err := app.projectService.ModifyModWith(app.cmdStack.Undo)
	if err != nil {
		app.onFailure("Undo", "", err)
	}
}

func (app *Application) tryRedo() {
	if !app.cmdStack.CanRedo() || app.modalActive() {
		return
	}
	imgui.ClearActiveID()
	err := app.projectService.ModifyModWith(app.cmdStack.Redo)
	if err != nil {
		app.onFailure("Redo", "", err)
	}
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

	style.SetColor(imgui.StyleColorTab, colorTripleLight(0.54))
	style.SetColor(imgui.StyleColorTabHovered, colorTripleLight(0.75))
	style.SetColor(imgui.StyleColorTabActive, colorTripleLight(1.0))

	style.SetColor(imgui.StyleColorSliderGrab, colorDoubleFull(1.0))
	style.SetColor(imgui.StyleColorButtonActive, colorDoubleFull(1.0))
	style.SetColor(imgui.StyleColorSeparatorHovered, colorDoubleFull(0.78))
	style.SetColor(imgui.StyleColorSeparatorActive, colorTripleLight(1.0))
}

func (app *Application) initSignalling() {
	app.eventDispatcher = event.NewDispatcher()
	app.cmdStack = new(cmd.Stack)
	app.txnBuilder.Commander = app
}

func (app *Application) initModel() {
	app.mod = world.NewMod(app.resourcesChanged, app.modReset)

	app.cp = text.DefaultCodepage()
	app.textLineCache = text.NewLineCache(app.cp, app.mod)
	app.textPageCache = text.NewPageCache(app.cp, app.mod)
	app.messagesCache = text.NewElectronicMessageCache(app.cp, app.mod)
	app.movieCache = movie.NewCache(app.cp, app.mod)
	app.soundEffectCache = sound.NewSoundCache(app.mod)

	for i := 0; i < archive.MaxLevels; i++ {
		app.levels[i] = level.NewLevel(ids.LevelResourcesStart, i, app.mod)
	}

	app.paletteCache = graphics.NewPaletteCache(app.gl, app.mod)
	app.textureCache = graphics.NewTextureCache(app.gl, app.mod)
	app.frameCache = graphics.NewFrameCache(app.gl)
	app.animationCache = bitmap.NewAnimationCache(app.mod)
}

func (app *Application) resourcesChanged(modifiedIDs []resource.ID, failedIDs []resource.ID) {
	app.textLineCache.InvalidateResources(modifiedIDs)
	app.textPageCache.InvalidateResources(modifiedIDs)
	app.messagesCache.InvalidateResources(modifiedIDs)
	app.movieCache.InvalidateResources(modifiedIDs)
	app.soundEffectCache.InvalidateResources(modifiedIDs)
	for _, lvl := range app.levels {
		lvl.InvalidateResources(modifiedIDs)
	}
	app.paletteCache.InvalidateResources(modifiedIDs)
	app.textureCache.InvalidateResources(modifiedIDs)
	app.animationCache.InvalidateResources(modifiedIDs)
}

func (app *Application) modReset() {
	app.cmdStack = new(cmd.Stack)
}

// nolint: lll
func (app *Application) initView() {
	textViewer := media.NewTextViewerService(app.textLineCache, app.textPageCache, app.mod)
	textSetter := media.NewTextSetterService(app.cp)
	audioViewer := media.NewAudioViewerService(app.movieCache, app.mod)
	audioSetter := media.NewAudioSetterService()
	movieViewer := media.NewMovieViewerService(app.movieCache, app.mod)
	movieSetter := media.NewMovieSetterService(app.cp)
	soundEffectViewer := media.NewSoundViewerService(app.soundEffectCache, app.mod)
	soundEffectSetter := media.NewSoundSetterService()
	soundEffectService := undoable.NewSoundEffectService(edit.NewSoundEffectService(soundEffectViewer, soundEffectSetter), app)
	augmentedTextService := undoable.NewAugmentedTextService(edit.NewAugmentedTextService(textViewer, textSetter, audioViewer, audioSetter), app)
	movieService := undoable.NewMovieService(edit.NewMovieService(app.cp, movieViewer, movieSetter), app)

	app.projectService = edit.NewProjectService(&app.txnBuilder, app.mod)
	app.gameStateService = edit.NewGameStateService(&app.txnBuilder)

	app.projectView = project.NewView(app.projectService, &app.modalState, app.GuiScale, &app.txnBuilder)
	app.archiveView = archives.NewArchiveView(app.gameStateService, app.mod, app.textLineCache, app.cp, app.GuiScale, app)
	app.levelControlView = levels.NewControlView(app.mod, app.GuiScale, app.textLineCache, app.textureCache, app, &app.eventQueue, app.eventDispatcher)
	app.levelTilesView = levels.NewTilesView(app.mod, app.GuiScale, app.textLineCache, app.textureCache, app, &app.eventQueue, app.eventDispatcher)
	app.levelObjectsView = levels.NewObjectsView(app.mod, app.GuiScale, app.textLineCache, app.textureCache, app, &app.eventQueue, app.eventDispatcher)
	app.messagesView = messages.NewMessagesView(app.mod, app.messagesCache, app.cp, app.movieCache, app.textureCache, &app.modalState, app.clipboard, app.GuiScale, app)
	app.textsView = texts.NewTextsView(augmentedTextService, &app.modalState, app.clipboard, app.GuiScale)
	app.bitmapsView = bitmaps.NewBitmapsView(app.mod, app.textureCache, app.paletteCache, &app.modalState, app.clipboard, app.GuiScale, app)
	app.texturesView = textures.NewTexturesView(app.mod, app.textLineCache, app.cp, app.textureCache, app.paletteCache, &app.modalState, app.clipboard, app.GuiScale, app)
	app.animationsView = animations.NewAnimationsView(app.mod, app.textureCache, app.paletteCache, app.animationCache, &app.modalState, app.GuiScale, app)
	app.moviesView = movies.NewMoviesView(app.mod, app.frameCache, movieService, &app.modalState, app.GuiScale, app)
	app.soundEffectsView = sounds.NewSoundEffectsView(soundEffectService, &app.modalState, app.GuiScale)
	app.objectsView = objects.NewView(app.mod, app.textLineCache, app.cp, app.textureCache, app.paletteCache, &app.modalState, app.clipboard, app.GuiScale, app)
	app.aboutView = about.NewView(app.clipboard, app.GuiScale, app.Version)
	app.licensesView = about.NewLicensesView(app.GuiScale)

	app.eventDispatcher.RegisterHandler(app.onLevelObjectRequestCreateEvent)
}

// Queue requests to perform the given command.
func (app *Application) Queue(command cmd.Command) {
	err := app.projectService.ModifyModWith(func(modder world.Modder) error {
		return app.cmdStack.Perform(command, modder)
	})
	if err != nil {
		app.onFailure("command", "", err)
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
			if imgui.MenuItemV("New Project", "Ctrl+N", false, true) {
				app.newProject()
			}
			if imgui.MenuItemV("Load Project", "", false, true) {
				app.loadProject()
			}
			if imgui.MenuItemV("Save Project", "", false, true) {
				app.saveProject()
			}
			imgui.Separator()
			if imgui.MenuItemV("Save Mod", "Ctrl+S", false, true) {
				app.projectView.StartSavingMod()
			}
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
			windowEntry("Messages", "F5", app.messagesView.WindowOpen())
			windowEntry("Texts", "", app.textsView.WindowOpen())
			windowEntry("Bitmaps", "", app.bitmapsView.WindowOpen())
			windowEntry("Textures", "", app.texturesView.WindowOpen())
			windowEntry("Animations", "", app.animationsView.WindowOpen())
			windowEntry("Movies", "", app.moviesView.WindowOpen())
			windowEntry("Sound Effects", "", app.soundEffectsView.WindowOpen())
			windowEntry("Game Objects", "", app.objectsView.WindowOpen())
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Help") {
			if imgui.MenuItem("About...") {
				app.aboutView.Show()
			}
			if imgui.MenuItem("Licenses") {
				app.licensesView.Show()
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

const settingsFileExtension = "hacked-project"

func projectFileTypes() []external.TypeInfo {
	return []external.TypeInfo{
		{
			Title:      "HackEd Project File (*." + settingsFileExtension + ")",
			Extensions: []string{settingsFileExtension},
		},
	}
}

func (app *Application) newProject() {
	app.saveStoredProject()
	app.restoreProjectState(projectState{}, "")
}

func (app *Application) loadProject() {
	app.saveStoredProject()

	external.LoadFile(&app.modalState, projectFileTypes(), func(filename string) error {
		state, err := projectStateFromFile(filename)
		if err != nil {
			return err
		}
		app.restoreProjectState(state, filename)
		return nil
	})
}

func (app *Application) saveStoredProject() {
	currentFilename := app.projectService.CurrentStateFilename()
	if len(currentFilename) > 0 {
		currentState := app.currentProjectState()
		_ = currentState.SaveTo(currentFilename)
	}
}

func (app *Application) saveProject() {
	currentState := app.currentProjectState()
	currentFilename := app.projectService.CurrentStateFilename()
	if (len(currentFilename) == 0) || (currentState.SaveTo(currentFilename) != nil) {
		external.SaveFile(&app.modalState, projectFileTypes(), func(filename string) error {
			completeFilename := filename
			dotExtension := "." + settingsFileExtension
			if !strings.HasSuffix(completeFilename, dotExtension) {
				completeFilename += dotExtension
			}
			err := currentState.SaveTo(completeFilename)
			if err != nil {
				return err
			}
			app.projectService.SetCurrentStateFilename(completeFilename)
			app.window.SetTitleSuffix(completeFilename)
			return nil
		})
	}
}

func (app *Application) currentProjectState() projectState {
	projectSettings := app.projectService.CurrentSettings()

	windowOpenByName := app.windowOpenByName()
	var openWindows []string
	for key, open := range windowOpenByName {
		if *open {
			openWindows = append(openWindows, key)
		}
	}
	activeLevelIndex := app.levelControlView.SelectedLevel()
	return projectState{
		Settings:         &projectSettings,
		OpenWindows:      openWindows,
		ActiveLevelIndex: &activeLevelIndex,
	}
}

func (app *Application) saveWorkspace() {
	currentProjectState := app.currentProjectState()
	currentProjectStateFilename := app.projectService.CurrentStateFilename()
	if len(currentProjectStateFilename) > 0 {
		err := currentProjectState.SaveTo(currentProjectStateFilename)
		if err != nil {
			currentProjectStateFilename = ""
		} else {
			currentProjectState = projectState{}
		}
	}

	lastProjectState := lastProjectState{
		projectState: currentProjectState,
		LastProject:  currentProjectStateFilename,
	}
	_ = lastProjectState.SaveTo(app.lastProjectStateConfigFilename())
}

func (app *Application) restoreWorkspace() {
	lastProjectState := lastProjectStateFromFile(app.lastProjectStateConfigFilename())

	projectStateToRestore := lastProjectState.projectState

	projectStateFilename := ""
	if len(lastProjectState.LastProject) > 0 {
		state, err := projectStateFromFile(lastProjectState.LastProject)
		if err == nil {
			projectStateToRestore = state
			projectStateFilename = lastProjectState.LastProject
		}
	}

	app.restoreProjectState(projectStateToRestore, projectStateFilename)
}

func (app *Application) restoreProjectState(state projectState, filename string) {
	var settings edit.ProjectSettings
	if state.Settings != nil {
		settings = *state.Settings
	}
	app.projectService.RestoreProject(settings, filename)

	windowOpenByName := app.windowOpenByName()
	for _, open := range windowOpenByName {
		*open = false
	}
	for _, key := range state.OpenWindows {
		open := windowOpenByName[key]
		if open != nil {
			*open = true
		}
	}
	if len(state.OpenWindows) == 0 {
		*app.projectView.WindowOpen() = true
	}

	activeLevel := world.StartingLevel
	if (state.ActiveLevelIndex != nil) &&
		(*state.ActiveLevelIndex >= 0) &&
		(*state.ActiveLevelIndex < len(app.levels)) {
		activeLevel = *state.ActiveLevelIndex
	}
	app.eventQueue.Event(levels.LevelSelectionSetEvent{Id: activeLevel})

	app.window.SetTitleSuffix(filename)
}

func (app *Application) updateAutoSave() {
	status := app.projectService.SaveStatus()
	app.window.SetProjectModified(status.FilesModified > 0)
	if status.SavePending && (status.SaveIn == 0) {
		status.ConfirmPendingSave()
		app.projectView.StartSavingMod()
	}
}

func (app *Application) windowStateConfigFilename() string {
	return filepath.Join(app.ConfigDir, "WindowState.json")
}

func (app *Application) lastProjectStateConfigFilename() string {
	return filepath.Join(app.ConfigDir, "LastProjectState.json")
}

func (app *Application) windowOpenByName() map[string]*bool {
	return map[string]*bool{
		"project":      app.projectView.WindowOpen(),
		"archive":      app.archiveView.WindowOpen(),
		"levelControl": app.levelControlView.WindowOpen(),
		"levelTiles":   app.levelTilesView.WindowOpen(),
		"levelObjects": app.levelObjectsView.WindowOpen(),
		"messages":     app.messagesView.WindowOpen(),
		"texts":        app.textsView.WindowOpen(),
		"bitmaps":      app.bitmapsView.WindowOpen(),
		"textures":     app.texturesView.WindowOpen(),
		"animations":   app.animationsView.WindowOpen(),
		"movies":       app.moviesView.WindowOpen(),
		"soundEffects": app.soundEffectsView.WindowOpen(),
		"gameObjects":  app.objectsView.WindowOpen(),
	}
}
