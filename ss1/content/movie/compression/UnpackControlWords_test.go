package compression_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inkyblackness/hacked/ss1/content/movie/compression"
)

func TestUnpackControlWordsFormatErrorForEmptyArray(t *testing.T) {
	_, err := compression.UnpackControlWords(nil)

	assert.Equal(t, compression.FormatError, err)
}

func TestUnpackControlWordsFormatErrorForTooSmallSizeField(t *testing.T) {
	_, err := compression.UnpackControlWords(make([]byte, 3))

	assert.Equal(t, compression.FormatError, err)
}

func TestUnpackControlWordsFormatErrorForSizeValueNotMultipleOfThree(t *testing.T) {
	_, err := compression.UnpackControlWords([]byte{0x02, 0x00, 0x00, 0x00, 0xCC, 0xBB, 0xAA, 0x01})

	assert.Equal(t, compression.FormatError, err)
}

func TestUnpackControlWordsEmptyResultForNoEntries(t *testing.T) {
	words, err := compression.UnpackControlWords([]byte{0x00, 0x00, 0x00, 0x00})

	require.Nil(t, err)
	assert.Equal(t, 0, len(words))
}

func TestUnpackControlWordsSingleEntry(t *testing.T) {
	words, err := compression.UnpackControlWords([]byte{0x03, 0x00, 0x00, 0x00, 0xCC, 0xBB, 0xAA, 0x01})

	require.Nil(t, err)
	assert.Equal(t, []compression.ControlWord{compression.ControlWord(0xAABBCC)}, words)
}

func TestUnpackControlWordsMultipleEntry(t *testing.T) {
	words, err := compression.UnpackControlWords([]byte{0x09, 0x00, 0x00, 0x00, 0xCC, 0xBB, 0xAA, 0x02, 0x33, 0x22, 0x11, 0x01})

	require.Nil(t, err)
	assert.Equal(t, []compression.ControlWord{compression.ControlWord(0xAABBCC), compression.ControlWord(0xAABBCC), compression.ControlWord(0x112233)}, words)
}

func TestUnpackControlWordsErrorForTooManyUnpacked(t *testing.T) {
	_, err := compression.UnpackControlWords([]byte{0x03, 0x00, 0x00, 0x00, 0xCC, 0xBB, 0xAA, 0x02})

	require.NotNil(t, err)
}
