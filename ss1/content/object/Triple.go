package object

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
