package gui

import (
	"github.com/inkyblackness/imgui-go/v2"
)

// StepSliderIntV creates a SliderInt with additional buttons to make single digit steps.
// min and max parameters are inclusive. Returns true on change.
func StepSliderIntV(label string, value *int, min, max int, format string) bool {
	imgui.PushID(label)
	changed := false
	if imgui.Button("-") && (*value > min) {
		*value--
		changed = true
	}
	innerSpacing := imgui.CurrentStyle().ItemInnerSpacing()
	imgui.SameLineV(0, innerSpacing.X)
	if imgui.Button("+") && (*value < max) {
		*value++
		changed = true
	}
	imgui.SameLineV(0, innerSpacing.X)
	value32 := int32(*value)
	if imgui.SliderIntV(label, &value32, int32(min), int32(max), format) {
		switch {
		case int(value32) < min:
			*value = min
		case int(value32) > max:
			*value = max
		default:
			*value = int(value32)
		}
		changed = true
	}
	imgui.PopID()
	return changed
}

// StepSliderInt calls StepSliderIntV(label, value, min, max, "%d")
func StepSliderInt(label string, value *int, min, max int) bool {
	return StepSliderIntV(label, value, min, max, "%d")
}
