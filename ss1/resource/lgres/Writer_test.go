package lgres

import (
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial"

	"github.com/stretchr/testify/assert"
)

func TestNewWriterReturnsErrorForNilTarget(t *testing.T) {
	writer, err := NewWriter(nil)

	assert.Nil(t, writer, "writer should be nil")
	assert.Equal(t, errTargetNil, err)
}

func TestWriterFinishWithoutAddingResourcesCreatesValidFileWithoutResources(t *testing.T) {
	emptyFileData := emptyResourceFile()
	store := serial.NewByteStore()
	writer, err := NewWriter(store)
	assert.Nil(t, err, "no error expected creating writer")

	err = writer.Finish()
	assert.Nil(t, err, "no error expected finishing writer")
	assert.Equal(t, emptyFileData, store.Data())
}

func TestWriterFinishReturnsErrorWhenAlreadyFinished(t *testing.T) {
	writer, _ := NewWriter(serial.NewByteStore())

	writer.Finish()

	err := writer.Finish()
	assert.Equal(t, errWriterFinished, err)
}

func TestWriterUncompressedSingleBlockResourceCanBeWritten(t *testing.T) {
	data := []byte{0xAB, 0x01, 0xCD, 0x02, 0xEF}
	store := serial.NewByteStore()
	writer, _ := NewWriter(store)
	resourceWriter, err := writer.CreateResource(resource.ID(0x1234), resource.ContentType(0x0A), false)
	assert.Nil(t, err, "no error expected")
	resourceWriter.Write(data)
	writer.Finish()

	result := store.Data()

	var expected []byte
	expected = append(expected, data...)
	expected = append(expected, 0x00, 0x00, 0x00)       // alignment for directory
	expected = append(expected, 0x01, 0x00)             // resource count
	expected = append(expected, 0x80, 0x00, 0x00, 0x00) // offset to first resource
	expected = append(expected, 0x34, 0x12)             // resource ID
	expected = append(expected, 0x05, 0x00, 0x00)       // resource length (uncompressed)
	expected = append(expected, 0x00)                   // resource type (uncompressed, single-block)
	expected = append(expected, 0x05, 0x00, 0x00)       // resource length in file
	expected = append(expected, 0x0A)                   // content type
	assert.Equal(t, expected, result[resourceDirectoryFileOffsetPos+4:])
}

func TestWriterUncompressedCompoundResourceCanBeWritten(t *testing.T) {
	blockData1 := []byte{0xAB, 0x01, 0xCD}
	blockData2 := []byte{0x11, 0x22, 0x33, 0x44}
	store := serial.NewByteStore()
	writer, _ := NewWriter(store)
	resourceWriter, err := writer.CreateCompoundResource(resource.ID(0x5678), resource.ContentType(0x0B), false)
	assert.Nil(t, err, "no error expected")
	resourceWriter.CreateBlock().Write(blockData1)
	resourceWriter.CreateBlock().Write(blockData2)
	writer.Finish()

	result := store.Data()

	var expected []byte
	expected = append(expected, 0x02, 0x00)             // number of blocks
	expected = append(expected, 0x0E, 0x00, 0x00, 0x00) // offset to first block
	expected = append(expected, 0x11, 0x00, 0x00, 0x00) // offset to second block
	expected = append(expected, 0x15, 0x00, 0x00, 0x00) // size of resource
	expected = append(expected, blockData1...)
	expected = append(expected, blockData2...)
	expected = append(expected, 0x00, 0x00, 0x00)       // alignment for directory
	expected = append(expected, 0x01, 0x00)             // resource count
	expected = append(expected, 0x80, 0x00, 0x00, 0x00) // offset to first resource
	expected = append(expected, 0x78, 0x56)             // resource ID
	expected = append(expected, 0x15, 0x00, 0x00)       // resource length (uncompressed)
	expected = append(expected, 0x02)                   // resource type
	expected = append(expected, 0x15, 0x00, 0x00)       // resource length in file
	expected = append(expected, 0x0B)                   // content type
	assert.Equal(t, expected, result[resourceDirectoryFileOffsetPos+4:])
}

