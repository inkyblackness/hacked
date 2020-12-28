package compression

type word uint16

const (
	bitsPerWord = uint(14)

	endOfStream  = word(0x3FFF)
	reset        = word(0x3FFE)
	literalLimit = word(0x0100)
)
