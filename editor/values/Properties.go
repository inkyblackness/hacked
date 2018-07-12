package values

import (
	"fmt"

	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
)

// RenderUnifiedSliderInt renders a control for an unified integer value.
func RenderUnifiedSliderInt(readOnly, multiple bool, label string, unifier Unifier,
	intConverter func(Unifier) int, formatter func(int) string, min, max int, changeHandler func(int)) {

	var labelValue string
	var selectedString string
	selectedValue := -1
	if unifier.IsUnique() {
		selectedValue = intConverter(unifier)
		selectedString = formatter(selectedValue)
		labelValue = fmt.Sprintf(selectedString, selectedValue)
	} else if multiple {
		selectedString = "(multiple)"
		labelValue = selectedString
	}
	if readOnly {
		imgui.LabelText(label, labelValue)
	} else {
		if gui.StepSliderIntV(label, &selectedValue, min, max, selectedString) {
			changeHandler(selectedValue)
		}
	}
}

// RenderUnifiedCombo renders a control for a unified enumeration value.
func RenderUnifiedCombo(readOnly, multiple bool, label string, unifier Unifier,
	intConverter func(Unifier) int, formatter func(int) string, count int, changeHandler func(int)) {
	var selectedString string
	selectedIndex := -1
	if unifier.IsUnique() {
		selectedIndex = intConverter(unifier)
		selectedString = formatter(selectedIndex)
	} else if multiple {
		selectedString = "(multiple)"
	}
	if readOnly {
		imgui.LabelText(label, selectedString)
	} else {
		if imgui.BeginCombo(label, selectedString) {
			for i := 0; i < count; i++ {
				entryString := formatter(i)

				if imgui.SelectableV(entryString, i == selectedIndex, 0, imgui.Vec2{}) {
					changeHandler(i)
				}
			}
			imgui.EndCombo()
		}
	}
}

// RenderUnifiedCheckboxCombo renders a control for a unified boolean value.
func RenderUnifiedCheckboxCombo(readOnly, multiple bool, label string, unifier Unifier, changeHandler func(bool)) {
	selectedString := ""
	if unifier.IsUnique() {
		selectedValue := unifier.Unified().(bool)
		if selectedValue {
			selectedString = "Yes"
		} else {
			selectedString = "No"
		}
	} else if multiple {
		selectedString = "(multiple)"
	}
	if readOnly {
		imgui.LabelText(label, selectedString)
	} else {
		if imgui.BeginCombo(label, selectedString) {
			for i, text := range []string{"No", "Yes"} {
				if imgui.SelectableV(text, text == selectedString, 0, imgui.Vec2{}) {
					changeHandler(i != 0)
				}
			}
			imgui.EndCombo()
		}
	}
}
