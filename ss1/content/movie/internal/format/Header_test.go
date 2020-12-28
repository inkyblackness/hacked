package format_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/movie/internal/format"
)

func TestHeaderSerializesToProperLength(t *testing.T) {
	source := bytes.NewReader(make([]byte, 0x200))
	var header format.Header

	_ = binary.Read(source, binary.LittleEndian, &header)
	curPos, _ := source.Seek(0, io.SeekCurrent)

	assert.Equal(t, int64(format.HeaderSize), curPos)
}
