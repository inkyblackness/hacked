package object

// Bitmap3D describes how a bitmap object should be represented in 3D world.
type Bitmap3D uint16

const (
	bitmap3DBitmapNumberLowMask  uint16 = 0x03FF
	bitmap3DBitmapNumberHighMask uint16 = 0x8000
	bitmap3DBitmapNumberMask            = bitmap3DBitmapNumberLowMask | bitmap3DBitmapNumberHighMask
	bitmap3DAnimationMask        uint16 = 0x0400
	bitmap3DRepeatMask           uint16 = 0x0800
	bitmap3DFrameNumberMask      uint16 = 0x7000

	// Bitmap3DBitmapNumberLimit is the maximum amount of bitmaps.
	Bitmap3DBitmapNumberLimit uint16 = 0x07FF
	// Bitmap3DFrameNumberLimit is the maximum amount of frames.
	Bitmap3DFrameNumberLimit uint16 = 0x0007
)

// WithBitmapNumber returns a Bitmap3D with the given value.
func (bmp Bitmap3D) WithBitmapNumber(value uint16) Bitmap3D {
	bitmapValue := value & bitmap3DBitmapNumberLowMask
	if value >= Bitmap3DBitmapNumberLimit {
		bitmapValue = bitmap3DBitmapNumberMask
	} else if value > bitmap3DBitmapNumberLowMask {
		bitmapValue |= bitmap3DBitmapNumberHighMask
	}
	return Bitmap3D((uint16(bmp) & ^bitmap3DBitmapNumberMask) | bitmapValue)
}

// BitmapNumber returns the number of bitmaps.
func (bmp Bitmap3D) BitmapNumber() uint16 {
	return ((uint16(bmp) & bitmap3DBitmapNumberHighMask) >> 5) | (uint16(bmp) & bitmap3DBitmapNumberLowMask)
}

// WithFrameNumber returns a Bitmap3D with the given value.
func (bmp Bitmap3D) WithFrameNumber(value uint16) Bitmap3D {
	frameValue := (value % (Bitmap3DFrameNumberLimit + 1)) << 12
	return Bitmap3D((uint16(bmp) & ^bitmap3DFrameNumberMask) | frameValue)
}

// FrameNumber returns the frame number.
func (bmp Bitmap3D) FrameNumber() uint16 {
	return (uint16(bmp) & bitmap3DFrameNumberMask) >> 12
}

// WithAnimation returns a Bitmap3D with the given value.
func (bmp Bitmap3D) WithAnimation(value bool) Bitmap3D {
	animationValue := uint16(0)
	if value {
		animationValue = bitmap3DAnimationMask
	}
	return Bitmap3D((uint16(bmp) & ^bitmap3DAnimationMask) | animationValue)
}

// Animation returns whether this object is animated.
func (bmp Bitmap3D) Animation() bool {
	return (uint16(bmp) & bitmap3DAnimationMask) != 0
}

// WithRepeat returns a Bitmap3D with the given value.
func (bmp Bitmap3D) WithRepeat(value bool) Bitmap3D {
	repeatValue := uint16(0)
	if value {
		repeatValue = bitmap3DRepeatMask
	}
	return Bitmap3D((uint16(bmp) & ^bitmap3DRepeatMask) | repeatValue)
}

// Repeat returns whether this object repeats its animation.
func (bmp Bitmap3D) Repeat() bool {
	return (uint16(bmp) & bitmap3DRepeatMask) != 0
}
