package voc

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadReturnsErrorOnNil(t *testing.T) {
	_, err := Load(nil)

	assert.Errorf(t, err, "source is nil")
}

func newHeader() *bytes.Buffer {
	writer := bytes.NewBufferString(fileHeader)
	version := uint16(0x010A)
	headerSize := uint16(0x001A)

	_ = binary.Write(writer, binary.LittleEndian, headerSize)
	_ = binary.Write(writer, binary.LittleEndian, version)
	versionValidity := uint16(^version + uint16(0x1234))
	_ = binary.Write(writer, binary.LittleEndian, versionValidity)

	return writer
}

func TestLoadReturnsErrorOnInvalidVersion(t *testing.T) {
	writer := newHeader()

	writer.Write([]byte{0x00}) // Terminator

	data := writer.Bytes()
	data[24] = 0x00
	source := bytes.NewReader(data)
	_, err := Load(source)

	assert.Errorf(t, err, "Version validity failed: 0x1129 != 0x1100")
}

func TestLoadReturnsErrorOnValidButEmptyFile(t *testing.T) {
	writer := newHeader()

	writer.Write([]byte{0x00}) // Terminator

	source := bytes.NewReader(writer.Bytes())
	_, err := Load(source)

	assert.Errorf(t, err, "No audio found")
}

func TestLoadReturnsSoundDataOnSampleData(t *testing.T) {
	writer := newHeader()

	writer.Write([]byte{0x01})             // sound data
	writer.Write([]byte{0x03, 0x00, 0x00}) // block size
	writer.Write([]byte{0x64, 0x00})       // divisor, sound type
	writer.Write([]byte{0x80})             // one sample

	writer.Write([]byte{0x00}) // Terminator

	source := bytes.NewReader(writer.Bytes())
	data, err := Load(source)

	require.Nil(t, err)
	assert.NotNil(t, data)
}

func TestLoadReturnsSoundDataWithSampleRate(t *testing.T) {
	writer := newHeader()

	writer.Write([]byte{0x01})             // sound data
	writer.Write([]byte{0x03, 0x00, 0x00}) // block size
	writer.Write([]byte{0x9C, 0x00})       // divisor, sound type
	writer.Write([]byte{0x80})             // one sample

	writer.Write([]byte{0x00}) // Terminator

	source := bytes.NewReader(writer.Bytes())
	data, err := Load(source)

	require.Nil(t, err)
	assert.Equal(t, float32(10000.0), data.SampleRate)
}

func TestLoadReturnsSoundDataWithSamples(t *testing.T) {
	writer := newHeader()
	samples := []byte{0x80, 0xFF, 0x00, 0xC0, 0x40, 0x7F, 0x81}

	writer.Write([]byte{0x01})             // sound data
	writer.Write([]byte{0x09, 0x00, 0x00}) // block size
	writer.Write([]byte{0x9C, 0x00})       // divisor, sound type
	writer.Write(samples)                  // samples

	writer.Write([]byte{0x00}) // Terminator

	source := bytes.NewReader(writer.Bytes())
	data, err := Load(source)

	require.Nil(t, err)
	assert.Equal(t, samples, data.Samples)
}
