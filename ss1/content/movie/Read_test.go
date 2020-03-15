package movie_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ss1/content/movie/internal/format"
)

func TestReadReturnsErrorOnNil(t *testing.T) {
	_, err := movie.Read(nil)

	assert.Errorf(t, err, "source is nil")
}

func TestReadReturnsContainerOnEmptyFile(t *testing.T) {
	buffer := bytes.NewBufferString(format.Tag)
	buffer.Write(make([]byte, 0x100+0x300-len(format.Tag)))
	emptyFile := buffer.Bytes()
	source := bytes.NewReader(emptyFile)
	container, _ := movie.Read(source)

	assert.NotNil(t, container)
}

func TestReadReturnsErrorOnMissingTag(t *testing.T) {
	emptyFile := make([]byte, 0x100+0x300)
	source := bytes.NewReader(emptyFile)
	_, err := movie.Read(source)

	assert.Errorf(t, err, "Not a MOVI format")
}

func TestReadReturnsContainerWithBasicPropertiesSet(t *testing.T) {
	buffer := bytes.NewBufferString(format.Tag)
	buffer.Write(make([]byte, 0x100+0x300-len(format.Tag)))
	emptyFile := buffer.Bytes()

	emptyFile[0x10] = 0x80
	emptyFile[0x11] = 0x40
	emptyFile[0x12] = 0x03
	emptyFile[0x18] = 0x80
	emptyFile[0x19] = 0x02
	emptyFile[0x1A] = 0xE0
	emptyFile[0x1B] = 0x01
	emptyFile[0x26] = 0x22
	emptyFile[0x27] = 0x56

	source := bytes.NewReader(emptyFile)
	container, _ := movie.Read(source)

	assert.Equal(t, uint16(640), container.VideoWidth)
	assert.Equal(t, uint16(480), container.VideoHeight)
	assert.Equal(t, uint16(22050), container.AudioSampleRate)
	assert.Equal(t, 0, len(container.Entries))
}

func TestReadReturnsContainerWithDataEntriesExceptTerminator(t *testing.T) {
	testData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	buffer := bytes.NewBufferString(format.Tag)
	buffer.Write(make([]byte, 0x100+0x300-len(format.Tag)))
	buffer.Write(make([]byte, 0xC00))
	buffer.Write(testData)
	raw := buffer.Bytes()

	raw[0x04] = 2
	// size of index table
	raw[0x08] = 0x00
	raw[0x09] = 0x0C

	// index entry 0
	raw[0x0400+3] = 0x02
	raw[0x0400+4] = 0x00
	raw[0x0400+5] = 0x10
	// index entry 1
	raw[0x0408+3] = 0x00
	raw[0x0408+4] = byte(len(testData))
	raw[0x0408+5] = 0x10

	source := bytes.NewReader(raw)
	container, _ := movie.Read(source)

	assert.Equal(t, 1, len(container.Entries))
	assert.Equal(t, movie.DataTypeAudio, container.Entries[0].Data.Type())
	assert.Equal(t, testData, container.Entries[0].Data.Bytes())
}
