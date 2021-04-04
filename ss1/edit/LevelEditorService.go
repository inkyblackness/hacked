package edit

import (
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/world"
)

// LevelEditorService provides level editing functionality based on the currently selected level and content.
type LevelEditorService struct {
	registry cmd.Registry

	levels         *EditableLevels
	levelSelection *LevelSelectionService
}

// NewLevelEditorService returns a new instance.
func NewLevelEditorService(registry cmd.Registry, levels *EditableLevels, levelSelection *LevelSelectionService) *LevelEditorService {
	return &LevelEditorService{
		registry:       registry,
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

// ChangeObjects modifies the basic properties of objects.
func (service *LevelEditorService) ChangeObjects(modifier func(*level.ObjectMainEntry)) error {
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
