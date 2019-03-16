package compression

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControlWordPackedReturnsPackedControlWord(t *testing.T) {
	word := ControlWord(0x00FFAA55)

	assert.Equal(t, PackedControlWord(0x80FFAA55), word.Packed(128))
}

func TestControlWordPackedClearsHighestByteBeforeSettingCount(t *testing.T) {
	word := ControlWord(0x55FFAA55)

	assert.Equal(t, PackedControlWord(0xC0FFAA55), word.Packed(0xC0))
}

func TestControlWordCountReturnsBits20To23(t *testing.T) {
	word := ControlWord(0x00F00000)

	assert.Equal(t, int(15), word.Count())
}

func TestControlWordIsLongOffsetReturnsTrueForCount0(t *testing.T) {
	word := ControlWord(0x00000000)

	assert.True(t, word.IsLongOffset())
}

func TestControlWordIsLongOffsetReturnsFalseForCountGreater0(t *testing.T) {
	word := ControlWord(0x00300000)

	assert.False(t, word.IsLongOffset())
}

func TestControlWordLongOffsetReturnsBits00To19(t *testing.T) {
	word := ControlWord(0xFFFA6665)

	assert.Equal(t, uint32(0xA6665), word.LongOffset())
}

func TestControlWordParameterReturnsBits00To16(t *testing.T) {
	word := ControlWord(0xFFFF1665)

	assert.Equal(t, uint32(0x11665), word.Parameter())
}

func TestControlWordTypeReturnsBits17To19(t *testing.T) {
	word := ControlWord(0xFFFAFFFF)

	assert.Equal(t, ControlType(5), word.Type())
}
