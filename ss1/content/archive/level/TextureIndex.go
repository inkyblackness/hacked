package level

const (
	// DefaultTextureAtlasSize is the amount of textures a level can hold.
	DefaultTextureAtlasSize = 54
	// FloorCeilingTextureLimit describes the amount of textures available for floors/ceilings.
	FloorCeilingTextureLimit = 32
)

// TextureIndex identifies one game texture.
type TextureIndex int16

// TextureAtlas is a selection of textures for the level.
type TextureAtlas []TextureIndex
