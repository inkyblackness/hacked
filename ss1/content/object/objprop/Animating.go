package objprop

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

var animationGenerics = interpreters.New().
	With("FrameTime", 0, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int) string {
		return fmt.Sprintf("%3.0f millisec - raw: %d", (float64(value)*900)/255.0, value)
	})).
	With("Flags", 1, 1).As(interpreters.Bitfield(map[uint32]string{0x01: "Emit Light"}))

var explosionAnimation = interpreters.New().
	With("FrameExplode", 0, 1)

func initAnimating() {
	objClass := object.Class(11)

	genericDescriptions[objClass] = animationGenerics

	setSpecific(objClass, 2, explosionAnimation)
}
