package level_test

import (
	"encoding/binary"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"

	"github.com/stretchr/testify/assert"
)

func TestObjectCrossReferenceEntrySerializedSize(t *testing.T) {
	var entry level.ObjectCrossReferenceEntry
	size := binary.Size(&entry)
	assert.Equal(t, level.ObjectCrossReferenceEntrySize, size)
}
