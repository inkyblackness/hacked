package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

var baseItem = interpreters.New()

var paperItem = baseItem.
	With("PaperId", 2, 1)

var briefcaseItem = baseItem.
	With("ObjectIndex1", 2, 2).As(interpreters.ObjectID()).
	With("ObjectIndex2", 4, 2).As(interpreters.ObjectID()).
	With("ObjectIndex3", 6, 2).As(interpreters.ObjectID()).
	With("ObjectIndex4", 8, 2).As(interpreters.ObjectID())

var corpseItem = baseItem.
	With("Unknown0000", 0, 2).As(interpreters.SpecialValue("Unknown")).
	With("ObjectIndex1", 2, 2).As(interpreters.ObjectID()).
	With("ObjectIndex2", 4, 2).As(interpreters.ObjectID()).
	With("ObjectIndex3", 6, 2).As(interpreters.ObjectID()).
	With("ObjectIndex4", 8, 2).As(interpreters.ObjectID())

var severedHeadItem = baseItem.
	With("ImageIndex", 2, 1)

var accessLevelMasks = map[uint32]string{
	0x00000001: "None",
	0x00000002: "Generic1",
	0x00000004: "Generic2",
	0x00000008: "Generic3",
	0x00000010: "Generic4",
	0x00000020: "Generic5",
	0x00000040: "Generic6",
	0x00000080: "Generic7",

	0x00000100: "Group1",
	0x00000200: "Group2",
	0x00000400: "Group3",
	0x00000800: "Group4",
	0x00001000: "Group5",
	0x00002000: "Group6",
	0x00004000: "Group7",
	0x00008000: "Group8",

	0x00010000: "Group9",
	0x00020000: "Group10",
	0x00040000: "Group11",
	0x00080000: "Group12",
	0x00100000: "Group13",
	0x00200000: "Group14",
	0x00400000: "Group15",
	0x00800000: "Group16",

	0x01000000: "Personal1",
	0x02000000: "Personal2",
	0x04000000: "Personal3",
	0x08000000: "Personal4",
	0x10000000: "Personal5",
	0x20000000: "Personal6",
	0x40000000: "Personal7"}

var accessMaskDescription = interpreters.Bitfield(accessLevelMasks)

var accessCardItem = baseItem.
	With("Ignored0000", 0, 2).As(interpreters.SpecialValue("Ignored")).
	With("AccessMask", 2, 4).As(accessMaskDescription)

var securityIDModuleItem = baseItem.
	With("AccessMask", 2, 4).As(accessMaskDescription)

var cyberInfoNodeItem = baseItem.
	With("TextIndex", 2, 1)

var cyberRestorative = baseItem.
	With("RestorationAmount", 2, 1)

var cyberDefenseMine = baseItem.
	With("DamageAmount", 2, 1)

var cyberBarricade = baseItem.
	With("Size", 2, 1).
	With("Height", 3, 1).
	With("Color", 6, 1)

func initSmallStuff() interpreterRetriever {

	junk := newInterpreterEntry(baseItem)
	junk.set(2, newInterpreterLeaf(paperItem))
	junk.set(7, newInterpreterLeaf(briefcaseItem))

	dead := newInterpreterEntry(baseItem)
	corpses := newInterpreterLeaf(corpseItem)
	severedHeads := newInterpreterLeaf(severedHeadItem)
	dead.set(0, corpses)
	dead.set(1, corpses)
	dead.set(2, corpses)
	dead.set(3, corpses)
	dead.set(4, corpses)
	dead.set(5, corpses)
	dead.set(6, corpses)
	dead.set(7, corpses)
	dead.set(13, severedHeads)
	dead.set(14, severedHeads)

	class := newInterpreterEntry(baseItem)
	class.set(0, junk)
	class.set(2, dead)
	class.set(4, newInterpreterLeaf(accessCardItem))

	return class
}

func initCyberspaceSmallStuff() interpreterRetriever {
	cyberspaceItems := newInterpreterEntry(baseItem)
	infoNodes := newInterpreterLeaf(cyberInfoNodeItem)
	cyberspaceItems.set(1, newInterpreterLeaf(cyberRestorative))
	cyberspaceItems.set(2, newInterpreterLeaf(cyberDefenseMine))
	cyberspaceItems.set(3, newInterpreterLeaf(securityIDModuleItem))
	cyberspaceItems.set(6, infoNodes)
	cyberspaceItems.set(8, infoNodes)
	cyberspaceItems.set(9, newInterpreterLeaf(cyberBarricade))

	class := newInterpreterEntry(baseItem)
	class.set(5, cyberspaceItems)

	return class
}
