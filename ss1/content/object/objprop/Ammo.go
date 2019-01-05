package objprop

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

var ammoClipGenerics = interpreters.New().
	Refining("BasicWeapon", 0, 8, basicWeapon, interpreters.Always).
	With("CartrigeSize", 8, 1).
	With("BulletMass", 9, 1).
	With("BulletSpeed", 10, 2).As(interpreters.RangedValue(-10000, +10000)).
	With("Range", 12, 1).
	With("RecoilForce", 13, 1)

func initAmmo() {
	objClass := object.Class(1)

	genericDescriptions[objClass] = ammoClipGenerics
}
