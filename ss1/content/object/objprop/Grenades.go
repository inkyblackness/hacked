package objprop

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

var grenadeGenerics = interpreters.New().
	Refining("BasicWeapon", 0, 8, basicWeapon, interpreters.Always).
	With("Touchiness", 8, 1).
	With("BlastRadius", 9, 1).
	With("BlastCoreRange", 10, 1).
	With("BlastDamage", 11, 1).
	With("AttackMass", 12, 1).
	With("Flags", 13, 1).As(interpreters.Bitfield(map[uint32]string{
	0x01: "Explode on contact",
	0x04: "Explode by timer",
	0x08: "Explode by vicinity"}))

var timedExplosives = interpreters.New().
	With("MinimumTime", 0, 1).
	With("MaximumTime", 1, 1).
	With("RandomFactor", 2, 1)

func initGrenades() {
	objClass := object.Class(3)

	genericDescriptions[objClass] = grenadeGenerics

	setSpecific(objClass, 1, timedExplosives)
}
