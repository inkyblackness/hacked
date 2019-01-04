package object_test

import (
	"fmt"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBitmap3DValueRange(t *testing.T) {
	for _, animation := range []bool{false, true} {
		for _, repeat := range []bool{false, true} {
			for frameNumber := uint16(0); frameNumber <= object.Bitmap3DFrameNumberLimit; frameNumber++ {
				for _, bitmapNumber := range []uint16{
					0,
					object.Bitmap3DBitmapNumberLimit / 4,
					object.Bitmap3DBitmapNumberLimit / 2,
					object.Bitmap3DBitmapNumberLimit / 2,
					object.Bitmap3DBitmapNumberLimit/2 + 4,
					object.Bitmap3DBitmapNumberLimit - 1,
					object.Bitmap3DBitmapNumberLimit,
				} {
					key := fmt.Sprintf("%v:%v:%v:%v", animation, repeat, frameNumber, bitmapNumber)
					value := object.Bitmap3D(0).
						WithAnimation(animation).
						WithRepeat(repeat).
						WithFrameNumber(frameNumber).
						WithBitmapNumber(bitmapNumber)

					assert.Equal(t, animation, value.Animation(), fmt.Sprintf("Animation wrong for %v", key))
					assert.Equal(t, repeat, value.Repeat(), fmt.Sprintf("Repeat wrong for %v", key))
					assert.Equal(t, frameNumber, value.FrameNumber(), fmt.Sprintf("Frame number wrong for %v", key))
					assert.Equal(t, bitmapNumber, value.BitmapNumber(), fmt.Sprintf("Bitmap number wrong for %v", key))
				}
			}
		}
	}
}
