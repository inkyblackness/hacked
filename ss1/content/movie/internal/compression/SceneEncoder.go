package compression

import (
	"math/bits"
	"sort"
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
	for vTile := 0; vTile < e.vTiles; vTile++ {
		vStart := vTile * e.stride * TileSideLength
		for hTile := 0; hTile < e.hTiles; hTile++ {
			delta.tiles = append(delta.tiles, e.deltaTile(vStart+(hTile*TileSideLength), frame))
		}
	}
	e.deltas = append(e.deltas, delta)

	e.copyLastFrame(frame)
	return nil
}

func (e *SceneEncoder) deltaTile(offset int, frame []byte) tileDelta {
	var delta tileDelta
	for y := 0; y < TileSideLength; y++ {
		start := offset + (y * e.stride)
		for x := 0; x < TileSideLength; x++ {
			pixel := frame[start+x]
			if pixel != e.lastFrame[start+x] {
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
func (e *SceneEncoder) Encode() (words []ControlWord, paletteLookup []byte, frames []EncodedFrame, err error) {
	var paletteLookupWriter PaletteLookupWriter
	var wordSequencer ControlWordSequencer
	tileColorOpsPerFrame := make([][]TileColorOp, len(e.deltas))
	frames = make([]EncodedFrame, len(e.deltas))

	for frameIndex := 0; frameIndex < len(e.deltas); frameIndex++ {
		outFrame := &frames[frameIndex]
		delta := e.deltas[frameIndex]

		for _, tile := range delta.tiles {
			var op TileColorOp
			pal, mask := maskOfTile(tile)
			palSize := len(pal)
			// TODO: determine how to implement skip (and skip row)...
			switch {
			case palSize == 1 && (pal[0] == 0x00):
				op.Type = CtrlSkip
			case palSize == 1:
				op.Type = CtrlColorTile2ColorsStatic
				op.Offset = uint32(pal[0])<<8 | uint32(pal[0])
			case palSize == 2 && mask == 0xAA:
				op.Type = CtrlColorTile2ColorsStatic
				op.Offset = uint32(pal[1])<<8 | uint32(pal[0])
			case palSize == 2 && mask == 0x55:
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
				paletteIndex := paletteLookupWriter.Write(pal)
				op.Type = CtrlColorTile4ColorsMasked
				op.Offset = paletteIndex
				outFrame.Maskstream = writeMaskstream(outFrame.Maskstream, 4, mask)
			case palSize <= 8:
				paletteIndex := paletteLookupWriter.Write(pal)
				op.Type = CtrlColorTile8ColorsMasked
				op.Offset = paletteIndex
				outFrame.Maskstream = writeMaskstream(outFrame.Maskstream, 6, mask)
			default:
				paletteIndex := paletteLookupWriter.Write(pal)
				op.Type = CtrlColorTile16ColorsMasked
				op.Offset = paletteIndex
				outFrame.Maskstream = writeMaskstream(outFrame.Maskstream, 8, mask)
			}

			err = wordSequencer.Add(op)
			if err != nil {
				return nil, nil, nil, err
			}
			tileColorOpsPerFrame[frameIndex] = append(tileColorOpsPerFrame[frameIndex], op)
		}
	}

	sequence, err := wordSequencer.Sequence()
	if err != nil {
		return nil, nil, nil, err
	}
	sequence.HTiles = uint32(e.hTiles)
	words = sequence.ControlWords()
	paletteLookup = paletteLookupWriter.Buffer
	for frameIndex, ops := range tileColorOpsPerFrame {
		frames[frameIndex].Bitstream, err = sequence.BitstreamFor(ops)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return
}

func maskOfTile(tile tileDelta) (pal []byte, mask uint64) {
	sorted := tile
	sort.Slice(sorted[:], func(a, b int) bool { return sorted[a] < sorted[b] })

	usedColors := 1
	last := sorted[PixelPerTile-1]
	for index := PixelPerTile - 1; index > 0; index-- {
		prev := index - 1
		if sorted[prev] == last {
			copy(sorted[prev:], sorted[index:])
		} else {
			usedColors++
		}
		last = sorted[prev]
	}

	var mapped [256]int
	for index, value := range sorted[0:usedColors] {
		mapped[value] = index
	}

	bitSize := uint(bits.Len(uint(usedColors - 1)))
	for index := PixelPerTile - 1; index >= 0; index-- {
		mask <<= bitSize
		mask |= uint64(mapped[tile[index]])
	}

	return sorted[0:usedColors], mask
}

func writeMaskstream(s []byte, bytes int, mask uint64) []byte {
	result := make([]byte, len(s), len(s)+bytes)
	copy(result, s)
	for b := 0; b < bytes; b++ {
		result = append(result, byte(mask>>(uint(b)*8)))
	}
	return result
}
