package movies

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"time"

	"github.com/asticode/go-astisub"
	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ss1/edit/undoable"
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

	frameCache    *graphics.FrameCache
	frameCacheKey graphics.FrameCacheKey

	movieService undoable.MovieService

	modalStateMachine gui.ModalStateMachine
	guiScale          float32
	commander         cmd.Commander

	model viewModel
}

// NewMoviesView returns a new instance.
func NewMoviesView(mod *world.Mod, frameCache *graphics.FrameCache,
	movieService undoable.MovieService,
	modalStateMachine gui.ModalStateMachine, guiScale float32, commander cmd.Commander) *View {
	view := &View{
		mod: mod,

		frameCache:    frameCache,
		frameCacheKey: frameCache.AllocateKey(),

		movieService: movieService,

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
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 800 * view.guiScale, Y: 300 * view.guiScale}, imgui.ConditionFirstUseEver)
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
					view.model.currentScene = 0
					view.model.currentFrame = 0
					view.model.frameTimeFraction = -1
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
			imgui.LabelText("Language", "(not localized)")
		}

		imgui.Separator()

		view.renderProperties()

		imgui.PopItemWidth()
	}
	imgui.EndChild()

	scenes := view.movieService.Video(view.model.currentKey)
	imgui.SameLine()
	imgui.BeginGroup()
	imgui.BeginChildV("Scenes", imgui.Vec2{X: 200 * view.guiScale, Y: -60 * view.guiScale}, true,
		imgui.WindowFlagsHorizontalScrollbar|imgui.WindowFlagsAlwaysVerticalScrollbar)
	for index, scene := range scenes {
		sceneInfo := fmt.Sprintf("Scene %02d (%d frames)", index, len(scene.Frames))
		if imgui.SelectableV(sceneInfo, view.model.currentScene == index, 0, imgui.Vec2{}) {
			view.model.currentScene = index
			view.model.currentFrame = 0
			view.model.frameTimeFraction = -1
		}
	}
	imgui.EndChild()
	if imgui.Button("Up") {
		view.requestMoveSceneEarlier()
	}
	imgui.SameLine()
	if imgui.Button("Down") {
		view.requestMoveSceneLater()
	}
	imgui.SameLine()
	if imgui.Button("Remove") {
		view.requestRemoveScene()
	}
	if imgui.Button("Import") {
		view.requestImportScene("")
	}
	imgui.SameLine()
	if imgui.Button("Export") {
		view.requestExportScene()
	}
	imgui.EndGroup()
	imgui.SameLine()
	if imgui.BeginChildV("Frames", imgui.Vec2{X: -1, Y: 0}, false, 0) {
		var scene *movie.Scene
		if (view.model.currentScene >= 0) && (view.model.currentScene < len(scenes)) {
			scene = &scenes[view.model.currentScene]
		}
		if scene != nil {
			var frame *movie.Frame
			if (view.model.currentFrame >= 0) && (view.model.currentFrame < len(scene.Frames)) {
				frame = &scene.Frames[view.model.currentFrame]
			}
			if frame != nil {
				// This code updates the texture every render cycle. In case of performance loss, this is a point to optimize.
				view.frameCache.SetTexture(view.frameCacheKey, 600, 300, frame.Pixels, &scene.Palette)

				render.FrameImage("Frame", view.frameCache, view.frameCacheKey,
					imgui.Vec2{
						X: float32(movie.HighResDefaultWidth) * view.guiScale,
						Y: float32(movie.HighResDefaultHeight) * view.guiScale,
					})
			}

			imgui.PushItemWidth(-150 * view.guiScale)
			gui.StepSliderInt("Frame Index", &view.model.currentFrame, 0, len(scene.Frames)-1)
			imgui.PopItemWidth()
			if (frame != nil) && (frame.DisplayTime < (time.Second / 4)) {
				imgui.Separator()
				if view.model.frameTimeFraction < 0 {
					view.model.frameTimeFraction = int((frame.DisplayTime * 0x10000) / time.Second)
				}
				imgui.PushItemWidth(-150 * view.guiScale)
				gui.StepSliderInt("Frame time fraction", &view.model.frameTimeFraction, 0, 0x3FFF)
				imgui.PopItemWidth()
				if imgui.Button("Set for scene") {
					displayTime := (time.Duration(view.model.frameTimeFraction) * time.Second) / 0x10000
					view.requestSetFramesDisplayTime(displayTime)
				}
			}
		}
	}
	imgui.EndChild()
}

func (view *View) renderProperties() {
	view.renderAudioProperties()
	view.renderSubtitlesProperties()

	if view.movieService.SizeWarning(view.model.currentKey) {
		view.renderSizeWarning()
	}
}

