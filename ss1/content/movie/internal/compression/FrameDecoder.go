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
		controlWords:      builder.controlWords,
	}

	copy(decoder.paletteLookupList[:listLen], builder.paletteLookupList)

	return decoder
}

// Decode reads the provided streams to paint a frame.
func (decoder *FrameDecoder) Decode(bitstreamData []byte, maskstreamData []byte) error {
	bitstream := NewBitstreamReader(bitstreamData)
	maskstream := NewMaskstreamReader(maskstreamData)

	for vTile := 0; vTile < decoder.verticalTiles && !bitstream.Exhausted(); vTile++ {
		lastControl := ControlWordOf(0, CtrlUnknown, 0)
		for hTile := 0; hTile < decoder.horizontalTiles && !bitstream.Exhausted(); hTile++ {
			control, err := decoder.readNextControlWord(bitstream)
			if err != nil {
				return err
			}

			if control.Type() == CtrlRepeatPrevious {
				if hTile == 0 {
					return errCannotRepeatWordOnFirstTileOfRow
				}
				control = lastControl
			}
			if control.Type() == CtrlUnknown {
				return errUnknownControl
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
	return nil
}

func (decoder *FrameDecoder) readNextControlWord(bitstream *BitstreamReader) (ControlWord, error) {
	availableWords := uint32(len(decoder.controlWords))
	controlIndex := bitstream.Read(12)
	if controlIndex > availableWords {
		return ControlWordOf(0, CtrlUnknown, 0),
			wordIndexOutOfRangeError{
				Index:     controlIndex,
				Available: availableWords,
			}
	}

	control := decoder.controlWords[controlIndex]
	if control.IsLongOffset() {
		longOffsetCount := 0
		bitstream.Advance(8)
		for control.IsLongOffset() {
			longOffsetCount++
			if longOffsetCount > 100 {
				return ControlWordOf(0, CtrlUnknown, 0), errTooManyLongOffsets
			}
			bitstream.Advance(4)
			offset := bitstream.Read(4)
			controlIndex = control.LongOffset() + offset
			if controlIndex > availableWords {
				return ControlWordOf(0, CtrlUnknown, 0),
					wordIndexOutOfRangeError{
						Index:     controlIndex,
						Available: availableWords,
					}
			}
			control = decoder.controlWords[controlIndex]
		}
	}
	bitstream.Advance(control.Count())
	return control, nil
}

func (decoder *FrameDecoder) colorTile(hTile, vTile int, control ControlWord, maskstream *MaskstreamReader) {
	param := control.Parameter()

	switch control.Type() {
	case CtrlColorTile2ColorsStatic:
		decoder.colorer(hTile, vTile, []byte{byte(param & 0xFF), byte(param >> 8 & 0xFF)}, 0xAAAA, 1)
	case CtrlColorTile2ColorsMasked:
		decoder.colorer(hTile, vTile, []byte{byte(param & 0xFF), byte(param >> 8 & 0xFF)}, maskstream.Read(2), 1)
	case CtrlColorTile4ColorsMasked:
		decoder.colorer(hTile, vTile, decoder.paletteLookupList[param:param+4], maskstream.Read(4), 2)
	case CtrlColorTile8ColorsMasked:
		decoder.colorer(hTile, vTile, decoder.paletteLookupList[param:param+8], maskstream.Read(6), 3)
	case CtrlColorTile16ColorsMasked:
		decoder.colorer(hTile, vTile, decoder.paletteLookupList[param:param+16], maskstream.Read(8), 4)
	}
}
