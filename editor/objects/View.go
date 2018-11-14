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
		imgui.PushItemWidth(-200 * view.guiScale)
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
	values.RenderUnifiedSliderInt(readOnly, false, "Hardness", hardnessUnifier, intIdentity, intFormat, 0, 255,
		func(newValue int) {
			view.requestSetObjectProperties(func(prop *object.Properties) {
				prop.Common.Hardness = byte(newValue)
			})
		})
}