func (view *View) renderAudioProperties() {
	imgui.PushID("audio")
	imgui.Separator()
	sound := view.currentSound()
	imgui.LabelText("Audio", fmt.Sprintf("%.2f sec", sound.Duration()))
	if imgui.Button("Import") {
		view.requestImportAudio()
	}
	if !sound.Empty() {
		imgui.SameLine()
		if imgui.Button("Export") {
			view.requestExportAudio(sound)
		}
		imgui.SameLine()
		if imgui.Button("Clear") {
			view.requestClearAudio()
		}
	}
	imgui.PopID()
}

func (view *View) renderSubtitlesProperties() {
	imgui.PushID("subtitles")
	imgui.Separator()
	if imgui.BeginCombo("Sub Language", view.model.currentSubtitleLang.String()) {
		languages := resource.Languages()
		for _, lang := range languages {
			if imgui.SelectableV(lang.String(), lang == view.model.currentSubtitleLang, 0, imgui.Vec2{}) {
				view.model.currentSubtitleLang = lang
			}
		}
		imgui.EndCombo()
	}
	sub := view.currentSubtitles()
	imgui.Text(fmt.Sprintf("%d lines", len(sub.Entries)))
	if imgui.Button("Import") {
		view.requestImportSubtitles()
	}
	if len(sub.Entries) > 0 {
		imgui.SameLine()
		if imgui.Button("Export") {
			view.requestExportSubtitles()
		}
		imgui.SameLine()
		if imgui.Button("Clear") {
			view.requestClearSubtitles()
		}
	}
	imgui.PopID()
}

func (view *View) renderSizeWarning() {
	imgui.Separator()
	imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
	imgui.Text("Size warning - reduce content!")
	imgui.PopStyleColor()
	if imgui.IsItemHovered() {
		imgui.BeginTooltip()
		imgui.SetTooltip("The movie has reached a size larger than officially supported.\n" +
			"While the game might still show the movie, a future update may stop supporting it.")
		imgui.EndTooltip()
	}
}

func (view *View) restoreFunc() func() {
	return view.restoreFuncWithScene(view.model.currentScene)
}

func (view *View) restoreFuncWithScene(oldScene int) func() {
	oldKey := view.model.currentKey
	oldSubtitlesLang := view.model.currentSubtitleLang
	oldFrame := view.model.currentFrame

	return func() {
		view.model.restoreFocus = true
		view.model.currentKey = oldKey
		view.model.currentSubtitleLang = oldSubtitlesLang
		view.model.currentScene = oldScene
		view.model.currentFrame = oldFrame
		view.model.frameTimeFraction = -1
	}
}

func (view *View) currentSound() audio.L8 {
	return view.movieService.Audio(view.model.currentKey)
}

func (view *View) requestExportAudio(sound audio.L8) {
	filename := fmt.Sprintf("%s_%s.wav", knownMovies[view.model.currentKey.ID].title, view.model.currentKey.Lang.String())

	external.ExportAudio(view.modalStateMachine, filename, sound)
}

func (view *View) requestImportAudio() {
	external.ImportAudio(view.modalStateMachine, func(sound audio.L8) {
		view.movieService.RequestSetAudio(view.model.currentKey, sound, view.restoreFunc())
	})
}

func (view *View) requestClearAudio() {
	view.movieService.RequestSetAudio(view.model.currentKey, audio.L8{}, view.restoreFunc())
}

func (view *View) currentSubtitles() movie.SubtitleList {
	return view.movieService.Subtitles(view.model.currentKey, view.model.currentSubtitleLang)
}

func (view View) requestExportSubtitles() {
	filename := fmt.Sprintf("%s_%s.srt", knownMovies[view.model.currentKey.ID].title, view.model.currentSubtitleLang.String())
	info := "File to be written: " + filename
	var exportTo func(string)
	currentSubtitles := view.currentSubtitles()

	exportTo = func(dirname string) {
		writer, err := os.Create(filepath.Join(dirname, filename))
		if err != nil {
			external.Export(view.modalStateMachine, "Could not create file.\n"+info, exportTo, true)
			return
		}
		defer func() { _ = writer.Close() }()

		sub := astisub.NewSubtitles()

		var lastItem *astisub.Item
		for _, entry := range currentSubtitles.Entries {
			var item astisub.Item
			var line astisub.Line
			line.Items = append(line.Items, astisub.LineItem{Text: entry.Text})
			item.Lines = []astisub.Line{line}
			item.StartAt = entry.Timestamp
			item.EndAt = item.StartAt
			if lastItem != nil {
				lastItem.EndAt = item.StartAt
			}
			lastItem = &item
			sub.Items = append(sub.Items, lastItem)
		}

		err = sub.WriteToSRT(writer)
		if err != nil {
			external.Export(view.modalStateMachine, "Could not export subtitles.\n"+info, exportTo, true)
			return
		}
	}

	external.Export(view.modalStateMachine, info, exportTo, false)
}

