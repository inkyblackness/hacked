package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

var baseContainer = interpreters.New()

var standardContainer = baseContainer.
	With("ObjectID1", 0, 2).As(interpreters.ObjectID()).
	With("ObjectID2", 2, 2).As(interpreters.ObjectID()).
	With("ObjectID3", 4, 2).As(interpreters.ObjectID()).
	With("ObjectID4", 6, 2).As(interpreters.ObjectID())

var crate = standardContainer.
	With("Width", 8, 1).
	With("Depth", 9, 1).
	With("Height", 10, 1).
	With("TopBottomTexture", 11, 1).
	With("SideTexture", 12, 1)

func initContainers() interpreterRetriever {
	standardContainers := newInterpreterLeaf(standardContainer)
	crates := newInterpreterLeaf(crate)

	class := newInterpreterEntry(baseContainer)
	class.set(0, crates)
	class.set(1, standardContainers)
	class.set(2, standardContainers)
	class.set(3, standardContainers)
	class.set(4, standardContainers)
	class.set(5, standardContainers)
	class.set(6, standardContainers)

	return class
}
