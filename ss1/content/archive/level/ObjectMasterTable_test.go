package level_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
)

func TestObjectMasterEntrySerializedSize(t *testing.T) {
	var entry level.ObjectMasterEntry
	size := binary.Size(&entry)
	assert.Equal(t, level.ObjectMasterEntrySize, size)
}

func TestObjectMasterTableCanBeSerializedWithCoder(t *testing.T) {
	table := level.ObjectMasterTable([]level.ObjectMasterEntry{{}, {}})
	buf := bytes.NewBuffer(nil)
	encoder := serial.NewEncoder(buf)
	encoder.Code(table)
	data := buf.Bytes()
	assert.Equal(t, level.ObjectMasterEntrySize*2, len(data))
}

func TestObjectMasterTableResetClearsInUseFlags(t *testing.T) {
	table := level.ObjectMasterTable([]level.ObjectMasterEntry{{InUse: 1}, {InUse: 1}})
	table.Reset()

	assert.Equal(t, byte(0), table[0].InUse, "table[0].InUse should be zero.")
	assert.Equal(t, byte(0), table[1].InUse, "table[1].InUse should be zero.")
}

func TestObjectMasterTableInitializesLists(t *testing.T) {
	table := level.ObjectMasterTable([]level.ObjectMasterEntry{{Next: 20, CrossReferenceTableIndex: 10}, {Next: 30}, {Next: 40}})
	table.Reset()

	assert.Equal(t, level.ObjectID(1), table[0].Next, "table[0].Next should be 1, the first free entry.")
	assert.Equal(t, int16(0), table[0].CrossReferenceTableIndex, "table[0].CrossReferenceTableIndex should be 0, the used list is empty.")
	assert.Equal(t, level.ObjectID(2), table[1].Next, "table[1].Next should be 2, the second free entry.")
	assert.Equal(t, level.ObjectID(0), table[2].Next, "table[2].Next should be 0, the start entry.")
}

func TestDefaultObjectMasterTable(t *testing.T) {
	table := level.DefaultObjectMasterTable()

	assert.Equal(t, 872, len(table), "Table length mismatch")
	assert.Equal(t, level.ObjectID(1), table[0].Next, "table[0].Next should be 1, the first free entry.")
}
