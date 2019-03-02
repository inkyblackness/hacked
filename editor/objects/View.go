package objects

import (
	"fmt"
	"strings"

	"github.com/inkyblackness/imgui-go"

	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/object/objprop"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
)

// View provides edit controls for game objects.
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

// NewView returns a new instance.
func NewView(mod *world.Mod, textCache *text.Cache, cp text.Codepage,
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
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 800 * view.guiScale, Y: 600 * view.guiScale}, imgui.ConditionOnce)
		if imgui.BeginV("Game Objects", view.WindowOpen(), imgui.WindowFlagsNoCollapse|imgui.WindowFlagsHorizontalScrollbar) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *View) renderContent() {
	if imgui.BeginChildV("Properties", imgui.Vec2{X: -330 * view.guiScale, Y: 0}, false, imgui.WindowFlagsHorizontalScrollbar) {
		imgui.PushItemWidth(-260 * view.guiScale)
		classString := func(class object.Class) string {
			return fmt.Sprintf("%2d: %v", int(class), class)
		}
		if imgui.BeginCombo("Object Class", classString(view.model.currentObject.Class)) {
			for _, class := range object.Classes() {
				if imgui.SelectableV(classString(class), class == view.model.currentObject.Class, 0, imgui.Vec2{}) {
					view.model.currentObject = object.TripleFrom(int(class), 0, 0)
					view.model.currentBitmap = 0
				}
			}
			imgui.EndCombo()
		}
		if imgui.BeginCombo("Object Type", view.tripleName(view.model.currentObject)) {
			allTypes := view.mod.ObjectProperties().TriplesInClass(view.model.currentObject.Class)
			for _, triple := range allTypes {
				if imgui.SelectableV(view.tripleName(triple), triple == view.model.currentObject, 0, imgui.Vec2{}) {
					view.model.currentObject = triple
					view.model.currentBitmap = 0
				}
			}
			imgui.EndCombo()
		}

		readOnly := !view.mod.HasModifyableObjectProperties()
		properties, propErr := view.mod.ObjectProperties().ForObject(view.model.currentObject)

		imgui.Separator()
		bitmapLimit := 0
		if propErr == nil {
			bitmapLimit = int(properties.Common.Bitmap3D.FrameNumber() + 2)
		}
		gui.StepSliderInt("Bitmap", &view.model.currentBitmap, 0, bitmapLimit)
		render.TextureSelector("BitmapSelector", -1, view.guiScale, bitmapLimit+1, view.model.currentBitmap,
			view.imageCache, view.currentBitmapKeyFor, func(int) string { return "" }, func(newValue int) {
				view.model.currentBitmap = newValue
			})

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
		view.renderText(readOnly, "Long Name",
			view.objectName(view.model.currentObject, view.model.currentLang, true),
			func(newValue string) {
				view.requestSetObjectName(view.model.currentObject, true, newValue)
			})
		view.renderText(readOnly, "Short Name",
			view.objectName(view.model.currentObject, view.model.currentLang, false),
			func(newValue string) {
				view.requestSetObjectName(view.model.currentObject, false, newValue)
			})

		if propErr == nil {
			if imgui.TreeNodeV("Common Properties", imgui.TreeNodeFlagsDefaultOpen|imgui.TreeNodeFlagsFramed) {
				view.renderCommonProperties(readOnly, properties)
				imgui.TreePop()
			}
			if imgui.TreeNodeV("Generic Properties", imgui.TreeNodeFlagsFramed) {
				view.renderGenericProperties(readOnly, properties)
				imgui.TreePop()
			}
			if imgui.TreeNodeV("Specific Properties", imgui.TreeNodeFlagsFramed) {
				view.renderSpecificProperties(readOnly, properties)
				imgui.TreePop()
			}
		} else {
			imgui.Text("(properties unavailable)")
		}

		imgui.PopItemWidth()
	}
	imgui.EndChild()
	imgui.SameLine()

	imgui.BeginGroup()
	view.renderObjectBitmap()
	imgui.EndGroup()
}

