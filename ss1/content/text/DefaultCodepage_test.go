package text_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/text"

	"github.com/stretchr/testify/assert"
)

func TestDefaultCodepageEncode(t *testing.T) {
	result := text.DefaultCodepage().Encode("ä")

	assert.Equal(t, []byte{132, 0x00}, result)
}

func TestDefaultCodepageDecode(t *testing.T) {
	result := text.DefaultCodepage().Decode([]byte{212, 225, 0x00})

	assert.Equal(t, "Èß", result)
}
