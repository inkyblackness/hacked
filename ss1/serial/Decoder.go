package serial

import (
	"encoding/binary"
	"io"
)

// Decoder is for decoding from a reader.
type Decoder struct {
	source     io.Reader
	firstError error
	offset     uint32
}

// NewDecoder creates a new Decoder from given source.
func NewDecoder(source io.Reader) *Decoder {
	return &Decoder{source: source}
}

// FirstError returns the error this Decoder encountered the first time.
func (coder *Decoder) FirstError() error {
	return coder.firstError
}

// Code serializes the given value in little endian format using binary.Read().
func (coder *Decoder) Code(value interface{}) {
	if coder.firstError != nil {
		return
	}
	if codable, isCodable := value.(Codable); isCodable {
		codable.Code(coder)
	} else {
		coder.firstError = binary.Read(coder, binary.LittleEndian, value)
	}
}

// Read reads the next bytes from the underlying source.
func (coder *Decoder) Read(data []byte) (read int, err error) {
	if coder.firstError != nil {
		return 0, coder.firstError
	}
	read, err = coder.source.Read(data)
	coder.offset += uint32(read)

	isErrEOF := err == io.EOF
	expectedAmountReturned := read == len(data)
	errorCanBeIgnored := isErrEOF && expectedAmountReturned
	if (coder.firstError == nil) && !errorCanBeIgnored {
		coder.firstError = err
	}
	return
}
