package compression_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"

	"github.com/stretchr/testify/assert"
)

func TestMaskstreamReaderReadPanicsForMoreThan8Bytes(t *testing.T) {
	reader := compression.NewMaskstreamReader([]byte{})

	assert.Panicsf(t, func() { reader.Read(9) }, "Limit of byte count: 8")
}

func TestMaskstreamReaderReadPanicsForLessThan0Bytes(t *testing.T) {
	reader := compression.NewMaskstreamReader([]byte{})

	assert.Panicsf(t, func() { reader.Read(-1) }, "Minimum byte count: 0")
}

func TestMaskstreamReaderReadReturnsValueOfRequestedByteSize(t *testing.T) {
	reader := compression.NewMaskstreamReader([]byte{0xAF})

	result := reader.Read(1)

	assert.Equal(t, uint64(0xAF), result)
}

func TestMaskstreamReaderReadIntegerFromSourceInLittleEndianOrder(t *testing.T) {
	reader := compression.NewMaskstreamReader([]byte{0x11, 0x22})

	result := reader.Read(2)

	assert.Equal(t, uint64(0x2211), result)
}

func TestMaskstreamReaderReadCanProvideUpTo64BitValues(t *testing.T) {
	reader := compression.NewMaskstreamReader([]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88})

	result := reader.Read(8)

	assert.Equal(t, uint64(0x8877665544332211), result)
}

func TestMaskstreamReaderReadFillsMissingBytesWithZero(t *testing.T) {
	reader := compression.NewMaskstreamReader([]byte{0xAA, 0xBB})

	result := reader.Read(8)

	assert.Equal(t, uint64(0x00BBAA), result)
}

func TestMaskstreamReaderReadAdvancesCurrentPosition(t *testing.T) {
	reader := compression.NewMaskstreamReader([]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88})

	reader.Read(2)
	result := reader.Read(3)

	assert.Equal(t, uint64(0x554433), result)
}
