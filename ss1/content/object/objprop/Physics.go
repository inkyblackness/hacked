package objprop

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

var projectileGenerics = interpreters.New().
	With("Flags", 0, 1).As(interpreters.Bitfield(map[uint32]string{
	0x01: "EmitLight",
	0x02: "BounceOffWalls",
	0x04: "BounceOffObjects",
	0x08: "BounceOffProjectiles"}))

var cyberProjectiles = interpreters.New().
	Refining("ColorScheme", 0, 6, cyberColorScheme, interpreters.Always)

func initPhysics() {
	objClass := object.Class(2)

	genericDescriptions[objClass] = projectileGenerics

	setSpecificByType(objClass, 1, 9, cyberProjectiles)
	setSpecificByType(objClass, 1, 10, cyberProjectiles)
	setSpecificByType(objClass, 1, 11, cyberProjectiles)
	setSpecificByType(objClass, 1, 12, cyberProjectiles)
	setSpecificByType(objClass, 1, 13, cyberProjectiles)
}
