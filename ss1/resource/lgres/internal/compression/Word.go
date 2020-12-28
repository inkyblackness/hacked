package compression

// Word is a compression stream entry.
type Word uint16

const (
	bitsPerWord = uint(14)

	// EndOfStream indicates the last word in the compressed stream.
	EndOfStream = Word(0x3FFF)
	// Reset indicates a reset of the dictionary.
	Reset        = Word(0x3FFE)
	literalLimit = Word(0x0100)
)
