package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

var realWorldEntries *interpreterEntry
var cyberspaceEntries *interpreterEntry

var realWorldExtras *interpreterEntry
var cyberspaceExtras *interpreterEntry

var extraDefault = interpreters.New().
	With("CurrentFrame", 1, 1).
	With("TimeRemainder", 2, 1)

var extraIced = interpreters.New().
	With("ICE-presence", 1, 1).
	With("ICE-level", 3, 1)

var extraIcedFixtures = interpreters.New().
	With("PanelName", 0, 1).
	With("ICE-presence", 1, 1).
	With("ICE-level", 3, 1)

var extraFixtures = interpreters.New().
	With("PanelName", 0, 1)

var extraSurfaces = interpreters.New().
	With("Index", 1, 1)

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

	realWorldExtras = newInterpreterEntry(extraDefault)
	realWorldExtras.set(int(object.ClassFixture), newInterpreterLeaf(extraFixtures))
	extraBigStuff := newInterpreterEntry(extraDefault)
	simpleSurfaces := newInterpreterLeaf(extraSurfaces)
	surfaces := newInterpreterEntry(extraDefault)
	surfaces.set(1, simpleSurfaces)
	surfaces.set(4, simpleSurfaces)
	extraBigStuff.set(2, surfaces)
	realWorldExtras.set(int(object.ClassBigStuff), extraBigStuff)
	cyberspaceExtras = newInterpreterEntry(interpreters.New())
	cyberspaceExtras.set(int(object.ClassSoftware), newInterpreterLeaf(extraIced))
	cyberspaceExtras.set(int(object.ClassBigStuff), newInterpreterLeaf(extraIced))
	cyberspaceExtras.set(int(object.ClassSmallStuff), newInterpreterLeaf(extraIced))
	cyberspaceExtras.set(int(object.ClassFixture), newInterpreterLeaf(extraIcedFixtures))
}
