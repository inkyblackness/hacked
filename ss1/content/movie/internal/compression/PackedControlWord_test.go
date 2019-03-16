package compression_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
)

func TestPackedControlWordTimesReturnsHighestByte(t *testing.T) {
	packed := compression.PackedControlWord(0xFF706050)

	assert.Equal(t, int(0xFF), packed.Times())
}

func TestPackedControlWordValueReturnsTheControlWord(t *testing.T) {
	packed := compression.PackedControlWord(0xFF706050)

	assert.Equal(t, compression.ControlWord(0x00706050), packed.Value())
}
