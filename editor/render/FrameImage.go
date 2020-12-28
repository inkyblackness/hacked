package render

import (
	"math"

	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/editor/graphics"
)

// FrameImage renders an image centered and fitted within the given size.
func FrameImage(label string, frameCache *graphics.FrameCache, key graphics.FrameCacheKey, size imgui.Vec2) {
	textureID := TextureIDForBitmapFrame(key)

	imgui.PushStyleColor(imgui.StyleColorChildBg, imgui.Vec4{X: 0, Y: 0, Z: 0, W: 1})
	imgui.PushStyleVarVec2(imgui.StyleVarWindowPadding, imgui.Vec2{X: 0, Y: 0})
	if imgui.BeginChildV(label, size, false,
		imgui.WindowFlagsNoNav|imgui.WindowFlagsNoInputs|imgui.WindowFlagsNoScrollWithMouse|
			imgui.WindowFlagsNoScrollbar) {
		texture := frameCache.Texture(key)
		if texture != nil {
			var uv imgui.Vec2
			uv.X, uv.Y = texture.UV()
			width, height := texture.Size()

			scaleFactor := float32(math.Min(float64(size.X/width), float64(size.Y/height)))
			imageSize := imgui.Vec2{X: width * scaleFactor, Y: height * scaleFactor}

			bufferSize := imgui.Vec2{X: (size.X - imageSize.X) / 2, Y: (size.Y - imageSize.Y) / 2}
			imgui.SetCursorPos(bufferSize)

			imgui.ImageV(textureID, imageSize, imgui.Vec2{}, uv,
				imgui.Vec4{X: 1, Y: 1, Z: 1, W: 1}, imgui.Vec4{X: 0, Y: 0, Z: 0, W: 0})
		}
	}
	imgui.EndChild()
	imgui.PopStyleVar()
	imgui.PopStyleColor()
}
