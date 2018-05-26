package resource

// ContentType identifies how resource data shall be interpreted.
type ContentType byte

const (
	// Palette refers to color tables.
	Palette = ContentType(0x00)
	// Text refers to texts.
	Text = ContentType(0x01)
	// Bitmap refers to images.
	Bitmap = ContentType(0x02)
	// Font refers to font descriptions.
	Font = ContentType(0x03)
	// Animation refers to graphical animations.
	Animation = ContentType(0x04)
	// Sound refers to audio samples.
	Sound = ContentType(0x07)
	// Geometry refers to 3D models.
	Geometry = ContentType(0x0F)
	// Movie refers to audio logs and cutscenes.
	Movie = ContentType(0x11)
	// Archive refers to archive data.
	Archive = ContentType(0x30)
)
