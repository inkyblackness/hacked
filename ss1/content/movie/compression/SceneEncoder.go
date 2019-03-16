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
	frames = make([]EncodedFrame, len(e.deltas))
	for frameIndex := 0; frameIndex < len(e.deltas); frameIndex++ {
		outFrame := &frames[frameIndex]
		delta := e.deltas[frameIndex]
		var bitstream BitstreamWriter

		for _, tile := range delta.tiles {
			paletteIndex := paletteLookupWriter.Write(tile[:])

			controlIndex := len(words)
			words = append(words, ControlWordOf(12, CtrlColorTile16ColorsMasked, paletteIndex))
			outFrame.Maskstream = writeMaskstream(outFrame.Maskstream, 8, 0xFEDCBA9876543210)
			bitstream.Write(12, uint32(controlIndex))
		}

		outFrame.Bitstream = bitstream.Buffer()
	}
	paletteLookup = paletteLookupWriter.Buffer

	// This needs to be done in several stages:
	// - gather a list of all necessary resulting tile operations (final palette index & op) over all frames
	// - these requests need also to be stored per frame
	// - after all control operations are requested, create a control word snapshot
	// - this snapshot has all the low-level control words created, as well as a mapping index from
	//   requested control operation to low-level index
	// - iterate again per frame, and find low-level index for requested operation - or, rather
	// - create bitstream out of requested list of control operations.

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
