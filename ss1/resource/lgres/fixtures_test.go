package lgres

import (
	"bytes"
	"encoding/binary"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial"
)

func emptyResourceFile() []byte {
	buf := bytes.NewBufferString(headerString)
	headerTrailer := make([]byte, resourceDirectoryFileOffsetPos-buf.Len())
	headerTrailer[0] = commentTerminator

	binary.Write(buf, binary.LittleEndian, headerTrailer)
	dictionaryOffset := uint32(buf.Len() + 4)
	binary.Write(buf, binary.LittleEndian, &dictionaryOffset)

	numberOfResources := uint16(0)
	firstResourceOffset := uint32(buf.Len())

	binary.Write(buf, binary.LittleEndian, &numberOfResources)
	binary.Write(buf, binary.LittleEndian, &firstResourceOffset)

	return buf.Bytes()
}

var exampleResourceIDSingleBlockResource = resource.ID(0x4000)
var exampleResourceIDSingleBlockResourceCompressed = resource.ID(0x1000)
var exampleResourceIDCompoundResource = resource.ID(0x2000)
var exampleResourceIDCompoundResourceCompressed = resource.ID(0x5000)

func exampleResourceFile() []byte {
	store := serial.NewByteStore()
	writer, _ := NewWriter(store)

	resource1, _ := writer.CreateResource(exampleResourceIDSingleBlockResource, resource.ContentType(0x01), false)
	resource1.Write([]byte{0x01, 0x01, 0x01})
	resource2, _ := writer.CreateResource(exampleResourceIDSingleBlockResourceCompressed, resource.ContentType(0x02), true)
	resource2.Write([]byte{0x02, 0x02})
	resource3, _ := writer.CreateCompoundResource(exampleResourceIDCompoundResource, resource.ContentType(0x03), false)
	resource3.CreateBlock().Write([]byte{0x30, 0x30, 0x30, 0x30})
	resource3.CreateBlock().Write([]byte{0x31, 0x31, 0x31})
	resource4, _ := writer.CreateCompoundResource(exampleResourceIDCompoundResourceCompressed, resource.ContentType(0x04), true)
	resource4.CreateBlock().Write([]byte{0x40, 0x40})
	resource4.CreateBlock().Write([]byte{0x41, 0x41, 0x41, 0x41})
	resource4.CreateBlock().Write([]byte{0x42})
	writer.Finish()

	return store.Data()
}
