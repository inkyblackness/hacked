package textures

import (
	"fmt"
	"math"

	"github.com/inkyblackness/imgui-go/v2"

	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
)

// View provides edit controls for textures.
type View struct {
	mod          *world.Mod
	textCache    *text.Cache
	cp           text.Codepage
	imageCache   *graphics.TextureCache
	paletteCache *graphics.PaletteCache

	modalStateMachine gui.ModalStateMachine
	clipboard         external.Clipboard
	guiScale          float32
	commander         cmd.Commander

	model viewModel
}

// NewTexturesView returns a new instance.
func NewTexturesView(mod *world.Mod, textCache *text.Cache, cp text.Codepage,
	imageCache *graphics.TextureCache, paletteCache *graphics.PaletteCache,
	modalStateMachine gui.ModalStateMachine,
	clipboard external.Clipboard, guiScale float32, commander cmd.Commander) *View {
	view := &View{
		mod:          mod,
		textCache:    textCache,
		cp:           cp,
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
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 800 * view.guiScale, Y: 300 * view.guiScale}, imgui.ConditionFirstUseEver)
		if imgui.BeginV("Textures", view.WindowOpen(), imgui.WindowFlagsNoCollapse|imgui.WindowFlagsHorizontalScrollbar) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *View) renderContent() {
	if imgui.BeginChildV("Properties", imgui.Vec2{X: -300 * view.guiScale, Y: 0}, false, imgui.WindowFlagsHorizontalScrollbar) {
		imgui.PushItemWidth(-200 * view.guiScale)

		gui.StepSliderInt("Index", &view.model.currentIndex, 0, world.MaxWorldTextures-1)

		render.TextureSelector("###IndexBitmap", -1, view.guiScale, world.MaxWorldTextures,
			view.model.currentIndex, view.imageCache,
			func(index int) resource.Key {
				return view.indexedResourceKey(ids.LargeTextures, index)
			},
			view.textureTooltip,
			func(newValue int) {
				view.model.currentIndex = newValue
			})

		readOnly := !view.mod.HasModifyableTextureProperties()

		imgui.Separator()
		if imgui.BeginCombo("Language", view.model.currentLang.String()) {
			languages := resource.Languages()
			for _, lang := range languages {
				if imgui.SelectableV(lang.String(), lang == view.model.currentLang, 0, imgui.Vec2{}) {
					view.model.currentLang = lang
				}
			}
			imgui.EndCombo()
		}

		nameKey := resource.KeyOf(ids.TextureNames, view.model.currentLang, view.model.currentIndex)
		name, _ := view.textCache.Text(nameKey)
		view.renderText(readOnly, "Name", name, view.requestSetTextureName)

		useKey := resource.KeyOf(ids.TextureUsages, view.model.currentLang, view.model.currentIndex)
		use, _ := view.textCache.Text(useKey)
		view.renderText(readOnly, "Use", use, view.requestSetTextureUsage)

		imgui.Separator()
		view.renderTextureProperties(readOnly)

		imgui.PopItemWidth()
	}
	imgui.EndChild()
	imgui.SameLine()

	imgui.BeginGroup()
	view.renderTextureSample("Large", ids.LargeTextures, 128, "Large")
	view.renderTextureSample("Medium", ids.MediumTextures, 64, "Medium")
	view.renderTextureSample("Small", ids.SmallTextures, 32, "Small")
	view.renderTextureSample("Icon", ids.IconTextures, 16, "Icon")
	imgui.EndGroup()
}

func (view *View) renderText(readOnly bool, label string, value string, changeCallback func(string)) {
	imgui.LabelText(label, value)
	view.clipboardPopup(readOnly, label, value, changeCallback)
}

func (view *View) clipboardPopup(readOnly bool, label string, value string, changeCallback func(string)) {
	if imgui.BeginPopupContextItemV(label+"-Popup", 1) {
		if imgui.Selectable("Copy to Clipboard") {
			view.clipboard.SetString(value)
		}
		if !readOnly && imgui.Selectable("Copy from Clipboard") {
			newValue, err := view.clipboard.String()
			if err == nil {
				changeCallback(newValue)
			}
		}
		imgui.EndPopup()
	}
}

