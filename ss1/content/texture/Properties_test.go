package texture_test

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/texture"
)

func TestPropertiesSize(t *testing.T) {
	size := binary.Size(texture.Properties{})
	assert.Equal(t, texture.PropertiesSize, size)
}
