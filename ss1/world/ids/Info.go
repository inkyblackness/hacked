package ids

import (
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"
	"github.com/inkyblackness/hacked/ss1/content/objprop"
	"github.com/inkyblackness/hacked/ss1/resource"
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
	levelInfo := func(lvl int, lvlResID int, compressed bool) ResourceInfo {
		resID := LevelResourcesStart.Plus(lvl*lvlids.PerLevel + lvlResID)
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
	for lvl := 0; lvl < archive.MaxLevels; lvl++ {
		register(levelInfo(lvl, lvlids.MapVersionNumber, false))
		register(levelInfo(lvl, lvlids.ObjectVersionNumber, false))
		register(levelInfo(lvl, lvlids.Information, true))

		register(levelInfo(lvl, lvlids.TileMap, true))
		register(levelInfo(lvl, lvlids.Schedules, false))
		register(levelInfo(lvl, lvlids.TextureMap, false))
		register(levelInfo(lvl, lvlids.MasterObjectTable, true))
		register(levelInfo(lvl, lvlids.ObjectCrossRefTable, true))
		for class := 0; class < objprop.ClassCount; class++ {
			register(levelInfo(lvl, lvlids.ObjectClassTablesStart+class, true))
			register(levelInfo(lvl, lvlids.ObjectDefaultTablesStart+class, true))
		}

		register(levelInfo(lvl, lvlids.SavefileVersion, false))
		register(levelInfo(lvl, lvlids.Unused41, false))

		register(levelInfo(lvl, lvlids.TextureAnimations, true))
		register(levelInfo(lvl, lvlids.SurveillanceSources, true))
		register(levelInfo(lvl, lvlids.SurveillanceSurrogates, true))
		register(levelInfo(lvl, lvlids.Variables, true))
		register(levelInfo(lvl, lvlids.MapNotes, true))
		register(levelInfo(lvl, lvlids.MapNotesPointer, false))

		register(levelInfo(lvl, lvlids.Unknown48, true))
		register(levelInfo(lvl, lvlids.Unknown49, true))
		register(levelInfo(lvl, lvlids.Unknown50, false))

		register(levelInfo(lvl, lvlids.LoopConfiguration, true))

		register(levelInfo(lvl, lvlids.Unknown52, false))
		register(levelInfo(lvl, lvlids.Unknown53, true))
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
