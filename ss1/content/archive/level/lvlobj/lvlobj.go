package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

var realWorldEntries *interpreterEntry
var cyberspaceEntries *interpreterEntry

var realWorldExtras *interpreterEntry
var cyberspaceExtras *interpreterEntry

var extraIced = interpreters.New().
	With("ICE-presence", 1, 1).
	With("ICE-level", 3, 1)

var extraIcedPanels = interpreters.New().
	With("PanelName", 0, 1).
	With("ICE-presence", 1, 1).
	With("ICE-level", 3, 1)

var extraPanels = interpreters.New().
	With("PanelName", 0, 1)

func init() {

	software := initSoftware()
	traps := initTraps()
	critters := initCritters()

	realWorldEntries = newInterpreterEntry(interpreters.New())
	realWorldEntries.set(int(object.ClassGun), initGuns())
	realWorldEntries.set(int(object.ClassAmmo), newInterpreterEntry(interpreters.New())) // have no data
	realWorldEntries.set(int(object.ClassPhysics), newInterpreterEntry(basePhysics))
	realWorldEntries.set(int(object.ClassGrenade), initGrenades())
	realWorldEntries.set(int(object.ClassDrug), newInterpreterEntry(interpreters.New())) // have no data
	realWorldEntries.set(int(object.ClassHardware), newInterpreterEntry(baseHardware))
	realWorldEntries.set(int(object.ClassSoftware), software)
	realWorldEntries.set(int(object.ClassBigStuff), initBigStuff())
	realWorldEntries.set(int(object.ClassSmallStuff), initSmallStuff())
	realWorldEntries.set(int(object.ClassFixture), initFixtures())
	realWorldEntries.set(int(object.ClassDoor), initDoors())
	realWorldEntries.set(int(object.ClassAnimating), newInterpreterEntry(baseAnimation))
	realWorldEntries.set(int(object.ClassTrap), traps)
	realWorldEntries.set(int(object.ClassContainer), initContainers())
	realWorldEntries.set(int(object.ClassCritter), critters)

	cyberspaceEntries = newInterpreterEntry(interpreters.New())
	cyberspaceEntries.set(int(object.ClassSoftware), software)
	cyberspaceEntries.set(int(object.ClassBigStuff), initCyberspaceBigStuff())
	cyberspaceEntries.set(int(object.ClassSmallStuff), initCyberspaceSmallStuff())
	cyberspaceEntries.set(int(object.ClassFixture), initCyberspaceFixtures())
	cyberspaceEntries.set(int(object.ClassTrap), traps)
	cyberspaceEntries.set(int(object.ClassCritter), critters)

	realWorldExtras = newInterpreterEntry(interpreters.New())
	realWorldExtras.set(9, newInterpreterLeaf(extraPanels))
	cyberspaceExtras = newInterpreterEntry(interpreters.New())
	cyberspaceExtras.set(6, newInterpreterLeaf(extraIced))
	cyberspaceExtras.set(7, newInterpreterLeaf(extraIced))
	cyberspaceExtras.set(8, newInterpreterLeaf(extraIced))
	cyberspaceExtras.set(9, newInterpreterLeaf(extraIcedPanels))
}
