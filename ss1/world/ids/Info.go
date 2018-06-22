package ids

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/resource/lgres"
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
	for _, info := range infoList {
		count := info.EndID.Value() - info.StartID.Value()
		for offset := 0; offset < int(count); offset++ {
			infoByID[info.StartID.Plus(offset)] = info
		}
	}
}

var infoByID = map[resource.ID]ResourceInfo{}

var infoList = []ResourceInfo{
	{IconTextures, IconTextures.Plus(1), resource.Bitmap, true, false, true, 300, lgres.Texture},
	{SmallTextures, SmallTextures.Plus(1), resource.Bitmap, true, false, true, 300, lgres.Texture},

	{PaperTextsStart, PaperTextsStart.Plus(16), resource.Text, true, false, false, 16, lgres.CybStrng},

	{TrapMessageTexts, TrapMessageTexts.Plus(1), resource.Text, true, false, true, 512, lgres.CybStrng},
	{WordTexts, WordTexts.Plus(1), resource.Text, true, false, true, 512, lgres.CybStrng},
	{PanelNameTexts, PanelNameTexts.Plus(1), resource.Text, true, false, true, 256, lgres.CybStrng},
	{LogCategoryTexts, LogCategoryTexts.Plus(1), resource.Text, true, false, true, 16, lgres.CybStrng},
	{VariousMessageTexts, VariousMessageTexts.Plus(1), resource.Text, true, false, true, 256, lgres.CybStrng},
	{ScreenMessageTexts, ScreenMessageTexts.Plus(1), resource.Text, true, false, true, 120, lgres.CybStrng},
	{InfoNodeMessageTexts, InfoNodeMessageTexts.Plus(1), resource.Text, true, false, true, 256, lgres.CybStrng},
	{AccessCardNameTexts, AccessCardNameTexts.Plus(1), resource.Text, true, false, true, 32 * 2, lgres.CybStrng},
	{DataletMessageTexts, DataletMessageTexts.Plus(1), resource.Text, true, false, true, 256, lgres.CybStrng},
}
