package animations

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/gif"
	"os"
	"path/filepath"

	"github.com/inkyblackness/imgui-go"

	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial/rle"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
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
	mod            *world.Mod
	imageCache     *graphics.TextureCache
	paletteCache   *graphics.PaletteCache
	animationCache *bitmap.AnimationCache

	modalStateMachine gui.ModalStateMachine
	guiScale          float32
	commander         cmd.Commander

	model viewModel
}

// NewAnimationsView returns a new instance.
func NewAnimationsView(mod *world.Mod, imageCache *graphics.TextureCache, paletteCache *graphics.PaletteCache,
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

	if imgui.BeginChildV("Frames", imgui.Vec2{X: -1, Y: 0}, false, 0) {
		anim, hasAnim, _ := view.currentAnimation()

		if hasAnim {
			frameKey := resource.KeyOf(anim.ResourceID, resource.LangAny, view.model.currentFrame)
			if view.cacheFrame(frameKey) {
				render.TextureImage("Frame", view.imageCache, frameKey,
					imgui.Vec2{X: float32(anim.Width) * view.guiScale, Y: float32(anim.Height) * view.guiScale})
				if imgui.Button("Export") {
					view.requestExport()
				}
			}
			imgui.SameLine()
			if imgui.Button("Import") {
				view.requestImport()
			}
		}
	}
	imgui.EndChild()
}

func (view *View) cacheFrame(key resource.Key) bool {
	var err error
	var lastKey *resource.Key
	for index := 0; (index <= key.Index) && (err == nil); index++ {
		nextKey := resource.KeyOf(key.ID, key.Lang, index)
		_, err = view.imageCache.TextureReferenced(nextKey, lastKey)
		lastKey = &nextKey
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

func (view *View) requestImport() {
	info := "File must be an animated GIF file.\nIdeally, it matches the game palette 1:1,\nothers are mapped closest fitting."
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

		palette, err := view.paletteCache.Palette(0)
		if err != nil {
			external.Import(view.modalStateMachine, "Can not import image without having a palette loaded.\n"+info, types, fileHandler, true)
			return
		}
		anim := bitmap.Animation{
			Width:      int16(data.Config.Width),
			Height:     int16(data.Config.Height),
			ResourceID: view.model.currentKey.ID.Plus(view.model.currentKey.Index).Plus(-12),
			IntroFlag:  0,
		}
		var frames [][]byte
		rawPalette := palette.Palette()

		highestBitShift := func(value int16) (result byte) {
			if value != 0 {
				for (value >> result) != 1 {
					result++
				}
			}
			return
		}

		var prevFrame []byte
		for index, img := range data.Image {
			if (img.Bounds().Max.X == data.Config.Width) && (img.Bounds().Max.Y == data.Config.Height) {
				bitmapper := bitmap.NewBitmapper(&rawPalette)
				bmp := bitmapper.Map(img)
				entry := bitmap.AnimationEntry{
					FirstFrame: byte(index),
					LastFrame:  byte(index),
					FrameTime:  int16(data.Delay[index] * 10),
				}
				anim.Entries = append(anim.Entries, entry)
				bmp.Header.Type = bitmap.TypeCompressed8Bit
				bmp.Header.WidthFactor = highestBitShift(bmp.Header.Width)
				bmp.Header.HeightFactor = highestBitShift(bmp.Header.Height)
				bmp.Header.Area = [4]int16{0, 0, anim.Width, anim.Height}
				bmp.Header.Stride = uint16(bmp.Header.Width)

				buf := bytes.NewBuffer(nil)
				_ = binary.Write(buf, binary.LittleEndian, &bmp.Header)
				_ = rle.Compress(buf, bmp.Pixels, prevFrame)
				prevFrame = bmp.Pixels
				frames = append(frames, buf.Bytes())
			}
		}

		view.requestSetAnimation(anim, frames)
	}

	external.Import(view.modalStateMachine, info, types, fileHandler, false)
}

func (view *View) requestExport() {
	anim, hasAnim, _ := view.currentAnimation()
	if hasAnim {
		filename := fmt.Sprintf("Anim%04X.gif", int(view.model.currentKey.ID.Plus(view.model.currentKey.Index)))
		view.exportTo(filename, anim)
	}
}

func (view *View) exportTo(filename string, anim bitmap.Animation) {
	info := "File to be written: " + filename
	var exportTo func(string)

	exportTo = func(dirname string) {
		writer, err := os.Create(filepath.Join(dirname, filename))
		if err != nil {
			external.Export(view.modalStateMachine, "Could not create file.\n"+info, exportTo, true)
			return
		}
		defer func() { _ = writer.Close() }()

		palTex, err := view.paletteCache.Palette(0)
		if err != nil {
			external.Export(view.modalStateMachine, "Could not create file. No palette loaded.\n"+info, exportTo, true)
			return
		}

		colorPalette := palTex.Palette().ColorPalette(false)
		data := gif.GIF{
			Config: image.Config{
				Width:      int(anim.Width),
				Height:     int(anim.Height),
				ColorModel: colorPalette,
			},
			LoopCount: -1,
		}

		imageRect := image.Rect(0, 0, data.Config.Width, data.Config.Height)
		frameIndex := 0
		for _, entry := range anim.Entries {
			for frameIndex <= int(entry.LastFrame) {
				frameKey := resource.KeyOf(anim.ResourceID, resource.LangAny, frameIndex)
				if !view.cacheFrame(frameKey) {
					external.Export(view.modalStateMachine, "Failed to cache frame.\n"+info, exportTo, true)
					return
				}
				frameTex, _ := view.imageCache.Texture(frameKey)
				frameImg := image.NewPaletted(imageRect, colorPalette)
				frameImg.Pix = frameTex.PixelData()
				data.Image = append(data.Image, frameImg)
				data.Delay = append(data.Delay, int(entry.FrameTime)/10)
				frameIndex++
			}
		}

		err = gif.EncodeAll(writer, &data)
		if err != nil {
			external.Export(view.modalStateMachine, info, exportTo, true)
			return
		}
	}

	external.Export(view.modalStateMachine, info, exportTo, false)
}

func (view *View) requestSetAnimation(newAnim bitmap.Animation, newFrames [][]byte) {
	encodeAnim := func(anim bitmap.Animation) []byte {
		buf := bytes.NewBuffer(nil)
		_ = bitmap.WriteAnimation(buf, anim)
		return buf.Bytes()
	}
	command := setAnimationCommand{
		model:        &view.model,
		animationKey: view.model.currentKey,
		newAnimation: encodeAnim(newAnim),
		newFrames:    newFrames,
		framesID:     newAnim.ResourceID,
	}
	oldAnim, oldExisting, _ := view.currentAnimation()
	if oldExisting {
		command.oldAnimation = encodeAnim(oldAnim)
	}
	command.oldFrames = view.mod.ModifiedBlocks(view.model.currentKey.Lang, newAnim.ResourceID)
	view.commander.Queue(command)
}
