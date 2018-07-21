package wav

import (
	"encoding/binary"
)

type waveFormatType uint16

const (
	waveFormatTypePcm = 1
)

type waveFormat struct {
	FormatType     waveFormatType
	Channels       uint16
	SamplesPerSec  uint32
	AvgBytesPerSec uint32
	BlockAlign     uint16
}

type waveFormatExtension struct {
	BitsPerSample uint16
	ExtensionSize uint16
}

type formatHeader struct {
	base      waveFormat
	extension waveFormatExtension
}

func (header formatHeader) size() uint32 {
	return uint32(binary.Size(&header.base) + binary.Size(&header.extension))
}
