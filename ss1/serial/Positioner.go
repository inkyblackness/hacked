package serial

// Positioner specifies a type that knows about a position which can be modified.
type Positioner interface {
	// CurPos gets the current position in the data.
	CurPos() uint32
	// SetCurPos sets the current position in the data.
	SetCurPos(offset uint32)
}
