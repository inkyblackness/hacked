package movie

// DataType identifies entries
type DataType byte

const (
	// endOfMedia marks the last entry
	endOfMedia = DataType(0)
	// LowResVideo for low resolution (low compression) video
	LowResVideo = DataType(0x21)
	// HighResVideo for high resolution (high compression) video
	HighResVideo = DataType(0x79)
	// Audio marks an audio entry.
	Audio = DataType(2)
	// Subtitle control
	Subtitle = DataType(3)
	// Palette data
	Palette = DataType(4)
	// PaletteReset is a zero-byte entry immediately before a Palette entry.
	PaletteReset = DataType(0x4C)
	// PaletteLookupList for high compression video.
	PaletteLookupList = DataType(5)
	// ControlDictionary for high compression video
	ControlDictionary = DataType(0x0D)
)
