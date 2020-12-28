package world

import "github.com/inkyblackness/hacked/ss1"

const (
	// StartingLevel identifies the level a new game is started in by default.
	StartingLevel = 1
	// StartingTileX identifies the default X position of the protagonist.
	StartingTileX = 30
	// StartingTileY identifies the default Y position of the protagonist.
	StartingTileY = 22

	// MaxWorldTextures is the limit of how many textures the engine supports.
	// Note that this value is equal to that of the resource limits in package ids. It is actually based on them.
	MaxWorldTextures = 293
)

const (
	// TexturePropertiesFilename specifies the lowercase name of the file containing texture properties.
	TexturePropertiesFilename = "textprop.dat"

	// ObjectPropertiesFilename specifies the lowercase name of the file containing object properties.
	ObjectPropertiesFilename = "objprop.dat"
)

const (
	errNoResourcesFound ss1.StringError = "resource unknown"
)
