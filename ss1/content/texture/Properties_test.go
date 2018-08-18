package texture_test

import (
	"encoding/binary"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/stretchr/testify/assert"
)

func TestPropertiesSize(t *testing.T) {
	size := binary.Size(texture.Properties{})
	assert.Equal(t, texture.PropertiesSize, size)
}
