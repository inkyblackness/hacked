package levels

import (
	"fmt"
	"math"
	"strings"

	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/hacked/ui/opengl"
	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/numbers"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

const shouldRenderRotationView = false

// ObjectsView is for object properties.
type ObjectsView struct {
	gameObjects     *edit.GameObjectsService
	levelSelection  *edit.LevelSelectionService
	editor          *edit.LevelEditorService
	varInfoProvider archive.GameVariableInfoProvider
	textCache       *text.Cache
	textureCache    *graphics.TextureCache

	guiScale float32
	registry cmd.Registry

	textureRenderer *render.TextureRenderer
	orientationView *render.OrientationView

	model objectsViewModel
}

// NewObjectsView returns a new instance.
func NewObjectsView(gameObjects *edit.GameObjectsService,
	editor *edit.LevelEditorService,
	levelSelection *edit.LevelSelectionService,
	varInfoProvider archive.GameVariableInfoProvider,
	guiScale float32, textCache *text.Cache, textureCache *graphics.TextureCache,
	registry cmd.Registry, gl opengl.OpenGL) *ObjectsView {
	textureRenderer := render.NewTextureRenderer(gl)
	viewMatrix := mgl.LookAt(0.35, 0.35, -1, 0, 0, 0, 1.0, 0.0, 0.0)
	context := render.Context{
		OpenGL:           gl,
		ViewMatrix:       &viewMatrix,
		ProjectionMatrix: mgl.Ortho(-0.5, 0.5, 0.5, -0.5, -10, 10.0),
	}
	// Set up an orientation view. The orientation parameter would match the
	// map as seen from above, and the rotation to match how the objects rotate in world.
	orientationView := render.NewOrientationView(context,
		mgl.Ident4().
			Mul4(mgl.HomogRotate3D(math.Pi/2, mgl.Vec3{0.0, 0.0, 1.0})).
			Mul4(mgl.HomogRotate3D(math.Pi, mgl.Vec3{1.0, 0.0, 0.0})),
		mgl.Vec3{1.0, 1.0, -1.0})

	view := &ObjectsView{
		gameObjects:     gameObjects,
		levelSelection:  levelSelection,
		editor:          editor,
		varInfoProvider: varInfoProvider,
		textCache:       textCache,
		textureCache:    textureCache,

		guiScale: guiScale,
		registry: registry,

		textureRenderer: textureRenderer,
		orientationView: orientationView,

		model: freshObjectsViewModel(),
	}
	return view
}

// WindowOpen returns the flag address, to be used with the main menu.
func (view *ObjectsView) WindowOpen() *bool {
	return &view.model.windowOpen
}

// Render renders the view.
func (view *ObjectsView) Render() {
	if view.model.restoreFocus {
		imgui.SetNextWindowFocus()
		view.model.restoreFocus = false
		view.model.windowOpen = true
	}
	if view.model.windowOpen {
		lvl := view.editor.Level()
		objects := view.editor.Objects()
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 500 * view.guiScale}, imgui.ConditionFirstUseEver)
		title := fmt.Sprintf("Level Objects, %d selected", len(objects))
		readOnly := view.editor.IsReadOnly()
		if readOnly {
			title += hintReadOnly
		}
		if imgui.BeginV(title+"###Level Objects", view.WindowOpen(), imgui.WindowFlagsHorizontalScrollbar|imgui.WindowFlagsAlwaysVerticalScrollbar) {
			view.renderContent(lvl, objects, readOnly)
		}
		imgui.End()
	}
}