func (view *View) requestImportSubtitles() {
	info := "File must be an .SRT file."
	types := []external.TypeInfo{{Title: "Subtitle files (*.srt)", Extensions: []string{"srt"}}}
	var fileHandler func(string)

	fileHandler = func(filename string) {
		reader, err := os.Open(filename)
		if err != nil {
			external.Import(view.modalStateMachine, "Could not open file.\n"+info, types, fileHandler, true)
			return
		}
		defer func() { _ = reader.Close() }()

		subtitles, err := astisub.ReadFromSRT(reader)
		if err != nil {
			external.Import(view.modalStateMachine, "File not recognized as SRT.\n"+info, types, fileHandler, true)
			return
		}
		var newSubtitles movie.SubtitleList
		for _, item := range subtitles.Items {
			var newEntry movie.Subtitle
			newEntry.Timestamp = item.StartAt
			for _, line := range item.Lines {
				for _, lineItem := range line.Items {
					if len(newEntry.Text) > 0 {
						newEntry.Text += "\n"
					}
					newEntry.Text += lineItem.Text
				}
			}
			newSubtitles.Entries = append(newSubtitles.Entries, newEntry)
		}

		view.movieService.RequestSetSubtitles(view.model.currentKey, view.model.currentSubtitleLang,
			newSubtitles, view.restoreFunc())
	}

	external.Import(view.modalStateMachine, info, types, fileHandler, false)
}

func (view *View) requestClearSubtitles() {
	view.movieService.RequestSetSubtitles(view.model.currentKey, view.model.currentSubtitleLang,
		movie.SubtitleList{}, view.restoreFunc())
}

func (view *View) requestImportScene(returningInfo string) {
	info := fmt.Sprintf("File must be an animated GIF file in size %dx%d.",
		movie.HighResDefaultWidth, movie.HighResDefaultHeight)
	types := []external.TypeInfo{{Title: "Animation files (*.gif)", Extensions: []string{"gif"}}}
	var fileHandler func(string)

	fileHandler = func(filename string) {
		reader, err := os.Open(filename)
		if err != nil {
			external.Import(view.modalStateMachine, "Could not open file.\n"+info, types, fileHandler, true)
			return
		}
		defer func() { _ = reader.Close() }()
		data, err := gif.DecodeAll(reader)
		if err != nil {
			external.Import(view.modalStateMachine, "File not recognized as GIF.\n"+info, types, fileHandler, true)
			return
		}

		if (data.Config.Width != movie.HighResDefaultWidth) || (data.Config.Height != movie.HighResDefaultHeight) {
			external.Import(view.modalStateMachine, info, types, fileHandler, true)
			return
		}

		var scene movie.Scene
		scene.Frames = make([]movie.Frame, len(data.Image))
		var palette bitmap.Palette
		if gifPalette, isPal := data.Config.ColorModel.(color.Palette); isPal {
			for index, clr := range gifPalette {
				r, g, b, _ := clr.RGBA()
				palette[index].Red = byte(r >> 8)
				palette[index].Green = byte(g >> 8)
				palette[index].Blue = byte(b >> 8)
			}
		}
		framebuffer := make([]byte, data.Config.Width*data.Config.Height)
		framebufferSnapshot := func() []byte {
			buf := make([]byte, len(framebuffer))
			copy(buf, framebuffer)
			return buf
		}

		var usedColors [256]bool
		for _, img := range data.Image {
			for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
				for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
					colorIndex := img.ColorIndexAt(x, y)
					usedColors[colorIndex] = true
				}
			}
		}
		var indexMap [256]byte
		for i := 0; i < 256; i++ {
			indexMap[byte(i)] = byte(i)
		}
		// Ensure color index 0 is not used.
		if usedColors[0] && (data.BackgroundIndex != 0) {
			if data.BackgroundIndex == 255 {
				indexMap[0] = palette.IndexClosestTo(palette[0], []byte{0, 255})
			} else {
				temp := palette[0]
				palette[0] = palette[data.BackgroundIndex]
				palette[data.BackgroundIndex] = temp
				indexMap[0] = data.BackgroundIndex
				indexMap[data.BackgroundIndex] = 0
			}
		}
		palette[0] = bitmap.RGB{Red: 0x00, Green: 0x00, Blue: 0x00} // default color for background
		// In case the last index is in use by the animation, use a similar color as this entry is reserved for subtitles.
		if usedColors[255] && (data.BackgroundIndex != 255) {
			indexMap[255] = palette.IndexClosestTo(palette[255], []byte{0, 255})
		}
		palette[255] = bitmap.RGB{Red: 0x9A, Green: 0x35, Blue: 0x35} // default color in intro

		scene.Palette = palette
		for index, img := range data.Image {
			for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
				for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
					colorIndex := img.ColorIndexAt(x, y)
					if colorIndex != data.BackgroundIndex {
						framebuffer[y*data.Config.Width+x] = indexMap[colorIndex]
					}
				}
			}

			scene.Frames[index].DisplayTime = time.Duration(data.Delay[index]) * 10 * time.Millisecond
			scene.Frames[index].Pixels = framebufferSnapshot()
		}

		view.compressAndAddScene(scene, data.Config.Width, data.Config.Height)
	}

	external.Import(view.modalStateMachine, returningInfo+info, types, fileHandler, false)
}

