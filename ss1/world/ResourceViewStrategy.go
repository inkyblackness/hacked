package world

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// ResourceViewStrategy returns a strategy that is typical for the game.
func ResourceViewStrategy() resource.ViewStrategy {
	return defaultResources{}
}

type defaultResources struct{}

func (def defaultResources) IsCompoundList(id resource.ID) bool {
	// TODO: extend & finalize these lists
	isGameTextureList := (id == ids.IconTextures) || (id == ids.SmallTextures)
	isTextLinesList := (id == ids.TrapMessageTexts) || (id == ids.WordTexts) || (id == ids.PanelNameTexts) ||
		(id == ids.LogCategoryTexts) || (id == ids.VariousMessageTexts) || (id == ids.ScreenMessageTexts) ||
		(id == ids.InfoNodeMessageTexts) || (id == ids.AccessCardNameTexts) || (id == ids.DataletMessageTexts)
	return isGameTextureList || isTextLinesList
}
