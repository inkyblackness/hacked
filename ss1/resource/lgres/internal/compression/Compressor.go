package compression

import (
	"io"

	"github.com/inkyblackness/hacked/ss1/serial"
)

type compressor struct {
	coder  serial.Coder
	writer *WordWriter

	dictBuffer     dictEntryBuffer
	overtime       int
	dictionary     *dictEntry
	dictionarySize int
	curEntry       *dictEntry
}

// NewCompressor creates a new compressor instance over a writer.
func NewCompressor(target io.Writer) io.WriteCloser {
	coder := serial.NewEncoder(target)
	obj := &compressor{
		coder:          coder,
		writer:         NewWordWriter(coder),
		dictionary:     rootDictEntry(),
		dictionarySize: 0,
		overtime:       0}

	obj.resetDictionary()

	return obj
}

func (obj *compressor) resetDictionary() {
	obj.dictionarySize = 0
	for i := 0; i < 0x100; i++ {
		obj.dictionary.Add(byte(i), Word(i), obj.dictBuffer.entry(Word(i)))
	}
	obj.curEntry = obj.dictionary
}

func (obj *compressor) Close() error {
	obj.writer.Write(obj.curEntry.key)
	obj.writer.Close()

	return obj.coder.FirstError()
}

func (obj *compressor) Write(p []byte) (n int, err error) {
	for _, input := range p {
		obj.addByte(input)
	}

	return len(p), obj.coder.FirstError()
}

func (obj *compressor) addByte(value byte) {
	nextEntry := obj.curEntry.next[int(value)]
	if nextEntry != nil {
		obj.curEntry = nextEntry
	} else {
		obj.writer.Write(obj.curEntry.key)

		key := Word(int(literalLimit) + obj.dictionarySize)
		if key < Reset {
			obj.curEntry.Add(value, key, obj.dictBuffer.entry(key))
			obj.dictionarySize++
		} else {
			obj.onKeySaturation()
		}

		obj.curEntry = obj.dictionary.next[value]
	}
}

func (obj *compressor) onKeySaturation() {
	obj.overtime++
	if obj.overtime > 1000 {
		obj.writer.Write(Reset)
		obj.resetDictionary()
		obj.overtime = 0
	}
}