func TestWriterUncompressedCompoundResourceCanBeWrittenWithPaddingForSpecialID(t *testing.T) {
	blockData1 := []byte{0xAB, 0x01, 0xCD}
	blockData2 := []byte{0x11, 0x22, 0x33, 0x44}
	store := serial.NewByteStore()
	writer, _ := NewWriter(store)
	resourceWriter, err := writer.CreateCompoundResource(resource.ID(0x08FD), resource.ContentType(0x0B), false)
	assert.Nil(t, err, "no error expected")
	resourceWriter.CreateBlock().Write(blockData1)
	resourceWriter.CreateBlock().Write(blockData2)
	writer.Finish()

	result := store.Data()

	var expected []byte
	expected = append(expected, 0x02, 0x00)             // number of blocks
	expected = append(expected, 0x10, 0x00, 0x00, 0x00) // offset to first block
	expected = append(expected, 0x13, 0x00, 0x00, 0x00) // offset to second block
	expected = append(expected, 0x17, 0x00, 0x00, 0x00) // size of resource
	expected = append(expected, 0x00, 0x00)             // padding
	expected = append(expected, blockData1...)
	expected = append(expected, blockData2...)
	expected = append(expected, 0x00)                   // alignment for directory
	expected = append(expected, 0x01, 0x00)             // resource count
	expected = append(expected, 0x80, 0x00, 0x00, 0x00) // offset to first resource
	expected = append(expected, 0xFD, 0x08)             // resource ID
	expected = append(expected, 0x17, 0x00, 0x00)       // resource length (uncompressed)
	expected = append(expected, 0x02)                   // resource type
	expected = append(expected, 0x17, 0x00, 0x00)       // resource length in file
	expected = append(expected, 0x0B)                   // content type
	assert.Equal(t, expected, result[resourceDirectoryFileOffsetPos+4:])
}

func TestWriterCompressedSingleBlockResourceCanBeWritten(t *testing.T) {
	data := []byte{0x01, 0x02, 0x01, 0x02}
	store := serial.NewByteStore()
	writer, _ := NewWriter(store)
	resourceWriter, err := writer.CreateResource(resource.ID(0x1122), resource.ContentType(0x0C), true)
	assert.Nil(t, err, "no error expected")
	resourceWriter.Write(data)
	writer.Finish()

	result := store.Data()

	var expected []byte
	// 0000 0000|0000 0100|0000 0000|0010 0000|0100 0000|0011 1111|1111 1111
	expected = append(expected, 0x00, 0x04, 0x00, 0x20, 0x40, 0x3F, 0xFF, 0x00) // 14bit words 0x0001 0x0002 0x0100 0x3FFF + trailing 0x00
	expected = append(expected)                                                 // alignment for directory
	expected = append(expected, 0x01, 0x00)                                     // resource count
	expected = append(expected, 0x80, 0x00, 0x00, 0x00)                         // offset to first resource
	expected = append(expected, 0x22, 0x11)                                     // resource ID
	expected = append(expected, 0x04, 0x00, 0x00)                               // resource length (uncompressed)
	expected = append(expected, 0x01)                                           // resource type
	expected = append(expected, 0x08, 0x00, 0x00)                               // resource length in file
	expected = append(expected, 0x0C)                                           // content type
	assert.Equal(t, expected, result[resourceDirectoryFileOffsetPos+4:])
}

func TestWriterCompressedCompoundResourceCanBeWritten(t *testing.T) {
	blockData1 := []byte{0x01, 0x02, 0x01, 0x02}
	blockData2 := []byte{0x01, 0x02, 0x01, 0x02}
	store := serial.NewByteStore()
	writer, _ := NewWriter(store)
	resourceWriter, err := writer.CreateCompoundResource(resource.ID(0x5544), resource.ContentType(0x09), true)
	assert.Nil(t, err, "no error expected")
	resourceWriter.CreateBlock().Write(blockData1)
	resourceWriter.CreateBlock().Write(blockData2)
	writer.Finish()

	result := store.Data()

	var expected []byte
	expected = append(expected, 0x02, 0x00)                               // number of blocks
	expected = append(expected, 0x0E, 0x00, 0x00, 0x00)                   // offset to first block
	expected = append(expected, 0x12, 0x00, 0x00, 0x00)                   // offset to second block
	expected = append(expected, 0x16, 0x00, 0x00, 0x00)                   // size of resource
	expected = append(expected, 0x00, 0x04, 0x00, 0x20, 0x40)             // compressed data, part 1
	expected = append(expected, 0x01, 0x02, 0x00, 0x0B, 0xFF, 0xF0, 0x00) // compressed data, part 2
	expected = append(expected, 0x00, 0x00)                               // alignment for directory
	expected = append(expected, 0x01, 0x00)                               // resource count
	expected = append(expected, 0x80, 0x00, 0x00, 0x00)                   // offset to first resource
	expected = append(expected, 0x44, 0x55)                               // resource ID
	expected = append(expected, 0x16, 0x00, 0x00)                         // resource length (uncompressed)
	expected = append(expected, 0x03)                                     // resource type
	expected = append(expected, 0x1A, 0x00, 0x00)                         // resource length in file
	expected = append(expected, 0x09)                                     // content type
	assert.Equal(t, expected, result[resourceDirectoryFileOffsetPos+4:])
}
