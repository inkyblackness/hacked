package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

var baseGrenade = interpreters.New().
	With("Unknown0000", 0, 2).
	With("State", 2, 2).As(interpreters.EnumValue(map[uint32]string{0: "Inert", 1: "Thrown Live", 5: "Landed Live"})).
	With("TimerTime", 4, 2).As(interpreters.RangedValue(0, 32767))

func initGrenades() interpreterRetriever {
	timedExplosives := newInterpreterEntry(baseGrenade)

	timedExplosives.set(2, newInterpreterLeaf(interpreters.New())) // Object explosion - not encountered

	class := newInterpreterEntry(baseGrenade)
	class.set(1, timedExplosives)

	return class
}
