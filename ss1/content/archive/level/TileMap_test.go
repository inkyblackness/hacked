package level_test

import (
	"encoding/binary"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"

	"github.com/stretchr/testify/assert"
)

func TestTileMapEntrySerializedSize(t *testing.T) {
	var entry level.TileMapEntry
	size := binary.Size(&entry)
	assert.Equal(t, 16, size, "Size mismatch")
}
