package movie

// SubtitleHeaderSize is the size, in bytes, of the header structure
const SubtitleHeaderSize = 16

// SubtitleHeader is the header structure of a subtitle data
type SubtitleHeader struct {
	// Control specifies how to interpret the string content
	Control SubtitleControl

	Unknown0004 byte
	Unknown0005 [11]byte
}
