package edit

import (
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// EditableLevels is a list of levels that can be modified.
type EditableLevels struct {
	list [archive.MaxLevels]*level.Level
}

// NewEditableLevels returns a new instance for given mod.
func NewEditableLevels(mod *world.Mod) *EditableLevels {
	var levels EditableLevels
	for i := 0; i < len(levels.list); i++ {
		levels.list[i] = level.NewLevel(ids.LevelResourcesStart, i, mod)
	}
	return &levels
}

// InvalidateResources informs all levels of changed resource identifier.
func (levels *EditableLevels) InvalidateResources(modifiedIDs []resource.ID) {
	for _, lvl := range levels.list {
		lvl.InvalidateResources(modifiedIDs)
	}
}

// Level returns the pointer to the given level. The function panics if the index is invalid.
func (levels *EditableLevels) Level(index int) *level.Level {
	if !levels.IsLevelAvailable(index) {
		panic("Invalid level index")
	}
	return levels.list[index]
}

// IsLevelAvailable returns true if the given index is a valid level.
func (levels *EditableLevels) IsLevelAvailable(index int) bool {
	return (index >= 0) && (index < len(levels.list))
}

// IsObjectInUse returns true if the given object ID is in use in given level.
func (levels *EditableLevels) IsObjectInUse(levelIndex int, id level.ObjectID) bool {
	return false
}

// IsTileOnMap returns true if the given tile position points to a valid tile in the given level.
func (levels *EditableLevels) IsTileOnMap(levelIndex int, x, y int) bool {
	return false
}
