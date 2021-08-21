package compression

// FrameDecoderBuilder implements the builder pattern for creating new instances of FrameDecoder.
// This builder can be reused during the processing of a movie. Every Build() call creates a new instance.
type FrameDecoderBuilder struct {
	width  int
	height int

	colorer TileColorFunction

	paletteLookupList []byte
	controlWords      []ControlWord
}

// NewFrameDecoderBuilder returns a new instance of a builder with given initial values.
func NewFrameDecoderBuilder(width, height int) *FrameDecoderBuilder {
	return &FrameDecoderBuilder{
		width:  width,
		height: height,
	}
}

// Build creates a new instance of a decoder with the most recent parameters.
func (builder *FrameDecoderBuilder) Build() *FrameDecoder {
	return newFrameDecoder(builder)
}

// ForStandardFrame registers the frame buffer and sets the standard coloring method.
func (builder *FrameDecoderBuilder) ForStandardFrame(frame []byte, stride int) {
	builder.colorer = StandardTileColorer(frame, stride)
}

// WithPaletteLookupList registers the new list.
func (builder *FrameDecoderBuilder) WithPaletteLookupList(list []byte) {
	builder.paletteLookupList = list
}

// WithControlWords registers the new word dictionary.
func (builder *FrameDecoderBuilder) WithControlWords(words []ControlWord) {
	builder.controlWords = words
}