func (view *ObjectsView) renderContent(lvl *level.Level, objects []*level.ObjectMainEntry, readOnly bool) {
	classUnifier := values.NewUnifier()
	typeUnifier := values.NewUnifier()
	zUnifier := values.NewUnifier()
	tileXUnifier := values.NewUnifier()
	fineXUnifier := values.NewUnifier()
	tileYUnifier := values.NewUnifier()
	fineYUnifier := values.NewUnifier()
	rotationXUnifier := values.NewUnifier()
	rotationYUnifier := values.NewUnifier()
	rotationZUnifier := values.NewUnifier()
	hitpointsUnifier := values.NewUnifier()

	for _, obj := range objects {
		classUnifier.Add(object.TripleFrom(int(obj.Class), 0, 0))
		typeUnifier.Add(obj.Triple())
		zUnifier.Add(obj.Z)
		tileXUnifier.Add(obj.X.Tile())
		fineXUnifier.Add(obj.X.Fine())
		tileYUnifier.Add(obj.Y.Tile())
		fineYUnifier.Add(obj.Y.Fine())
		rotationXUnifier.Add(int32(obj.XRotation))
		rotationYUnifier.Add(int32(obj.YRotation))
		rotationZUnifier.Add(int32(obj.ZRotation))
		hitpointsUnifier.Add(obj.Hitpoints)
	}

	imgui.PushItemWidth(-150 * view.guiScale)
	columns, rows, levelHeight := lvl.Size()

	objectHeightFormatter := objectHeightFormatterFor(levelHeight)

	if !readOnly {
		classString := func(class object.Class) string {
			active, capacity := lvl.ObjectClassStats(class)
			if class == object.ClassSmallStuff {
				// The engine backs up inventory items in the level during a level transition.
				// For this it requires space. If there is none, bugs do appear.
				// This limit could apply for several classes, yet only small stuff is being used.
				practicalLimit := 0
				if capacity > level.InventorySize {
					practicalLimit = capacity - level.InventorySize
				}
				return fmt.Sprintf("%2d: %v -- %d/%d (%d)", int(class), class, active, practicalLimit, capacity)
			}
			return fmt.Sprintf("%2d: %v -- %d/%d", int(class), class, active, capacity)
		}
		if imgui.BeginCombo("New Object Class", classString(view.editor.NewObjectTriple().Class)) {
			for _, class := range object.Classes() {
				if imgui.SelectableV(classString(class), class == view.editor.NewObjectTriple().Class, 0, imgui.Vec2{}) {
					view.editor.SetNewObjectTriple(object.TripleFrom(int(class), 0, 0))
				}
			}
			imgui.EndCombo()
		}
		if imgui.BeginCombo("New Object Type", view.tripleName(view.editor.NewObjectTriple())) {
			allTypes := view.gameObjects.AllProperties().TriplesInClass(view.editor.NewObjectTriple().Class)
			bitmapInfo := view.gameObjects.BitmapInfo()
			lineHeight := imgui.TextLineHeight()
			iconSize := imgui.Vec2{X: lineHeight * view.guiScale, Y: lineHeight * view.guiScale}
			for index, triple := range allTypes {
				info := bitmapInfo[triple]
				iconKey := resource.KeyOf(ids.ObjectBitmaps, resource.LangAny, info.Start+info.IconRecommendation)
				imgui.PushIDInt(index)
				imgui.BeginGroup()
				render.TextureImage("Icon", view.textureCache, iconKey, iconSize)
				imgui.SameLine()
				if imgui.SelectableV(view.tripleName(triple), triple == view.editor.NewObjectTriple(), 0, imgui.Vec2{}) {
					view.editor.SetNewObjectTriple(triple)
				}
				imgui.EndGroup()
				imgui.PopID()
			}
			imgui.EndCombo()
		}
		if imgui.Button("Delete Selected") {
			view.requestDeleteObjects()
		}
		imgui.Separator()
	}

	selectedIDs := view.levelSelection.CurrentSelectedObjects()
	switch {
	case len(selectedIDs) == 1:
		imgui.LabelText("ID", fmt.Sprintf("%3d", int(selectedIDs[0])))
	case len(selectedIDs) > 1:
		imgui.LabelText("ID", "(multiple)")
	default:
		imgui.LabelText("ID", "")
	}

	objTypeRenderer := values.ObjectTypeControlRenderer{
		Meta:      view.gameObjects.AllProperties(),
		TextCache: view.textCache,
	}
	objTypeRenderer.Render(readOnly, "Object Type", classUnifier, typeUnifier,
		func(u values.Unifier) object.Triple { return u.Unified().(object.Triple) },
		func(newValue object.Triple) {
			view.requestChangeObjects("ChangeBaseType", func(entry *level.ObjectMainEntry) {
				entry.Subclass = newValue.Subclass
				entry.Type = newValue.Type
			})
		})

	if imgui.TreeNodeV("Base Properties", imgui.TreeNodeFlagsDefaultOpen|imgui.TreeNodeFlagsFramed) {
		values.RenderUnifiedSliderInt(readOnly, "Z", zUnifier,
			func(u values.Unifier) int { return int(u.Unified().(level.HeightUnit)) },
			objectHeightFormatter,
			0, 0xFF,
			func(newValue int) {
				view.requestChangeObjects("ChangeBaseZ", func(entry *level.ObjectMainEntry) { entry.Z = level.HeightUnit(newValue) })
			})
		values.RenderUnifiedSliderInt(readOnly, "Tile X", tileXUnifier,
			func(u values.Unifier) int { return int(u.Unified().(byte)) },
			func(value int) string { return "%d" },
			0, columns-1,
			func(newValue int) {
				view.requestChangeObjects("ChangeBaseTileX", func(entry *level.ObjectMainEntry) { entry.X = level.CoordinateAt(byte(newValue), entry.X.Fine()) })
			})
		values.RenderUnifiedSliderInt(readOnly, "Fine X", fineXUnifier,
			func(u values.Unifier) int { return int(u.Unified().(byte)) },
			func(value int) string { return "%d" },
			0, level.FineCoordinatesPerTileSide-1,
			func(newValue int) {
				view.requestChangeObjects("ChangeBaseFineX", func(entry *level.ObjectMainEntry) { entry.X = level.CoordinateAt(entry.X.Tile(), byte(newValue)) })
			})
		values.RenderUnifiedSliderInt(readOnly, "Tile Y", tileYUnifier,
			func(u values.Unifier) int { return int(u.Unified().(byte)) },
			func(value int) string { return "%d" },
			0, rows-1,
			func(newValue int) {
				view.requestChangeObjects("ChangeBaseTileY", func(entry *level.ObjectMainEntry) { entry.Y = level.CoordinateAt(byte(newValue), entry.Y.Fine()) })
			})
		values.RenderUnifiedSliderInt(readOnly, "Fine Y", fineYUnifier,
			func(u values.Unifier) int { return int(u.Unified().(byte)) },
			func(value int) string { return "%d" },
			0, level.FineCoordinatesPerTileSide-1,
			func(newValue int) {
				view.requestChangeObjects("ChangeFineY", func(entry *level.ObjectMainEntry) { entry.Y = level.CoordinateAt(entry.Y.Tile(), byte(newValue)) })
			})

		values.RenderUnifiedRotation(readOnly, "Rotation X", rotationXUnifier,
			0, 0x0FF,
			values.RotationInfo{Horizontal: true, Positive: false, Clockwise: false},
			func(newValue int) {
				view.requestChangeObjects("ChangeBaseRotationX", func(entry *level.ObjectMainEntry) { entry.XRotation = level.RotationUnit(newValue) })
			})
		values.RenderUnifiedRotation(readOnly, "Rotation Y", rotationYUnifier,
			0, 0x0FF,
			values.RotationInfo{Horizontal: false, Positive: true, Clockwise: true},
			func(newValue int) {
				view.requestChangeObjects("ChangeBaseRotationY", func(entry *level.ObjectMainEntry) { entry.YRotation = level.RotationUnit(newValue) })
			})
		values.RenderUnifiedRotation(readOnly, "Rotation Z", rotationZUnifier,
			0, 0x0FF,
			values.RotationInfo{Horizontal: false, Positive: false, Clockwise: true},
			func(newValue int) {
				view.requestChangeObjects("ChangeBaseRotationZ", func(entry *level.ObjectMainEntry) { entry.ZRotation = level.RotationUnit(newValue) })
			})

		if shouldRenderRotationView {
			var rotation mgl.Vec3
			if rotationXUnifier.IsUnique() {
				rotation[0] = float32(rotationXUnifier.Unified().(int32)*360) / 256.0
			}
			if rotationYUnifier.IsUnique() {
				rotation[1] = float32(rotationYUnifier.Unified().(int32)*360) / 256.0
			}
			if rotationZUnifier.IsUnique() {
				rotation[2] = float32(rotationZUnifier.Unified().(int32)*360) / 256.0
			}
			view.textureRenderer.Render(func() { view.orientationView.Render(rotation) })
			imgui.ImageV(gui.TextureIDForColorTexture(view.textureRenderer.Handle()),
				imgui.Vec2{X: 200 * view.guiScale, Y: 200 * view.guiScale},
				imgui.Vec2{X: 0.0, Y: 0.0}, imgui.Vec2{X: 1.0, Y: 1.0},
				imgui.Vec4{X: 1, Y: 1, Z: 1, W: 1.0}, imgui.Vec4{X: 0, Y: 0, Z: 0, W: 0})
		}

		values.RenderUnifiedSliderInt(readOnly, "Hitpoints", hitpointsUnifier,
			func(u values.Unifier) int { return int(u.Unified().(int16)) },
			func(value int) string { return "%d" },
			0, 10000,
			func(newValue int) {
				view.requestChangeObjects("ChangeBaseHitpoints", func(entry *level.ObjectMainEntry) { entry.Hitpoints = int16(newValue) })
			})
		imgui.TreePop()
	}
	if imgui.TreeNodeV("Extra Properties", imgui.TreeNodeFlagsFramed) {
		view.renderProperties(lvl, objects, readOnly,
			func(obj *level.ObjectMainEntry) *interpreters.Instance { return lvl.ObjectExtraData(obj) })
		imgui.TreePop()
	}
	if imgui.TreeNodeV("Class Properties", imgui.TreeNodeFlagsFramed) {
		view.renderProperties(lvl, objects, readOnly,
			func(obj *level.ObjectMainEntry) *interpreters.Instance { return lvl.ObjectClassData(obj) })
		view.renderBlockPuzzleControl(lvl, objects, readOnly)
		imgui.TreePop()
	}

	imgui.PopItemWidth()
}

