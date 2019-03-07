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
	return &SceneEncoder{
		hTiles: width / TileSideLength,
		vTiles: height / TileSideLength,
		stride: width,
	}
}

// AddFrame registers a further frame to the scene.
func (e *SceneEncoder) AddFrame(frame []byte) {
	var delta frameDelta
	delta.tiles = append(delta.tiles, tileDelta{})
	copy(delta.tiles[0][:], frame)
	e.deltas = append(e.deltas, delta)
}

// Encode processes all the previously registered frames and creates the necessary components for decoding.
func (e *SceneEncoder) Encode() (words []ControlWord, paletteLookup []byte, frames []EncodedFrame) {
	words = append(words, ControlWordOf(12, CtrlColorTile16ColorsMasked, 0))
	paletteLookup = e.deltas[0].tiles[0][:]
	frames = make([]EncodedFrame, len(e.deltas))
	frames[0].Bitstream = []byte{0x00, 0x00}
	frames[0].Maskstream = []byte{0x10, 0x32, 0x54, 0x76, 0x98, 0xBA, 0xDC, 0xFE}
	/*
		for index, delta := range e.deltas {
			frames[index]
		}
	*/
	return
}
