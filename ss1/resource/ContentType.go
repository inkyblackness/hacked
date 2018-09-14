package resource

import "fmt"

// ContentType identifies how resource data shall be interpreted.
type ContentType byte

var contentTypeNames = map[ContentType]string{
	Palette:   "Palette",
	Text:      "Text",
	Bitmap:    "Bitmap",
	Font:      "Font",
	Animation: "Animation",
	Sound:     "Sound",
	Geometry:  "Geometry",
	Movie:     "Movie",
	Archive:   "Archive",
}

// String returns the textual representation of the type.
func (t ContentType) String() string {
	s, existing := contentTypeNames[t]
	if existing {
		return s
	}
	return fmt.Sprintf("Unknown%02X", int(t))
}

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
