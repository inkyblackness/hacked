package compression

import (
	"fmt"

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

type tooManyWordsError struct {
	Total               uint32
	MaximumDirect       uint32
	BitstreamIndexLimit uint32
}

func (err tooManyWordsError) Error() string {
	return fmt.Sprintf("too many words: total=%v, maxDirect=%v, indexLimit=%v", err.Total, err.MaximumDirect, err.BitstreamIndexLimit)
}

type wordIndexOutOfRangeError struct {
	Index     uint32
	Available uint32
}

func (err wordIndexOutOfRangeError) Error() string {
	return fmt.Sprintf("control word index out of range: %v/%v", err.Index, err.Available)
}

type paletteLookupTooBigError struct {
	Size int
}

func (err paletteLookupTooBigError) Error() string {
	return fmt.Sprintf("palette lookup is too big: %vB", err.Size)
}