func (view *ObjectsView) renderProperties(lvl *level.Level, objects []*level.ObjectMainEntry, readOnly bool,
	interpreterRetriever func(*level.ObjectMainEntry) *interpreters.Instance) {
	propertyUnifier := make(map[string]*values.Unifier)
	propertyDescribers := make(map[string]func(*interpreters.Simplifier))
	var propertyOrder []string
	describer := func(interpreter *interpreters.Instance, key string) func(simpl *interpreters.Simplifier) {
		return func(simpl *interpreters.Simplifier) { interpreter.Describe(key, simpl) }
	}

	var unifyInterpreter func(string, *interpreters.Instance, bool, map[string]bool)
	unifyInterpreter = func(path string, interpreter *interpreters.Instance, first bool, thisKeys map[string]bool) {
		for _, key := range interpreter.Keys() {
			fullPath := path + key
			thisKeys[fullPath] = true
			if unifier, existing := propertyUnifier[fullPath]; existing || first {
				if !existing {
					u := values.NewUnifier()
					unifier = &u
					propertyUnifier[fullPath] = unifier
					propertyDescribers[fullPath] = describer(interpreter, key)
					propertyOrder = append(propertyOrder, fullPath)
				}
				unifier.Add(int32(interpreter.Get(key)))
			}
		}
		for _, key := range interpreter.ActiveRefinements() {
			unifyInterpreter(path+key+".", interpreter.Refined(key), first, thisKeys)
		}
	}

	for index, obj := range objects {
		interpreter := interpreterRetriever(obj)

		thisKeys := make(map[string]bool)
		unifyInterpreter("", interpreter, index == 0, thisKeys)
		{
			var toRemove []string
			for previousKey := range propertyUnifier {
				if !thisKeys[previousKey] {
					toRemove = append(toRemove, previousKey)
				}
			}
			for _, key := range toRemove {
				delete(propertyUnifier, key)
			}
		}
	}

	lastTitle := ""
	for _, key := range propertyOrder {
		if unifier, existing := propertyUnifier[key]; existing {
			subKeys := strings.Split(key, ".")
			baseKey := ""
			if len(subKeys) > 1 {
				baseKey = subKeys[0]
			}
			if baseKey != lastTitle {
				imgui.Separator()
				imgui.Text(baseKey + ":")
				lastTitle = baseKey
			}
			view.renderPropertyControl(lvl, readOnly, key, *unifier, propertyDescribers[key],
				func(modifier func(uint32) uint32) {
					view.requestPropertiesChange(interpreterRetriever, key, modifier) // nolint: scopelint
				})
		}
	}
}

