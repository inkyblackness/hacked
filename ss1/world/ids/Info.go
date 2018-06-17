package ids

import "github.com/inkyblackness/hacked/ss1/resource"

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
	{IconTextures, IconTextures.Plus(1), resource.Bitmap, true, false, true},
	{SmallTextures, SmallTextures.Plus(1), resource.Bitmap, true, false, true},

	{PaperTextsStart, PaperTextsStart.Plus(256), resource.Text, true, false, false},

	{TrapMessageTexts, TrapMessageTexts.Plus(1), resource.Text, true, false, true},
	{WordTexts, WordTexts.Plus(1), resource.Text, true, false, true},
	{PanelNameTexts, PanelNameTexts.Plus(1), resource.Text, true, false, true},
	{LogCategoryTexts, LogCategoryTexts.Plus(1), resource.Text, true, false, true},
	{VariousMessageTexts, VariousMessageTexts.Plus(1), resource.Text, true, false, true},
	{ScreenMessageTexts, ScreenMessageTexts.Plus(1), resource.Text, true, false, true},
	{InfoNodeMessageTexts, InfoNodeMessageTexts.Plus(1), resource.Text, true, false, true},
	{AccessCardNameTexts, AccessCardNameTexts.Plus(1), resource.Text, true, false, true},
	{DataletMessageTexts, DataletMessageTexts.Plus(1), resource.Text, true, false, true},
}
