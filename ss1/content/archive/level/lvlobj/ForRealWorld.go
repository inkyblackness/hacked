package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

// ForRealWorld returns an interpreter instance that handles the level class
// data of the specified object - in real world.
func ForRealWorld(triple object.Triple, data []byte) *interpreters.Instance {
	return realWorldEntries.specialize(int(triple.Class)).specialize(int(triple.Subclass)).specialize(int(triple.Type)).instance(data)
}

// RealWorldExtra returns an interpreter instance that handles the level object extra
// data of the specified object - in real world.
func RealWorldExtra(triple object.Triple, data []byte) *interpreters.Instance {
	return realWorldExtras.specialize(int(triple.Class)).specialize(int(triple.Subclass)).specialize(int(triple.Type)).instance(data)
}
