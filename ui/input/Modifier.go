package input

// Modifier is a collection of currently pressed modifier keys.
type Modifier uint32

// Constants for common modifier.
const (
	ModNone    = Modifier(0x00000000)
	ModShift   = Modifier(0x00000001)
	ModControl = Modifier(0x00000002)
	ModAlt     = Modifier(0x00000004)
	ModSuper   = Modifier(0x00000008)
)

// With returns a combination of this modifier and the given one.
func (mod Modifier) With(other Modifier) Modifier {
	return Modifier(uint32(mod) | uint32(other))
}

// Without returns the result of this modifier without the given one.
func (mod Modifier) Without(other Modifier) Modifier {
	return Modifier(uint32(mod) & ^uint32(other))
}

// Has returns true if the given modifier is included in this set.
// This modifier can have more keys set than the requested and still return true.
func (mod Modifier) Has(other Modifier) bool {
	return (uint32(mod) | uint32(other)) == uint32(mod)
}