func (view *ObjectsView) renderPropertyControl(lvl *level.Level, readOnly bool,
	fullKey string, unifier values.Unifier, describer func(*interpreters.Simplifier),
	updater func(func(uint32) uint32)) {
	keys := strings.Split(fullKey, ".")
	key := keys[len(keys)-1]
	label := key + "###" + fullKey
	_, _, levelHeight := lvl.Size()
	tileHeightFormatter := tileHeightFormatterFor(levelHeight)
	objectHeightFormatter := objectHeightFormatterFor(levelHeight)
	moveTileHeightFormatter := moveTileHeightFormatterFor(levelHeight)

	objTypeRenderer := values.ObjectTypeControlRenderer{
		Meta:      view.gameObjects.AllProperties(),
		TextCache: view.textCache,
	}
	simplifier := values.StandardSimplifier(readOnly, fullKey, unifier, updater, objTypeRenderer)

	simplifier.SetObjectIDHandler(func() {
		values.RenderUnifiedSliderInt(readOnly, label, unifier,
			func(u values.Unifier) int { return int(u.Unified().(int32)) },
			func(value int) string { return "%d" },
			0, lvl.ObjectCapacity(),
			func(newValue int) {
				updater(func(oldValue uint32) uint32 { return uint32(newValue) })
			})
	})

	addVariableKey := func() {
		variableTypes := []string{"boolean", "Integer"}
		typeMask := 0x1000

		typeLabel := key + "-Type###" + fullKey + "-Type"
		values.RenderUnifiedCombo(readOnly, typeLabel, unifier,
			func(u values.Unifier) int {
				value := u.Unified().(int32)
				result := 0
				if (value & int32(typeMask)) != 0 {
					result = 1
				}
				return result
			},
			func(index int) string { return variableTypes[index] },
			len(variableTypes),
			func(newIndex int) {
				setMask := uint32(0)
				if newIndex > 0 {
					setMask |= uint32(typeMask)
				}
				updater(func(oldValue uint32) uint32 {
					return (oldValue & uint32(0xE000)) | setMask
				})
			})

		valueLabel := key + "-Value###" + fullKey + "-Value"
		switch {
		case unifier.IsUnique():
			key := unifier.Unified().(int32)
			isIntegerVar := (key & int32(typeMask)) != 0
			limit := archive.BooleanVarCount
			if isIntegerVar {
				limit = archive.IntegerVarCount
			}
			values.RenderUnifiedSliderInt(readOnly, valueLabel, unifier,
				func(u values.Unifier) int { return int(u.Unified().(int32)) & 0x1FF },
				func(value int) string {
					if isIntegerVar {
						return fmt.Sprintf("%%d: %s", view.varInfoProvider.IntegerVariable(value).Name)
					}
					return fmt.Sprintf("%%d: %s", view.varInfoProvider.BooleanVariable(value).Name)
				},
				0, limit-1,
				func(newValue int) {
					updater(func(oldValue uint32) uint32 {
						result := oldValue & ^uint32(0x01FF)
						result |= uint32(newValue) & uint32(0x01FF)
						return result
					})
				})

		case !unifier.IsUnique():
			imgui.LabelText(valueLabel, "(multiple)")
		default:
			imgui.LabelText(valueLabel, "")
		}
	}

	simplifier.SetSpecialHandler("VariableKey", addVariableKey)
	simplifier.SetSpecialHandler("VariableCondition", func() {
		addVariableKey()

		comparisons := []string{
			"Var == Val",
			"Var < Val",
			"Var <= Val",
			"Var > Val",
			"Var >= Val",
			"Var != Val",
		}

		comboLabel := key + "-Check###" + fullKey + "-Value"

		values.RenderUnifiedCombo(readOnly, comboLabel, unifier,
			func(u values.Unifier) int {
				key := unifier.Unified().(int32)
				return int((key >> 13) & 0x7)
			},
			func(index int) string { return comparisons[index] },
			len(comparisons),
			func(newIndex int) {
				updater(func(oldValue uint32) uint32 {
					result := oldValue & ^uint32(0xE000)
					result |= uint32(newIndex<<13) & uint32(0xE000)
					return result
				})
			})
	})
	simplifier.SetSpecialHandler("LockVariable", func() {
		limit := archive.BooleanVarCount
		valueLabel := key + "###" + fullKey
		values.RenderUnifiedSliderInt(readOnly, valueLabel, unifier,
			func(u values.Unifier) int { return int(u.Unified().(int32)) & 0x1FF },
			func(value int) string {
				return fmt.Sprintf("%%d: %s", view.varInfoProvider.BooleanVariable(value).Name)
			},
			0, limit-1,
			func(newValue int) {
				updater(func(oldValue uint32) uint32 {
					result := oldValue & ^uint32(0x01FF)
					result |= uint32(newValue) & uint32(0x01FF)
					return result
				})
			})
	})

	simplifier.SetSpecialHandler("BinaryCodedDecimal", func() {
		values.RenderUnifiedSliderInt(readOnly, label, unifier,
			func(u values.Unifier) int { return int(numbers.FromBinaryCodedDecimal(uint16(u.Unified().(int32)))) },
			func(value int) string { return "%03d" },
			0, 999,
			func(newValue int) {
				updater(func(oldValue uint32) uint32 { return uint32(numbers.ToBinaryCodedDecimal(uint16(newValue))) })
			})
	})

	simplifier.SetSpecialHandler("LevelTexture", func() {
		atlas := lvl.TextureAtlas()
		selectedIndex := -1
		if unifier.IsUnique() {
			selectedIndex = int(unifier.Unified().(int32))
		}

		values.RenderUnifiedSliderInt(readOnly, key+" (atlas index)###"+fullKey+"-atlas", unifier,
			func(u values.Unifier) int { return int(u.Unified().(int32)) },
			func(value int) string { return "%d" },
			0, len(atlas)-1,
			func(newValue int) {
				updater(func(oldValue uint32) uint32 { return uint32(newValue) })
			})
		render.TextureSelector(label, -1, view.guiScale, len(atlas), selectedIndex,
			view.textureCache,
			func(index int) resource.Key {
				return resource.KeyOf(ids.LargeTextures.Plus(int(atlas[index])), resource.LangAny, 0)
			},
			func(index int) string { return view.textureName(int(atlas[index])) },
			func(index int) {
				if !readOnly {
					updater(func(oldValue uint32) uint32 { return uint32(index) })
				}
			})
	})
	simplifier.SetSpecialHandler("MaterialOrLevelTexture", func() {
		types := []string{"Material", "Level Texture"}
		selectedType := -1
		selectedIndex := -1
		if unifier.IsUnique() {
			value := int(unifier.Unified().(int32))
			selectedType = (value >> 7) & 1
			selectedIndex = value & 0x7F
		}

		values.RenderUnifiedCombo(readOnly, key+"###"+fullKey+"-Type", unifier,
			func(u values.Unifier) int { return selectedType },
			func(index int) string { return types[index] },
			len(types),
			func(newIndex int) {
				updater(func(oldValue uint32) uint32 {
					newValue := uint32(0)
					if newIndex > 0 {
						newValue = 0x80
					}
					return newValue
				})
			})
		selectorLabel := key + "-Texture###" + fullKey + "-Texture"
		if selectedType == 0 {
			resInfo, _ := ids.Info(ids.ObjectMaterialBitmaps)
			render.TextureSelector(selectorLabel, -1, view.guiScale, resInfo.MaxCount, selectedIndex,
				view.textureCache,
				func(index int) resource.Key {
					return resource.KeyOf(ids.ObjectMaterialBitmaps.Plus(index), resource.LangAny, 0)
				},
				func(index int) string { return fmt.Sprintf("%3d", index) },
				func(index int) {
					if !readOnly {
						updater(func(oldValue uint32) uint32 { return uint32(index) })
					}
				})
		} else {
			atlas := lvl.TextureAtlas()
			render.TextureSelector(selectorLabel, -1, view.guiScale, len(atlas), selectedIndex,
				view.textureCache,
				func(index int) resource.Key {
					return resource.KeyOf(ids.LargeTextures.Plus(int(atlas[index])), resource.LangAny, 0)
				},
				func(index int) string { return view.textureName(int(atlas[index])) },
				func(index int) {
					if !readOnly {
						updater(func(oldValue uint32) uint32 { return uint32(0x80 | index) })
					}
				})
		}
	})

	simplifier.SetSpecialHandler("TileHeight", func() {
		values.RenderUnifiedSliderInt(readOnly, label, unifier,
			func(u values.Unifier) int { return int(u.Unified().(int32)) },
			tileHeightFormatter,
			0, int(level.TileHeightUnitMax-1),
			func(newValue int) {
				updater(func(oldValue uint32) uint32 { return uint32(newValue) })
			})
	})
	simplifier.SetSpecialHandler("ObjectHeight", func() {
		values.RenderUnifiedSliderInt(readOnly, label, unifier,
			func(u values.Unifier) int { return int(u.Unified().(int32)) },
			objectHeightFormatter,
			0, 255,
			func(newValue int) {
				updater(func(oldValue uint32) uint32 { return uint32(newValue) })
			})
	})
	simplifier.SetSpecialHandler("MoveTileHeight", func() {
		values.RenderUnifiedSliderInt(readOnly, label, unifier,
			func(u values.Unifier) int { return int(u.Unified().(int32)) },
			moveTileHeightFormatter,
			0, 0x0FFF,
			func(newValue int) {
				updater(func(oldValue uint32) uint32 { return uint32(newValue) })
			})
	})

	simplifier.SetSpecialHandler("TileType", func() {
		values.RenderUnifiedCombo(readOnly, label, unifier,
			func(u values.Unifier) int { return int(u.Unified().(int32)) },
			func(index int) string { return level.TileType(index).String() },
			len(level.TileTypes()),
			func(newValue int) {
				updater(func(oldValue uint32) uint32 { return uint32(newValue) })
			})
	})

	describer(simplifier)
}

