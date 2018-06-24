package level

const (
	// TextureAnimationCount describes how many entries of texture animations a level can have.
	TextureAnimationCount = 4
)

// TextureAnimationEntry describes one entry of the texture animation table.
type TextureAnimationEntry struct {
	FrameTime         uint16
	CurrentFrameTime  uint16
	CurrentFrameIndex byte
	FrameCount        byte
	LoopType          TextureAnimationLoopType
}

// TextureAnimationLoopType describes how a texture animation loop should advance.
type TextureAnimationLoopType byte

const (
	// TextureAnimationForward has the texture animation run in a linear loop.
	TextureAnimationForward = TextureAnimationLoopType(0x00)
	// TextureAnimationForthAndBack reverses the sequence at every end.
	TextureAnimationForthAndBack = TextureAnimationLoopType(0x01)
	// TextureAnimationBackAndForth reverses the sequence at every end.
	TextureAnimationBackAndForth = TextureAnimationLoopType(0x81)
)
