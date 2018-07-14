package lvlobj

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

// ForCyberspace returns an interpreter instance that handles the level class
// data of the specified object - in cyberspace.
func ForCyberspace(triple object.Triple, data []byte) *interpreters.Instance {
	return cyberspaceEntries.specialize(int(triple.Class)).specialize(int(triple.Subclass)).specialize(int(triple.Type)).instance(data)
}

// CyberspaceExtra returns an interpreter instance that handles the level object extra
// data of the specified object - in cybperspace.
func CyberspaceExtra(triple object.Triple, data []byte) *interpreters.Instance {
	return cyberspaceExtras.specialize(int(triple.Class)).specialize(int(triple.Subclass)).specialize(int(triple.Type)).instance(data)
}
