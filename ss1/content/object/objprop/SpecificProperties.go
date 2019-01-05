package objprop

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

// SpecificProperties returns an interpreter specific for the given object class and subclass.
func SpecificProperties(triple object.Triple, data []byte) *interpreters.Instance {
	desc := specificDescriptions[triple]
	if desc == nil {
		desc = specificDescriptions[object.TripleFrom(int(triple.Class), int(triple.Subclass), anyObjectType)]
	}
	if desc == nil {
		desc = interpreters.New()
	}
	return desc.For(data)
}
