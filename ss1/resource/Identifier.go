package resource

import "fmt"

// ID represents an integer key of resources.
type ID uint16

// Value returns the numerical value of the identifier.
func (id ID) Value() uint16 {
	return uint16(id)
}

// Plus adds the given offset and returns the resulting ID.
func (id ID) Plus(offset int) ID {
	return ID(int(id) + offset)
}

func (id ID) String() string {
	return fmt.Sprintf("%04X", uint16(id))
}
