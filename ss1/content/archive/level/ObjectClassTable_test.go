package level_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
)

func TestNewObjectClassEntrySetsData(t *testing.T) {
	entry := level.NewObjectClassEntry(20)
	assert.Equal(t, 20, len(entry.Data))
}

func TestObjectClassEntryCanBeSerializedWithCoder(t *testing.T) {
	entry := level.NewObjectClassEntry(10)
	buf := bytes.NewBuffer(nil)
	encoder := serial.NewEncoder(buf)
	encoder.Code(&entry)
	data := buf.Bytes()
	assert.Equal(t, 10+6, len(data))
}

func TestObjectClassTableCanBeSerializedWithCoder(t *testing.T) {
	table := level.ObjectClassTable([]level.ObjectClassEntry{level.NewObjectClassEntry(5), level.NewObjectClassEntry(5)})
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

func aRandomObjectClassEntry() level.ObjectClassEntry {
	decoder := serial.NewDecoder(rand.Reader)
	var dataSize byte
	decoder.Code(&dataSize)
	entry := level.NewObjectClassEntry(int(dataSize) + 1)
	decoder.Code(&entry)
	return entry
}
