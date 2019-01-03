package level_test

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestObjectCrossReferenceTableAllocate(t *testing.T) {
	tt := []int{0, 1, 2, 3, 100}

	for _, tc := range tt {
		table := make(level.ObjectCrossReferenceTable, tc)
		table.Reset()
		possible := tc - 1

		for attempt := 0; attempt < possible; attempt++ {
			index := table.Allocate()
			assert.NotEqual(t, 0, index, "could not allocate at attempt %d for size %d", attempt, tc)
		}
		last := table.Allocate()
		assert.Equal(t, 0, last, "table was not exhausted although it should be")
	}
}

func TestObjectCrossReferenceTableRelease(t *testing.T) {
	stats := func(table level.ObjectCrossReferenceTable) (used, free int) {
		for next := table[0].NextInTile; next != 0; next = table[next].NextInTile {
			free++
		}
		used = len(table) - 1 - free
		return
	}

	table := make(level.ObjectCrossReferenceTable, 10)
	table.Reset()
	var allocated []int
	for i := 0; i < 4; i++ {
		id := table.Allocate()
		allocated = append(allocated, id)
	}
	used, free := stats(table)
	require.Equal(t, 4, used, "invalid amount of used entries")
	require.Equal(t, 5, free, "invalid amount of free entries")

	for _, id := range allocated {
		table.Release(id)
	}

	used, free = stats(table)
	assert.Equal(t, 0, used, "invalid amount of used entries after release")
	assert.Equal(t, 9, free, "invalid amount of free entries after release")

	table.Release(0)
	table.Release(20)
	for i := 0; i < len(table)-1; i++ {
		id := table.Allocate()
		assert.NotEqual(t, level.ObjectID(0), id, "should have been able to re-allocate")
	}
}

func aRandomObjectCrossReferenceEntry() level.ObjectCrossReferenceEntry {
	decoder := serial.NewDecoder(rand.Reader)
	var dataSize byte
	decoder.Code(&dataSize)
	var entry level.ObjectCrossReferenceEntry
	decoder.Code(&entry)
	return entry
}
