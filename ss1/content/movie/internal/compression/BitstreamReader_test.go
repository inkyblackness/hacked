package compression_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"

	"github.com/stretchr/testify/assert"
)

func TestBitstreamReadPanicsForMoreThan32Bits(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0x11, 0x22, 0x33, 0x44})

	assert.Panicsf(t, func() { reader.Read(33) }, "Limit of bit count: 32")
}

func TestBitstreamReadReturnsValueOfRequestedBitSize(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xAF})

	result := reader.Read(3)

	assert.Equal(t, uint32(5), result)
}

func TestBitstreamRepeatedReadReturnsSameValue(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xAF})

	result1 := reader.Read(3)
	result2 := reader.Read(3)

	assert.Equal(t, uint32(5), result1)
	assert.Equal(t, result1, result2)
}

func TestBitstreamReadReturnsZeroesForBitsBeyondEndA(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xAF})

	result := reader.Read(9)

	assert.Equal(t, uint32(0x15E), result)
}

func TestBitstreamReadReturnsZeroesForBitsBeyondEndB(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xBF})

	result := reader.Read(32)

	assert.Equal(t, uint32(0xBF000000), result)
}

func TestBitstreamAdvancePanicsForNegativeValues(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0x11, 0x22, 0x33, 0x44})

	assert.Panicsf(t, func() { reader.Advance(-10) }, "Can only advance forward")
}

func TestBitstreamAdvanceLetsReadFurtherBits(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xAF})

	reader.Advance(2)
	result := reader.Read(4)

	assert.Equal(t, uint32(0x0B), result)
}

func TestBitstreamAdvanceToEndIsPossible(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xAF})

	reader.Advance(8)
	result := reader.Read(8)

	assert.Equal(t, uint32(0), result)
}

func TestBitstreamInternalBufferDoesNotLoseData(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0x7F, 0xFF, 0xFF, 0xFF, 0x80})

	reader.Advance(1)
	result := reader.Read(32)

	assert.Equal(t, uint32(0xFFFFFFFF), result)
}

func TestBitstreamAdvanceCanJumpToLastBit(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0x00, 0x00, 0x01})

	reader.Advance(23)
	result := reader.Read(1)

	assert.Equal(t, uint32(1), result)
}

func TestBitstreamReadAdvanceBeyondFirstRead(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xFF, 0x00, 0xFA})

	reader.Read(10)
	reader.Advance(20)
	result := reader.Read(4)

	assert.Equal(t, uint32(0x0A), result)
}

func TestBitstreamReadAdvanceWithinFirstRead(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xFF, 0x00, 0xFA})

	reader.Read(10)
	reader.Advance(4)
	result := reader.Read(20)

	assert.Equal(t, uint32(0xF00FA), result)
}

func TestBitstreamReadOfZeroBitsIsPossibleMidStream(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xFF})

	reader.Read(4)
	result := reader.Read(0)

	assert.Equal(t, uint32(0), result)
}

func TestBitstreamReadOfZeroBitsIsPossibleAtEnd(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xFF})

	reader.Advance(8)
	result := reader.Read(0)

	assert.Equal(t, uint32(0), result)
}

func TestBitstreamReadOfZeroBitsIsPossibleWithEmptySource(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{})

	result := reader.Read(0)

	assert.Equal(t, uint32(0), result)
}

func TestBitstreamAdvanceBeyondEndIsPossible(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xFF})

	reader.Advance(9)
	result := reader.Read(32)

	assert.Equal(t, uint32(0), result)
}

func TestBitstreamExhaustedReturnsFalseForAvailableBits(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xFF})

	result := reader.Exhausted()

	assert.False(t, result)
}

func TestBitstreamExhaustedReturnsFalseForStillOneAvailableBit(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xFF})

	reader.Advance(7)
	result := reader.Exhausted()

	assert.False(t, result)
}

func TestBitstreamExhaustedReturnsTrueAfterAdvancingToEnd(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xFF, 0xFF})

	reader.Advance(16)
	result := reader.Exhausted()

	assert.True(t, result)
}

func TestBitstreamExhaustedReturnsTrueAfterAdvancingBeyondEnd(t *testing.T) {
	reader := compression.NewBitstreamReader([]byte{0xFF, 0xFF})

	reader.Advance(30)
	result := reader.Exhausted()

	assert.True(t, result)
}
