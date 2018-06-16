package ids

import "github.com/inkyblackness/hacked/ss1/resource"

// ResourceInfo describes a group of resources with their default serialization properties.
type ResourceInfo struct {
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
	info, existing := infoById[id]
	return info, existing
}

var infoById = map[resource.ID]ResourceInfo{
	PaperTextsStart: {resource.Text, true, false, false},

	TrapMessageTexts:     {resource.Text, true, false, true},
	WordTexts:            {resource.Text, true, false, true},
	PanelNameTexts:       {resource.Text, true, false, true},
	LogCategoryTexts:     {resource.Text, true, false, true},
	VariousMessageTexts:  {resource.Text, true, false, true},
	ScreenMessageTexts:   {resource.Text, true, false, true},
	InfoNodeMessageTexts: {resource.Text, true, false, true},
	AccessCardNameTexts:  {resource.Text, true, false, true},
	DataletMessageTexts:  {resource.Text, true, false, true},
}
