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
