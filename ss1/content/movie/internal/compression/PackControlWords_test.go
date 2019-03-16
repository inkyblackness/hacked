package compression_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inkyblackness/hacked/ss1/content/movie/internal/compression"
)

func TestControlWordsEmptyArrayResultsInZeroCount(t *testing.T) {
	data := compression.PackControlWords(nil)

	assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00}, data)
}

func TestControlWordsCountIsInLengthTimesThree(t *testing.T) {
	data := compression.PackControlWords([]compression.ControlWord{
		compression.ControlWord(0),
		compression.ControlWord(1),
	})

	assert.Equal(t, []byte{0x06, 0x00, 0x00, 0x00}, data[:4])
}

func TestControlWordsSingleItem(t *testing.T) {
	data := compression.PackControlWords([]compression.ControlWord{compression.ControlWord(0x00BBCCDD)})

	assert.Equal(t, []byte{0x03, 0x00, 0x00, 0x00, 0xDD, 0xCC, 0xBB, 0x01}, data)
}

func TestControlWordsMultiItemArray(t *testing.T) {
	data := compression.PackControlWords([]compression.ControlWord{
		compression.ControlWord(0x00BBCCDD),
		compression.ControlWord(0x00112233),
	})

	assert.Equal(t, []byte{0x06, 0x00, 0x00, 0x00, 0xDD, 0xCC, 0xBB, 0x01, 0x33, 0x22, 0x11, 0x01}, data)
}

func TestControlWordsIdenticalWordsArePacked(t *testing.T) {
	data := compression.PackControlWords([]compression.ControlWord{
		compression.ControlWord(0x00BBCCDD),
		compression.ControlWord(0x00BBCCDD),
	})

	require.Equal(t, 8, len(data))
	assert.Equal(t, []byte{0xDD, 0xCC, 0xBB, 0x02}, data[4:8])
}

func TestControlWordsCountIsResetForFurtherWords(t *testing.T) {
	data := compression.PackControlWords([]compression.ControlWord{
		compression.ControlWord(0x00BBCCDD),
		compression.ControlWord(0x00BBCCDD),
		compression.ControlWord(0x00112233),
	})

	require.Equal(t, 12, len(data))
	assert.Equal(t, []byte{0xDD, 0xCC, 0xBB, 0x02, 0x33, 0x22, 0x11, 0x01}, data[4:12])
}

func TestControlWordsMaximumCountIs255(t *testing.T) {
	words := make([]compression.ControlWord, 260)
	for i := 0; i < len(words); i++ {
		words[i] = compression.ControlWord(0x00112233)
	}
	data := compression.PackControlWords(words)

	require.Equal(t, 12, len(data))
	assert.Equal(t, []byte{0x0C, 0x03, 0x00, 0x00, 0x33, 0x22, 0x11, 0xFF, 0x33, 0x22, 0x11, 0x05}, data)
}
