package levels

import (
	"fmt"
	"sort"
	"strings"

	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/event"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlobj"
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/imgui-go"
)

// ObjectsView is for object properties.
type ObjectsView struct {
	mod       *model.Mod
	textCache *text.Cache

	guiScale      float32
	commander     cmd.Commander
	eventListener event.Listener

	model objectsViewModel
}

// NewObjectsView returns a new instance.
func NewObjectsView(mod *model.Mod, guiScale float32, textCache *text.Cache,
	commander cmd.Commander, eventListener event.Listener, eventRegistry event.Registry) *ObjectsView {
	view := &ObjectsView{
		mod:       mod,
		textCache: textCache,

		guiScale:      guiScale,
		commander:     commander,
		eventListener: eventListener,

		model: freshObjectsViewModel(),
	}
	view.model.selectedObjects.registerAt(eventRegistry)
	return view
}

// WindowOpen returns the flag address, to be used with the main menu.
func (view *ObjectsView) WindowOpen() *bool {
	return &view.model.windowOpen
}

// Render renders the view.
func (view *ObjectsView) Render(lvl *level.Level) {
	if view.model.restoreFocus {
		imgui.SetNextWindowFocus()
		view.model.restoreFocus = false
		view.model.windowOpen = true
	}
	if view.model.windowOpen {
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 500 * view.guiScale}, imgui.ConditionOnce)
		title := fmt.Sprintf("Level Objects, %d selected", len(view.model.selectedObjects.list))
		readOnly := !view.editingAllowed(lvl.ID())
		if readOnly {
			title += " (read-only)"
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

	for _, id := range view.model.selectedObjects.list {
		obj := lvl.Object(id)
		if obj != nil {
			objectIDUnifier.Add(id)
			classUnifier.Add(object.TripleFrom(int(obj.Class), 0, 0))
			typeUnifier.Add(object.TripleFrom(int(obj.Class), int(obj.Subclass), int(obj.Type)))
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
	multiple := len(view.model.selectedObjects.list) > 1
	columns, rows, levelHeight := lvl.Size()

	objectHeightFormatter := func(value int) string {
		tileHeight, err := levelHeight.ValueFromObjectHeight(level.HeightUnit(value))
		tileHeightString := "???"
		if err == nil {
			tileHeightString = fmt.Sprintf("%2.3f", tileHeight)
		}
		return tileHeightString + " tile(s) - raw: %d"
	}
	rotationFormatter := func(value int) string {
		return fmt.Sprintf("%.3f degrees  - raw: %%d", level.RotationUnit(value).ToDegrees())
	}

	if multiple {
		imgui.LabelText("ID", "(multiple)")
	} else if objectIDUnifier.IsUnique() {
		imgui.LabelText("ID", fmt.Sprintf("%3d", int(objectIDUnifier.Unified().(level.ObjectID))))
	} else {
		imgui.LabelText("ID", "")
	}
	view.renderTypeCombo(readOnly, multiple, "Object Type", classUnifier, typeUnifier,
		func(u values.Unifier) object.Triple { return u.Unified().(object.Triple) },
		func(newValue object.Triple) {
			view.requestBaseChange(lvl, func(entry *level.ObjectMasterEntry) {
				entry.Subclass = newValue.Subclass
				entry.Type = newValue.Type
			})
		})

	if imgui.TreeNodeV("Base Properties", imgui.TreeNodeFlagsDefaultOpen|imgui.TreeNodeFlagsFramed) {

		values.RenderUnifiedSliderInt(readOnly, multiple, "Z", zUnifier,
			func(u values.Unifier) int { return int(u.Unified().(level.HeightUnit)) },
			objectHeightFormatter,
			0, 0xFF,
			func(newValue int) {
				view.requestBaseChange(lvl, func(entry *level.ObjectMasterEntry) { entry.Z = level.HeightUnit(newValue) })
			})
		values.RenderUnifiedSliderInt(readOnly, multiple, "Tile X", tileXUnifier,
			func(u values.Unifier) int { return int(u.Unified().(byte)) },
			func(value int) string { return "%d" },
			0, columns-1,
			func(newValue int) {
				view.requestBaseChange(lvl, func(entry *level.ObjectMasterEntry) { entry.X = level.CoordinateAt(byte(newValue), entry.X.Fine()) })
			})
		values.RenderUnifiedSliderInt(readOnly, multiple, "Fine X", fineXUnifier,
			func(u values.Unifier) int { return int(u.Unified().(byte)) },
			func(value int) string { return "%d" },
			0, 0xFF,
			func(newValue int) {
				view.requestBaseChange(lvl, func(entry *level.ObjectMasterEntry) { entry.X = level.CoordinateAt(entry.X.Tile(), byte(newValue)) })
			})
		values.RenderUnifiedSliderInt(readOnly, multiple, "Tile Y", tileYUnifier,
			func(u values.Unifier) int { return int(u.Unified().(byte)) },
			func(value int) string { return "%d" },
			0, rows-1,
			func(newValue int) {
				view.requestBaseChange(lvl, func(entry *level.ObjectMasterEntry) { entry.Y = level.CoordinateAt(byte(newValue), entry.Y.Fine()) })
			})
		values.RenderUnifiedSliderInt(readOnly, multiple, "Fine Y", fineYUnifier,
			func(u values.Unifier) int { return int(u.Unified().(byte)) },
			func(value int) string { return "%d" },
			0, 0xFF,
			func(newValue int) {
				view.requestBaseChange(lvl, func(entry *level.ObjectMasterEntry) { entry.Y = level.CoordinateAt(entry.Y.Tile(), byte(newValue)) })
			})
		values.RenderUnifiedSliderInt(readOnly, multiple, "Rotation X", rotationXUnifier,
			func(u values.Unifier) int { return int(u.Unified().(level.RotationUnit)) },
			rotationFormatter,
			0, 0xFF,
			func(newValue int) {
				view.requestBaseChange(lvl, func(entry *level.ObjectMasterEntry) { entry.XRotation = level.RotationUnit(newValue) })
			})
		values.RenderUnifiedSliderInt(readOnly, multiple, "Rotation Y", rotationYUnifier,
			func(u values.Unifier) int { return int(u.Unified().(level.RotationUnit)) },
			rotationFormatter,
			0, 0xFF,
			func(newValue int) {
				view.requestBaseChange(lvl, func(entry *level.ObjectMasterEntry) { entry.YRotation = level.RotationUnit(newValue) })
			})
		values.RenderUnifiedSliderInt(readOnly, multiple, "Rotation Z", rotationZUnifier,
			func(u values.Unifier) int { return int(u.Unified().(level.RotationUnit)) },
			rotationFormatter,
			0, 0xFF,
			func(newValue int) {
				view.requestBaseChange(lvl, func(entry *level.ObjectMasterEntry) { entry.ZRotation = level.RotationUnit(newValue) })
			})

		values.RenderUnifiedSliderInt(readOnly, multiple, "Hitpoints", hitpointsUnifier,
			func(u values.Unifier) int { return int(u.Unified().(int16)) },
			func(value int) string { return "%d" },
			0, 10000,
			func(newValue int) {
				view.requestBaseChange(lvl, func(entry *level.ObjectMasterEntry) { entry.Hitpoints = int16(newValue) })
			})

		imgui.TreePop()
	}
	if imgui.TreeNodeV("Extra Properties", imgui.TreeNodeFlagsFramed) {
		view.renderProperties(lvl, readOnly,
			func(id level.ObjectID, entry *level.ObjectMasterEntry) []byte { return entry.Extra[:] },
			view.extraInterpreterFactory(lvl))
		imgui.TreePop()
	}
	if imgui.TreeNodeV("Class Properties", imgui.TreeNodeFlagsFramed) {
		view.renderProperties(lvl, readOnly,
			func(id level.ObjectID, entry *level.ObjectMasterEntry) []byte { return lvl.ObjectClassData(id) },
			view.classInterpreterFactory(lvl))
		imgui.TreePop()
	}

	imgui.PopItemWidth()
}

func (view *ObjectsView) renderProperties(lvl *level.Level, readOnly bool,
	dataRetriever func(level.ObjectID, *level.ObjectMasterEntry) []byte,
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

	for index, id := range view.model.selectedObjects.list {
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

	multiple := len(view.model.selectedObjects.list) > 1
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
			view.renderPropertyControl(lvl, readOnly, multiple, key, *unifier, propertyDescribers[key],
				func(modifier func(uint32) uint32) {
					view.requestPropertiesChange(lvl, dataRetriever, interpreterFactory, key, modifier)
				})
		}
	}
}

func (view *ObjectsView) renderPropertyControl(lvl *level.Level, readOnly bool, multiple bool,
	fullKey string, unifier values.Unifier, describer func(*interpreters.Simplifier),
	updater func(func(uint32) uint32)) {
	keys := strings.Split(fullKey, ".")
	key := keys[len(keys)-1]
	label := key + "###" + fullKey

	simplifier := interpreters.NewSimplifier(func(minValue, maxValue int64, formatter interpreters.RawValueFormatter) {
		values.RenderUnifiedSliderInt(readOnly, multiple, label, unifier,
			func(u values.Unifier) int { return int(u.Unified().(int32)) },
			func(value int) string { return "%d" },
			int(minValue), int(maxValue),
			func(newValue int) {
				updater(func(oldValue uint32) uint32 { return uint32(newValue) })
			})
	})

	simplifier.SetEnumValueHandler(func(enumValues map[uint32]string) {
		valueKeys := make([]uint32, 0, len(enumValues))
		for valueKey := range enumValues {
			valueKeys = append(valueKeys, valueKey)
		}
		sort.Slice(valueKeys, func(indexA, indexB int) bool { return valueKeys[indexA] < valueKeys[indexB] })

		values.RenderUnifiedCombo(readOnly, multiple, label, unifier,
			func(u values.Unifier) int {
				unifiedValue := uint32(u.Unified().(int32))
				for index, valueKey := range valueKeys {
					if valueKey == unifiedValue {
						return index
					}
				}
				return -1
			},
			func(index int) string { return enumValues[valueKeys[index]] },
			len(valueKeys),
			func(newIndex int) {
				updater(func(oldValue uint32) uint32 { return valueKeys[newIndex] })
			})
	})

	simplifier.SetBitfieldHandler(func(maskNames map[uint32]string) {
		masks := make([]uint32, 0, len(maskNames))
		for mask := range maskNames {
			masks = append(masks, mask)
		}
		sort.Slice(masks, func(indexA, indexB int) bool { return masks[indexA] < masks[indexB] })

		addMaskedItem := func(mask uint32) {
			maxValue := mask
			shift := 0
			maskedLabel := key + "." + maskNames[mask] + "###" + label + "-" + maskNames[mask]

			for (maxValue & 1) == 0 {
				shift++
				maxValue >>= 1
			}

			if maxValue == 1 {
				booleanUnifier := values.NewUnifier()
				if unifier.IsUnique() {
					booleanUnifier.Add((uint32(unifier.Unified().(int32)) & mask) != 0)
				}
				values.RenderUnifiedCheckboxCombo(readOnly, multiple, maskedLabel, booleanUnifier,
					func(newValue bool) {
						updater(func(oldValue uint32) uint32 {
							result := oldValue & ^mask
							if newValue {
								result |= mask
							}
							return result
						})
					})
			} else {
				values.RenderUnifiedSliderInt(readOnly, multiple, maskedLabel, unifier,
					func(u values.Unifier) int { return int((uint32(u.Unified().(int32)) & mask) >> uint32(shift)) },
					func(value int) string { return "%d" },
					0, int(maxValue),
					func(newValue int) {
						updater(func(oldValue uint32) uint32 {
							return (oldValue & ^mask) | (uint32(newValue) << uint32(shift))
						})
					})
			}
		}

		for _, mask := range masks {
			addMaskedItem(mask)
		}
	})

	simplifier.SetSpecialHandler("ObjectTriple", func() {
		var classNames [object.ClassCount]string
		for index, class := range object.Classes() {
			classNames[index] = class.String()
		}
		tripleResolver := func(u values.Unifier) object.Triple { return object.TripleFromInt(int(u.Unified().(int32))) }
		values.RenderUnifiedCombo(readOnly, multiple, key+"-Class###"+fullKey+"-Class", unifier,
			func(u values.Unifier) int {
				triple := tripleResolver(u)
				return int(triple.Class)
			},
			func(value int) string { return fmt.Sprintf("%2d: %v", value, object.Class(value)) },
			object.ClassCount,
			func(newValue int) {
				triple := object.TripleFrom(newValue, 0, 0)
				updater(func(oldValue uint32) uint32 { return uint32(triple.Int()) })
			})

		view.renderTypeCombo(readOnly, multiple, key+"###"+fullKey+"-Type", unifier, unifier,
			tripleResolver,
			func(newValue object.Triple) {
				updater(func(oldValue uint32) uint32 { return uint32(newValue.Int()) })
			})
	})

	simplifier.SetObjectIDHandler(func() {
		values.RenderUnifiedSliderInt(readOnly, multiple, label, unifier,
			func(u values.Unifier) int { return int(u.Unified().(int32)) },
			func(value int) string { return "%d" },
			0, int(lvl.ObjectLimit()),
			func(newValue int) {
				updater(func(oldValue uint32) uint32 { return uint32(newValue) })
			})
	})

	simplifier.SetSpecialHandler("Mistake", func() {})
	simplifier.SetSpecialHandler("Ignored", func() {})
	simplifier.SetSpecialHandler("Unknown", func() {})

	describer(simplifier)
}

func (view *ObjectsView) renderTypeCombo(readOnly, multiple bool, label string,
	classUnifier values.Unifier, typeUnifier values.Unifier,
	tripleResolver func(values.Unifier) object.Triple,
	changeHandler func(object.Triple)) {
	if classUnifier.IsUnique() {
		objectProperties := view.mod.ObjectProperties()
		class := tripleResolver(classUnifier).Class
		triples := objectProperties.TriplesInClass(class)
		selectedIndex := -1
		if typeUnifier.IsUnique() {
			triple := tripleResolver(typeUnifier)
			for index, availableTriple := range triples {
				if availableTriple == triple {
					selectedIndex = index
				}
			}
			if selectedIndex < 0 {
				selectedIndex = len(triples)
				triples = append(triples, triple)
			}
		}
		values.RenderUnifiedCombo(readOnly, multiple, label, typeUnifier,
			func(u values.Unifier) int { return selectedIndex },
			func(value int) string {
				triple := triples[value]
				suffix := "???"
				linearIndex := objectProperties.TripleIndex(triple)
				if linearIndex >= 0 {
					key := resource.KeyOf(ids.ObjectLongNames, resource.LangDefault, linearIndex)
					objName, err := view.textCache.Text(key)
					if err == nil {
						suffix = objName
					}
				}
				return triple.String() + ": " + suffix
			},
			len(triples),
			func(newValue int) {
				triple := triples[newValue]
				changeHandler(triple)
			})
	} else if multiple {
		imgui.LabelText(label, "(multiple classes)")
	} else {
		imgui.LabelText(label, "")
	}
}

func (view *ObjectsView) editingAllowed(id int) bool {
	gameStateData := view.mod.ModifiedBlocks(resource.LangAny, ids.GameState)
	isSavegame := (len(gameStateData) == 1) && (len(gameStateData[0]) == archive.GameStateSize) && (gameStateData[0][0x009C] > 0)
	moddedLevel := len(view.mod.ModifiedBlocks(resource.LangAny, ids.LevelResourcesStart.Plus(lvlids.PerLevel*id+lvlids.FirstUsed))) > 0

	return moddedLevel && !isSavegame
}

func (view *ObjectsView) requestBaseChange(lvl *level.Level, modifier func(*level.ObjectMasterEntry)) {
	objectIDs := view.model.selectedObjects.list
	for _, id := range objectIDs {
		obj := lvl.Object(id)
		if obj != nil {
			modifier(obj)
		}
	}

	view.patchLevel(lvl, objectIDs)
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
	dataRetriever func(level.ObjectID, *level.ObjectMasterEntry) []byte,
	interpreterFactory lvlobj.InterpreterFactory,
	key string, modifier func(uint32) uint32) {
	objectIDs := view.model.selectedObjects.list
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

func (view *ObjectsView) patchLevel(lvl *level.Level, objectIDs []level.ObjectID) {

	command := patchLevelDataCommand{
		restoreState: func() {
			view.model.restoreFocus = true
			view.setSelectedLevel(lvl.ID())
			view.setSelectedObjects(objectIDs)
		},
	}

	newDataSet := lvl.EncodeState()
	for id, newData := range newDataSet {
		if len(newData) > 0 {
			resourceID := ids.LevelResourcesStart.Plus(lvlids.PerLevel*lvl.ID() + id)
			patch, changed, err := view.mod.CreateBlockPatch(resource.LangAny, resourceID, 0, newData)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				// TODO how to handle this? We're not expecting this, so crash and burn?
			} else if changed {
				command.patches = append(command.patches, patch)
			}
		}
	}

	view.commander.Queue(command)
}

func (view *ObjectsView) setSelectedLevel(id int) {
	view.eventListener.Event(LevelSelectionSetEvent{id: id})
}

func (view *ObjectsView) setSelectedObjects(objectIDs []level.ObjectID) {
	view.eventListener.Event(ObjectSelectionSetEvent{objects: objectIDs})
}
