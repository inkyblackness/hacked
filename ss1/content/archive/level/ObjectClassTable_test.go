package level_test

import (
	"bytes"
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
