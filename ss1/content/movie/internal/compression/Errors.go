package compression

import (
	"github.com/inkyblackness/hacked/ss1"
)

const (
	// ErrFormat indicates errors in the compression format.
	ErrFormat ss1.StringError = "format error"

	errCannotRepeatWordOnFirstTileOfRow ss1.StringError = "cannot repeat word on first tile of row"
	errUnknownControl                   ss1.StringError = "found unknown control"
	errTooManyLongOffsets               ss1.StringError = "too many long offsets"
	errInvalidByteCount                 ss1.StringError = "invalid byte count"
	errInvalidFrameSize                 ss1.StringError = "invalid frame size"
)
