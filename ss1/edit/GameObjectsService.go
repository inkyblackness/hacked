package edit

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/world"
)

// GameObjectBitmapInfo contains indices about bitmaps for a particual game object.
type GameObjectBitmapInfo struct {
	// Start is the index for the first object bitmap in ids.ObjectBitmaps.
	Start int
	// Count is the amount of bitmaps in total the object has.
	Count int
	// IconRecommendation is the offset to start for using icons.
	IconRecommendation int
}

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

// BitmapInfo returns information on all bitmaps.
func (service *GameObjectsService) BitmapInfo() map[object.Triple]GameObjectBitmapInfo {
	properties := service.AllProperties()
	result := make(map[object.Triple]GameObjectBitmapInfo)
	start := 0
	properties.Iterate(func(triple object.Triple, prop *object.Properties) bool {
		numExtra := int(prop.Common.Bitmap3D.FrameNumber())
		info := GameObjectBitmapInfo{
			Start:              start,
			Count:              3 + numExtra,
			IconRecommendation: 1,
		}
		if triple.Class != object.ClassTrap {
			info.IconRecommendation += 2
		}
		start += info.Count
		result[triple] = info
		return true
	})
	return result
}
