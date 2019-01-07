package objprop

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

var critterAttackInfo = interpreters.New().
	With("DamageType", 0x0000, 1).As(damageType).
	With("SpecialDamageType", 0x0001, 1).As(specialDamageType).
	With("DamageModifier", 0x0004, 2).As(interpreters.RangedValue(0, 500)).
	With("OffenceValue", 0x0006, 1).
	With("Penetration", 0x0007, 1).
	With("AttackMass", 0x0008, 1).
	With("AttackVelocity", 0x0009, 2).
	With("Accuracy", 0x000B, 1).
	With("AttackRange", 0x000C, 1).
	With("ReloadTime", 0x000D, 2).As(interpreters.RangedValue(0, 1000)).
	With("Projectile", 0x0011, 4).As(interpreters.SpecialValue("ObjectTriple"))

var critterGenerics = interpreters.New().
	Refining("PrimaryAttack", 0x0001, 21, critterAttackInfo, interpreters.Always).
	Refining("SecondaryAttack", 0x0016, 21, critterAttackInfo, interpreters.Always).
	With("Perception", 0x002B, 1).
	With("Defence", 0x002C, 1).
	With("ProjectileSourceHeightOffset", 0x002D, 1).As(interpreters.RangedValue(-128, 127)).
	With("Flags", 0x002E, 1).As(interpreters.Bitfield(map[uint32]string{
	0x01: "Flying"})).
	With("AnimationSpeed", 0x003B, 1).
	With("AttackSoundIndex", 0x003C, 1).
	With("NearSoundIndex", 0x003D, 1).
	With("HurtSoundIndex", 0x003E, 1).
	With("DeathSoundIndex", 0x003F, 1).
	With("NoticeSoundIndex", 0x0040, 1).
	With("Corpse", 0x0041, 4).As(interpreters.SpecialValue("ObjectTriple")).
	With("ViewCount", 0x0045, 1).
	With("SecondaryAttackProbability", 0x0046, 1).
	With("DisruptProbability", 0x0047, 1).
	With("TreasureType", 0x0048, 1).As(interpreters.EnumValue(map[uint32]string{
	0x00: "No treasure",
	0x01: "Humanoid",
	0x02: "Drone",
	0x03: "Assassin",
	0x04: "Warrior cyborg",
	0x05: "Flier bot",
	0x06: "Security 1 bot",
	0x07: "Exec bot",
	0x08: "Cyborg enforcer",
	0x09: "Security 2 bot",
	0x0A: "Elite cyborg",
	0x0B: "Standard corpse",
	0x0C: "Loot-oriented corpse",
	0x0D: "Electro-stuff treasure",
	0x0E: "Serv-bot"})).
	With("HitEffect", 0x0049, 1).As(interpreters.EnumValue(map[uint32]string{0: "meat", 1: "plant", 2: "robot", 3: "cyborg", 4: "generic"})).
	With("AttackKeyFrame", 0x004A, 1).As(interpreters.RangedValue(0, 10))

var cyberCritters = interpreters.New().
	With("Color0", 0, 1).
	With("Color1", 1, 1).
	With("Color2", 2, 1).
	With("AltColor0", 3, 1).
	With("AltColor1", 4, 1).
	With("AltColor2", 5, 1)

func initCritters() {
	objClass := object.Class(14)

	genericDescriptions[objClass] = critterGenerics

	setSpecific(objClass, 3, cyberCritters)
}
