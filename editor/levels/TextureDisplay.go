package levels

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
)

// TextureDisplay is an enumeration which texture to display in a 2D map view.
type TextureDisplay int

// TextureDisplay constants are listed below.
const (
	TextureDisplayFloor   TextureDisplay = 0
	TextureDisplayWall    TextureDisplay = 1
	TextureDisplayCeiling TextureDisplay = 2
)

// String returns a textual representation.
func (display TextureDisplay) String() string {
	switch display {
	case TextureDisplayFloor:
		return "Floor"
	case TextureDisplayWall:
		return "Wall"
	case TextureDisplayCeiling:
		return "Ceiling"
	default:
		return fmt.Sprintf("Unknown%d", int(display))
	}
}

// TextureDisplays returns all TextureDisplay constants.
func TextureDisplays() []TextureDisplay {
	return []TextureDisplay{TextureDisplayFloor, TextureDisplayWall, TextureDisplayCeiling}
}

// TextureDisplayFunc is one that resolves texture properties.
type TextureDisplayFunc func(*level.TileMapEntry) (atlasIndex level.AtlasIndex, textureRotations int)

// FloorTexture returns the properties for the floor.
func FloorTexture(tile *level.TileMapEntry) (atlasIndex level.AtlasIndex, textureRotations int) {
	return tile.TextureInfo.FloorTextureIndex(), tile.Floor.TextureRotations()
}

// CeilingTexture returns the properties for the ceiling.
func CeilingTexture(tile *level.TileMapEntry) (atlasIndex level.AtlasIndex, textureRotations int) {
	return tile.TextureInfo.CeilingTextureIndex(), tile.Ceiling.TextureRotations()
}

// WallTexture returns the properties for the floor.
func WallTexture(tile *level.TileMapEntry) (atlasIndex level.AtlasIndex, textureRotations int) {
	return tile.TextureInfo.WallTextureIndex(), 0
}

// Func returns the display func for the current display setting.
func (display TextureDisplay) Func() TextureDisplayFunc {
	if display == TextureDisplayFloor {
		return FloorTexture
	} else if display == TextureDisplayCeiling {
		return CeilingTexture
	}
	return WallTexture
}
