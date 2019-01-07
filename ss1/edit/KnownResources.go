package edit

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// TextInfo describes one kind of text.
type TextInfo struct {
	ID    resource.ID
	Title string
}

// TextInfoList is a set of TextInfo.
type TextInfoList []TextInfo

// Title returns the title property of the identified info.
func (list TextInfoList) Title(id resource.ID) string {
	for _, info := range list {
		if info.ID == id {
			return info.Title
		}
	}
	return fmt.Sprintf("??? %v", id)
}

var knownTexts = TextInfoList{
	{ID: ids.TrapMessageTexts, Title: "Trap Messages"},
	{ID: ids.WordTexts, Title: "Words"},
	{ID: ids.LogCategoryTexts, Title: "Log Categories"},
	{ID: ids.VariousMessageTexts, Title: "Various Messages"},
	{ID: ids.ScreenMessageTexts, Title: "Screen Messages"},
	{ID: ids.InfoNodeMessageTexts, Title: "Info Node Message Texts (8/5/6)"},
	{ID: ids.AccessCardNameTexts, Title: "Access Card Names"},
	{ID: ids.DataletMessageTexts, Title: "Datalet Messages (8/5/8)"},
	{ID: ids.PaperTextsStart, Title: "Papers"},
	{ID: ids.PanelNameTexts, Title: "Panel Names"},
}

// KnownTexts returns a set of known texts.
func KnownTexts() TextInfoList {
	return knownTexts
}
