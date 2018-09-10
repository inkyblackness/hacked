package bitmaps

import (
	"fmt"

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
	bitmapType       bitmap.Type
	bitmapFlags      bitmap.Flag
}

var knownBitmapTypes = map[resource.ID]bitmapInfo{
	ids.MfdDataBitmaps:        {title: "MFD Data Images", languageSpecific: true, bitmapType: bitmap.TypeCompressed8Bit, bitmapFlags: bitmap.FlagTransparent},
	ids.ObjectMaterialBitmaps: {title: "Object Materials", languageSpecific: false, bitmapType: bitmap.TypeFlat8Bit, bitmapFlags: 0},
	ids.ObjectTextureBitmaps:  {title: "Object Textures", languageSpecific: false, bitmapType: bitmap.TypeFlat8Bit, bitmapFlags: 0},
	ids.IconBitmaps:           {title: "Wall Icons", languageSpecific: false, bitmapType: bitmap.TypeCompressed8Bit, bitmapFlags: bitmap.FlagTransparent},
	ids.GraffitiBitmaps:       {title: "Graffiti", languageSpecific: false, bitmapType: bitmap.TypeCompressed8Bit, bitmapFlags: bitmap.FlagTransparent},
}

var knownBitmapTypesOrder = []resource.ID{
	ids.MfdDataBitmaps,
	ids.ObjectMaterialBitmaps,
	ids.ObjectTextureBitmaps,
	ids.IconBitmaps,
	ids.GraffitiBitmaps,
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
			view.requestClear(selectedType)
		}
		imgui.SameLine()
		if imgui.Button("Import") {
			view.requestImport(selectedType)
		}
		if err == nil {
			imgui.SameLine()
			if imgui.Button("Export") {
				view.requestExport()
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

func (view *View) requestExport() {
	key := view.currentResourceKey()
	texture, err := view.imageCache.Texture(key)
	if err != nil {
		return
	}
	palette, err := view.paletteCache.Palette(0)
	if err != nil {
		return
	}
	rawPalette := palette.Palette()
	filename := fmt.Sprintf("%05d_%03d_%s.png", key.ID.Value(), key.Index, key.Lang.String())
	width, height := texture.Size()
	bmp := bitmap.Bitmap{
		Header: bitmap.Header{
			Width:  int16(width),
			Height: int16(height),
		},
		Pixels:  texture.PixelData(),
		Palette: &rawPalette,
	}

	external.ExportImage(view.modalStateMachine, filename, bmp)
}

func (view *View) requestImport(bmpInfo bitmapInfo) {
	paletteRetriever := func() (bitmap.Palette, error) {
		palette, err := view.paletteCache.Palette(0)
		if err != nil {
			return bitmap.Palette{}, err
		}
		return palette.Palette(), nil
	}
	external.ImportImage(view.modalStateMachine, paletteRetriever, func(bmp bitmap.Bitmap) {
		view.requestSetBitmap(bmp, bmpInfo)
	})
}

func (view *View) requestClear(bmpInfo bitmapInfo) {
	bmp := bitmap.Bitmap{
		Header: bitmap.Header{
			Width:  1,
			Height: 1,
		},
		Pixels: []byte{0x00},
	}
	view.requestSetBitmap(bmp, bmpInfo)
}

func (view *View) requestSetBitmap(bmp bitmap.Bitmap, bmpInfo bitmapInfo) {
	highestBitShift := func(value int16) (result byte) {
		if value != 0 {
			for (value >> result) != 1 {
				result++
			}
		}
		return
	}

	bmp.Header.Flags = bmpInfo.bitmapFlags
	bmp.Header.Type = bmpInfo.bitmapType
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
