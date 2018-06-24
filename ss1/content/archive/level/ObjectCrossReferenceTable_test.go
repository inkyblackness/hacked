package level_test

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
)

func TestObjectCrossReferenceEntrySerializedSize(t *testing.T) {
	var entry level.ObjectCrossReferenceEntry
	size := binary.Size(&entry)
	assert.Equal(t, level.ObjectCrossReferenceEntrySize, size)
}

func TestObjectCrossReferenceTableCanBeSerializedWithCoder(t *testing.T) {
	table := level.ObjectCrossReferenceTable([]level.ObjectCrossReferenceEntry{{}, {}})
	buf := bytes.NewBuffer(nil)
	encoder := serial.NewEncoder(buf)
	encoder.Code(table)
	data := buf.Bytes()
	assert.Equal(t, level.ObjectCrossReferenceEntrySize*2, len(data))
}

func TestObjectCrossReferenceTableResetInitializesChainPointers(t *testing.T) {
	table := level.ObjectCrossReferenceTable([]level.ObjectCrossReferenceEntry{
		aRandomObjectCrossReferenceEntry(), aRandomObjectCrossReferenceEntry(), aRandomObjectCrossReferenceEntry()})
	table.Reset()

	start := table[0]
	assert.Equal(t, level.ObjectID(0), start.ObjectID, "start.ObjectID should be 0, it is unused.")
	assert.Equal(t, int16(1), start.NextInTile, "start.NextInTile should be 1, the first free.")
	free0 := table[1]
	assert.Equal(t, level.ObjectID(0), free0.ObjectID, "free0.ObjectID should be 0, it is unused.")
	assert.Equal(t, int16(2), free0.NextInTile, "free0.NextInTile should be 2, the second free.")
	free1 := table[2]
	assert.Equal(t, level.ObjectID(0), free1.ObjectID, "free1.ObjectID should be 0, it is unused.")
	assert.Equal(t, int16(0), free1.NextInTile, "free1.NextInTile should be 0, the start.")
}

func TestDefaultObjectCrossReferenceTable(t *testing.T) {
	table := level.DefaultObjectCrossReferenceTable()

	assert.Equal(t, 1600, len(table), "Wrong default size")
	assert.Equal(t, int16(1), table[0].NextInTile, "start entry should point to first free.")
}

func aRandomObjectCrossReferenceEntry() level.ObjectCrossReferenceEntry {
	decoder := serial.NewDecoder(rand.Reader)
	var dataSize byte
	decoder.Code(&dataSize)
	var entry level.ObjectCrossReferenceEntry
	decoder.Code(&entry)
	return entry
}
