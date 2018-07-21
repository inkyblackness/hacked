package movie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeFromRawA(t *testing.T) {
	result := timeFromRaw(byte(4), uint16(0x8000))

	assert.Equal(t, float32(4.5), result)
}

func TestTimeFromRawB(t *testing.T) {
	result := timeFromRaw(byte(6), uint16(0))

	assert.Equal(t, float32(6.0), result)
}

func TestTimeToRawA(t *testing.T) {
	second, fraction := timeToRaw(7.25)

	assert.Equal(t, byte(7), second)
	assert.Equal(t, uint16(0x4000), fraction)
}

func TestTimeToRawB(t *testing.T) {
	second, fraction := timeToRaw(255.75)

	assert.Equal(t, byte(255), second)
	assert.Equal(t, uint16(0xC000), fraction)
}