func (view *View) renderTextureProperties(readOnly bool) {
	list := view.mod.TextureProperties()
	distanceUnifier := values.NewUnifier()
	climbableUnifier := values.NewUnifier()
	transparencyControlUnifier := values.NewUnifier()
	animationGroupUnifier := values.NewUnifier()
	animationIndexUnifier := values.NewUnifier()
	if view.model.currentIndex < len(list) {
		properties := list[view.model.currentIndex]
		distanceUnifier.Add(int(properties.DistanceModifier))
		climbableUnifier.Add(properties.Climbable != 0)
		transparencyControlUnifier.Add(int(properties.TransparencyControl))
		animationGroupUnifier.Add(int(properties.AnimationGroup))
		animationIndexUnifier.Add(int(properties.AnimationIndex))
	}

	values.RenderUnifiedSliderInt(readOnly, false, "Distance Modifier", distanceUnifier,
		func(u values.Unifier) int { return u.Unified().(int) },
		func(value int) string { return "%d" },
		math.MinInt16, math.MaxInt16,
		func(newValue int) {
			view.requestChangeProperties(func(prop *texture.Properties) {
				prop.DistanceModifier = int16(newValue)
			})
		})

	values.RenderUnifiedCheckboxCombo(readOnly, false, "Climbable", climbableUnifier,
		func(newValue bool) {
			view.requestChangeProperties(func(prop *texture.Properties) {
				prop.Climbable = 0
				if newValue {
					prop.Climbable = 1
				}
			})
		})

	values.RenderUnifiedCombo(readOnly, false, "Transparency Control", transparencyControlUnifier,
		func(u values.Unifier) int { return u.Unified().(int) },
		func(index int) string { return texture.TransparencyControl(index).String() },
		len(texture.TransparencyControls()),
		func(newValue int) {
			view.requestChangeProperties(func(prop *texture.Properties) {
				prop.TransparencyControl = texture.TransparencyControl(newValue)
			})
		})

	values.RenderUnifiedSliderInt(readOnly, false, "Animation Group", animationGroupUnifier,
		func(u values.Unifier) int { return u.Unified().(int) },
		func(value int) string { return "%d" },
		0, 3,
		func(newValue int) {
			view.requestChangeProperties(func(prop *texture.Properties) {
				prop.AnimationGroup = byte(newValue)
			})
		})
	values.RenderUnifiedSliderInt(readOnly, false, "Animation Index", animationIndexUnifier,
		func(u values.Unifier) int { return u.Unified().(int) },
		func(value int) string { return "%d" },
		0, 3,
		func(newValue int) {
			view.requestChangeProperties(func(prop *texture.Properties) {
				prop.AnimationIndex = byte(newValue)
			})
		})
}

