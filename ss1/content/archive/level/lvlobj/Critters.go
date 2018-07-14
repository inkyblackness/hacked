package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

var baseCritter = interpreters.New().
	With("PendingZRotation", 0, 4).
	With("ForwardVelocityFraction", 4, 2).
	With("ForwardVelocity", 6, 2).
	With("Hastiness", 0x0A, 2).
	With("StateTimeout", 0x0C, 2).As(interpreters.RangedValue(0, 300)).
	With("Unknown000E", 0x0E, 2).As(interpreters.SpecialValue("Unknown")).
	With("Unknown0010", 0x10, 2).As(interpreters.SpecialValue("Unknown")).
	With("RoamingState", 0x14, 1).
	With("PrimaryState", 0x15, 1).As(interpreters.EnumValue(map[uint32]string{
	0: "docile",
	1: "cautious",
	2: "hostile",
	3: "cautious (?)",
	4: "attacking",
	5: "sleeping",
	6: "tranquilized",
	7: "confused"})).
	With("SecondaryState", 0x16, 1).
	With("TertiaryState", 0x17, 1).
	With("AICoordinates0018", 0x18, 2).As(interpreters.SpecialValue("Unknown")).
	With("AICoordinates001A", 0x1A, 2).As(interpreters.SpecialValue("Unknown")).
	With("AICoordinates001C", 0x1C, 2).As(interpreters.SpecialValue("Unknown")).
	With("Always255", 0x1E, 1).As(interpreters.RangedValue(0, 255)).
	With("AICoordinates001F", 0x1F, 1).As(interpreters.SpecialValue("Unknown")).
	With("LootObjectIndex1", 0x20, 2).As(interpreters.ObjectID()).
	With("LootObjectIndex2", 0x22, 2).As(interpreters.ObjectID()).
	With("Unknown0026", 0x26, 2).As(interpreters.SpecialValue("Unknown"))

func initCritters() interpreterRetriever {
	class := newInterpreterEntry(baseCritter)

	return class
}
