package object_test

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBitmap3DValueRange(t *testing.T) {
	for _, animation := range []bool{false, true} {
		for _, repeat := range []bool{false, true} {
			for frameNumber := uint16(0); frameNumber <= object.Bitmap3DFrameNumberLimit; frameNumber++ {
				for _, bitmapNumber := range []uint16{0, object.Bitmap3DBitmapNumberLimit, object.Bitmap3DBitmapNumberLimit / 2} {
					value := object.Bitmap3D(0).
						WithAnimation(animation).
						WithRepeat(repeat).
						WithFrameNumber(frameNumber).
						WithBitmapNumber(bitmapNumber)

					assert.Equal(t, animation, value.Animation(), "Animation wrong")
					assert.Equal(t, repeat, value.Repeat(), "Repeat wrong")
					assert.Equal(t, frameNumber, value.FrameNumber(), "Frame number wrong")
					assert.Equal(t, bitmapNumber, value.BitmapNumber(), "Bitmap number wrong")
				}
			}
		}
	}
}
