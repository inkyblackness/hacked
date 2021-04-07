package values

import (
	"fmt"
	"image/color"
	"math"

	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/ui/gui"
)

const (
	hintMultiple = "(multiple)"
	hintYes      = "Yes"
	hintNo       = "No"
)

// RenderUnifiedSliderInt renders a control for an unified integer value.
func RenderUnifiedSliderInt(readOnly bool, label string, unifier Unifier,
	intConverter func(Unifier) int, formatter func(int) string, min, max int, changeHandler func(int)) {
	var labelValue string
	var selectedString string
	selectedValue := -1
	if unifier.IsUnique() {
		selectedValue = intConverter(unifier)
		selectedString = formatter(selectedValue)
		labelValue = fmt.Sprintf(selectedString, selectedValue)
	} else if !unifier.IsEmpty() {
		selectedString = hintMultiple
		labelValue = selectedString
	}
	if readOnly {
		imgui.LabelText(label, labelValue)
	} else if gui.StepSliderIntV(label, &selectedValue, min, max, selectedString) {
		changeHandler(selectedValue)
	}
}

// RenderUnifiedCombo renders a control for a unified enumeration value.
func RenderUnifiedCombo(readOnly bool, label string, unifier Unifier,
	intConverter func(Unifier) int, formatter func(int) string, count int, changeHandler func(int)) {
	var selectedString string
	selectedIndex := -1
	if unifier.IsUnique() {
		selectedIndex = intConverter(unifier)
		selectedString = formatter(selectedIndex)
	} else if !unifier.IsEmpty() {
		selectedString = hintMultiple
	}
	if readOnly {
		imgui.LabelText(label, selectedString)
	} else if imgui.BeginCombo(label, selectedString) {
		for i := 0; i < count; i++ {
			entryString := formatter(i)

			if imgui.SelectableV(entryString, i == selectedIndex, 0, imgui.Vec2{}) {
				changeHandler(i)
			}
		}
		imgui.EndCombo()
	}
}

// RenderUnifiedCheckboxCombo renders a control for a unified boolean value.
func RenderUnifiedCheckboxCombo(readOnly bool, label string, unifier Unifier, changeHandler func(bool)) {
	selectedString := ""
	if unifier.IsUnique() {
		selectedValue := unifier.Unified().(bool)
		if selectedValue {
			selectedString = hintYes
		} else {
			selectedString = hintNo
		}
	} else if !unifier.IsEmpty() {
		selectedString = hintMultiple
	}
	if readOnly {
		imgui.LabelText(label, selectedString)
	} else if imgui.BeginCombo(label, selectedString) {
		for i, text := range []string{hintNo, hintYes} {
			if imgui.SelectableV(text, text == selectedString, 0, imgui.Vec2{}) {
				changeHandler(i != 0)
			}
		}
		imgui.EndCombo()
	}
}

// RotationInfo describes how the rotation is aligned in its zero position.
type RotationInfo struct {
	Horizontal bool
	Positive   bool
	Clockwise  bool
}

// RenderUnifiedRotation renders a control for rotation value.
func RenderUnifiedRotation(readOnly bool, label string, unifier Unifier,
	min, max int, info RotationInfo, changeHandler func(int)) {
	calcValuePercent := func(value int) float64 {
		valueRange := (max - min) + 1
		return float64(value) / float64(valueRange)
	}
	RenderUnifiedSliderInt(readOnly, label, unifier,
		func(u Unifier) int {
			unifiedValue := u.Unified().(int32)
			return int(unifiedValue)
		},
		func(value int) string {
			valuePercent := calcValuePercent(value)
			result := fmt.Sprintf("%3.02f°  - raw: %%d", valuePercent*360.0)
			return result
		},
		min, max,
		changeHandler)

	if (imgui.IsItemFocused() || imgui.IsItemActive() || imgui.IsItemHovered()) && unifier.IsUnique() {
		valuePercent := calcValuePercent(int(unifier.Unified().(int32)))
		valueRadian := valuePercent * (math.Pi * 2)

		lineColor := imgui.Packed(color.RGBA{R: 0x21, G: 0xFF, B: 0x43, A: 0xFF}) // color copied from style.
		lineThickness := float32(4.0)
		imgui.BeginTooltip()
		imgui.Dummy(imgui.Vec2{X: imgui.TextLineHeightWithSpacing() * 2, Y: imgui.TextLineHeightWithSpacing() * 2})
		dl := imgui.WindowDrawList()
		winTopLeft := imgui.WindowPos()
		winSize := imgui.WindowSize()
		center := winTopLeft.Plus(winSize.Times(0.5))
		circleRadius := ((winSize.X / 2) * 5) / 6
		dl.AddCircleV(center, circleRadius, lineColor, 100, lineThickness)
		// This code is specialized for Hacker-Yaw. To re-use it for other rotations, the target needs to be provided.
		targetX := 0.0
		targetY := 0.0
		if info.Horizontal {
			targetX = 1.0
			if !info.Positive {
				targetX = -1.0
			}
		} else {
			targetY = -1.0
			if !info.Positive {
				targetY = 1.0
			}
		}
		var target imgui.Vec2
		if info.Clockwise {
			target.X = center.X - float32((float64(circleRadius)*targetX)*math.Cos(valueRadian)+(float64(circleRadius)*targetY)*math.Sin(valueRadian))
			target.Y = center.Y + float32((float64(circleRadius)*targetY)*math.Cos(valueRadian)-(float64(circleRadius)*targetX)*math.Sin(valueRadian))
		} else {
			target.X = center.X + float32((float64(circleRadius)*targetX)*math.Cos(valueRadian)-(float64(circleRadius)*targetY)*math.Sin(valueRadian))
			target.Y = center.Y - float32((float64(circleRadius)*targetY)*math.Cos(valueRadian)+(float64(circleRadius)*targetX)*math.Sin(valueRadian))
		}
		dl.AddLineV(center, target, lineColor, lineThickness)
		imgui.EndTooltip()
	}
}
