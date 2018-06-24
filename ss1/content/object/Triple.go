package object

import "fmt"

// Class describes a general category of objects.
type Class byte

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

// String returns the textual representation of the triple as "cl/s/ty" string.
func (triple Triple) String() string {
	return fmt.Sprintf("%2d/%d/%2d", triple.Class, triple.Subclass, triple.Type)
}
