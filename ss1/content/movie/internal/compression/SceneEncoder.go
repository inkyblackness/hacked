package compression

import (
	"context"
)

// EncodedFrame contains the streams of one compressed frame.
type EncodedFrame struct {
	Bitstream  []byte
	Maskstream []byte
}

type tileDelta [PixelPerTile]byte

type frameDelta struct {
	tiles []tileDelta
}

// SceneEncoder encodes an entire scene of bitmaps sharing the same palette.
type SceneEncoder struct {
	hTiles     int
	vTiles     int
	lineStride int
	tileStride int

	lastFrame []byte
	deltas    []frameDelta
}

// NewSceneEncoder returns a new instance.
func NewSceneEncoder(width, height int) *SceneEncoder {
	e := &SceneEncoder{
		hTiles:     width / TileSideLength,
		vTiles:     height / TileSideLength,
		lineStride: width,
	}
	e.tileStride = e.lineStride * TileSideLength
	e.lastFrame = make([]byte, e.vTiles*TileSideLength*e.lineStride)
	return e
}

// AddFrame registers a further frame to the scene.
func (e *SceneEncoder) AddFrame(frame []byte) error {
	if len(frame) != len(e.lastFrame) {
		return errInvalidFrameSize
	}
	var delta frameDelta
	isFirstFrame := len(e.deltas) == 0
	vStart := 0
	for vTile := 0; vTile < e.vTiles; vTile++ {
		tileStart := vStart
		for hTile := 0; hTile < e.hTiles; hTile++ {
			delta.tiles = append(delta.tiles, e.deltaTile(isFirstFrame, tileStart, frame))
			tileStart += TileSideLength
		}
		vStart += e.tileStride
	}
	e.deltas = append(e.deltas, delta)
	copy(e.lastFrame, frame)
	return nil
}

func (e *SceneEncoder) deltaTile(isFirstFrame bool, offset int, frame []byte) tileDelta {
	var delta tileDelta
	for y := 0; y < TileSideLength; y++ {
		start := offset + (y * e.lineStride)
		for x := 0; x < TileSideLength; x++ {
			pixel := frame[start+x]
			if isFirstFrame || (pixel != e.lastFrame[start+x]) {
				delta[y*TileSideLength+x] = pixel
			}
		}
	}
	return delta
}

// Encode processes all the previously registered frames and creates the necessary components for decoding.
func (e *SceneEncoder) Encode(ctx context.Context) (
	words []ControlWord, paletteLookupBuffer []byte, frames []EncodedFrame, err error) {
	var wordSequencer ControlWordSequencer
	tileColorOpsPerFrame := make([][]TileColorOp, len(e.deltas))
	paletteLookup, err := e.createPaletteLookup(ctx)
	if err != nil {
		return
	}

	paletteLookupBuffer = paletteLookup.Buffer()
	if len(paletteLookupBuffer) > 0x1FFFF {
		err = paletteLookupTooBigError{Size: len(paletteLookupBuffer)}
		return
	}

	frames = make([]EncodedFrame, len(e.deltas))
	for frameIndex := 0; frameIndex < len(e.deltas); frameIndex++ {
		var maskstreamWriter MaskstreamWriter
		outFrame := &frames[frameIndex]
		delta := e.deltas[frameIndex]

		lastOp := TileColorOp{Type: CtrlUnknown}
		for tileIndex, tile := range delta.tiles {
			var op TileColorOp
			paletteIndex, pal, mask := paletteLookup.Lookup(tile)
			palSize := len(pal)

			switch {
			case palSize == 1 && (pal[0] == 0x00):
				op.Type = CtrlSkip
			case palSize == 1:
				op.Type = CtrlColorTile2ColorsStatic
				op.Offset = uint32(pal[0])<<8 | uint32(pal[0])
			case palSize == 2 && mask == 0xAAAA && (pal[0] != 0x00) && (pal[1] != 0x00):
				op.Type = CtrlColorTile2ColorsStatic
				op.Offset = uint32(pal[1])<<8 | uint32(pal[0])
			case palSize == 2 && mask == 0x5555 && (pal[0] != 0x00) && (pal[1] != 0x00):
				op.Type = CtrlColorTile2ColorsStatic
				op.Offset = uint32(pal[0])<<8 | uint32(pal[1])
			case palSize <= 2:
				op.Type = CtrlColorTile2ColorsMasked
				if palSize == 2 {
					op.Offset = uint32(pal[1])
					op.Offset <<= 8
				}
				if palSize > 0 {
					op.Offset |= uint32(pal[0])
				}

				_ = maskstreamWriter.Write(2, mask)
			case palSize <= 4:
				op.Type = CtrlColorTile4ColorsMasked
				op.Offset = uint32(paletteIndex)
				_ = maskstreamWriter.Write(4, mask)
			case palSize <= 8:
				op.Type = CtrlColorTile8ColorsMasked
				op.Offset = uint32(paletteIndex)
				_ = maskstreamWriter.Write(6, mask)
			default:
				op.Type = CtrlColorTile16ColorsMasked
				op.Offset = uint32(paletteIndex)
				_ = maskstreamWriter.Write(8, mask)
			}

			if op.Type != CtrlSkip && (tileIndex%e.hTiles) != 0 && lastOp == op {
				op = TileColorOp{Type: CtrlRepeatPrevious}
			} else {
				lastOp = op
			}

			err = wordSequencer.Add(op)
			if err != nil {
				return nil, nil, nil, err
			}
			tileColorOpsPerFrame[frameIndex] = append(tileColorOpsPerFrame[frameIndex], op)
		}
		outFrame.Maskstream = maskstreamWriter.Buffer
	}

	wordSequence, err := wordSequencer.Sequence()
	if err != nil {
		return nil, nil, nil, err
	}
	wordSequence.HTiles = uint32(e.hTiles)
	words = wordSequence.ControlWords()
	for frameIndex, ops := range tileColorOpsPerFrame {
		frames[frameIndex].Bitstream, err = wordSequence.BitstreamFor(ops)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	return
}

func (e *SceneEncoder) createPaletteLookup(ctx context.Context) (PaletteLookup, error) {
	var paletteLookupGenerator PaletteLookupGenerator
	for _, delta := range e.deltas {
		for _, tile := range delta.tiles {
			paletteLookupGenerator.Add(tile)
		}
	}
	return paletteLookupGenerator.Generate(ctx)
}