func (view *View) renderText(readOnly bool, label string, value string, changeCallback func(string)) {
	imgui.LabelText(label, value)
	view.clipboardPopup(readOnly, label, value, changeCallback)
}

func (view *View) tripleName(triple object.Triple) string {
	return triple.String() + ": " + view.objectName(triple, resource.LangDefault, true)
}

func (view *View) objectName(triple object.Triple, lang resource.Language, longName bool) string {
	result := "???"
	linearIndex := view.mod.ObjectProperties().TripleIndex(triple)
	if linearIndex >= 0 {
		nameID := ids.ObjectShortNames
		if longName {
			nameID = ids.ObjectLongNames
		}
		key := resource.KeyOf(nameID, lang, linearIndex)
		objName, err := view.textCache.Text(key)
		if err == nil {
			result = objName
		}
	}
	return result
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

func (view *View) requestSetObjectName(triple object.Triple, longName bool, newValue string) {
	linearIndex := view.mod.ObjectProperties().TripleIndex(triple)
	if linearIndex >= 0 {
		id := ids.ObjectShortNames
		if longName {
			id = ids.ObjectLongNames
		}
		key := resource.KeyOf(id, view.model.currentLang, linearIndex)
		oldValue, _ := view.textCache.Text(key)

		if oldValue != newValue {
			command := setObjectTextCommand{
				model:   &view.model,
				triple:  view.model.currentObject,
				bitmap:  view.model.currentBitmap,
				key:     key,
				oldData: view.cp.Encode(oldValue),
				newData: view.cp.Encode(text.Blocked(newValue)[0]),
			}
			view.commander.Queue(command)
		}
	}
}

func (view *View) requestSetObjectProperties(modifier func(*object.Properties)) {
	command := setObjectPropertiesCommand{
		model:  &view.model,
		triple: view.model.currentObject,
		bitmap: view.model.currentBitmap,
	}
	currentProp, err := view.mod.ObjectProperties().ForObject(command.triple)
	if err != nil {
		return
	}
	command.oldProperties = currentProp.Clone()
	command.newProperties = currentProp.Clone()
	modifier(&command.newProperties)
	view.commander.Queue(command)
}

func (view *View) renderCommonProperties(readOnly bool, properties *object.Properties) {
	intIdentity := func(u values.Unifier) int { return u.Unified().(int) }
	intFormat := func(value int) string { return "%d" }

	if imgui.TreeNode("Flags") {
		renderFlag := func(flag object.CommonFlag) {
			flagUnifier := values.NewUnifier()
			flagUnifier.Add(properties.Common.Flags.Has(flag))
			values.RenderUnifiedCheckboxCombo(readOnly, false, flag.String(), flagUnifier, func(newValue bool) {
				view.requestSetObjectProperties(func(prop *object.Properties) {
					if newValue {
						prop.Common.Flags = prop.Common.Flags.With(flag)
					} else {
						prop.Common.Flags = prop.Common.Flags.Without(flag)
					}
				})
			})
		}
		for _, flag := range object.CommonFlags() {
			renderFlag(flag)
		}

		imgui.TreePop()
	}

	lightTypes := object.LightTypes()
	lightTypeUnifier := values.NewUnifier()
	lightTypeUnifier.Add(int(properties.Common.Flags.LightType()))
	values.RenderUnifiedCombo(readOnly, false, "LightType", lightTypeUnifier, intIdentity,
		func(value int) string {
			return object.LightType(value).String()
		}, len(lightTypes), func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.Flags = prop.Common.Flags.WithLightType(object.LightType(newValue))
			})
		})

	useModeUnifier := values.NewUnifier()
	useModeUnifier.Add(int(properties.Common.Flags.UseMode()))
	values.RenderUnifiedCombo(readOnly, false, "UseMode", useModeUnifier, intIdentity,
		func(value int) string {
			return object.UseMode(value).String()
		}, 4, func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.Flags = prop.Common.Flags.WithUseMode(object.UseMode(newValue))
			})
		})

	massUnifier := values.NewUnifier()
	massUnifier.Add(int(properties.Common.Mass))
	values.RenderUnifiedSliderInt(readOnly, false, "Mass", massUnifier, intIdentity, intFormat, -1, 5000,
		func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.Mass = int32(newValue)
			})
		})

	hitpointsUnifier := values.NewUnifier()
	hitpointsUnifier.Add(int(properties.Common.Hitpoints))
	values.RenderUnifiedSliderInt(readOnly, false, "Hitpoints", hitpointsUnifier, intIdentity, intFormat, 0, 10000,
		func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.Hitpoints = int16(newValue)
			})
		})

	armorUnifier := values.NewUnifier()
	armorUnifier.Add(int(properties.Common.Armor))
	values.RenderUnifiedSliderInt(readOnly, false, "Armor", armorUnifier, intIdentity, intFormat, 0, 255,
		func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.Armor = byte(newValue)
			})
		})

	renderTypeUnifier := values.NewUnifier()
	renderTypeUnifier.Add(int(properties.Common.RenderType))
	values.RenderUnifiedCombo(readOnly, false, "Render Type", renderTypeUnifier, intIdentity,
		func(value int) string { return object.RenderType(value).String() },
		len(object.RenderTypes()), func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.RenderType = object.RenderType(newValue)
			})
		})

	physicsModelUnifier := values.NewUnifier()
	physicsModelUnifier.Add(int(properties.Common.PhysicsModel))
	values.RenderUnifiedCombo(readOnly, false, "Physics Model", physicsModelUnifier, intIdentity,
		func(value int) string { return object.PhysicsModel(value).String() },
		len(object.PhysicsModels()), func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.PhysicsModel = object.PhysicsModel(newValue)
			})
		})

	hardnessUnifier := values.NewUnifier()
	hardnessUnifier.Add(int(properties.Common.Hardness))
	values.RenderUnifiedSliderInt(readOnly, false, "Hardness", hardnessUnifier, intIdentity, intFormat, 0, object.HardnessLimit,
		func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.Hardness = byte(newValue)
			})
		})

	physicsXRUnifier := values.NewUnifier()
	physicsXRUnifier.Add(int(properties.Common.PhysicsXR))
	values.RenderUnifiedSliderInt(readOnly, false, "Physics XR", physicsXRUnifier, intIdentity, intFormat, 0, object.PhysicsXRLimit,
		func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.PhysicsXR = byte(newValue)
			})
		})

	physicsZUnifier := values.NewUnifier()
	physicsZUnifier.Add(int(properties.Common.PhysicsZ))
	values.RenderUnifiedSliderInt(readOnly, false, "Physics Z", physicsZUnifier, intIdentity, intFormat, 0, 255,
		func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.PhysicsZ = byte(newValue)
			})
		})

	if imgui.TreeNode("Vulnerabilities") {
		renderVulnerability := func(damageType object.DamageType) {
			damageUnifier := values.NewUnifier()
			damageUnifier.Add(properties.Common.Vulnerabilities.Has(damageType))
			values.RenderUnifiedCheckboxCombo(readOnly, false, damageType.String(), damageUnifier, func(newValue bool) {
				view.requestSetObjectProperties(func(prop *object.Properties) {
					if newValue {
						prop.Common.Vulnerabilities = prop.Common.Vulnerabilities.With(damageType)
					} else {
						prop.Common.Vulnerabilities = prop.Common.Vulnerabilities.Without(damageType)
					}
				})
			})
		}
		for _, damageType := range object.DamageTypes() {
			renderVulnerability(damageType)
		}

		primaryUnifier := values.NewUnifier()
		primaryUnifier.Add(properties.Common.SpecialVulnerabilities.PrimaryValue())
		values.RenderUnifiedSliderInt(readOnly, false, "Primary (double dmg)", primaryUnifier, intIdentity, intFormat, 0, object.SpecialDamageTypeLimit,
			func(newValue int) {
				view.requestSetObjectProperties(func(prop *object.Properties) {
					prop.Common.SpecialVulnerabilities = prop.Common.SpecialVulnerabilities.WithPrimaryValue(newValue)
				})
			})

		superUnifier := values.NewUnifier()
		superUnifier.Add(properties.Common.SpecialVulnerabilities.SuperValue())
		values.RenderUnifiedSliderInt(readOnly, false, "Super (quad dmg)", superUnifier, intIdentity, intFormat, 0, object.SpecialDamageTypeLimit,
			func(newValue int) {
				view.requestSetObjectProperties(func(prop *object.Properties) {
					prop.Common.SpecialVulnerabilities = prop.Common.SpecialVulnerabilities.WithSuperValue(newValue)
				})
			})

		imgui.TreePop()
	}

	defenseUnifier := values.NewUnifier()
	defenseUnifier.Add(int(properties.Common.Defense))
	values.RenderUnifiedSliderInt(readOnly, false, "Defense", defenseUnifier, intIdentity,
		func(value int) string {
			if value == object.DefenseNoCriticals {
				return "%d - (no criticals)"
			}
			return "%d"
		}, 0, 255,
		func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.Defense = byte(newValue)
			})
		})

	toughnessUnifier := values.NewUnifier()
	toughnessUnifier.Add(int(properties.Common.Toughness))
	values.RenderUnifiedSliderInt(readOnly, false, "Toughness", toughnessUnifier, intIdentity,
		func(value int) string {
			if value == object.ToughnessNoDamage {
				return "(no damage) -- raw: %d"
			}
			return fmt.Sprintf("%d:1 dmg -- raw: %%d", 1<<uint(value))
		},
		0, 7,
		func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.Toughness = byte(newValue)
			})
		})

	mfdOrMeshIDUnifier := values.NewUnifier()
	mfdOrMeshIDUnifier.Add(int(properties.Common.MfdOrMeshID))
	values.RenderUnifiedSliderInt(readOnly, false, "MFD/Mesh ID", mfdOrMeshIDUnifier, intIdentity, intFormat, 0, 1000,
		func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.MfdOrMeshID = uint16(newValue)
			})
		})

	bitmap3DBitmapNumUnifier := values.NewUnifier()
	bitmap3DBitmapNumUnifier.Add(int(properties.Common.Bitmap3D.BitmapNumber()))
	values.RenderUnifiedSliderInt(readOnly, false, "Bitmap Number", bitmap3DBitmapNumUnifier, intIdentity, intFormat, 0, int(object.Bitmap3DBitmapNumberLimit),
		func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.Bitmap3D = prop.Common.Bitmap3D.WithBitmapNumber(uint16(newValue))
			})
		})
	bitmap3DFrameNumUnifier := values.NewUnifier()
	bitmap3DFrameNumUnifier.Add(int(properties.Common.Bitmap3D.FrameNumber()))
	values.RenderUnifiedSliderInt(readOnly, false, "Frame Number", bitmap3DFrameNumUnifier, intIdentity, intFormat, 0, int(object.Bitmap3DFrameNumberLimit),
		func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.Bitmap3D = prop.Common.Bitmap3D.WithFrameNumber(uint16(newValue))
			})
		})
	bitmap3DAnimUnifier := values.NewUnifier()
	bitmap3DAnimUnifier.Add(properties.Common.Bitmap3D.Animation())
	values.RenderUnifiedCheckboxCombo(readOnly, false, "Animation", bitmap3DAnimUnifier,
		func(newValue bool) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.Bitmap3D = prop.Common.Bitmap3D.WithAnimation(newValue)
			})
		})
	bitmap3DRepeatUnifier := values.NewUnifier()
	bitmap3DRepeatUnifier.Add(properties.Common.Bitmap3D.Repeat())
	values.RenderUnifiedCheckboxCombo(readOnly, false, "Repeat", bitmap3DRepeatUnifier,
		func(newValue bool) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.Bitmap3D = prop.Common.Bitmap3D.WithRepeat(newValue)
			})
		})

	destroyEffectValueUnifier := values.NewUnifier()
	destroyEffectValueUnifier.Add(int(properties.Common.DestroyEffect.Value()))
	values.RenderUnifiedSliderInt(readOnly, false, "DestroyEffect Value", destroyEffectValueUnifier,
		intIdentity, intFormat, 0, int(object.DestroyEffectValueLimit),
		func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.DestroyEffect = prop.Common.DestroyEffect.WithValue(byte(newValue))
			})
		})
	destroyEffectPlaySoundUnifier := values.NewUnifier()
	destroyEffectPlaySoundUnifier.Add(properties.Common.DestroyEffect.PlaySound())
	values.RenderUnifiedCheckboxCombo(readOnly, false, "DestroyEffect PlaySound", destroyEffectPlaySoundUnifier,
		func(newValue bool) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.DestroyEffect = prop.Common.DestroyEffect.WithSound(newValue)
			})
		})
	destroyEffectShowExplosionUnifier := values.NewUnifier()
	destroyEffectShowExplosionUnifier.Add(properties.Common.DestroyEffect.ShowExplosion())
	values.RenderUnifiedCheckboxCombo(readOnly, false, "DestroyEffect ShowExplosion", destroyEffectShowExplosionUnifier,
		func(newValue bool) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.DestroyEffect = prop.Common.DestroyEffect.WithExplosion(newValue)
			})
		})
}

