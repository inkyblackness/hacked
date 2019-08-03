package compression

import (
	"fmt"
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
	hTiles int
	vTiles int
	stride int

	lastFrame []byte
	deltas    []frameDelta
}

// NewSceneEncoder returns a new instance.
func NewSceneEncoder(width, height int) *SceneEncoder {
	e := &SceneEncoder{
		hTiles: width / TileSideLength,
		vTiles: height / TileSideLength,
		stride: width,
	}
	e.lastFrame = make([]byte, e.vTiles*TileSideLength*e.stride)
	return e
}

// AddFrame registers a further frame to the scene.
func (e *SceneEncoder) AddFrame(frame []byte) error {
	var delta frameDelta
	isFirstFrame := len(e.deltas) == 0
	for vTile := 0; vTile < e.vTiles; vTile++ {
		vStart := vTile * e.stride * TileSideLength
		for hTile := 0; hTile < e.hTiles; hTile++ {
			delta.tiles = append(delta.tiles, e.deltaTile(isFirstFrame, vStart+(hTile*TileSideLength), frame))
		}
	}
	e.deltas = append(e.deltas, delta)

	e.copyLastFrame(frame)
	return nil
}

func (e *SceneEncoder) deltaTile(isFirstFrame bool, offset int, frame []byte) tileDelta {
	var delta tileDelta
	for y := 0; y < TileSideLength; y++ {
		start := offset + (y * e.stride)
		for x := 0; x < TileSideLength; x++ {
			pixel := frame[start+x]
			if isFirstFrame || (pixel != e.lastFrame[start+x]) {
				delta[y*TileSideLength+x] = pixel
			}
		}
	}
	return delta
}

func (e *SceneEncoder) copyLastFrame(frame []byte) {
	hPixel := e.hTiles * TileSideLength
	for y := 0; y < e.vTiles*TileSideLength; y++ {
		start := y * e.stride
		copy(e.lastFrame[start:start+hPixel], frame[start:])
	}
}

// Encode processes all the previously registered frames and creates the necessary components for decoding.
func (e *SceneEncoder) Encode() (words []ControlWord, paletteLookupBuffer []byte, frames []EncodedFrame, err error) {
	var wordSequencer ControlWordSequencer
	tileColorOpsPerFrame := make([][]TileColorOp, len(e.deltas))
	frames = make([]EncodedFrame, len(e.deltas))

	var paletteLookupGenerator PaletteLookupGenerator
	for frameIndex := 0; frameIndex < len(e.deltas); frameIndex++ {
		delta := e.deltas[frameIndex]

		for _, tile := range delta.tiles {
			paletteLookupGenerator.Add(tile)
		}
	}
	paletteLookup := paletteLookupGenerator.Generate()
	paletteLookupBuffer = paletteLookup.Buffer()
	if len(paletteLookupBuffer) > 0x1FFFF {
		err = fmt.Errorf("palette lookup is too big: %vB", len(paletteLookupBuffer))
		return
	}

	for frameIndex := 0; frameIndex < len(e.deltas); frameIndex++ {
		outFrame := &frames[frameIndex]
		delta := e.deltas[frameIndex]
		controlStatistics := make(map[int]int)

		lastOp := TileColorOp{Type: CtrlUnknown}
		for tileIndex, tile := range delta.tiles {
			var op TileColorOp
			paletteIndex, pal, mask := paletteLookup.Lookup(tile)
			palSize := len(pal)
			// TODO: determine how to implement skip (and skip row)...
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
				outFrame.Maskstream = writeMaskstream(outFrame.Maskstream, 2, mask)
			case palSize <= 4:
				op.Type = CtrlColorTile4ColorsMasked
				op.Offset = uint32(paletteIndex)
				outFrame.Maskstream = writeMaskstream(outFrame.Maskstream, 4, mask)
			case palSize <= 8:
				op.Type = CtrlColorTile8ColorsMasked
				op.Offset = uint32(paletteIndex)
				outFrame.Maskstream = writeMaskstream(outFrame.Maskstream, 6, mask)
			default:
				op.Type = CtrlColorTile16ColorsMasked
				op.Offset = uint32(paletteIndex)
				outFrame.Maskstream = writeMaskstream(outFrame.Maskstream, 8, mask)
			}
			/* todo: see if this skipping could be moved to op sequencer */
			if op.Type != CtrlSkip && (tileIndex%e.hTiles) != 0 && lastOp == op {
				op = TileColorOp{Type: CtrlRepeatPrevious}
			} else {
				lastOp = op
			}
			/**/
			controlStatistics[int(op.Type)]++
			err = wordSequencer.Add(op)
			if err != nil {
				return nil, nil, nil, err
			}
			tileColorOpsPerFrame[frameIndex] = append(tileColorOpsPerFrame[frameIndex], op)
		}
		fmt.Printf("Out %v control statistics: %v\n", frameIndex, controlStatistics)
	}

	sequence, err := wordSequencer.Sequence()
	if err != nil {
		return nil, nil, nil, err
	}
	sequence.HTiles = uint32(e.hTiles)
	words = sequence.ControlWords()
	paletteLookupBuffer = paletteLookup.Buffer()
	for frameIndex, ops := range tileColorOpsPerFrame {
		frames[frameIndex].Bitstream, err = sequence.BitstreamFor(ops)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	return
}

func writeMaskstream(s []byte, bytes int, mask uint64) []byte {
	result := make([]byte, len(s), len(s)+bytes)
	copy(result, s)
	for b := 0; b < bytes; b++ {
		result = append(result, byte(mask>>(uint(b)*8)))
	}
	return result
}
