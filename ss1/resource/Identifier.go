package resource

import "fmt"

// ID represents an integer key of resources.
type ID uint16

// Value returns the numerical value of the identifier.
func (id ID) Value() uint16 {
	return uint16(id)
}

func (id ID) String() string {
	return fmt.Sprintf("%04X", uint16(id))
}