func (view *View) compressAndAddScene(scene movie.Scene, width, height int) {
	view.modalStateMachine.SetState(&compressingStartState{
		machine:  view.modalStateMachine,
		view:     view,
		width:    width,
		height:   height,
		input:    scene,
		listener: view.onCompressionResult,
	})
}

func (view *View) onCompressionResult(result compressionResult) {
	switch typedResult := result.(type) {
	case compressionAborted:
	case compressionFinished:
		view.requestAddScene(typedResult.scene)
	case compressionFailed:
		view.requestImportScene("Could not compress. Follow recommendations and retry.\n" +
			"Technical details:\n" + typedResult.err.Error() + "\n\n")
	}
}

func (view *View) requestAddScene(scene movie.HighResScene) {
	view.movieService.RequestAddScene(view.model.currentKey, scene, view.restoreFunc())
}

func (view *View) requestExportScene() {
	filename := fmt.Sprintf("%s_Scene%02d_%s.gif",
		knownMovies[view.model.currentKey.ID].title,
		view.model.currentScene,
		view.model.currentKey.Lang.String())
	info := "File to be written: " + filename
	var exportTo func(string)

	scenes := view.movieService.Video(view.model.currentKey)
	var scene *movie.Scene
	if view.model.currentScene >= 0 && view.model.currentScene < len(scenes) {
		scene = &scenes[view.model.currentScene]
	}

	if (scene == nil) || len(scene.Frames) == 0 {
		return
	}

	exportTo = func(dirname string) {
		writer, err := os.Create(filepath.Join(dirname, filename))
		if err != nil {
			external.Export(view.modalStateMachine, "Could not create file.\n"+info, exportTo, true)
			return
		}
		defer func() { _ = writer.Close() }()

		colorPalette := scene.Palette.ColorPalette(false)
		data := gif.GIF{
			Config: image.Config{
				Width:      600,
				Height:     300,
				ColorModel: colorPalette,
			},
			LoopCount: -1,
		}

		imageRect := image.Rect(0, 0, data.Config.Width, data.Config.Height)
		for _, frame := range scene.Frames {
			frameImg := image.NewPaletted(imageRect, colorPalette)
			frameImg.Pix = frame.Pixels
			frameImg.Stride = data.Config.Width
			data.Image = append(data.Image, frameImg)
			data.Delay = append(data.Delay, int(((frame.DisplayTime/time.Millisecond)+5)/10))
		}

		err = gif.EncodeAll(writer, &data)
		if err != nil {
			external.Export(view.modalStateMachine, info, exportTo, true)
			return
		}
	}

	external.Export(view.modalStateMachine, info, exportTo, false)
}

func (view *View) requestMoveSceneEarlier() {
	scenes := view.movieService.Video(view.model.currentKey)
	if (view.model.currentScene > 0) && (view.model.currentScene < len(scenes)) {
		view.movieService.RequestMoveSceneEarlier(view.model.currentKey, view.model.currentScene,
			view.restoreFuncWithScene(view.model.currentScene-1))
	}
}

func (view *View) requestMoveSceneLater() {
	scenes := view.movieService.Video(view.model.currentKey)
	if (view.model.currentScene >= 0) && (view.model.currentScene < (len(scenes) - 1)) {
		view.movieService.RequestMoveSceneLater(view.model.currentKey, view.model.currentScene,
			view.restoreFuncWithScene(view.model.currentScene+1))
	}
}

func (view *View) requestRemoveScene() {
	scenes := view.movieService.Video(view.model.currentKey)
	if (view.model.currentScene >= 0) && (view.model.currentScene < len(scenes)) {
		view.movieService.RequestRemoveScene(view.model.currentKey, view.model.currentScene, view.restoreFunc())
	}
}

func (view *View) requestSetFramesDisplayTime(displayTime time.Duration) {
	view.movieService.RequestSetSceneFramesDisplayTime(view.model.currentKey, view.model.currentScene, displayTime, view.restoreFunc())
}
