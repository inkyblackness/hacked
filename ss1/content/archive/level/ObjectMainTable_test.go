package level_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObjectMainEntrySerializedSize(t *testing.T) {
	var entry level.ObjectMainEntry
	size := binary.Size(&entry)
	assert.Equal(t, level.ObjectMainEntrySize, size)
}

func TestObjectMainTableCanBeSerializedWithCoder(t *testing.T) {
	table := level.ObjectMainTable([]level.ObjectMainEntry{{}, {}})
	buf := bytes.NewBuffer(nil)
	encoder := serial.NewEncoder(buf)
	encoder.Code(table)
	data := buf.Bytes()
	assert.Equal(t, level.ObjectMainEntrySize*2, len(data))
}

func TestObjectMainTableResetClearsInUseFlags(t *testing.T) {
	table := level.ObjectMainTable([]level.ObjectMainEntry{{InUse: 1}, {InUse: 1}})
	table.Reset()

	assert.Equal(t, byte(0), table[0].InUse, "table[0].InUse should be zero.")
	assert.Equal(t, byte(0), table[1].InUse, "table[1].InUse should be zero.")
}

func TestObjectMainTableInitializesLists(t *testing.T) {
	table := level.ObjectMainTable([]level.ObjectMainEntry{{Next: 20, CrossReferenceTableIndex: 10}, {Next: 30}, {Next: 40}})
	table.Reset()

	assert.Equal(t, level.ObjectID(1), table[0].Next, "table[0].Next should be 1, the first free entry.")
	assert.Equal(t, int16(0), table[0].CrossReferenceTableIndex, "table[0].CrossReferenceTableIndex should be 0, the used list is empty.")
	assert.Equal(t, level.ObjectID(2), table[1].Next, "table[1].Next should be 2, the second free entry.")
	assert.Equal(t, level.ObjectID(0), table[2].Next, "table[2].Next should be 0, the start entry.")
}

func TestDefaultObjectMainTable(t *testing.T) {
	table := level.DefaultObjectMainTable()

	assert.Equal(t, 872, len(table), "Table length mismatch")
	assert.Equal(t, level.ObjectID(1), table[0].Next, "table[0].Next should be 1, the first free entry.")
}

func TestObjectMainTableAllocate(t *testing.T) {
	tt := []int{0, 1, 2, 3, 100}

	for _, tc := range tt {
		table := make(level.ObjectMainTable, tc)
		table.Reset()
		possible := tc - 1

		if possible < 0 {
			possible = 0
		}
		assert.Equal(t, possible, table.Capacity(), "table should have proper capacity")
		for attempt := 0; attempt < possible; attempt++ {
			assert.Equal(t, attempt, table.AllocatedCount(), "table should have expected allocation count")
			id := table.Allocate()
			assert.NotEqual(t, level.ObjectID(0), id, "could not allocate at attempt %d for size %d", attempt, tc)
		}
		assert.Equal(t, possible, table.AllocatedCount(), "table should be full")
		last := table.Allocate()
		assert.Equal(t, level.ObjectID(0), last, "table was not exhausted although it should be")
	}
}

func TestObjectMainTableRelease(t *testing.T) {
	stats := func(table level.ObjectMainTable) (used, free int) {
		for i := 1; i < len(table); i++ {
			entry := &table[i]
			if entry.InUse != 0 {
				used++
			} else {
				free++
			}
		}
		return
	}

	table := make(level.ObjectMainTable, 10)
	table.Reset()
	var allocated []level.ObjectID
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
