package level

import "fmt"

// WallTexturePattern describes how a wall texture should be applied for a tile.
type WallTexturePattern byte

// String returns the textual information.
func (pattern WallTexturePattern) String() string {
	switch pattern {
	case WallTexturePatternRegular:
		return "Regular"
	case WallTexturePatternFlipHorizontal:
		return "FlipHorizontal"
	case WallTexturePatternFlipAlternating:
		return "FlipAlternating"
	case WallTexturePatternFlipAlternatingInverted:
		return "FlipAlternatingInverted"
	default:
		return fmt.Sprintf("Unknown%02X", int(pattern))
	}
}

// WallTexturePattern constants
const (
	WallTexturePatternRegular                 WallTexturePattern = 0
	WallTexturePatternFlipHorizontal          WallTexturePattern = 1
	WallTexturePatternFlipAlternating         WallTexturePattern = 2
	WallTexturePatternFlipAlternatingInverted WallTexturePattern = 3
)

// WallTexturePatterns returns all available texture patterns.
func WallTexturePatterns() []WallTexturePattern {
	return []WallTexturePattern{
		WallTexturePatternRegular,
		WallTexturePatternFlipHorizontal,
		WallTexturePatternFlipAlternating,
		WallTexturePatternFlipAlternatingInverted,
	}
}
