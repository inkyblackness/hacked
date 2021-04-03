package edit

import (
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// EditableLevels is a list of levels that can be modified.
type EditableLevels struct {
	registry cmd.Registry
	mod      *world.Mod
	list     [archive.MaxLevels]*level.Level
}

// NewEditableLevels returns a new instance for given mod.
func NewEditableLevels(registry cmd.Registry, mod *world.Mod) *EditableLevels {
	var levels EditableLevels
	levels.registry = registry
	levels.mod = mod
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

// IsLevelReadOnly returns true if the identified level cannot be modified.
func (levels *EditableLevels) IsLevelReadOnly(index int) bool {
	isInMod := len(levels.mod.ModifiedBlocks(resource.LangAny, ids.LevelResourcesStart.Plus(lvlids.PerLevel*index+lvlids.FirstUsed))) > 0
	return !levels.IsLevelAvailable(index) || !isInMod
}

// IsLevelAvailable returns true if the given index is a valid level.
func (levels *EditableLevels) IsLevelAvailable(index int) bool {
	return (index >= 0) && (index < len(levels.list))
}

// IsObjectInUse returns true if the given object ID is in use in given level.
func (levels *EditableLevels) IsObjectInUse(levelIndex int, id level.ObjectID) bool {
	if !levels.IsLevelAvailable(levelIndex) {
		return false
	}
	lvl := levels.list[levelIndex]
	obj := lvl.Object(id)
	return (obj != nil) && (obj.InUse != 0)
}

// IsTileOnMap returns true if the given tile position points to a valid tile in the given level.
func (levels *EditableLevels) IsTileOnMap(levelIndex int, pos level.TilePosition) bool {
	if !levels.IsLevelAvailable(levelIndex) {
		return false
	}
	lvl := levels.list[levelIndex]
	return lvl.Tile(pos) != nil
}

// CommitLevelChanges applies any level data change compared to the underlying mod.
func (levels *EditableLevels) CommitLevelChanges(index int) error {
	lvl := levels.Level(index)
	newDataSet := lvl.EncodeState()
	var patches []world.BlockPatch
	for id, newData := range &newDataSet {
		if len(newData) > 0 {
			resourceID := ids.LevelResourcesStart.Plus(lvlids.PerLevel*lvl.ID() + id)
			patch, changed, err := levels.mod.CreateBlockPatch(resource.LangAny, resourceID, 0, newData)
			if err != nil {
				return err
			} else if changed {
				patches = append(patches, patch)
			}
		}
	}

	return levels.registry.Register(cmd.Named("CommitLevelChanges"),
		cmd.Forward(forwardTask(patches)), cmd.Reverse(reverseTask(patches)))
}

func forwardTask(patches []world.BlockPatch) cmd.Task {
	return func(modder world.Modder) error {
		for _, patch := range patches {
			modder.PatchResourceBlock(resource.LangAny, patch.ID, patch.BlockIndex, patch.BlockLength, patch.ForwardData)
		}
		return nil
	}
}

func reverseTask(patches []world.BlockPatch) cmd.Task {
	return func(modder world.Modder) error {
		for _, patch := range patches {
			modder.PatchResourceBlock(resource.LangAny, patch.ID, patch.BlockIndex, patch.BlockLength, patch.ReverseData)
		}
		return nil
	}
}
