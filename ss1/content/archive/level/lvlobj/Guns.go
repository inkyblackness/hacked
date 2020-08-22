package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

const (
	// EnergyWeaponMaxTemperature is the maximum value possible for heat.
	EnergyWeaponMaxTemperature = 100
	// EnergyWeaponMaxEnergy is the maximum value for energy.
	// The original source does not explicitly state this value, it has been determined by looking at
	// a savegame with the energy set to max. A closely related number in the source is
	// in the macro "mfd_charge_units_per_pixel", which scales to 69.
	EnergyWeaponMaxEnergy = 59
)

var baseGun = interpreters.New()

var energyWeapon = baseGun.
	With("Charge", 0, 1).As(interpreters.Bitfield(
	map[uint32]string{
		0x80: "Overload",
		0x7F: "Energy",
	})).
	With("Temperature", 1, 1).As(interpreters.RangedValue(0, EnergyWeaponMaxTemperature))

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
