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
	isTextLinesList := id == ids.TrapMessages
	return isGameTextureList || isTextLinesList
}
