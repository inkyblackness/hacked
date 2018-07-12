package levels

import (
	"fmt"

	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/event"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/imgui-go"
)

// ObjectsView is for object properties.
type ObjectsView struct {
	mod *model.Mod

	guiScale      float32
	commander     cmd.Commander
	eventListener event.Listener

	model objectsViewModel
}

// NewObjectsView returns a new instance.
func NewObjectsView(mod *model.Mod, guiScale float32, commander cmd.Commander, eventListener event.Listener, eventRegistry event.Registry) *ObjectsView {
	view := &ObjectsView{
		mod:           mod,
		guiScale:      guiScale,
		commander:     commander,
		eventListener: eventListener,
		model:         freshObjectsViewModel(),
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
		if imgui.BeginV(title+"###Level Objects", view.WindowOpen(), 0) {
			view.renderContent(lvl, readOnly)
		}
		imgui.End()
	}
}

func (view *ObjectsView) renderContent(lvl *level.Level, readOnly bool) {
	objectIDUnifier := values.NewUnifier()
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
		objectIDUnifier.Add(id)
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

	imgui.PushItemWidth(-250 * view.guiScale)
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
		return fmt.Sprintf("%.3f degrees  - raw: %v", level.RotationUnit(value).ToDegrees(), value)
	}

	if multiple {
		imgui.LabelText("ID", "(multiple)")
	} else if objectIDUnifier.IsUnique() {
		imgui.LabelText("ID", fmt.Sprintf("%3d", int(objectIDUnifier.Unified().(level.ObjectID))))
	} else {
		imgui.LabelText("ID", "")
	}

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

		imgui.TreePop()
	}
	if imgui.TreeNodeV("Class Properties", imgui.TreeNodeFlagsFramed) {
		imgui.TreePop()
	}

	imgui.PopItemWidth()
}

func (view *ObjectsView) editingAllowed(id int) bool {
	gameStateData := view.mod.ModifiedBlocks(resource.LangAny, ids.GameState)
	isSavegame := (len(gameStateData) == 1) && (len(gameStateData[0]) == archive.GameStateSize) && (gameStateData[0][0x009C] > 0)
	moddedLevel := len(view.mod.ModifiedBlocks(resource.LangAny, ids.LevelResourcesStart.Plus(lvlids.PerLevel*id+lvlids.FirstUsed))) > 0

	return moddedLevel && !isSavegame
}

func (view *ObjectsView) requestBaseChange(lvl *level.Level, modifier func(*level.ObjectMasterEntry)) {
	view.changeObjectMaster(lvl, view.model.selectedObjects.list, modifier)
}

func (view *ObjectsView) changeObjectMaster(lvl *level.Level, objectIDs []level.ObjectID, modifier func(*level.ObjectMasterEntry)) {
	for _, id := range objectIDs {
		obj := lvl.Object(id)
		modifier(obj)
	}

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
