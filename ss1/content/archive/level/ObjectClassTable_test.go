package level_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewObjectClassEntrySetsData(t *testing.T) {
	entry := level.NewObjectClassEntry(20)
	assert.Equal(t, 20, len(entry.Data))
}

func TestObjectClassEntryCanBeSerializedWithCoder(t *testing.T) {
	entry := level.NewObjectClassEntry(10)
	buf := bytes.NewBuffer(nil)
	encoder := serial.NewEncoder(buf)
	encoder.Code(entry)
	data := buf.Bytes()
	assert.Equal(t, 10+6, len(data))
}

func TestObjectClassTableCanBeSerializedWithCoder(t *testing.T) {
	table := level.ObjectClassTable([]level.ObjectClassEntry{*level.NewObjectClassEntry(5), *level.NewObjectClassEntry(5)})
	buf := bytes.NewBuffer(nil)
	encoder := serial.NewEncoder(buf)
	encoder.Code(table)
	data := buf.Bytes()
	assert.Equal(t, (5+6)*2, len(data))
}

func TestObjectClassTableCanAllocateData(t *testing.T) {
	table := level.ObjectClassTable([]level.ObjectClassEntry{{}, {}})
	table.AllocateData(4)
	assert.Equal(t, 4, len(table[0].Data))
	assert.Equal(t, 4, len(table[1].Data))
}

func TestObjectClassTableResetClearsEntryData(t *testing.T) {
	table := level.ObjectClassTable([]level.ObjectClassEntry{aRandomObjectClassEntry(), aRandomObjectClassEntry()})
	table.Reset()

	for index, entry := range table {
		expected := make([]byte, len(entry.Data))
		assert.Equal(t, expected, entry.Data, fmt.Sprintf("Data not zero for %d", index))
	}
}

func TestObjectClassTableResetInitializesChainPointers(t *testing.T) {
	table := level.ObjectClassTable([]level.ObjectClassEntry{aRandomObjectClassEntry(), aRandomObjectClassEntry(), aRandomObjectClassEntry()})
	table.Reset()

	start := table[0]
	assert.Equal(t, int16(1), start.Next, "start.Next should be 1, the first free entry.")
	assert.Equal(t, level.ObjectID(0), start.ObjectID, "start.ObjectID should be 0, the list is empty.")
	free0 := table[1]
	assert.Equal(t, int16(2), free0.Next, "free0.Next should be 2, the second free entry.")
	free1 := table[2]
	assert.Equal(t, int16(0), free1.Next, "free1.Next should be 0, the start entry.")
}

func TestDefaultObjectClassTable(t *testing.T) {
	table := level.DefaultObjectClassTable(0)
	assert.Equal(t, 16, len(table), "16 entries for class 0 expected")
	assert.Equal(t, 2, len(table[0].Data), "2 bytes for class 0 data expected")
	assert.Equal(t, int16(1), table[0].Next, "free list expected to be initialized")
}

func TestObjectClassTableAllocate(t *testing.T) {
	tt := []int{0, 1, 2, 3, 100}

	for _, tc := range tt {
		table := make(level.ObjectClassTable, tc)
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

func TestObjectClassTableRelease(t *testing.T) {
	stats := func(table level.ObjectClassTable) (used, free int) {
		for next := int16(table[0].ObjectID); next != 0; next = table[next].Next {
			used++
		}
		for next := table[0].Next; next != 0; next = table[next].Next {
			free++
		}
		return
	}

	table := make(level.ObjectClassTable, 10)
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

func aRandomObjectClassEntry() level.ObjectClassEntry {
	decoder := serial.NewDecoder(rand.Reader)
	var dataSize byte
	decoder.Code(&dataSize)
	entry := level.NewObjectClassEntry(int(dataSize) + 1)
	decoder.Code(entry)
	return *entry
}
