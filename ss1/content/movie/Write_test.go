package movie_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ss1/content/text"
)

func TestWriteOfEmptyContainerCreatesMinimumSizeData(t *testing.T) {
	var container movie.Container
	buffer := bytes.NewBuffer(nil)

	err := movie.Write(buffer, container, text.DefaultCodepage())
	require.Nil(t, err)
	assert.Equal(t, 0x0800, len(buffer.Bytes()))
}

func TestWriteCanSaveEmptyContainer(t *testing.T) {
	var container movie.Container
	buffer := bytes.NewBuffer(nil)

	err := movie.Write(buffer, container, text.DefaultCodepage())
	require.Nil(t, err)

	_, err = movie.Read(bytes.NewReader(buffer.Bytes()), text.DefaultCodepage())

	require.Nil(t, err)
}

func TestWriteSavesAudio(t *testing.T) {
	dataBytes := []byte{0x01, 0x02, 0x03}
	var container movie.Container
	container.Audio.Sound.SampleRate = 22050.0
	container.Audio.Sound.Samples = dataBytes
	buffer := bytes.NewBuffer(nil)

	err := movie.Write(buffer, container, text.DefaultCodepage())
	require.Nil(t, err)

	result, err := movie.Read(bytes.NewReader(buffer.Bytes()), text.DefaultCodepage())

	require.Nil(t, err)
	require.NotNil(t, result)
	assert.Equal(t, dataBytes, result.Audio.Sound.Samples)
}
