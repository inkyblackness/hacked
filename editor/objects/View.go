package objects

import (
	"fmt"

	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/external"
	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
)

// View provides edit controls for game objects.
type View struct {
	mod          *model.Mod
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
func NewView(mod *model.Mod, textCache *text.Cache, cp text.Codepage,
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
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 600 * view.guiScale, Y: 600 * view.guiScale}, imgui.ConditionOnce)
		if imgui.BeginV("Game Objects", view.WindowOpen(), imgui.WindowFlagsNoCollapse|imgui.WindowFlagsHorizontalScrollbar) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *View) renderContent() {
	if imgui.BeginChildV("Properties", imgui.Vec2{X: -100 * view.guiScale, Y: 0}, false, imgui.WindowFlagsHorizontalScrollbar) {
		imgui.PushItemWidth(-250 * view.guiScale)
		classString := func(class object.Class) string {
			return fmt.Sprintf("%2d: %v", int(class), class)
		}
		if imgui.BeginCombo("Object Class", classString(view.model.currentObject.Class)) {
			for _, class := range object.Classes() {
				if imgui.SelectableV(classString(class), class == view.model.currentObject.Class, 0, imgui.Vec2{}) {
					view.model.currentObject = object.TripleFrom(int(class), 0, 0)
				}
			}
			imgui.EndCombo()
		}
		if imgui.BeginCombo("Object Type", view.tripleName(view.model.currentObject)) {
			allTypes := view.mod.ObjectProperties().TriplesInClass(view.model.currentObject.Class)
			for _, triple := range allTypes {
				if imgui.SelectableV(view.tripleName(triple), triple == view.model.currentObject, 0, imgui.Vec2{}) {
					view.model.currentObject = triple
				}
			}
			imgui.EndCombo()
		}

		readOnly := !view.mod.HasModifyableObjectProperties()

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

		properties, err := view.mod.ObjectProperties().ForObject(view.model.currentObject)
		if err == nil {
			if imgui.TreeNodeV("Common Properties", imgui.TreeNodeFlagsDefaultOpen|imgui.TreeNodeFlagsFramed) {
				view.renderCommonProperties(readOnly, properties)
				imgui.TreePop()
			}
			if imgui.TreeNodeV("Generic Properties", imgui.TreeNodeFlagsFramed) {
				imgui.Text("(not yet)")
				imgui.TreePop()
			}
			if imgui.TreeNodeV("Specific Properties", imgui.TreeNodeFlagsFramed) {
				imgui.Text("(not yet)")
				imgui.TreePop()
			}
		}

		imgui.PopItemWidth()
	}
	imgui.EndChild()
	//imgui.SameLine()

	//imgui.BeginGroup()
	// view.renderObjectBitmap()
	//imgui.EndGroup()
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
			damageUnifier.Add(bool(properties.Common.Vulnerabilities.Has(damageType)))
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
