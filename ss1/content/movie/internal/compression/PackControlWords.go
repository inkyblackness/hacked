package compression

import (
	"bytes"
	"encoding/binary"
)

// PackControlWords packs and encodes the given words into a byte stream.
func PackControlWords(words []ControlWord) []byte {
	buf := bytes.NewBuffer(nil)
	controlBytes := uint32(len(words) * bytesPerControlWord)

	_ = binary.Write(buf, binary.LittleEndian, controlBytes)
	if controlBytes > 0 {
		lastPacked := words[0].Packed(1)

		for _, word := range words[1:] {
			if (lastPacked.Value() == word) && (lastPacked.Times() < 0xFF) {
				lastPacked = word.Packed(lastPacked.Times() + 1)
			} else {
				_ = binary.Write(buf, binary.LittleEndian, uint32(lastPacked))
				lastPacked = word.Packed(1)
			}
		}
		_ = binary.Write(buf, binary.LittleEndian, uint32(lastPacked))
	}

	return buf.Bytes()
}