func (view *View) renderGenericProperties(readOnly bool, properties *object.Properties) {
	readInterpreter := objprop.GenericProperties(view.model.currentObject.Class, properties.Generic)
	view.createPropertyControls(readOnly, readInterpreter, func(key string, modifier func(uint32) uint32) {
		view.requestSetObjectProperties(func(prop *object.Properties) {
			writeInterpreter := objprop.GenericProperties(view.model.currentObject.Class, prop.Generic)
			view.setInterpreterValueKeyed(writeInterpreter, key, modifier)
		})
	})
}

func (view *View) renderSpecificProperties(readOnly bool, properties *object.Properties) {
	readInterpreter := objprop.SpecificProperties(view.model.currentObject, properties.Specific)
	view.createPropertyControls(readOnly, readInterpreter, func(key string, modifier func(uint32) uint32) {
		view.requestSetObjectProperties(func(prop *object.Properties) {
			writeInterpreter := objprop.SpecificProperties(view.model.currentObject, prop.Specific)
			view.setInterpreterValueKeyed(writeInterpreter, key, modifier)
		})
	})
}

func (view *View) createPropertyControls(readOnly bool, rootInterpreter *interpreters.Instance,
	updater func(string, func(uint32) uint32)) {
	objTypeRenderer := values.ObjectTypeControlRenderer{
		Meta:      view.mod.ObjectProperties(),
		TextCache: view.textCache,
	}

	var processInterpreter func(string, *interpreters.Instance)
	processInterpreter = func(path string, interpreter *interpreters.Instance) {
		for _, key := range interpreter.Keys() {
			fullKey := path + key
			unifier := values.NewUnifier()
			unifier.Add(int32(interpreter.Get(key)))
			simplifier := values.StandardSimplifier(readOnly, false, fullKey, unifier,
				func(modifier func(uint32) uint32) {
					updater(fullKey, modifier)
				}, objTypeRenderer)

			interpreter.Describe(key, simplifier)
		}

		for _, key := range interpreter.ActiveRefinements() {
			fullKey := path + key
			if len(fullKey) > 0 {
				imgui.Separator()
				imgui.Text(fullKey + ":")
			}
			processInterpreter(fullKey+".", interpreter.Refined(key))
		}
	}
	processInterpreter("", rootInterpreter)
}

