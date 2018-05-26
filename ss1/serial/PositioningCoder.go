package serial

// PositioningCoder is a coder that also knows about positioning.
// Any error state regarding repositioning will be provided via FirstError().
// SetCurPos will do nothing if the coder is already in error state.
type PositioningCoder interface {
	Positioner
	Coder
}
