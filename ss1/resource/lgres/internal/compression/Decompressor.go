package compression

import (
	"io"

	"github.com/inkyblackness/hacked/ss1/serial"
)

type decompressor struct {
	coder  serial.Coder
	reader *WordReader

	isEndOfStream  bool
	dictBuffer     dictEntryBuffer
	dictionary     *dictEntry
	dictionarySize int
	lastEntry      *dictEntry
	lookup         []*dictEntry

	scratch  []byte
	leftover []byte
}

// NewDecompressor creates a new decompressor instance over a reader.
func NewDecompressor(source io.Reader) io.Reader {
	coder := serial.NewDecoder(source)
	obj := &decompressor{
		coder:      coder,
		reader:     NewWordReader(coder),
		dictionary: rootDictEntry(),
		scratch:    make([]byte, 1024)}
	obj.resetDictionary()

	return obj
}

func (obj *decompressor) resetDictionary() {
	obj.dictionarySize = 0
	obj.lookup = make([]*dictEntry, 1024)
	obj.dictionary = rootDictEntry()
	for i := 0; i < 0x100; i++ {
		entry := obj.dictionary.Add(byte(i), Word(i), obj.dictBuffer.entry(Word(i)))
		obj.lookup[Word(i)] = entry
	}
	obj.lastEntry = obj.dictionary
}

func (obj *decompressor) Read(p []byte) (n int, err error) {
	requested := len(p)

	for n < requested && !obj.isEndOfStream && obj.coder.FirstError() == nil {
		n += obj.takeFromLeftover(p[n:])
		if n < requested {
			obj.readNextWord()
			n += obj.takeFromLeftover(p[n:])
		}
	}
	err = obj.coder.FirstError()
	if err == nil && obj.isEndOfStream {
		err = io.EOF
	}

	return
}

func (obj *decompressor) takeFromLeftover(target []byte) (provided int) {
	requested := len(target)
	available := len(obj.leftover)

	if available > 0 && requested > 0 {
		provided = available
		if provided > requested {
			provided = requested
		}
		copy(target[0:provided], obj.leftover)
		obj.leftover = obj.leftover[provided:]
	}

	return
}

func (obj *decompressor) readNextWord() {
	nextWord := obj.reader.Read()

	if obj.lastEntry.depth > len(obj.scratch) {
		obj.scratch = make([]byte, len(obj.scratch)+1024)
	}
	obj.leftover = obj.scratch[:obj.lastEntry.depth]
	obj.lastEntry.Data(obj.leftover)
	switch nextWord {
	case EndOfStream:
		obj.isEndOfStream = true
	case Reset:
		obj.resetDictionary()
	default:
		var nextEntry *dictEntry
		if int(nextWord) < len(obj.lookup) {
			nextEntry = obj.lookup[int(nextWord)]
		}

		switch {
		case nextEntry != nil:
			if obj.lastEntry.depth > 0 {
				obj.addToDictionary(nextEntry.FirstByte())
			}
			obj.lastEntry = nextEntry
		case nextWord >= literalLimit:
			nextValue := obj.lastEntry.FirstByte()
			obj.addToDictionary(nextValue)
			obj.lastEntry = obj.lastEntry.next[nextValue]
		default:
			nextValue := byte(nextWord)
			obj.addToDictionary(nextValue)
			obj.lastEntry = obj.dictionary.next[nextValue]
		}
	}
}

func (obj *decompressor) addToDictionary(value byte) {
	key := Word(int(literalLimit) + obj.dictionarySize)
	nextEntry := obj.lastEntry.Add(value, key, obj.dictBuffer.entry(key))
	if int(key) >= len(obj.lookup) {
		newLookup := make([]*dictEntry, len(obj.lookup)+1024)
		copy(newLookup, obj.lookup)
		obj.lookup = newLookup
	}
	obj.lookup[key] = nextEntry
	obj.dictionarySize++
}
