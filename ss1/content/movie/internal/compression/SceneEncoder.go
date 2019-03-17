package compression

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
			paletteIndex := paletteLookupWriter.Write(tile[:])
			op := TileColorOp{Type: CtrlColorTile16ColorsMasked, Offset: paletteIndex}
			outFrame.Maskstream = writeMaskstream(outFrame.Maskstream, 8, 0xFEDCBA9876543210)

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

func writeMaskstream(s []byte, bytes int, mask uint64) []byte {
	result := make([]byte, len(s), len(s)+bytes)
	copy(result, s)
	for b := 0; b < bytes; b++ {
		result = append(result, byte(mask>>(uint(b)*8)))
	}
	return result
}
