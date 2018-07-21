package compression

// FrameDecoder is for decoding compressed frames with the help of data streams.
// A new instance of a decoder is created with a FrameDecoderBuilder.
type FrameDecoder struct {
	horizontalTiles int
	verticalTiles   int

	colorer TileColorFunction

	paletteLookupList []byte
	controlWords      []ControlWord
}

func newFrameDecoder(builder *FrameDecoderBuilder) *FrameDecoder {
	listLen := len(builder.paletteLookupList)
	decoder := &FrameDecoder{
		horizontalTiles:   builder.width / TileSideLength,
		verticalTiles:     builder.height / TileSideLength,
		colorer:           builder.colorer,
		paletteLookupList: make([]byte, listLen+16),
		controlWords:      builder.controlWords}

	copy(decoder.paletteLookupList[:listLen], builder.paletteLookupList)

	return decoder
}

// Decode reads the provided streams to paint a frame.
func (decoder *FrameDecoder) Decode(bitstreamData []byte, maskstreamData []byte) {
	bitstream := NewBitstreamReader(bitstreamData)
	maskstream := NewMaskstreamReader(maskstreamData)
	lastControl := ControlWord(0)

	for vTile := 0; vTile < decoder.verticalTiles && !bitstream.Exhausted(); vTile++ {
		for hTile := 0; hTile < decoder.horizontalTiles && !bitstream.Exhausted(); hTile++ {
			control := decoder.readNextControlWord(bitstream)

			if control.Type() == CtrlUnknown {
				panic("Unknown control in use")
			} else if control.Type() == CtrlRepeatPrevious {
				control = lastControl
			}

			if control.Type() == CtrlSkip {
				skipCount := bitstream.Read(5)
				bitstream.Advance(5)
				if skipCount == 0x1F {
					hTile = decoder.horizontalTiles
				} else {
					hTile += int(skipCount)
				}
			} else {
				decoder.colorTile(hTile, vTile, control, maskstream)
			}

			lastControl = control
		}
	}
}

func (decoder *FrameDecoder) readNextControlWord(bitstream *BitstreamReader) ControlWord {
	controlIndex := bitstream.Read(12)
	control := decoder.controlWords[controlIndex]

	if control.IsLongOffset() {
		bitstream.Advance(8)
		for control.IsLongOffset() {
			bitstream.Advance(4)
			offset := bitstream.Read(4)
			controlIndex = control.LongOffset() + offset
			control = decoder.controlWords[controlIndex]
		}
	}
	bitstream.Advance(control.Count())

	return control
}

func (decoder *FrameDecoder) colorTile(hTile, vTile int, control ControlWord, maskstream *MaskstreamReader) {
	param := control.Parameter()

	switch control.Type() {
	case CtrlColorTile2ColorsStatic:
		{
			decoder.colorer(hTile, vTile, []byte{byte(param & 0xFF), byte(param >> 8 & 0xFF)}, 0xAAAA, 1)
		}
	case CtrlColorTile2ColorsMasked:
		{
			decoder.colorer(hTile, vTile, []byte{byte(param & 0xFF), byte(param >> 8 & 0xFF)}, maskstream.Read(2), 1)
		}
	case CtrlColorTile4ColorsMasked:
		{
			decoder.colorer(hTile, vTile, decoder.paletteLookupList[param:param+4], maskstream.Read(4), 2)
		}
	case CtrlColorTile8ColorsMasked:
		{
			decoder.colorer(hTile, vTile, decoder.paletteLookupList[param:param+8], maskstream.Read(6), 3)
		}
	case CtrlColorTile16ColorsMasked:
		{
			decoder.colorer(hTile, vTile, decoder.paletteLookupList[param:param+16], maskstream.Read(8), 4)
		}
	}
}
