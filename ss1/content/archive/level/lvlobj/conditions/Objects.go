package conditions

import (
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

var objectType = interpreters.New().
	With("ObjectType", 0, 3).As(interpreters.SpecialValue("ObjectTriple"))

var objectID = interpreters.New().
	With("ObjectID", 0, 2).As(interpreters.ObjectID())

// ObjectType returns a condition description for object types.
func ObjectType() *interpreters.Description {
	return objectType
}

// ObjectID returns a condition description for object indices.
func ObjectID() *interpreters.Description {
	return objectID
}