func (view *ObjectsView) renderBlockPuzzleControl(lvl *level.Level, objects []*level.ObjectMainEntry, readOnly bool) {
	var blockPuzzleData *interpreters.Instance

	if len(objects) == 1 {
		interpreter := lvl.ObjectClassData(objects[0])
		isBlockPuzzle := interpreter.Refined("Puzzle").Get("Type") == 0x10
		if isBlockPuzzle {
			blockPuzzleData = interpreter.Refined("Puzzle").Refined("Block")
		}
	}
	if blockPuzzleData != nil {
		blockPuzzleDataID := level.ObjectID(blockPuzzleData.Get("StateStoreObjectID"))
		dataObj := lvl.Object(blockPuzzleDataID)
		blockPuzzleDataState := lvl.ObjectClassData(dataObj).Raw()
		expectedTriple := object.TripleFrom(int(object.ClassTrap), 0, 1)

		imgui.Separator()
		if (dataObj != nil) && (dataObj.InUse != 0) &&
			(dataObj.Triple() == expectedTriple) &&
			(len(blockPuzzleDataState) >= (16 + 6)) {
			blockLayout := blockPuzzleData.Get("Layout")
			blockWidth := int((blockLayout >> 20) & 7)
			blockHeight := int((blockLayout >> 24) & 7)
			state := level.NewBlockPuzzleState(blockPuzzleDataState[6:6+16], blockHeight, blockWidth)
			startRow := 1 + (7-blockHeight)/2
			startColumn := 1 + (7-blockWidth)/2
			var cells [9][9]string

			placeConnector := func(side, offset int, text string) {
				xOffsets := []int{offset, offset, -1, blockWidth}
				yOffsets := []int{-1, blockHeight, offset, offset}
				y := startRow + yOffsets[side]
				x := startColumn + xOffsets[side]

				if (x >= 0) && (x < 9) && (y >= 0) && (y < 9) {
					cells[y][x] = text
				}
			}
			stateMapping := []string{" ", "X", "+", "(+)", "F", "(F)", "H", "(H)"}

			for row := 0; row < blockHeight; row++ {
				for col := 0; col < blockWidth; col++ {
					value := state.CellValue(row, col)
					cells[startRow+row][startColumn+col] = stateMapping[value]
				}
			}
			placeConnector(int((blockLayout>>7)&3), int((blockLayout>>4)&7), "S")
			placeConnector(int((blockLayout>>15)&3), int((blockLayout>>12)&7), "D")

			cellSize := imgui.Vec2{X: imgui.TextLineHeightWithSpacing() * 1.5, Y: imgui.TextLineHeightWithSpacing() * 1.5}
			for x := 0; x < 9; x++ {
				imgui.BeginGroup()
				for y := 0; y < 9; y++ {
					cellText := cells[y][x]
					imgui.PushID(fmt.Sprintf("%d:%d", x, y))
					if len(cellText) > 0 {
						if imgui.ButtonV(cellText, cellSize) && !readOnly {
							clickRow := y - startRow
							clickCol := x - startColumn
							if (clickRow >= 0) && (clickRow < blockHeight) && (clickCol >= 0) && (clickCol < blockWidth) {
								oldValue := state.CellValue(clickRow, clickCol)
								state.SetCellValue(clickRow, clickCol, (8+oldValue+1)%8)

								// The following seems to do nothing, yet the data is still synchronized.
								// This appears to be a bad design - maybe find a better way to design the concept
								// of block puzzle references...
								view.requestAction("ChangeBlockPuzzle", func() error {
									return view.editor.ChangeObjects(func(entry *level.ObjectMainEntry) {})
								})
							}
						}
					} else {
						imgui.Dummy(cellSize)
					}
					imgui.PopID()
				}
				imgui.EndGroup()
				imgui.SameLine()
			}
		} else {
			imgui.Text("No proper state store found for block puzzle!\n'StateStoreObjectID' must refer to a " + expectedTriple.String() + " object.")
		}
	}
}

