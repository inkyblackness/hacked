package object

import "fmt"

// Class describes a general category of objects.
type Class byte

// String returns the textual representation.
func (c Class) String() string {
	if int(c) >= len(classNames) {
		return fmt.Sprintf("Unknown%02X", int(c))
	}
	return classNames[c]
}

// Object classes constants.
const (
	ClassGun        Class = 0
	ClassAmmo       Class = 1
	ClassPhysics    Class = 2
	ClassGrenade    Class = 3
	ClassDrug       Class = 4
	ClassHardware   Class = 5
	ClassSoftware   Class = 6
	ClassBigStuff   Class = 7
	ClassSmallStuff Class = 8
	ClassFixture    Class = 9
	ClassDoor       Class = 10
	ClassAnimating  Class = 11
	ClassTrap       Class = 12
	ClassContainer  Class = 13
	ClassCritter    Class = 14
)

var classNames = []string{
	"Gun",
	"Ammo",
	"Physics",
	"Grenade",
	"Drug",
	"Hardware",
	"Software",
	"BigStuff",
	"SmallStuff",
	"Fixture",
	"Door",
	"Animating",
	"Trap",
	"Container",
	"Critter",
}

// Classes returns a list of all classes.
func Classes() []Class {
	var classes [ClassCount]Class
	for class := Class(0); class < ClassCount; class++ {
		classes[class] = class
	}
	return classes[:]
}

// Subclass divides an object class.
type Subclass byte

// Type describes one specific object.
type Type byte

// Triple identifies one specific object by its full coordinate.
type Triple struct {
	Class    Class
	Subclass Subclass
	Type     Type
}

// TripleFrom returns a Triple instance with given values as coordinates.
func TripleFrom(class, subclass, objType int) Triple {
	return Triple{
		Class:    Class(class),
		Subclass: Subclass(subclass),
		Type:     Type(objType),
	}
}

// TripleFromInt returns a Triple instance based on a single integer value.
func TripleFromInt(value int) Triple {
	return Triple{
		Class:    Class((value >> 16) & 0xFF),
		Subclass: Subclass((value >> 8) & 0xFF),
		Type:     Type((value >> 0) & 0xFF),
	}
}

// String returns the textual representation of the triple as "cl/s/ty" string.
func (triple Triple) String() string {
	return fmt.Sprintf("%2d/%d/%2d", triple.Class, triple.Subclass, triple.Type)
}

// Int returns the triple as a single integer value.
func (triple Triple) Int() int {
	value := int(triple.Class) << 16
	value |= int(triple.Subclass) << 8
	value |= int(triple.Type) << 0
	return value
}
