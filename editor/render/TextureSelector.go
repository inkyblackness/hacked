package render

import (
	"fmt"

	"github.com/inkyblackness/imgui-go/v2"

	"github.com/inkyblackness/hacked/editor/graphics"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// TextureSelector renders a horizontal list of game textures, with selection capability.
func TextureSelector(label string, width float32, guiScale float32,
	count int, selectedIndex int, cache *graphics.TextureCache, keyResolver func(int) resource.Key,
	tooltipText func(int) string,
	changeCallback func(int)) {
	if imgui.BeginChildV(label, imgui.Vec2{X: width * guiScale, Y: 100 * guiScale}, true,
		imgui.WindowFlagsHorizontalScrollbar|imgui.WindowFlagsNoScrollWithMouse) {
		for i := 0; i < count; i++ {
			key := keyResolver(i)
			imgui.PushStyleVarVec2(imgui.StyleVarWindowPadding, imgui.Vec2{X: 0, Y: 0})
			if imgui.BeginChildV(fmt.Sprintf("%3d", i), imgui.Vec2{X: 80 * guiScale, Y: 64 * guiScale}, false,
				imgui.WindowFlagsNoNav|imgui.WindowFlagsNoScrollWithMouse) {
				imgui.BeginGroup()
				if imgui.SelectableV("", selectedIndex == i, 0, imgui.Vec2{X: 0, Y: 64 * guiScale}) {
					changeCallback(i)
				}
				imgui.SameLine()
				TextureImage(fmt.Sprintf("%3d", i), cache, key, imgui.Vec2{X: 64 * guiScale, Y: 64 * guiScale})
				imgui.EndGroup()
				if imgui.IsItemHovered() {
					text := tooltipText(i)
					if len(text) > 0 {
						imgui.SetTooltip(text)
					}
				}
			}
			imgui.EndChild()
			imgui.PopStyleVar()
			imgui.SameLine()
		}
	}
	imgui.EndChild()
}
