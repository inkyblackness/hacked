package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

var baseGun = interpreters.New()

var energyWeapon = baseGun.
	With("Charge", 0, 1).As(interpreters.RangedValue(0, 255)).
	With("Temperature", 1, 1).As(interpreters.RangedValue(0, 255))

var projectileWeapon = baseGun.
	With("AmmoType", 0, 1).As(interpreters.EnumValue(map[uint32]string{0: "Standard", 1: "Special"})).
	With("AmmoCount", 1, 1).As(interpreters.RangedValue(0, 255))

func initGuns() interpreterRetriever {
	projectileWeapons := newInterpreterLeaf(projectileWeapon)
	energyWeapons := newInterpreterLeaf(energyWeapon)

	class := newInterpreterEntry(baseGun)
	class.set(0, projectileWeapons)
	class.set(1, projectileWeapons)
	class.set(2, projectileWeapons)

	class.set(4, energyWeapons)
	class.set(5, energyWeapons)

	return class
}
