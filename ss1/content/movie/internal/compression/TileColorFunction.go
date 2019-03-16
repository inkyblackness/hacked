package compression

// TileColorFunction is one which shall color the pixel of a video frame tile using the provided values.
// mask is a packed field of 16 indices, each with a length of indexBitSize. These indices point into a
// value of lookupArray.
type TileColorFunction func(hTile int, vTile int, lookupArray []byte, mask uint64, indexBitSize uint64)