func (view *View) setInterpreterValueKeyed(instance *interpreters.Instance, key string, modifier func(uint32) uint32) {
	resolvedInterpreter := instance
	keys := strings.Split(key, ".")
	keyCount := len(keys)
	if len(keys) > 1 {
		for _, subKey := range keys[:keyCount-1] {
			resolvedInterpreter = resolvedInterpreter.Refined(subKey)
		}
	}
	resolvedInterpreter.Set(keys[keyCount-1], modifier(resolvedInterpreter.Get(keys[keyCount-1])))
}

func (view *View) renderObjectBitmap() {
	render.TextureImage("BitmapImage", view.imageCache, view.currentBitmapKey(), imgui.Vec2{X: 320 * view.guiScale, Y: 240 * view.guiScale})
	if imgui.Button("Clear") {
		view.requestClearBitmap()
	}
	imgui.SameLine()
	if imgui.Button("Import") {
		view.requestImportBitmap()
	}
	imgui.SameLine()
	if imgui.Button("Export") {
		view.requestExportBitmap()
	}
	if view.hasModCurrentBitmap() {
		imgui.SameLine()
		if imgui.Button("Remove") {
			view.requestSetBitmapData(nil)
		}
	}
}

func (view *View) requestClearBitmap() {
	bmp := bitmap.Bitmap{
		Header: bitmap.Header{
			Width:  1,
			Height: 1,
		},
		Pixels: []byte{0x00},
	}
	view.requestSetBitmap(bmp)
}

