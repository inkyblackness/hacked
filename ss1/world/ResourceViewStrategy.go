package world

import "github.com/inkyblackness/hacked/ss1/resource"

// ResourceViewStrategy defines how selected resources shall be viewed.
type ResourceViewStrategy interface {
	// IsCompoundList returns true for compound from where each contained block is a separate entity.
	// Separate entities are those that can be replaced without affecting others.
	// Examples are the small game textures and the list of object names.
	IsCompoundList(id resource.ID) bool
}

// StandardResourceViewStrategy returns a strategy that is typical for the game.
func StandardResourceViewStrategy() ResourceViewStrategy {
	return defaultResources{}
}

type defaultResources struct{}

func (def defaultResources) IsCompoundList(id resource.ID) bool {
	isGameTextureList := (id == resource.ID(0x004C)) || (id == resource.ID(0x004D))
	return isGameTextureList
}
