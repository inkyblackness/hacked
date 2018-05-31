package world

import "github.com/inkyblackness/hacked/ss1/resource"

// ResourceViewStrategy returns a strategy that is typical for the game.
func ResourceViewStrategy() resource.ResourceViewStrategy {
	return defaultResources{}
}

type defaultResources struct{}

func (def defaultResources) IsCompoundList(id resource.ID) bool {
	// TODO: extend & finalize these lists
	isGameTextureList := (id == resource.ID(0x004C)) || (id == resource.ID(0x004D))
	return isGameTextureList
}
