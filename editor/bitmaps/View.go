package bitmaps

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
)

type bitmapInfo struct {
	title            string
	languageSpecific bool
}

var knownBitmapTypes = map[resource.ID]bitmapInfo{
	ids.MfdDataBitmaps:        {"MFD Data Images", true},
	ids.ObjectMaterialBitmaps: {"Object Materials", false},
	ids.ObjectTextureBitmaps:  {"Object Textures", false},
}

var knownBitmapTypesOrder = []resource.ID{
	ids.MfdDataBitmaps,
	ids.ObjectMaterialBitmaps,
	ids.ObjectTextureBitmaps,
}

// View provides edit controls for bitmaps.
type View struct {
	mod          *model.Mod
	imageCache   *graphics.TextureCache
	paletteCache *graphics.PaletteCache

	modalStateMachine gui.ModalStateMachine
	clipboard         external.Clipboard
	guiScale          float32
	commander         cmd.Commander

	model viewModel
}

// NewBitmapsView returns a new instance.
func NewBitmapsView(mod *model.Mod, imageCache *graphics.TextureCache, paletteCache *graphics.PaletteCache,
	modalStateMachine gui.ModalStateMachine, clipboard external.Clipboard,
	guiScale float32, commander cmd.Commander) *View {
	view := &View{
		mod:          mod,
		imageCache:   imageCache,
		paletteCache: paletteCache,

		modalStateMachine: modalStateMachine,
		clipboard:         clipboard,
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
		if imgui.BeginV("Bitmaps", view.WindowOpen(), imgui.WindowFlagsNoCollapse|imgui.WindowFlagsHorizontalScrollbar) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *View) renderContent() {
	if imgui.BeginChildV("Properties", imgui.Vec2{X: 350 * view.guiScale, Y: 0}, false, 0) {
		imgui.PushItemWidth(-150 * view.guiScale)
		if imgui.BeginCombo("Bitmap Type", knownBitmapTypes[view.model.currentKey.ID].title) {
			for _, id := range knownBitmapTypesOrder {
				if imgui.SelectableV(knownBitmapTypes[id].title, id == view.model.currentKey.ID, 0, imgui.Vec2{}) {
					view.model.currentKey.ID = id
					view.model.currentKey.Index = 0
				}
			}
			imgui.EndCombo()
		}
		selectedType := knownBitmapTypes[view.model.currentKey.ID]
		if selectedType.languageSpecific {
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
			view.model.currentKey.Lang = resource.LangAny
		}

		info, _ := ids.Info(view.model.currentKey.ID)

		gui.StepSliderInt("Index", &view.model.currentKey.Index, 0, info.MaxCount-1)

		render.TextureSelector("###"+"IndexBitmap", -1, view.guiScale, info.MaxCount,
			view.model.currentKey.Index, view.imageCache,
			func(index int) resource.Key {
				return view.indexedResourceKey(index)
			},
			func(index int) string { return fmt.Sprintf("%d", index) },
			func(newValue int) {
				view.model.currentKey.Index = newValue
			})

		tex, err := view.imageCache.Texture(view.currentResourceKey())

		if imgui.Button("Clear") {
			view.requestClear()
		}
		imgui.SameLine()
		if imgui.Button("Import") {
			view.requestImport(false)
		}
		if err == nil {
			imgui.SameLine()
			if imgui.Button("Export") {
				view.requestExport(false)
			}
			if view.hasModCurrentBitmap() {
				imgui.SameLine()
				if imgui.Button("Remove") {
					view.requestSetBitmapData(nil)
				}
			}

			width, height := tex.Size()
			imgui.LabelText("Width", fmt.Sprintf("%d", int(width)))
			imgui.LabelText("Height", fmt.Sprintf("%d", int(height)))
		}

		imgui.PopItemWidth()
	}
	imgui.EndChild()
	imgui.SameLine()
	render.TextureImage("Big texture", view.imageCache, view.currentResourceKey(), imgui.Vec2{X: 320 * view.guiScale, Y: 240 * view.guiScale})
}

func (view *View) currentResourceKey() resource.Key {
	return view.indexedResourceKey(view.model.currentKey.Index)
}

func (view *View) indexedResourceKey(index int) resource.Key {
	key := view.model.currentKey
	info, _ := ids.Info(view.model.currentKey.ID)
	if !info.List {
		key.ID = key.ID.Plus(index)
		key.Index = 0
	} else {
		key.Index = index
	}
	return key
}

func (view *View) hasModCurrentBitmap() bool {
	key := view.currentResourceKey()
	return len(view.mod.ModifiedBlock(key.Lang, key.ID, key.Index)) > 0
}

func (view *View) requestExport(withError bool) {
	key := view.currentResourceKey()
	filename := fmt.Sprintf("%05d_%03d_%s.png", key.ID.Value(), key.Index, key.Lang.String())
	info := "File to be written: " + filename
	var exportTo func(string)

	exportTo = func(dirname string) {
		writer, err := os.Create(filepath.Join(dirname, filename))
		if err != nil {
			external.Export(view.modalStateMachine, "Could not create file.\n"+info, exportTo, true)
			return
		}
		defer func() { _ = writer.Close() }()

		texture, err := view.imageCache.Texture(key)
		if err != nil {
			external.Export(view.modalStateMachine, "Image not available.\n"+info, exportTo, true)
			return
		}
		palette, err := view.paletteCache.Palette(0)
		if err != nil {
			external.Export(view.modalStateMachine, "Palette not available.\n"+info, exportTo, true)
			return
		}

		width, height := texture.Size()
		imageRect := image.Rect(0, 0, int(width), int(height))
		imagePal := palette.Palette().ColorPalette(true)
		paletted := image.NewPaletted(imageRect, imagePal)
		paletted.Pix = texture.PixelData()
		err = png.Encode(writer, paletted)
		if err != nil {
			external.Export(view.modalStateMachine, info, exportTo, true)
			return
		}
	}

	external.Export(view.modalStateMachine, info, exportTo, withError)
}

func (view *View) requestImport(withError bool) {
	info := "File should be either a PNG or a GIF file."

	var fileHandler func(string)

	fileHandler = func(filename string) {
		palette, err := view.paletteCache.Palette(0)
		if err != nil {
			external.Export(view.modalStateMachine, "No palette loaded.\n"+info, fileHandler, true)
			return
		}
		reader, err := os.Open(filename)
		if err != nil {
			external.Import(view.modalStateMachine, "Could not open file.\n"+info, fileHandler, true)
			return
		}
		defer func() { _ = reader.Close() }()
		img, _, err := image.Decode(reader)
		if err != nil {
			external.Import(view.modalStateMachine, "File not recognized as image.\n"+info, fileHandler, true)
			return
		}

		bitmapper := bitmap.NewBitmapper(palette.Palette())
		bmp := bitmapper.Map(img)
		view.requestSetBitmap(bmp)
	}

	external.Import(view.modalStateMachine, info, fileHandler, false)
}

func (view *View) requestClear() {
	bmp := bitmap.Bitmap{
		Header: bitmap.Header{
			Width:  1,
			Height: 1,
		},
		Pixels: []byte{0x00},
	}
	view.requestSetBitmap(bmp)
}

func (view *View) requestSetBitmap(bmp bitmap.Bitmap) {
	highestBitShift := func(value int16) (result byte) {
		if value != 0 {
			for (value >> result) != 1 {
				result++
			}
		}
		return
	}

	bmp.Header.Flags = bitmap.FlagTransparent
	bmp.Header.Type = bitmap.TypeCompressed8Bit
	bmp.Header.WidthFactor = highestBitShift(bmp.Header.Width)
	bmp.Header.HeightFactor = highestBitShift(bmp.Header.Height)
	bmp.Header.Stride = uint16(bmp.Header.Width)
	data := bitmap.Encode(&bmp, 0)
	view.requestSetBitmapData(data)
}

func (view *View) requestSetBitmapData(newData []byte) {
	resourceKey := view.currentResourceKey()

	command := setBitmapCommand{
		displayKey: view.model.currentKey,
		model:      &view.model,

		resourceKey: resourceKey,
		oldData:     view.mod.ModifiedBlock(resourceKey.Lang, resourceKey.ID, resourceKey.Index),
		newData:     newData,
	}
	view.commander.Queue(command)
}
