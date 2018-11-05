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
	result := text.DefaultCodepage().Decode([]byte{144, 225, 0x00})

	assert.Equal(t, "Éß", result)
}

func TestDefaultCodepageMapsUnknownCharacterToQuestionMark(t *testing.T) {
	result := text.DefaultCodepage().Encode("„quoted”")

	assert.Equal(t, []byte{0x3F, 0x71, 0x75, 0x6F, 0x74, 0x65, 0x64, 0x3F, 0x00}, result)
}
