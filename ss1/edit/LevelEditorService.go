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
	lvl := service.Level()
	positions := service.levelSelection.CurrentSelectedTiles()
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

func (service *LevelEditorService) setSelectedTilesTask(positions []level.TilePosition) cmd.Task {
	return func(world.Modder) error {
		service.levelSelection.SetCurrentSelectedTiles(positions)
		return nil
	}
}