func (view *ObjectsView) requestChangeObjects(name string, modifier func(*level.ObjectMainEntry)) {
	view.requestAction(name, func() error {
		return view.editor.ChangeObjects(modifier)
	})
}

func (view *ObjectsView) requestPropertiesChange(interpreterRetriever func(*level.ObjectMainEntry) *interpreters.Instance,
	key string, modifier func(uint32) uint32) {
	subKeys := strings.Split(key, ".")
	valueIndex := len(subKeys) - 1

	view.requestAction("ChangeProperties", func() error {
		return view.editor.ChangeObjects(func(obj *level.ObjectMainEntry) {
			interpreter := interpreterRetriever(obj)
			for subIndex := 0; subIndex < valueIndex; subIndex++ {
				interpreter = interpreter.Refined(subKeys[subIndex])
			}
			subKey := subKeys[valueIndex]
			interpreter.Set(subKey, modifier(interpreter.Get(subKey)))
		})
	})
}

func (view *ObjectsView) requestDeleteObjects() {
	view.requestAction("DeleteObjects", view.editor.DeleteObjects)
}

func (view *ObjectsView) requestAction(name string, nested func() error) {
	err := view.registry.Register(cmd.Named(name),
		cmd.Forward(view.restoreFocusTask()),
		cmd.Nested(nested),
		cmd.Reverse(view.restoreFocusTask()))
	if err != nil {
		panic(err)
	}
}

func (view *ObjectsView) restoreFocusTask() cmd.Task {
	return func(modder world.Modder) error {
		view.model.restoreFocus = true
		return nil
	}
}

func (view *ObjectsView) tripleName(triple object.Triple) string {
	suffix := hintUnknown
	linearIndex := view.gameObjects.AllProperties().TripleIndex(triple)
	if linearIndex >= 0 {
		key := resource.KeyOf(ids.ObjectLongNames, resource.LangDefault, linearIndex)
		objName, err := view.textCache.Text(key)
		if err == nil {
			suffix = objName
		}
	}
	return triple.String() + ": " + suffix
}

func (view *ObjectsView) textureName(index int) string {
	key := resource.KeyOf(ids.TextureNames, resource.LangDefault, index)
	name, err := view.textCache.Text(key)
	suffix := ""
	if err == nil {
		suffix = ": " + name
	}
	return fmt.Sprintf("%3d", index) + suffix
}
