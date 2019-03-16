package format

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderSerializesToProperLength(t *testing.T) {
	source := bytes.NewReader(make([]byte, 0x200))
	var header Header

	_ = binary.Read(source, binary.LittleEndian, &header)
	curPos, _ := source.Seek(0, io.SeekCurrent)

	assert.Equal(t, int64(HeaderSize), curPos)
}
