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
