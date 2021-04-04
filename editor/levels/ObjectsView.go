package levels

import (
	"fmt"
	"strings"

	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlobj"
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

// ObjectsView is for object properties.
type ObjectsView struct {
	levels          *edit.EditableLevels
	levelSelection  *edit.LevelSelectionService
	editor          *edit.LevelEditorService
	varInfoProvider archive.GameVariableInfoProvider
	mod             *world.Mod
	textCache       *text.Cache
	textureCache    *graphics.TextureCache

	guiScale float32
	registry cmd.Registry

	model objectsViewModel
}

// NewObjectsView returns a new instance.
func NewObjectsView(editor *edit.LevelEditorService, levels *edit.EditableLevels, levelSelection *edit.LevelSelectionService,
	varInfoProvider archive.GameVariableInfoProvider,
	mod *world.Mod, guiScale float32, textCache *text.Cache, textureCache *graphics.TextureCache,
	registry cmd.Registry) *ObjectsView {
	view := &ObjectsView{
		levelSelection:  levelSelection,
		levels:          levels,
		editor:          editor,
		varInfoProvider: varInfoProvider,
		mod:             mod,
		textCache:       textCache,
		textureCache:    textureCache,

		guiScale: guiScale,
		registry: registry,

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
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 500 * view.guiScale}, imgui.ConditionFirstUseEver)
		title := fmt.Sprintf("Level Objects, %d selected", view.levelSelection.NumberOfSelectedObjects())
		readOnly := view.editor.IsReadOnly()
		if readOnly {
			title += hintReadOnly
		}
		if imgui.BeginV(title+"###Level Objects", view.WindowOpen(), imgui.WindowFlagsHorizontalScrollbar|imgui.WindowFlagsAlwaysVerticalScrollbar) {
			view.renderContent(lvl, readOnly)
		}
		imgui.End()
	}
}

func (view *ObjectsView) renderContent(lvl *level.Level, readOnly bool) {
	objectIDUnifier := values.NewUnifier()
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

	selectedObjectIDs := view.levelSelection.CurrentSelectedObjects()
	for _, id := range selectedObjectIDs {
		obj := lvl.Object(id)
		if obj != nil {
			objectIDUnifier.Add(id)
			classUnifier.Add(object.TripleFrom(int(obj.Class), 0, 0))
			typeUnifier.Add(obj.Triple())
			zUnifier.Add(obj.Z)
			tileXUnifier.Add(obj.X.Tile())
			fineXUnifier.Add(obj.X.Fine())
			tileYUnifier.Add(obj.Y.Tile())
			fineYUnifier.Add(obj.Y.Fine())
			rotationXUnifier.Add(obj.XRotation)
			rotationYUnifier.Add(obj.YRotation)
			rotationZUnifier.Add(obj.ZRotation)
			hitpointsUnifier.Add(obj.Hitpoints)
		}
	}

	imgui.PushItemWidth(-150 * view.guiScale)
	columns, rows, levelHeight := lvl.Size()

	objectHeightFormatter := objectHeightFormatterFor(levelHeight)
	rotationFormatter := func(value int) string {
		return fmt.Sprintf("%.3f degrees  - raw: %%d", level.RotationUnit(value).ToDegrees())
	}

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
		if imgui.BeginCombo("New Object Class", classString(view.model.newObjectTriple.Class)) {
			for _, class := range object.Classes() {
				if imgui.SelectableV(classString(class), class == view.model.newObjectTriple.Class, 0, imgui.Vec2{}) {
					view.model.newObjectTriple = object.TripleFrom(int(class), 0, 0)
				}
			}
			imgui.EndCombo()
		}
		if imgui.BeginCombo("New Object Type", view.tripleName(view.model.newObjectTriple)) {
			allTypes := view.mod.ObjectProperties().TriplesInClass(view.model.newObjectTriple.Class)
			for _, triple := range allTypes {
				if imgui.SelectableV(view.tripleName(triple), triple == view.model.newObjectTriple, 0, imgui.Vec2{}) {
					view.model.newObjectTriple = triple
				}
			}
			imgui.EndCombo()
		}
		if imgui.Button("Delete Selected") {
			view.requestDeleteObjects()
		}
		imgui.Separator()
	}

	switch {
	case objectIDUnifier.IsUnique():
		imgui.LabelText("ID", fmt.Sprintf("%3d", int(objectIDUnifier.Unified().(level.ObjectID))))
	case !objectIDUnifier.IsEmpty():
		imgui.LabelText("ID", "(multiple)")
	default:
		imgui.LabelText("ID", "")
	}

	objTypeRenderer := values.ObjectTypeControlRenderer{
		Meta:      view.mod.ObjectProperties(),
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
			0, 0xFF,
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
			0, 0xFF,
			func(newValue int) {
				view.requestChangeObjects("ChangeFineY", func(entry *level.ObjectMainEntry) { entry.Y = level.CoordinateAt(entry.Y.Tile(), byte(newValue)) })
			})
		values.RenderUnifiedSliderInt(readOnly, "Rotation X", rotationXUnifier,
			func(u values.Unifier) int { return int(u.Unified().(level.RotationUnit)) },
			rotationFormatter,
			0, 0xFF,
			func(newValue int) {
				view.requestChangeObjects("ChangeBaseRotationX", func(entry *level.ObjectMainEntry) { entry.XRotation = level.RotationUnit(newValue) })
			})
		values.RenderUnifiedSliderInt(readOnly, "Rotation Y", rotationYUnifier,
			func(u values.Unifier) int { return int(u.Unified().(level.RotationUnit)) },
			rotationFormatter,
			0, 0xFF,
			func(newValue int) {
				view.requestChangeObjects("ChangeBaseRotationY", func(entry *level.ObjectMainEntry) { entry.YRotation = level.RotationUnit(newValue) })
			})
		values.RenderUnifiedSliderInt(readOnly, "Rotation Z", rotationZUnifier,
			func(u values.Unifier) int { return int(u.Unified().(level.RotationUnit)) },
			rotationFormatter,
			0, 0xFF,
			func(newValue int) {
				view.requestChangeObjects("ChangeBaseRotationZ", func(entry *level.ObjectMainEntry) { entry.ZRotation = level.RotationUnit(newValue) })
			})

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
		view.renderProperties(lvl, readOnly,
			func(id level.ObjectID, entry *level.ObjectMainEntry) []byte { return entry.Extra[:] },
			view.extraInterpreterFactory(lvl))
		imgui.TreePop()
	}
	if imgui.TreeNodeV("Class Properties", imgui.TreeNodeFlagsFramed) {
		view.renderProperties(lvl, readOnly,
			func(id level.ObjectID, entry *level.ObjectMainEntry) []byte { return lvl.ObjectClassData(id) },
			view.classInterpreterFactory(lvl))
		view.renderBlockPuzzleControl(lvl, readOnly)
		imgui.TreePop()
	}

	imgui.PopItemWidth()
}

