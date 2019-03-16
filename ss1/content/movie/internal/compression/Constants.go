package compression

const (
	bytesPerControlWord = 3

	// TileSideLength is the number of pixel per side of a square frame tile
	TileSideLength = 4
	// PixelPerTile is the number of pixel within one frame tile
	PixelPerTile = TileSideLength * TileSideLength
)
