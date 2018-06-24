package ids

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/constants"
)

// ResourceInfo describes a group of resources with their default serialization properties.
type ResourceInfo struct {
	// StartID is the first ID of the resource block (inclusive).
	StartID resource.ID
	// EndID is the last ID of the resource block (exclusive).
	EndID resource.ID

	// ContentType describes how to interpret resource data.
	ContentType resource.ContentType
	// Compound indicates whether the resource has a variable amount of blocks. Simple resources always have one block.
	Compound bool
	// Compressed indicates that the resource data shall be stored in compressed form.
	Compressed bool

	// List is set for compound resources that have an atomic resource per block.
	List bool
	// MaxCount describes how many resources can be stored at maximum.
	MaxCount int

	// ResFile specifies the .res file in which this resource should be stored.
	ResFile resource.Filename
}

// Info returns the resource information for the identified resource.
func Info(id resource.ID) (ResourceInfo, bool) {
	info, existing := infoByID[id]
	return info, existing
}

func init() {
	register := func(info ResourceInfo) {
		count := info.EndID.Value() - info.StartID.Value()
		for offset := 0; offset < int(count); offset++ {
			infoByID[info.StartID.Plus(offset)] = info
		}
	}
	for _, info := range infoList {
		register(info)
	}
	levelInfo := func(lvl int, lvlResID LevelResourceID, compressed bool) ResourceInfo {
		resID := lvlResID.ForLevel(lvl)
		return ResourceInfo{
			StartID: resID,
			EndID:   resID.Plus(1),

			ContentType: resource.Archive,
			Compound:    false,
			Compressed:  compressed,

			List:     false,
			MaxCount: 1,

			ResFile: Archive,
		}
	}
	// In the following, the compression flags were set if the vanilla archive had
	// the respective resource compressed in at least one level.
	// Most of them were compressed uniform across all levels, some had exceptions.
	for lvl := 0; lvl < constants.MaxLevels; lvl++ {
		register(levelInfo(lvl, LvlMapVersionNumber, false))
		register(levelInfo(lvl, LvlObjectVersionNumber, false))
		register(levelInfo(lvl, LvlInformation, true))

		register(levelInfo(lvl, LvlTileMap, true))
		register(levelInfo(lvl, LvlSchedules, false))
		register(levelInfo(lvl, LvlTextureMap, false))
		register(levelInfo(lvl, LvlMasterObjectTable, true))
		register(levelInfo(lvl, LvlObjectCrossRefTable, true))
		for class := LevelResourceID(0); class < constants.ObjectClasses; class++ {
			register(levelInfo(lvl, LvlObjectClassTablesStart+class, true))
			register(levelInfo(lvl, LvlObjectDefaultTablesStart+class, true))
		}

		register(levelInfo(lvl, LvlSavefileVersion, false))
		register(levelInfo(lvl, LvlUnused41, false))

		register(levelInfo(lvl, LvlTextureAnimations, true))
		register(levelInfo(lvl, LvlSurveillanceSources, true))
		register(levelInfo(lvl, LvlSurveillanceSurrogates, true))
		register(levelInfo(lvl, LvlVariables, true))
		register(levelInfo(lvl, LvlMapNotes, true))
		register(levelInfo(lvl, LvlMapNotesPointer, false))

		register(levelInfo(lvl, LvlUnknown48, true))
		register(levelInfo(lvl, LvlUnknown49, true))
		register(levelInfo(lvl, LvlUnknown50, false))

		register(levelInfo(lvl, LvlLoopConfiguration, true))

		register(levelInfo(lvl, LvlUnknown52, false))
		register(levelInfo(lvl, LvlUnknown53, true))
	}
}

var infoByID = map[resource.ID]ResourceInfo{}

var infoList = []ResourceInfo{
	{IconTextures, IconTextures.Plus(1), resource.Bitmap, true, false, true, 300, Texture},
	{SmallTextures, SmallTextures.Plus(1), resource.Bitmap, true, false, true, 300, Texture},

	{PaperTextsStart, PaperTextsStart.Plus(16), resource.Text, true, false, false, 16, CybStrng},

	{TrapMessageTexts, TrapMessageTexts.Plus(1), resource.Text, true, false, true, 512, CybStrng},
	{WordTexts, WordTexts.Plus(1), resource.Text, true, false, true, 512, CybStrng},
	{PanelNameTexts, PanelNameTexts.Plus(1), resource.Text, true, false, true, 256, CybStrng},
	{LogCategoryTexts, LogCategoryTexts.Plus(1), resource.Text, true, false, true, 16, CybStrng},
	{VariousMessageTexts, VariousMessageTexts.Plus(1), resource.Text, true, false, true, 256, CybStrng},
	{ScreenMessageTexts, ScreenMessageTexts.Plus(1), resource.Text, true, false, true, 120, CybStrng},
	{InfoNodeMessageTexts, InfoNodeMessageTexts.Plus(1), resource.Text, true, false, true, 256, CybStrng},
	{AccessCardNameTexts, AccessCardNameTexts.Plus(1), resource.Text, true, false, true, 32 * 2, CybStrng},
	{DataletMessageTexts, DataletMessageTexts.Plus(1), resource.Text, true, false, true, 256, CybStrng},

	{ArchiveName, ArchiveName.Plus(1), resource.Archive, false, false, false, 1, Archive},
	{GameState, GameState.Plus(1), resource.Archive, false, true, false, 1, Archive},
}