func (view *View) requestImportBitmap() {
	paletteRetriever := func() (bitmap.Palette, error) {
		palette, err := view.paletteCache.Palette(0)
		if err != nil {
			return bitmap.Palette{}, err
		}
		return palette.Palette(), nil
	}
	external.ImportImage(view.modalStateMachine, paletteRetriever, func(bmp bitmap.Bitmap) {
		view.requestSetBitmap(bmp)
	})
}

func (view *View) requestExportBitmap() {
	key := view.currentBitmapKey()
	texture, err := view.imageCache.Texture(key)
	if err != nil {
		return
	}
	palette, err := view.paletteCache.Palette(0)
	if err != nil {
		return
	}
	rawPalette := palette.Palette()
	filename := fmt.Sprintf("%02d_%d_%02d-%02d.png",
		view.model.currentObject.Class, view.model.currentObject.Subclass, view.model.currentObject.Type, view.model.currentBitmap)
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

func (view *View) hasModCurrentBitmap() bool {
	key := view.currentBitmapKey()
	return len(view.mod.ModifiedBlock(key.Lang, key.ID, key.Index)) > 0
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
	bmp.Header.Type = bitmap.TypeFlat8Bit
	bmp.Header.WidthFactor = highestBitShift(bmp.Header.Width)
	bmp.Header.HeightFactor = highestBitShift(bmp.Header.Height)
	bmp.Header.Stride = uint16(bmp.Header.Width)
	data := bitmap.Encode(&bmp, 0)
	view.requestSetBitmapData(data)
}

func (view *View) requestSetBitmapData(newData []byte) {
	resourceKey := view.currentBitmapKey()

	command := setObjectBitmapCommand{
		model:  &view.model,
		triple: view.model.currentObject,
		bitmap: view.model.currentBitmap,

		resourceKey: resourceKey,
		oldData:     view.mod.ModifiedBlock(resourceKey.Lang, resourceKey.ID, resourceKey.Index),
		newData:     newData,
	}
	view.commander.Queue(command)
}

func (view *View) currentBitmapKey() resource.Key {
	return view.currentBitmapKeyFor(view.model.currentBitmap)
}

func (view *View) currentBitmapKeyFor(offset int) resource.Key {
	baseOffset := 1
	done := false
	view.mod.ObjectProperties().Iterate(func(triple object.Triple, prop *object.Properties) bool {
		if !done && triple == view.model.currentObject {
			done = true
		}
		if done {
			return false
		}
		numExtra := int(prop.Common.Bitmap3D.FrameNumber())
		baseOffset += 3 + numExtra
		return true
	})

	return resource.KeyOf(ids.ObjectBitmaps, resource.LangAny, baseOffset+offset)
}
