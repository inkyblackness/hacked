package compression

// StandardTileColorer returns a TileColorFunction for given frame buffer.
//
// This colorer iterates over all pixel of a tile. It sets the pixel value if the
// corresponding palette value is > 0x00.
func StandardTileColorer(frame []byte, stride int) TileColorFunction {
	return func(hTile int, vTile int, lookupArray []byte, mask uint64, indexBitSize uint64) {
		start := vTile*TileSideLength*stride + hTile*TileSideLength
		singleMask := ^(^uint64(0) << indexBitSize)

		for i := 0; i < PixelPerTile; i++ {
			pixelValue := lookupArray[(mask>>(indexBitSize*uint64(i)))&singleMask]

			if pixelValue != 0x00 {
				offset := start + (i % TileSideLength) + stride*(i/TileSideLength)

				frame[offset] = pixelValue
			}
		}
	}
}
