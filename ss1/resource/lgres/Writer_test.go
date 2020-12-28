package lgres_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/resource/lgres"
	"github.com/inkyblackness/hacked/ss1/resource/lgres/internal/format"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
)

func TestNewWriterReturnsErrorForNilTarget(t *testing.T) {
	writer, err := lgres.NewWriter(nil)

	assert.Nil(t, writer, "writer should be nil")
	assert.Equal(t, lgres.ErrTargetNil, err)
}

func TestWriterFinishWithoutAddingResourcesCreatesValidFileWithoutResources(t *testing.T) {
	emptyFileData := emptyResourceFile()
	store := serial.NewByteStore()
	writer, err := lgres.NewWriter(store)
	assert.Nil(t, err, "no error expected creating writer")

	err = writer.Finish()
	assert.Nil(t, err, "no error expected finishing writer")
	assert.Equal(t, emptyFileData, store.Data())
}

func TestWriterFinishReturnsErrorWhenAlreadyFinished(t *testing.T) {
	writer, _ := lgres.NewWriter(serial.NewByteStore())

	err := writer.Finish()
	assert.Nil(t, err, "no error expected finishing")

	err = writer.Finish()
	assert.Equal(t, lgres.ErrWriterFinished, err)
}

func TestWriterUncompressedSingleBlockResourceCanBeWritten(t *testing.T) {
	data := []byte{0xAB, 0x01, 0xCD, 0x02, 0xEF}
	store := serial.NewByteStore()
	writer, _ := lgres.NewWriter(store)
	resourceWriter, err := writer.CreateResource(resource.ID(0x1234), resource.ContentType(0x0A), false)
	assert.Nil(t, err, "no error expected creating resource")
	_, err = resourceWriter.Write(data)
	assert.Nil(t, err, "no error expected writing")
	err = writer.Finish()
	assert.Nil(t, err, "no error expected finishing")

	result := store.Data()

	var expected []byte
	expected = append(expected, data...)
	expected = append(expected,
		0x00, 0x00, 0x00, // alignment for directory
		0x01, 0x00, // resource count
		0x80, 0x00, 0x00, 0x00, // offset to first resource
		0x34, 0x12, // resource ID
		0x05, 0x00, 0x00, // resource length (uncompressed)
		0x00,             // resource type (uncompressed, single-block)
		0x05, 0x00, 0x00, // resource length in file
		0x0A) // content type
	assert.Equal(t, expected, result[format.ResourceDirectoryFileOffsetPos+4:])
}

func TestWriterUncompressedCompoundResourceCanBeWritten(t *testing.T) {
	blockData1 := []byte{0xAB, 0x01, 0xCD}
	blockData2 := []byte{0x11, 0x22, 0x33, 0x44}
	store := serial.NewByteStore()
	writer, _ := lgres.NewWriter(store)
	resourceWriter, err := writer.CreateCompoundResource(resource.ID(0x5678), resource.ContentType(0x0B), false)
	assert.Nil(t, err, "no error expected creating resource")
	_, _ = resourceWriter.CreateBlock().Write(blockData1)
	_, _ = resourceWriter.CreateBlock().Write(blockData2)
	err = writer.Finish()
	assert.Nil(t, err, "no error expected finishing")

	result := store.Data()

	var expected []byte
	expected = append(expected,
		0x02, 0x00, // number of blocks
		0x0E, 0x00, 0x00, 0x00, // offset to first block
		0x11, 0x00, 0x00, 0x00, // offset to second block
		0x15, 0x00, 0x00, 0x00) // size of resource
	expected = append(expected, blockData1...)
	expected = append(expected, blockData2...)
	expected = append(expected,
		0x00, 0x00, 0x00, // alignment for directory
		0x01, 0x00, // resource count
		0x80, 0x00, 0x00, 0x00, // offset to first resource
		0x78, 0x56, // resource ID
		0x15, 0x00, 0x00, // resource length (uncompressed)
		0x02,             // resource type
		0x15, 0x00, 0x00, // resource length in file
		0x0B) // content type
	assert.Equal(t, expected, result[format.ResourceDirectoryFileOffsetPos+4:])
}

