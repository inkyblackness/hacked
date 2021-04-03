package edit

import (
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/world"
)

type LevelEditorService struct {
	registry cmd.Registry

	levels         *EditableLevels
	levelSelection *LevelSelectionService
}

func NewLevelEditorService(registry cmd.Registry, levels *EditableLevels, levelSelection *LevelSelectionService) *LevelEditorService {
	return &LevelEditorService{
		registry:       registry,
		levels:         levels,
		levelSelection: levelSelection,
	}
}

func (service *LevelEditorService) IsReadOnly() bool {
	return service.levels.IsLevelReadOnly(service.levelSelection.CurrentLevelID())
}

func (service *LevelEditorService) Level() *level.Level {
	return service.levels.Level(service.levelSelection.CurrentLevelID())
}

func (service *LevelEditorService) Tiles() []*level.TileMapEntry {
	lvl := service.Level()
	positions := service.levelSelection.CurrentSelectedTiles()
	tiles := make([]*level.TileMapEntry, len(positions))
	for i, pos := range positions {
		tiles[i] = lvl.Tile(pos)
	}
	return tiles
}

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