func (view *View) renderTextureSample(label string, id resource.ID, sideLength float32, sizeID string) {
	if imgui.BeginChildV(label, imgui.Vec2{X: -1, Y: (128 + 7) * view.guiScale}, true, imgui.WindowFlagsNoScrollbar) {
		key := view.indexedResourceKey(id, view.model.currentIndex)
		render.TextureImage("Texture Bitmap", view.imageCache, key,
			imgui.Vec2{X: sideLength * view.guiScale, Y: sideLength * view.guiScale})

		imgui.SameLine()
		imgui.BeginGroup()

		tex, err := view.imageCache.Texture(key)

		if imgui.Button("Clear") {
			view.requestClear(id, view.model.currentIndex, int(sideLength))
		}
		imgui.SameLine()
		if imgui.Button("Import") {
			view.requestImport(id, view.model.currentIndex)
		}
		if err == nil {
			if imgui.Button("Export") {
				view.requestExport(id, view.model.currentIndex, sizeID)
			}
			if view.hasModCurrentBitmap() {
				imgui.SameLine()
				if imgui.Button("Remove") {
					view.requestSetBitmapData(id, view.model.currentIndex, nil)
				}
			}

			width, height := tex.Size()
			imgui.Text(fmt.Sprintf("%d x %d px", int(width), int(height)))
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

func (view *View) textureTooltip(index int) string {
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

func (view *View) requestChangeProperties(modifier func(*texture.Properties)) {
	list := view.mod.TextureProperties()
	if view.model.currentIndex < len(list) {
		command := setTexturePropertiesCommand{
			model:         &view.model,
			textureIndex:  view.model.currentIndex,
			oldProperties: list[view.model.currentIndex],
			newProperties: list[view.model.currentIndex],
		}
		modifier(&command.newProperties)
		view.commander.Queue(command)
	}
}

func (view *View) requestSetTextureName(value string) {
	view.requestSetTextureText(ids.TextureNames, value)
}

func (view *View) requestSetTextureUsage(value string) {
	view.requestSetTextureText(ids.TextureUsages, value)
}

func (view *View) requestSetTextureText(id resource.ID, newValue string) {
	key := resource.KeyOf(id, view.model.currentLang, view.model.currentIndex)
	oldValue, _ := view.textCache.Text(key)

	if oldValue != newValue {
		command := setTextureTextCommand{
			model:   &view.model,
			key:     key,
			oldData: view.cp.Encode(oldValue),
			newData: view.cp.Encode(text.Blocked(newValue)[0]),
		}
		view.commander.Queue(command)
	}
}

func (view *View) requestExport(id resource.ID, index int, sizeID string) {
	key := view.indexedResourceKey(id, index)
	tex, err := view.imageCache.Texture(key)
	if err != nil {
		return
	}
	palette, err := view.paletteCache.Palette(0)
	if err != nil {
		return
	}
	rawPalette := palette.Palette()
	filename := fmt.Sprintf("Texture_%03d_%s.png", index, sizeID)
	width, height := tex.Size()
	bmp := bitmap.Bitmap{
		Header: bitmap.Header{
			Width:  int16(width),
			Height: int16(height),
		},
		Pixels:  tex.PixelData(),
		Palette: &rawPalette,
	}

	external.ExportImage(view.modalStateMachine, filename, bmp)
}

func (view *View) requestImport(id resource.ID, index int) {
	paletteRetriever := func() (bitmap.Palette, error) {
		palette, err := view.paletteCache.Palette(0)
		if err != nil {
			return bitmap.Palette{}, err
		}
		return palette.Palette(), nil
	}

	external.ImportImage(view.modalStateMachine, paletteRetriever, func(bmp bitmap.Bitmap) {
		view.requestSetBitmap(id, index, bmp)
	})
}

func (view *View) requestClear(id resource.ID, index int, sideLength int) {
	bmp := bitmap.Bitmap{
		Header: bitmap.Header{
			Width:  int16(sideLength),
			Height: int16(sideLength),
		},
		Pixels: []byte{0x00},
	}
	view.requestSetBitmap(id, index, bmp)
}

func (view *View) requestSetBitmap(id resource.ID, index int, bmp bitmap.Bitmap) {
	highestBitShift := func(value int16) (result byte) {
		if value != 0 {
			for (value >> result) != 1 {
				result++
			}
		}
		return
	}

	bmp.Header.Flags = 0
	bmp.Header.Type = bitmap.TypeFlat8Bit
	bmp.Header.WidthFactor = highestBitShift(bmp.Header.Width)
	bmp.Header.HeightFactor = highestBitShift(bmp.Header.Height)
	bmp.Header.Stride = uint16(bmp.Header.Width)
	data := bitmap.Encode(&bmp, 0)
	view.requestSetBitmapData(id, index, data)
}

func (view *View) requestSetBitmapData(id resource.ID, index int, newData []byte) {
	resourceKey := view.indexedResourceKey(id, index)
	command := setTextureBitmapCommand{
		model:        &view.model,
		id:           id,
		textureIndex: index,
		oldData:      view.mod.ModifiedBlock(resource.LangAny, resourceKey.ID, resourceKey.Index),
		newData:      newData,
	}
	view.commander.Queue(command)
}
