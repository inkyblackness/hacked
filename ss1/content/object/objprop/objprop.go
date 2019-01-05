package objprop

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

var genericDescriptions map[object.Class]*interpreters.Description
var specificDescriptions map[object.Triple]*interpreters.Description

const anyObjectType = 0xFF

var damageType = interpreters.Bitfield(map[uint32]string{
	0x01: "Explosion",
	0x02: "Energy",
	0x04: "Magnetic",
	0x08: "Radiation",
	0x10: "Gas",
	0x20: "Tranquilizer",
	0x40: "Needle",
	0x80: "Bio",
})

var specialDamageType = interpreters.Bitfield(map[uint32]string{
	0x0F: "Primary (double)",
	0xF0: "Super (quad)",
})

func init() {
	genericDescriptions = make(map[object.Class]*interpreters.Description)
	specificDescriptions = make(map[object.Triple]*interpreters.Description)

	initGuns()
	initAmmo()
	initPhysics()
	initGrenades()
	initSmallStuff()
	initAnimating()
	initCritters()
}

func setSpecific(objClass object.Class, objSubclass int, desc *interpreters.Description) {
	specificDescriptions[object.TripleFrom(int(objClass), objSubclass, anyObjectType)] = desc
}

func setSpecificByType(objClass object.Class, objSubclass int, objType int, desc *interpreters.Description) {
	specificDescriptions[object.TripleFrom(int(objClass), objSubclass, objType)] = desc
}
