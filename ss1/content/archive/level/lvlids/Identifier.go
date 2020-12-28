package lvlids

const (
	// PerLevel is the amount of resources reserved per level.
	PerLevel = 100

	// FirstUsed identifies the first level-specific resource in a level-range.
	FirstUsed = 2
	// FirstUnused identifies the first unused level-specific resource in a level-range.
	FirstUnused = 54
)

// Level resource identifier are listed below.
const (
	MapVersionNumber    = 2
	ObjectVersionNumber = 3
	Information         = 4

	TileMap                  = 5
	Schedules                = 6
	TextureAtlas             = 7
	ObjectMainTable          = 8
	ObjectCrossRefTable      = 9
	ObjectClassTablesStart   = 10
	ObjectDefaultTablesStart = 25

	SavefileVersion = 40
	Unused41        = 41

	TextureAnimations      = 42
	SurveillanceSources    = 43
	SurveillanceSurrogates = 44
	Parameters             = 45
	MapNotes               = 46
	MapNotesPointer        = 47

	Unknown48 = 48
	Unknown49 = 49
	Unknown50 = 50

	LoopConfiguration = 51

	Unknown52        = 52
	HeightSemaphores = 53
)
