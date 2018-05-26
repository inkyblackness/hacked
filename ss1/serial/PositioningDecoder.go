package serial

import "io"

// PositioningDecoder is a Decoder with positioning capabilities.
type PositioningDecoder struct {
	Decoder

	seeker io.Seeker
}

// NewPositioningDecoder creates a new PositiongingDecoder from given reader.
func NewPositioningDecoder(source io.ReadSeeker) *PositioningDecoder {
	return &PositioningDecoder{Decoder: Decoder{source: source, offset: 0}, seeker: source}
}

// CurPos returns the current decoding position, in bytes.
func (coder *PositioningDecoder) CurPos() uint32 {
	return coder.offset
}

// SetCurPos changes the encoding position to the specified absolute offset.
func (coder *PositioningDecoder) SetCurPos(offset uint32) {
	if coder.firstError != nil {
		return
	}
	_, coder.firstError = coder.seeker.Seek(int64(offset), io.SeekStart)
	if coder.firstError != nil {
		return
	}
	coder.offset = offset
}
