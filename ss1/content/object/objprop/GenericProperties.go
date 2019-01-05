package objprop

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
)

// GenericProperties returns an interpreter specific for the given object class.
func GenericProperties(objClass object.Class, data []byte) *interpreters.Instance {
	desc := genericDescriptions[objClass]
	if desc == nil {
		desc = interpreters.New()
	}
	return desc.For(data)
}
