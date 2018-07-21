package compression

import (
	"bytes"
	"encoding/binary"
)

// UnpackControlWords reads from an encoded string of bytes a series of packed control words.
// If all is OK, the control words are returned. If the function returns an error, something
// could not be read/decoded properly.
func UnpackControlWords(data []byte) (words []ControlWord, err error) {
	if len(data) < 4 {
		return nil, FormatError
	}

	var controlBytes uint32
	reader := bytes.NewReader(data)

	err = binary.Read(reader, binary.LittleEndian, &controlBytes)
	if err != nil {
		return
	}
	if controlBytes%bytesPerControlWord != 0 {
		return nil, FormatError
	}
	wordCount := int(controlBytes / bytesPerControlWord)
	unpacked := 0

	words = make([]ControlWord, wordCount)
	for unpacked < wordCount {
		var packed PackedControlWord
		err = binary.Read(reader, binary.LittleEndian, &packed)
		if err != nil {
			return
		}
		times := packed.Times()
		if times > len(words)-unpacked {
			return nil, FormatError
		}
		for i := 0; i < times; i++ {
			words[unpacked+i] = packed.Value()
		}
		unpacked += times
	}
	return
}