func (view *ObjectsView) renderProperties(lvl *level.Level, readOnly bool,
	dataRetriever func(level.ObjectID, *level.ObjectMainEntry) []byte,
	interpreterFactory lvlobj.InterpreterFactory) {
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

	for index, id := range view.levelSelection.CurrentSelectedObjects() {
		obj := lvl.Object(id)
		if obj != nil {
			data := dataRetriever(id, obj)
			interpreter := interpreterFactory(obj.Triple(), data)

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
					view.requestPropertiesChange(lvl, dataRetriever, interpreterFactory, key, modifier) // nolint: scopelint
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
		Meta:      view.mod.ObjectProperties(),
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

func (view *ObjectsView) renderBlockPuzzleControl(lvl *level.Level, readOnly bool) {
	var blockPuzzleData *interpreters.Instance

	selectedObjectIDs := view.levelSelection.CurrentSelectedObjects()
	if len(selectedObjectIDs) == 1 {
		id := selectedObjectIDs[0]
		obj := lvl.Object(id)
		if (obj != nil) && (obj.InUse != 0) {
			triple := obj.Triple()
			classData := lvl.ObjectClassData(id)
			interpreter := view.classInterpreterFactory(lvl)(triple, classData)
			isBlockPuzzle := interpreter.Refined("Puzzle").Get("Type") == 0x10
			if isBlockPuzzle {
				blockPuzzleData = interpreter.Refined("Puzzle").Refined("Block")
			}
		}
	}
	if blockPuzzleData != nil {
		blockPuzzleDataID := level.ObjectID(blockPuzzleData.Get("StateStoreObjectID"))
		blockPuzzleDataState := lvl.ObjectClassData(blockPuzzleDataID)
		dataObj := lvl.Object(blockPuzzleDataID)
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
								view.patchLevel(lvl, view.levelSelection.CurrentSelectedObjects())
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

func (view *ObjectsView) extraInterpreterFactory(lvl *level.Level) lvlobj.InterpreterFactory {
	interpreterFactory := lvlobj.RealWorldExtra
	if lvl.IsCyberspace() {
		interpreterFactory = lvlobj.CyberspaceExtra
	}
	return interpreterFactory
}

func (view *ObjectsView) classInterpreterFactory(lvl *level.Level) lvlobj.InterpreterFactory {
	interpreterFactory := lvlobj.ForRealWorld
	if lvl.IsCyberspace() {
		interpreterFactory = lvlobj.ForCyberspace
	}
	return interpreterFactory
}

func (view *ObjectsView) requestPropertiesChange(lvl *level.Level,
	dataRetriever func(level.ObjectID, *level.ObjectMainEntry) []byte,
	interpreterFactory lvlobj.InterpreterFactory,
	key string, modifier func(uint32) uint32) {
	objectIDs := view.levelSelection.CurrentSelectedObjects()
	subKeys := strings.Split(key, ".")
	valueIndex := len(subKeys) - 1

	for _, id := range objectIDs {
		obj := lvl.Object(id)
		if obj != nil {
			data := dataRetriever(id, obj)
			interpreter := interpreterFactory(obj.Triple(), data)
			for subIndex := 0; subIndex < valueIndex; subIndex++ {
				interpreter = interpreter.Refined(subKeys[subIndex])
			}
			subKey := subKeys[valueIndex]
			interpreter.Set(subKey, modifier(interpreter.Get(subKey)))
		}
	}

	view.patchLevel(lvl, objectIDs)
}

// RequestCreateObject requests to create a new object of the currently selected type.
func (view *ObjectsView) RequestCreateObject(pos MapPosition) {
	if view.editor.IsReadOnly() {
		return
	}
	view.requestCreateObject(view.model.newObjectTriple, pos)
}

func (view *ObjectsView) requestCreateObject(triple object.Triple, pos MapPosition) {
	lvl := view.editor.Level()
	if !lvl.HasRoomForObjectOf(triple.Class) {
		return
	}
	placeObject := func(obj *level.ObjectMainEntry) {
		obj.X = pos.X
		obj.Y = pos.Y
	}
	view.requestAction("NewObject", func() error { return view.editor.NewObject(triple, placeObject) })
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

func (view *ObjectsView) patchLevel(lvl *level.Level, forwardObjectIDs []level.ObjectID) {
	oldLevelID := view.levelSelection.CurrentLevelID()
	reverseObjectIDs := view.levelSelection.CurrentSelectedObjects()

	err := view.registry.Register(cmd.Named("PatchLevel"),
		cmd.Forward(view.restoreFocusTask()),
		cmd.Forward(view.levelSelection.SetCurrentLevelIDTask(lvl.ID())),
		cmd.Reverse(view.setSelectedObjectsTask(reverseObjectIDs)),
		cmd.Nested(func() error { return view.levels.CommitLevelChanges(lvl.ID()) }),
		cmd.Forward(view.setSelectedObjectsTask(forwardObjectIDs)),
		cmd.Reverse(view.levelSelection.SetCurrentLevelIDTask(oldLevelID)),
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

func (view *ObjectsView) setSelectedObjectsTask(ids []level.ObjectID) cmd.Task {
	return func(modder world.Modder) error {
		view.levelSelection.SetCurrentSelectedObjects(ids)
		return nil
	}
}

func (view *ObjectsView) tripleName(triple object.Triple) string {
	suffix := hintUnknown
	linearIndex := view.mod.ObjectProperties().TripleIndex(triple)
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

// PlaceSelectedObjectsOnFloor puts all selected objects to sit on the floor.
func (view *ObjectsView) PlaceSelectedObjectsOnFloor() {
	lvl := view.editor.Level()
	_, _, height := lvl.Size()
	view.placeSelectedObjects(lvl, func(tile *level.TileMapEntry, pos level.FinePosition, objPivot float32) level.HeightUnit {
		floorHeight := tile.FloorTileHeightAt(pos, height)
		return height.ValueToObjectHeight(floorHeight + objPivot)
	})
}

// PlaceSelectedObjectsOnEyeLevel puts all selected objects to be at eye level (approximately).
func (view *ObjectsView) PlaceSelectedObjectsOnEyeLevel() {
	lvl := view.editor.Level()
	_, _, height := lvl.Size()
	view.placeSelectedObjects(lvl, func(tile *level.TileMapEntry, pos level.FinePosition, objPivot float32) level.HeightUnit {
		floorHeight := tile.FloorTileHeightAt(pos, height)
		return height.ValueToObjectHeight(floorHeight + 0.75 - objPivot)
	})
}

// PlaceSelectedObjectsOnCeiling puts all selected objects to hang from the ceiling.
func (view *ObjectsView) PlaceSelectedObjectsOnCeiling() {
	lvl := view.editor.Level()
	_, _, height := lvl.Size()
	view.placeSelectedObjects(lvl, func(tile *level.TileMapEntry, pos level.FinePosition, objPivot float32) level.HeightUnit {
		ceilingHeight := tile.CeilingTileHeightAt(pos, height)
		return height.ValueToObjectHeight(ceilingHeight - objPivot)
	})
}

func (view *ObjectsView) placeSelectedObjects(lvl *level.Level,
	atHeight func(tile *level.TileMapEntry, pos level.FinePosition, objPivot float32) level.HeightUnit) {
	view.requestChangeObjects("ChangeBasePlacement", func(obj *level.ObjectMainEntry) {
		var objPivot float32
		prop, err := view.mod.ObjectProperties().ForObject(obj.Triple())
		if err == nil {
			objPivot = object.Pivot(prop.Common)
		}
		tilePos := obj.TilePosition()
		tile := lvl.Tile(tilePos)
		if tile != nil {
			obj.Z = atHeight(tile, obj.FinePosition(), objPivot)
		}
	})
}
