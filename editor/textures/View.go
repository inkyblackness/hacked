package textures

import (
	"fmt"

	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
)

// View provides edit controls for textures.
type View struct {
	mod          *model.Mod
	textCache    *text.Cache
	imageCache   *graphics.TextureCache
	paletteCache *graphics.PaletteCache

	modalStateMachine gui.ModalStateMachine
	guiScale          float32
	commander         cmd.Commander

	model viewModel
}

// NewTexturesView returns a new instance.
func NewTexturesView(mod *model.Mod, textCache *text.Cache, imageCache *graphics.TextureCache, paletteCache *graphics.PaletteCache,
	modalStateMachine gui.ModalStateMachine,
	guiScale float32, commander cmd.Commander) *View {
	view := &View{
		mod:          mod,
		textCache:    textCache,
		imageCache:   imageCache,
		paletteCache: paletteCache,

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
		if imgui.BeginV("Textures", view.WindowOpen(), imgui.WindowFlagsNoCollapse|imgui.WindowFlagsHorizontalScrollbar) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *View) renderContent() {
	if imgui.BeginChildV("Properties", imgui.Vec2{X: 350 * view.guiScale, Y: 0}, false, 0) {
		imgui.PushItemWidth(-150 * view.guiScale)

		gui.StepSliderInt("Index", &view.model.currentIndex, 0, world.MaxWorldTextures-1)

		render.TextureSelector("###IndexBitmap", -1, view.guiScale, world.MaxWorldTextures,
			view.model.currentIndex, view.imageCache,
			func(index int) resource.Key {
				return view.indexedResourceKey(ids.LargeTextures, index)
			},
			func(index int) string { return view.textureName(index) },
			func(newValue int) {
				view.model.currentIndex = newValue
			})

		imgui.Separator()

		imgui.PopItemWidth()
	}
	imgui.EndChild()
	imgui.SameLine()

	imgui.BeginGroup()
	view.renderTextureSample("Large", ids.LargeTextures, 128)
	view.renderTextureSample("Medium", ids.MediumTextures, 64)
	view.renderTextureSample("Small", ids.SmallTextures, 32)
	view.renderTextureSample("Icon", ids.IconTextures, 16)
	imgui.EndGroup()
}

func (view *View) renderTextureSample(label string, id resource.ID, sideLength float32) {
	if imgui.BeginChildV(label, imgui.Vec2{X: -1, Y: 128 * view.guiScale}, true, imgui.WindowFlagsNoScrollbar) {
		key := view.indexedResourceKey(id, view.model.currentIndex)
		render.TextureImage("Texture Bitmap", view.imageCache, key,
			imgui.Vec2{X: sideLength * view.guiScale, Y: sideLength * view.guiScale})

		imgui.SameLine()
		imgui.BeginGroup()

		tex, err := view.imageCache.Texture(key)

		if imgui.Button("Clear") {
			//view.requestClear()
		}
		imgui.SameLine()
		if imgui.Button("Import") {
			//view.requestImport(false)
		}
		if err == nil {
			imgui.SameLine()
			if imgui.Button("Export") {
				//view.requestExport(false)
			}
			if view.hasModCurrentBitmap() {
				imgui.SameLine()
				if imgui.Button("Remove") {
					//view.requestSetBitmapData(nil)
				}
			}

			width, height := tex.Size()
			imgui.LabelText("Width", fmt.Sprintf("%d", int(width)))
			imgui.LabelText("Height", fmt.Sprintf("%d", int(height)))
		}

		imgui.EndGroup()
	}
	imgui.EndChild()
}

func (view *View) currentResourceKey() resource.Key {
	return view.indexedResourceKey(ids.LargeTextures, view.model.currentIndex)
}

func (view *View) indexedResourceKey(id resource.ID, index int) resource.Key {
	info, _ := ids.Info(id)
	if info.List {
		return resource.KeyOf(id, resource.LangAny, index)
	}
	return resource.KeyOf(id.Plus(index), resource.LangAny, 0)
}

func (view *View) textureName(index int) string {
	key := resource.KeyOf(ids.TextureNames, resource.LangDefault, index)
	name, err := view.textCache.Text(key)
	suffix := ""
	if err == nil {
		suffix = ": " + name
	}
	return fmt.Sprintf("%3d", index) + suffix
}

func (view *View) hasModCurrentBitmap() bool {
	key := view.currentResourceKey()
	return len(view.mod.ModifiedBlock(key.Lang, key.ID, key.Index)) > 0
}

func (view *View) requestExport(withError bool) {
	// TODO needs four textures to be exported?
	/*
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
	*/
}

func (view *View) requestImport(withError bool) {
	// TODO allow import of four bitmaps at once
	/*
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
	*/
}

func (view *View) requestClear() {
	/*
		bmp := bitmap.Bitmap{
			Header: bitmap.Header{
				Width:  1,
				Height: 1,
			},
			Pixels: []byte{0x00},
		}
		view.requestSetBitmap(bmp)
	*/
}
