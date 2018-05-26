package lgres

import (
	"errors"
	"io"
	"math"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/resource/lgres/compression"
	"github.com/inkyblackness/hacked/ss1/serial"
)

// Writer provides methods to write a new resource file from scratch.
// Resources have to be created sequentially. The writer does not support
// concurrent creation and modification of resources.
type Writer struct {
	encoder *serial.PositioningEncoder

	firstResourceOffset        uint32
	currentResourceStartOffset uint32
	currentResource            resourceWriter

	directory []*resourceDirectoryEntry
}

var errTargetNil = errors.New("target is nil")

// NewWriter returns a new Writer instance prepared to add resources.
// To finalize the created file, call Finish().
//
// This function will write initial information to the target and will return
// an error if the writer did. In such a case, the returned writer instance
// will produce invalid results and the state of the target is undefined.
func NewWriter(target io.WriteSeeker) (*Writer, error) {
	if target == nil {
		return nil, errTargetNil
	}

	encoder := serial.NewPositioningEncoder(target)
	writer := &Writer{encoder: encoder}
	writer.writeHeader()
	writer.firstResourceOffset = writer.encoder.CurPos()

	return writer, writer.encoder.FirstError()
}

var errWriterFinished = errors.New("writer is finished")

// CreateResource adds a new single-block resource to the current resource file.
// This resource is closed by creating another resource, or by finishing the writer.
func (writer *Writer) CreateResource(id resource.ID, contentType resource.ContentType,
	compressed bool) (*BlockWriter, error) {
	if writer.encoder == nil {
		return nil, errWriterFinished
	}

	writer.finishLastResource()
	if writer.encoder.FirstError() != nil {
		return nil, writer.encoder.FirstError()
	}

	var targetWriter io.Writer = serial.NewEncoder(writer.encoder)
	targetFinisher := func() {}
	resourceType := byte(0x00)
	if compressed {
		compressor := compression.NewCompressor(targetWriter)
		resourceType |= resourceTypeFlagCompressed
		targetWriter = compressor
		targetFinisher = func() { compressor.Close() } // nolint: errcheck
	}
	blockWriter := &BlockWriter{target: targetWriter, finisher: targetFinisher}
	writer.addNewResource(id, contentType, resourceType, blockWriter)

	return blockWriter, nil
}

// CreateCompoundResource adds a new compound resource to the current resource file.
// This resource is closed by creating another resource, or by finishing the writer.
func (writer *Writer) CreateCompoundResource(id resource.ID, contentType resource.ContentType,
	compressed bool) (*CompoundResourceWriter, error) {
	if writer.encoder == nil {
		return nil, errWriterFinished
	}

	writer.finishLastResource()
	if writer.encoder.FirstError() != nil {
		return nil, writer.encoder.FirstError()
	}

	resourceType := resourceTypeFlagCompound
	if compressed {
		resourceType |= resourceTypeFlagCompressed
	}
	resourceWriter := &CompoundResourceWriter{
		target:          serial.NewPositioningEncoder(writer.encoder),
		compressed:      compressed,
		dataPaddingSize: writer.dataPaddingSizeForCompoundResource(id)}
	writer.addNewResource(id, contentType, resourceType, resourceWriter)

	return resourceWriter, nil
}

// Finish finalizes the resource file. After calling this function, the
// writer becomes unusable.
func (writer *Writer) Finish() (err error) {
	if writer.encoder == nil {
		return errWriterFinished
	}

	writer.finishLastResource()

	directoryOffset := writer.encoder.CurPos()
	writer.encoder.SetCurPos(resourceDirectoryFileOffsetPos)
	writer.encoder.Code(directoryOffset)
	writer.encoder.SetCurPos(directoryOffset)
	writer.encoder.Code(uint16(len(writer.directory)))
	writer.encoder.Code(writer.firstResourceOffset)
	for _, entry := range writer.directory {
		writer.encoder.Code(entry)
	}

	err = writer.encoder.FirstError()
	writer.encoder = nil

	return
}

func (writer *Writer) writeHeader() {
	header := make([]byte, resourceDirectoryFileOffsetPos)
	for index, r := range headerString {
		header[index] = byte(r)
	}
	header[len(headerString)] = commentTerminator
	writer.encoder.Code(header)
	writer.encoder.Code(uint32(math.MaxUint32))
}

func (writer *Writer) addNewResource(id resource.ID, contentType resource.ContentType, resourceType byte, newResource resourceWriter) {
	entry := &resourceDirectoryEntry{ID: id.Value()}
	entry.setContentType(byte(contentType))
	entry.setResourceType(resourceType)
	writer.directory = append(writer.directory, entry)
	writer.currentResource = newResource
	writer.currentResourceStartOffset = writer.encoder.CurPos()
}

func (writer *Writer) finishLastResource() {
	if writer.currentResource != nil {
		currentEntry := writer.directory[len(writer.directory)-1]
		currentEntry.setUnpackedLength(writer.currentResource.finish())
		currentEntry.setPackedLength(writer.encoder.CurPos() - writer.currentResourceStartOffset)

		writer.currentResourceStartOffset = 0
		writer.currentResource = nil
	}
	writer.alignToBoundary()
}

func (writer *Writer) alignToBoundary() {
	extraBytes := writer.encoder.CurPos() % boundarySize
	if extraBytes > 0 {
		padding := make([]byte, boundarySize-extraBytes)
		writer.encoder.Code(padding)
	}
}

func (writer *Writer) dataPaddingSizeForCompoundResource(id resource.ID) (padding int) {
	// Some directories have a 2byte padding before the actual data
	idValue := id.Value()
	if (idValue >= 0x08FC) && (idValue <= 0x094B) { // all resources in obj3d.res
		padding = 2
	}
	return
}
