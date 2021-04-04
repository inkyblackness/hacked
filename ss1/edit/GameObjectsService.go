package edit

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/world"
)

// GameObjectsService provides information for game objects.
type GameObjectsService struct {
	mod *world.Mod
}

// NewGameObjectsService returns a new instance.
func NewGameObjectsService(mod *world.Mod) *GameObjectsService {
	return &GameObjectsService{
		mod: mod,
	}
}

// AllProperties returns the currently know properties.
func (service *GameObjectsService) AllProperties() object.PropertiesTable {
	return service.mod.ObjectProperties()
}

// AllProperties returns the currently know properties.
func (service *GameObjectsService) PropertiesFor(triple object.Triple) (object.Properties, error) {
	prop, err := service.AllProperties().ForObject(triple)
	if err != nil {
		return object.Properties{}, err
	}
	return *prop, err
}
