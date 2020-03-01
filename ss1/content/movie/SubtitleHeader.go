package movie

// SubtitleDefaultTextOffset is the typical offset, in bytes, to text content.
const SubtitleDefaultTextOffset = 16

// SubtitleHeader is the header structure of a subtitle data
type SubtitleHeader struct {
	// Control specifies how to interpret the string content
	Control SubtitleControl
	// TextOffset specifies the offset, in bytes from beginning of entry data, where the text starts.
	TextOffset uint32
}
