package compression_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/movie/compression"
)

func TestMaskstreamWriter(t *testing.T) {
	verifyMaskstreamWriter(t,
		[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		[]byte{},
		[]byte{0x01, 0x02},
		[]byte{0x01},
		[]byte{0x03, 0x04},
		[]byte{0x05, 0x06},
		[]byte{0x06, 0x07, 0x08})
}

func TestMaskstreamWriterOffset(t *testing.T) {
	tt := []struct {
		name     string
		buffer   []byte
		add      []byte
		expected uint32
	}{
		{name: "add at end if not found", buffer: []byte{0x22}, add: []byte{0x11}, expected: 1},
		{name: "find existing", buffer: []byte{0x11}, add: []byte{0x11}, expected: 0},
		{name: "add partially at end A", buffer: []byte{0x01, 0x02}, add: []byte{0x01, 0x02, 0x03}, expected: 0},
		{name: "add partially at end B", buffer: []byte{0x01, 0x02}, add: []byte{0x02, 0x03}, expected: 1},
	}

	for _, tc := range tt {
		td := tc
		t.Run(td.name, func(t *testing.T) {
			w := compression.PaletteLookupWriter{Buffer: td.buffer}
			result := w.Write(td.add)
			assert.Equal(t, td.expected, result)
		})
	}
}

func verifyMaskstreamWriter(t *testing.T, expected []byte, parts ...[]byte) {
	t.Helper()
	var w compression.PaletteLookupWriter
	for _, part := range parts {
		w.Write(part)
	}
	assert.Equal(t, expected, w.Buffer)
}