func TestWriterUncompressedCompoundResourceCanBeWrittenWithPaddingForSpecialID(t *testing.T) {
	blockData1 := []byte{0xAB, 0x01, 0xCD}
	blockData2 := []byte{0x11, 0x22, 0x33, 0x44}
	store := serial.NewByteStore()
	writer, _ := lgres.NewWriter(store)
	resourceWriter, err := writer.CreateCompoundResource(resource.ID(0x08FD), resource.ContentType(0x0B), false)
	assert.Nil(t, err, "no error expected creating resource")
	_, _ = resourceWriter.CreateBlock().Write(blockData1)
	_, _ = resourceWriter.CreateBlock().Write(blockData2)
	err = writer.Finish()
	assert.Nil(t, err, "no error expected finishing")

	result := store.Data()

	var expected []byte
	expected = append(expected,
		0x02, 0x00, // number of blocks
		0x10, 0x00, 0x00, 0x00, // offset to first block
		0x13, 0x00, 0x00, 0x00, // offset to second block
		0x17, 0x00, 0x00, 0x00, // size of resource
		0x00, 0x00) // padding
	expected = append(expected, blockData1...)
	expected = append(expected, blockData2...)
	expected = append(expected,
		0x00,       // alignment for directory
		0x01, 0x00, // resource count
		0x80, 0x00, 0x00, 0x00, // offset to first resource
		0xFD, 0x08, // resource ID
		0x17, 0x00, 0x00, // resource length (uncompressed)
		0x02,             // resource type
		0x17, 0x00, 0x00, // resource length in file
		0x0B) // content type
	assert.Equal(t, expected, result[format.ResourceDirectoryFileOffsetPos+4:])
}

func TestWriterCompressedSingleBlockResourceCanBeWritten(t *testing.T) {
	data := []byte{0x01, 0x02, 0x01, 0x02}
	store := serial.NewByteStore()
	writer, _ := lgres.NewWriter(store)
	resourceWriter, err := writer.CreateResource(resource.ID(0x1122), resource.ContentType(0x0C), true)
	assert.Nil(t, err, "no error expected creating resource")
	_, _ = resourceWriter.Write(data)
	err = writer.Finish()
	assert.Nil(t, err, "no error expected finishing")

	result := store.Data()

	var expected []byte
	// 0000 0000|0000 0100|0000 0000|0010 0000|0100 0000|0011 1111|1111 1111
	expected = append(expected,
		0x00, 0x04, 0x00, 0x20, 0x40, 0x3F, 0xFF, 0x00, // 14bit words 0x0001 0x0002 0x0100 0x3FFF + trailing 0x00
		// (no bytes) alignment for directory
		0x01, 0x00, // resource count
		0x80, 0x00, 0x00, 0x00, // offset to first resource
		0x22, 0x11, // resource ID
		0x04, 0x00, 0x00, // resource length (uncompressed)
		0x01,             // resource type
		0x08, 0x00, 0x00, // resource length in file
		0x0C) // content type
	assert.Equal(t, expected, result[format.ResourceDirectoryFileOffsetPos+4:])
}

func TestWriterCompressedCompoundResourceCanBeWritten(t *testing.T) {
	blockData1 := []byte{0x01, 0x02, 0x01, 0x02}
	blockData2 := []byte{0x01, 0x02, 0x01, 0x02}
	store := serial.NewByteStore()
	writer, _ := lgres.NewWriter(store)
	resourceWriter, err := writer.CreateCompoundResource(resource.ID(0x5544), resource.ContentType(0x09), true)
	assert.Nil(t, err, "no error expected creating resource")
	_, _ = resourceWriter.CreateBlock().Write(blockData1)
	_, _ = resourceWriter.CreateBlock().Write(blockData2)
	err = writer.Finish()
	assert.Nil(t, err, "no error expected finishing")

	result := store.Data()

	var expected []byte
	expected = append(expected,
		0x02, 0x00, // number of blocks
		0x0E, 0x00, 0x00, 0x00, // offset to first block
		0x12, 0x00, 0x00, 0x00, // offset to second block
		0x16, 0x00, 0x00, 0x00, // size of resource
		0x00, 0x04, 0x00, 0x20, 0x40, // compressed data, part 1
		0x01, 0x02, 0x00, 0x0B, 0xFF, 0xF0, 0x00, // compressed data, part 2
		0x00, 0x00, // alignment for directory
		0x01, 0x00, // resource count
		0x80, 0x00, 0x00, 0x00, // offset to first resource
		0x44, 0x55, // resource ID
		0x16, 0x00, 0x00, // resource length (uncompressed)
		0x03,             // resource type
		0x1A, 0x00, 0x00, // resource length in file
		0x09) // content type
	assert.Equal(t, expected, result[format.ResourceDirectoryFileOffsetPos+4:])
}
