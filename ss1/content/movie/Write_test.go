package movie

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inkyblackness/hacked/ss1/content/text"
)

func TestWriteOfEmptyContainerCreatesMinimumSizeData(t *testing.T) {
	var container Container
	buffer := bytes.NewBuffer(nil)

	err := Write(buffer, container, text.DefaultCodepage())
	require.Nil(t, err)
	assert.Equal(t, 0x0800, len(buffer.Bytes()))
}

func TestWriteCanSaveEmptyContainer(t *testing.T) {
	var container Container
	buffer := bytes.NewBuffer(nil)

	err := Write(buffer, container, text.DefaultCodepage())
	require.Nil(t, err)

	_, err = Read(bytes.NewReader(buffer.Bytes()), text.DefaultCodepage())

	require.Nil(t, err)
}

func TestWriteSavesAudio(t *testing.T) {
	dataBytes := []byte{0x01, 0x02, 0x03}
	var container Container
	container.Audio.Sound.SampleRate = 22050.0
	container.Audio.Sound.Samples = dataBytes
	buffer := bytes.NewBuffer(nil)

	err := Write(buffer, container, text.DefaultCodepage())
	require.Nil(t, err)

	result, err := Read(bytes.NewReader(buffer.Bytes()), text.DefaultCodepage())

	require.Nil(t, err)
	require.NotNil(t, result)
	assert.Equal(t, dataBytes, result.Audio.Sound.Samples)
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
