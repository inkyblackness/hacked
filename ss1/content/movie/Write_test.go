package movie

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteOfEmptyContainerCreatesMinimumSizeData(t *testing.T) {
	builder := NewContainerBuilder()
	container := builder.Build()
	buffer := bytes.NewBuffer(nil)

	err := Write(buffer, container)
	require.Nil(t, err)
	assert.Equal(t, 0x0800, len(buffer.Bytes()))
}

func TestWriteCanSaveEmptyContainer(t *testing.T) {
	builder := NewContainerBuilder()
	container := builder.Build()
	buffer := bytes.NewBuffer(nil)

	err := Write(buffer, container)
	require.Nil(t, err)

	result, err := Read(bytes.NewReader(buffer.Bytes()))

	require.Nil(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 0, result.EntryCount())
}

func TestWriteSavesEntries(t *testing.T) {
	dataBytes := []byte{0x01, 0x02, 0x03}
	builder := NewContainerBuilder()
	builder.AudioSampleRate(22050.0)
	builder.AddEntry(NewMemoryEntry(0.0, Audio, dataBytes))
	container := builder.Build()
	buffer := bytes.NewBuffer(nil)

	err := Write(buffer, container)
	require.Nil(t, err)

	result, err := Read(bytes.NewReader(buffer.Bytes()))

	require.Nil(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.EntryCount())
	assert.Equal(t, dataBytes, result.Entry(0).Data())
}

func TestIndexTableSizeFor_ExistingSizes(t *testing.T) {
	// These sample values are always the minimum and maximum amount of index entries
	// found for a given index size.

	assert.Equal(t, 0x0400, indexTableSizeFor(3))
	assert.Equal(t, 0x0400, indexTableSizeFor(127))

	assert.Equal(t, 0x0C00, indexTableSizeFor(130))
	assert.Equal(t, 0x0C00, indexTableSizeFor(218))

	assert.Equal(t, 0x1C00, indexTableSizeFor(738))
	assert.Equal(t, 0x1C00, indexTableSizeFor(755))

	assert.Equal(t, 0x3400, indexTableSizeFor(1475))
	assert.Equal(t, 0x3400, indexTableSizeFor(1523))
}
