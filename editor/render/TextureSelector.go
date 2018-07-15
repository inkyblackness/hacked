package render

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/imgui-go"
)

// TextureSelector renders a horizontal list of game textures, with selection capability.
func TextureSelector(label string, width float32, guiScale float32,
	count int, selectedIndex int, keyResolver func(int) resource.Key,
	tooltipText func(int) string,
	changeCallback func(int)) {
	if imgui.BeginChildV(label, imgui.Vec2{X: width * guiScale, Y: 110 * guiScale}, true,
		imgui.WindowFlagsHorizontalScrollbar|imgui.WindowFlagsNoScrollWithMouse) {
		for i := 0; i < count; i++ {
			key := keyResolver(i)
			textureID := TextureIDForBitmapTexture(key)
			if imgui.BeginChildV(fmt.Sprintf("%3d", i), imgui.Vec2{X: 80 * guiScale, Y: 80 * guiScale}, false, 0) {
				imgui.BeginGroup()
				if imgui.SelectableV("", selectedIndex == i, 0, imgui.Vec2{X: 0, Y: 70 * guiScale}) {
					changeCallback(i)
				}
				imgui.SameLine()
				imgui.Image(textureID, imgui.Vec2{X: 64 * guiScale, Y: 64 * guiScale})
				imgui.EndGroup()
				if imgui.IsItemHovered() {
					text := tooltipText(i)
					if len(text) > 0 {
						imgui.SetTooltip(text)
					}
				}
			}
			imgui.EndChild()
			imgui.SameLine()
		}
	}
	imgui.EndChild()
}
