package serial

import (
	"errors"
	"io"
)

// PositioningEncoder is an Encoder with positioning capabilities.
type PositioningEncoder struct {
	Encoder

	seeker io.Seeker
}

// NewPositioningEncoder returns an Encoder that also implements the Positioner interface.
// The new Encoder starts with its zero position at the current position in the writer.
func NewPositioningEncoder(target io.WriteSeeker) *PositioningEncoder {
	return &PositioningEncoder{Encoder: Encoder{target: target}, seeker: target}
}

// CurPos returns the current encoding position, in bytes.
func (coder *PositioningEncoder) CurPos() uint32 {
	return coder.offset
}

// SetCurPos changes the encoding position to the specified absolute offset.
func (coder *PositioningEncoder) SetCurPos(offset uint32) {
	if coder.firstError != nil {
		return
	}
	_, coder.firstError = coder.seeker.Seek(int64(offset)-int64(coder.offset), io.SeekCurrent)
	if coder.firstError != nil {
		return
	}
	coder.offset = offset
}

var errInvalidWhence = errors.New("seek: invalid whence")
var errSeekBeforeStart = errors.New("seek: seeking before start")

// Seek repositions the current encoding offset.
// This implementation does not support a whence value of io.SeekEnd.
func (coder *PositioningEncoder) Seek(offset int64, whence int) (int64, error) {
	var newPosition int64
	switch whence {
	default:
		return int64(coder.offset), errInvalidWhence
	case io.SeekStart:
		newPosition = offset
	case io.SeekCurrent:
		newPosition = int64(coder.offset) + offset
	}
	if newPosition < 0 {
		return int64(coder.offset), errSeekBeforeStart
	}

	coder.SetCurPos(uint32(newPosition))
	return int64(coder.offset), coder.firstError
}
