package objprop

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

var basicWeapon = interpreters.New().
	With("Damage", 0, 2).As(interpreters.RangedValue(0, 0x7FFF)).
	With("OffenceValue", 2, 1).
	With("DamageType", 3, 1).As(damageType).
	With("SpecialDamageType", 4, 1).As(specialDamageType).
	With("ArmorPenetration", 7, 1)

var weaponGenerics = interpreters.New().
	With("FireRate", 0, 1).
	With("AmmoType", 1, 1).As(interpreters.Bitfield(map[uint32]string{
	0x01: "AmmoType0",
	0x02: "AmmoType1",
	0x04: "AmmoType2",
	0x08: "AmmoType3",
	0xF0: "AmmoSubclass"}))

var projectileWeapons = interpreters.New().
	Refining("BasicWeapon", 0, 8, basicWeapon, interpreters.Always).
	With("ProjectileTravelSpeed", 8, 1).
	With("Projectile", 9, 4).As(interpreters.SpecialValue("ObjectTriple")).
	With("AttackMass", 13, 1).
	With("AttackSpeed", 14, 2).As(interpreters.RangedValue(-10000, +10000))

var meleeWeapons = interpreters.New().
	Refining("BasicWeapon", 0, 8, basicWeapon, interpreters.Always).
	With("EnergyUsage", 8, 1).
	With("AttackMass", 9, 1).
	With("Range", 10, 1).
	With("AttackSpeed", 11, 2).As(interpreters.RangedValue(-10000, +10000))

var energyBeamWeapons = interpreters.New().
	Refining("BasicWeapon", 0, 8, basicWeapon, interpreters.Always).
	With("MaxCharge", 8, 1).
	With("AttackMass", 9, 1).
	With("Range", 10, 1).
	With("AttackSpeed", 11, 2).As(interpreters.RangedValue(-10000, +10000))

var energyProjectileWeapons = interpreters.New().
	Refining("BasicWeapon", 0, 8, basicWeapon, interpreters.Always).
	With("MaxCharge", 8, 1).
	With("AttackMass", 9, 1).
	With("AttackSpeed", 10, 2).As(interpreters.RangedValue(-10000, +10000)).
	With("ProjectileTravelSpeed", 12, 1).
	With("Projectile", 13, 4).As(interpreters.SpecialValue("ObjectTriple")).
	With("Flags", 17, 1).As(interpreters.Bitfield(map[uint32]string{
	0x02: "Stun"}))

func initGuns() {
	objClass := object.Class(0)

	genericDescriptions[objClass] = weaponGenerics

	setSpecific(objClass, 2, projectileWeapons)
	setSpecific(objClass, 3, meleeWeapons)
	setSpecific(objClass, 4, energyBeamWeapons)
	setSpecific(objClass, 5, energyProjectileWeapons)
}
