package edit

import (
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/world"
)

// LevelEditorService provides level editing functionality based on the currently selected level and content.
type LevelEditorService struct {
	registry cmd.Registry

	gameObjects    *GameObjectsService
	levels         *EditableLevels
	levelSelection *LevelSelectionService

	newObjectTriple object.Triple
}

// NewLevelEditorService returns a new instance.
func NewLevelEditorService(registry cmd.Registry,
	gameObjects *GameObjectsService,
	levels *EditableLevels,
	levelSelection *LevelSelectionService) *LevelEditorService {
	return &LevelEditorService{
		registry:       registry,
		gameObjects:    gameObjects,
		levels:         levels,
		levelSelection: levelSelection,
	}
}

// IsReadOnly returns true if the currently selected level cannot be modified.
func (service *LevelEditorService) IsReadOnly() bool {
	return service.levels.IsLevelReadOnly(service.levelSelection.CurrentLevelID())
}

// Level returns the currently selected level.
func (service *LevelEditorService) Level() *level.Level {
	return service.levels.Level(service.levelSelection.CurrentLevelID())
}

// Tiles returns the list of currently selected tiles of the current level.
func (service *LevelEditorService) Tiles() []*level.TileMapEntry {
	lvl := service.Level()
	positions := service.levelSelection.CurrentSelectedTiles()
	tiles := make([]*level.TileMapEntry, len(positions))
	for i, pos := range positions {
		tiles[i] = lvl.Tile(pos)
	}
	return tiles
}

// HasSelectedTiles returns true if Tiles() returns at least one entry.
func (service *LevelEditorService) HasSelectedTiles() bool {
	return len(service.levelSelection.CurrentSelectedTiles()) > 0
}

// ClearTiles resets the selected tiles.
func (service *LevelEditorService) ClearTiles() error {
	return service.ChangeTiles(func(tile *level.TileMapEntry) {
		index := tile.FirstObjectIndex
		tile.Reset()
		tile.FirstObjectIndex = index
	})
}

// ChangeTiles performs a modification on all currently selected tiles and commits these changes to the repository.
func (service *LevelEditorService) ChangeTiles(modifier func(*level.TileMapEntry)) error {
	positions := service.levelSelection.CurrentSelectedTiles()
	if len(positions) == 0 {
		return nil
	}
	lvl := service.Level()
	for _, pos := range positions {
		modifier(lvl.Tile(pos))
	}
	levelID := lvl.ID()
	return service.registry.Register(cmd.Named("ChangeTiles"),
		cmd.Forward(service.levelSelection.SetCurrentLevelIDTask(levelID)),
		cmd.Forward(service.setSelectedTilesTask(positions)),
		cmd.Nested(func() error { return service.levels.CommitLevelChanges(levelID) }),
		cmd.Reverse(service.setSelectedTilesTask(positions)),
		cmd.Reverse(service.levelSelection.SetCurrentLevelIDTask(levelID)),
	)
}

// Objects returns the list of currently selected objects of the current level.
func (service *LevelEditorService) Objects() []*level.ObjectMainEntry {
	lvl := service.Level()
	ids := service.levelSelection.CurrentSelectedObjects()
	objects := make([]*level.ObjectMainEntry, len(ids))
	for i, id := range ids {
		objects[i] = lvl.Object(id)
	}
	return objects
}

// HasSelectedObjects returns true if a call to Objects() returns at least one item.
func (service *LevelEditorService) HasSelectedObjects() bool {
	return len(service.levelSelection.CurrentSelectedObjects()) > 0
}

// NewObjectTriple returns the identifier of implicitly created objects.
func (service *LevelEditorService) NewObjectTriple() object.Triple {
	return service.newObjectTriple
}

// SetNewObjectTriple sets the identifier for implicitly created objects.
func (service *LevelEditorService) SetNewObjectTriple(value object.Triple) {
	service.newObjectTriple = value
}

// CreateNewObject adds a new object to the level and selects it.
func (service *LevelEditorService) CreateNewObject(triple object.Triple, modifier level.ObjectMainEntryModifier) error {
	lvl := service.Level()
	id, err := lvl.NewObject(triple.Class)
	if err != nil {
		return err
	}
	obj := lvl.Object(id)
	obj.Subclass = triple.Subclass
	obj.Type = triple.Type
	service.resetHitpointsToDefault(obj)
	modifier(obj)
	service.placeObject(lvl, obj, service.atFloorLevelIn(lvl))
	lvl.UpdateObjectLocation(id)
	return service.patchLevelObjects(lvl, service.levelSelection.CurrentSelectedObjects(), []level.ObjectID{id})
}

// PlaceObjectsOnFloor puts all selected objects to sit on the floor.
func (service *LevelEditorService) PlaceObjectsOnFloor() error {
	lvl := service.Level()
	atHeight := service.atFloorLevelIn(lvl)
	return service.ChangeObjects(func(obj *level.ObjectMainEntry) { service.placeObject(lvl, obj, atHeight) })
}

