package level_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
)

func TestMasterObjectEntrySerializedSize(t *testing.T) {
	var entry level.MasterObjectEntry
	size := binary.Size(&entry)
	assert.Equal(t, level.MasterObjectEntrySize, size)
}

func TestMasterObjectTableCanBeSerializedWithCoder(t *testing.T) {
	table := level.MasterObjectTable([]level.MasterObjectEntry{{}, {}})
	buf := bytes.NewBuffer(nil)
	encoder := serial.NewEncoder(buf)
	encoder.Code(table)
	data := buf.Bytes()
	assert.Equal(t, level.MasterObjectEntrySize*2, len(data))
}
