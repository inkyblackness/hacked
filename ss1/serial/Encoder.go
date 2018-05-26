package serial

import (
	"encoding/binary"
	"io"
)

// Encoder implements the Coder interface to write to a writer.
// It also implements the Writer interface.
type Encoder struct {
	target     io.Writer
	firstError error
	offset     uint32
}

// NewEncoder creates and returns a fresh Encoder.
func NewEncoder(target io.Writer) *Encoder {
	return &Encoder{target: target}
}

// FirstError returns the error this Encoder encountered the first time.
func (coder *Encoder) FirstError() error {
	return coder.firstError
}

// Code serializes the given value in little endian format using binary.Write().
func (coder *Encoder) Code(value interface{}) {
	if coder.firstError != nil {
		return
	}
	if codable, isCodable := value.(Codable); isCodable {
		codable.Code(coder)
	} else {
		coder.firstError = binary.Write(coder, binary.LittleEndian, value)
	}
}

// Write serializes the given data to the contained writer.
// If the encoder is in the error state, the result of FirstError() is
// returned and nothing else is done.
func (coder *Encoder) Write(data []byte) (written int, err error) {
	if coder.firstError != nil {
		return 0, coder.firstError
	}
	written, err = coder.target.Write(data)
	coder.offset += uint32(written)
	if coder.firstError == nil {
		coder.firstError = err
	}
	return
}
