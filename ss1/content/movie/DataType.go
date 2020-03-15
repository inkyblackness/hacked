package movie

// DataType identifies entries
type DataType byte

const (
	// dataTypeEndOfMedia marks the last entry
	dataTypeEndOfMedia = DataType(0)
	// DataTypeLowResVideo for low resolution (low compression) video
	DataTypeLowResVideo = DataType(0x21)
	// DataTypeHighResVideo for high resolution (high compression) video
	DataTypeHighResVideo = DataType(0x79)
	// DataTypeAudio marks an audio entry.
	DataTypeAudio = DataType(2)
	// DataTypeSubtitle control
	DataTypeSubtitle = DataType(3)
	// Palette data
	DataTypePalette = DataType(4)
	// DataTypePaletteReset is a zero-byte entry immediately before a Palette entry.
	DataTypePaletteReset = DataType(0x4C)
	// DataTypePaletteLookupList for high compression video.
	DataTypePaletteLookupList = DataType(5)
	// DataTypeControlDictionary for high compression video
	DataTypeControlDictionary = DataType(0x0D)
)
