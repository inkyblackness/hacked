package level

import "fmt"

const (
	// TextureAnimationCount describes how many entries of texture animations a level can have.
	TextureAnimationCount = 4

	// TextureAnimationEntrySize is the size, in bytes, of one animation entry.
	TextureAnimationEntrySize = 7
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

// String returns the textual representation.
func (loopType TextureAnimationLoopType) String() string {
	switch loopType {
	case TextureAnimationForward:
		return "Forward"
	case TextureAnimationForthAndBack:
		return "Forth-And-Back"
	case TextureAnimationBackAndForth:
		return "Back-And-Forth"
	default:
		return fmt.Sprintf("Unknown%02X", int(loopType))
	}
}

const (
	// TextureAnimationForward has the texture animation run in a linear loop.
	TextureAnimationForward = TextureAnimationLoopType(0x00)
	// TextureAnimationForthAndBack reverses the sequence at every end.
	TextureAnimationForthAndBack = TextureAnimationLoopType(0x01)
	// TextureAnimationBackAndForth reverses the sequence at every end.
	TextureAnimationBackAndForth = TextureAnimationLoopType(0x81)
)

// TextureAnimationLoopTypes returns all constants.
func TextureAnimationLoopTypes() []TextureAnimationLoopType {
	return []TextureAnimationLoopType{
		TextureAnimationForward,
		TextureAnimationForthAndBack,
		TextureAnimationBackAndForth,
	}
}