// PlaceObjectsOnEyeLevel puts all selected objects to be at eye level (approximately).
func (service *LevelEditorService) PlaceObjectsOnEyeLevel() error {
	lvl := service.Level()
	atHeight := service.atEyeLevelIn(lvl)
	return service.ChangeObjects(func(obj *level.ObjectMainEntry) { service.placeObject(lvl, obj, atHeight) })
}

// PlaceObjectsOnCeiling puts all selected objects to hang from the ceiling.
func (service *LevelEditorService) PlaceObjectsOnCeiling() error {
	lvl := service.Level()
	atHeight := service.atCeilingLevelIn(lvl)
	return service.ChangeObjects(func(obj *level.ObjectMainEntry) { service.placeObject(lvl, obj, atHeight) })
}

// ChangeObjects modifies properties of selected objects. The modifier is called for each selected object.
func (service *LevelEditorService) ChangeObjects(modifier level.ObjectMainEntryModifier) error {
	objectIDs := service.levelSelection.CurrentSelectedObjects()
	if len(objectIDs) == 0 {
		return nil
	}
	lvl := service.Level()
	for _, id := range objectIDs {
		obj := lvl.Object(id)
		oldPosition := obj.TilePosition()
		modifier(obj)
		if oldPosition != obj.TilePosition() {
			lvl.UpdateObjectLocation(id)
		}
	}
	return service.patchLevelObjects(lvl, objectIDs, objectIDs)
}

// DeleteObjects deletes all currently selected objects, clearing the selection afterwards.
func (service *LevelEditorService) DeleteObjects() error {
	objectIDs := service.levelSelection.CurrentSelectedObjects()
	if len(objectIDs) == 0 {
		return nil
	}
	lvl := service.Level()
	for _, id := range objectIDs {
		lvl.DelObject(id)
	}
	return service.patchLevelObjects(lvl, objectIDs, nil)
}

func (service *LevelEditorService) patchLevelObjects(
	lvl *level.Level,
	reverseObjectIDs []level.ObjectID,
	forwardObjectIDs []level.ObjectID) error {
	levelID := lvl.ID()

	return service.registry.Register(cmd.Named("PatchLevelObjects"),
		cmd.Forward(service.levelSelection.SetCurrentLevelIDTask(levelID)),
		cmd.Reverse(service.setSelectedObjectsTask(reverseObjectIDs)),
		cmd.Nested(func() error { return service.levels.CommitLevelChanges(levelID) }),
		cmd.Forward(service.setSelectedObjectsTask(forwardObjectIDs)),
		cmd.Reverse(service.levelSelection.SetCurrentLevelIDTask(levelID)),
	)
}

func (service *LevelEditorService) setSelectedTilesTask(positions []level.TilePosition) cmd.Task {
	return func(world.Modder) error {
		service.levelSelection.SetCurrentSelectedTiles(positions)
		return nil
	}
}

func (service *LevelEditorService) setSelectedObjectsTask(ids []level.ObjectID) cmd.Task {
	return func(world.Modder) error {
		service.levelSelection.SetCurrentSelectedObjects(ids)
		return nil
	}
}

type heightCalculatorFunc func(tile *level.TileMapEntry, pos level.FinePosition, objPivot float32) level.HeightUnit

func (service *LevelEditorService) placeObject(lvl *level.Level, obj *level.ObjectMainEntry, atHeight heightCalculatorFunc) {
	var objPivot float32
	prop, err := service.gameObjects.PropertiesFor(obj.Triple())
	if err == nil {
		objPivot = object.Pivot(prop.Common)
	}
	tile := lvl.Tile(obj.TilePosition())
	if tile != nil {
		obj.Z = atHeight(tile, obj.FinePosition(), objPivot)
	}
}

func (service *LevelEditorService) atFloorLevelIn(lvl *level.Level) heightCalculatorFunc {
	_, _, height := lvl.Size()
	return func(tile *level.TileMapEntry, pos level.FinePosition, objPivot float32) level.HeightUnit {
		floorHeight := tile.FloorTileHeightAt(pos, height)
		return height.ValueToObjectHeight(floorHeight + objPivot)
	}
}

func (service *LevelEditorService) atEyeLevelIn(lvl *level.Level) heightCalculatorFunc {
	_, _, height := lvl.Size()
	return func(tile *level.TileMapEntry, pos level.FinePosition, objPivot float32) level.HeightUnit {
		floorHeight := tile.FloorTileHeightAt(pos, height)
		return height.ValueToObjectHeight(floorHeight + 0.75 - objPivot)
	}
}

func (service *LevelEditorService) atCeilingLevelIn(lvl *level.Level) heightCalculatorFunc {
	_, _, height := lvl.Size()
	return func(tile *level.TileMapEntry, pos level.FinePosition, objPivot float32) level.HeightUnit {
		floorHeight := tile.CeilingTileHeightAt(pos, height)
		return height.ValueToObjectHeight(floorHeight - objPivot)
	}
}

func (service *LevelEditorService) resetHitpointsToDefault(obj *level.ObjectMainEntry) {
	prop, err := service.gameObjects.PropertiesFor(obj.Triple())
	if err != nil {
		return
	}
	obj.Hitpoints = prop.Common.Hitpoints
}
