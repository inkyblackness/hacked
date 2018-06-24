package ids

import "github.com/inkyblackness/hacked/ss1/resource"

const (
	// MaxResourcesPerLevel describes how many resources a level could store.
	MaxResourcesPerLevel = 100

	// FirstUsedLevelResource identifies the first level-specific resource in a level-range.
	FirstUsedLevelResource = 2
	// FirstUnusedLevelResource identifies the first unused level-specific resource in a level-range.
	FirstUnusedLevelResource = 54
)

// LevelResourceID identifies a level-specific resource.
type LevelResourceID byte

// ForLevel returns the final resource identifier for a specific level.
func (id LevelResourceID) ForLevel(level int) resource.ID {
	return LevelResourcesStart.Plus((level * MaxResourcesPerLevel) + int(id))
}

// Level resource identifier
const (
	LvlMapVersionNumber    LevelResourceID = 2
	LvlObjectVersionNumber LevelResourceID = 3
	LvlInformation         LevelResourceID = 4

	LvlTileMap                  LevelResourceID = 5
	LvlSchedules                LevelResourceID = 6
	LvlTextureMap               LevelResourceID = 7
	LvlMasterObjectTable        LevelResourceID = 8
	LvlObjectCrossRefTable      LevelResourceID = 9
	LvlObjectClassTablesStart   LevelResourceID = 10
	LvlObjectDefaultTablesStart LevelResourceID = 25

	LvlSavefileVersion LevelResourceID = 40
	LvlUnused41        LevelResourceID = 41

	LvlTextureAnimations      LevelResourceID = 42
	LvlSurveillanceSources    LevelResourceID = 43
	LvlSurveillanceSurrogates LevelResourceID = 44
	LvlVariables              LevelResourceID = 45
	LvlMapNotes               LevelResourceID = 46
	LvlMapNotesPointer        LevelResourceID = 47

	LvlUnknown48 LevelResourceID = 48
	LvlUnknown49 LevelResourceID = 49
	LvlUnknown50 LevelResourceID = 50

	LvlLoopConfiguration LevelResourceID = 51

	LvlUnknown52 LevelResourceID = 52
	LvlUnknown53 LevelResourceID = 53
)
