package archives

import (
	"fmt"
	"math"
	"sort"

	"github.com/inkyblackness/imgui-go/v2"

	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ui/gui"
)

type editIntegerVariableDialog struct {
	machine gui.ModalStateMachine
	view    *View

	service     *edit.GameStateService
	index       int
	varOverride bool
	info        archive.GameVariableInfo

	selectedEnum int
	opened       bool
}

func (dialog *editIntegerVariableDialog) Render() {
	if !dialog.opened {
		dialog.opened = true
		dialog.varOverride = dialog.service.IntegerVariableOverride(dialog.index)
		dialog.info = dialog.service.IntegerVariable(dialog.index)
		imgui.OpenPopup("Edit integer variable")
	}

	imgui.SetNextWindowSize(imgui.Vec2{X: 600, Y: 0}.Times(dialog.view.guiScale))
	if imgui.BeginPopupModalV("Edit integer variable", nil,
		imgui.WindowFlagsNoSavedSettings|imgui.WindowFlagsAlwaysAutoResize) {
		dialog.renderControls()
		imgui.Separator()
		if imgui.Button("OK") {
			dialog.applyChanges()
			dialog.close()
		}
		imgui.SameLine()
		if imgui.Button("Cancel") {
			dialog.close()
		}
		imgui.EndPopup()
	} else {
		dialog.machine.SetState(nil)
	}
}

func (dialog *editIntegerVariableDialog) applyChanges() {
	if dialog.varOverride {
		dialog.service.SetIntegerVariable(dialog.index, dialog.info)
	} else {
		dialog.service.DefaultIntegerVariable(dialog.index)
	}
}

func (dialog *editIntegerVariableDialog) close() {
	dialog.machine.SetState(nil)
	imgui.CloseCurrentPopup()
}

func (dialog *editIntegerVariableDialog) renderControls() {
	imgui.PushItemWidth(-150 * dialog.view.guiScale)
	values.RenderUnifiedCheckboxCombo(false, false, "Override", values.UnifierFor(dialog.varOverride), func(newValue bool) {
		dialog.varOverride = newValue
	})

	textFlags := imgui.InputTextFlagsNoUndoRedo
	if !dialog.varOverride {
		textFlags |= imgui.InputTextFlagsReadOnly
	}
	imgui.InputTextV("Name", &dialog.info.Name, textFlags, nil)
	imgui.InputTextMultilineV("Description", &dialog.info.Description, imgui.Vec2{X: 0, Y: 100 * dialog.view.guiScale},
		textFlags, nil)

	imgui.Separator()

	valueLabels := map[bool]string{false: "Numeric", true: "Enumerated"}
	isEnumerated := len(dialog.info.ValueNames) > 0
	if imgui.BeginCombo("Values", valueLabels[isEnumerated]) {
		if imgui.SelectableV(valueLabels[false], !isEnumerated, 0, imgui.Vec2{}) {
			dialog.info.ValueNames = nil
		}
		if imgui.SelectableV(valueLabels[true], isEnumerated, 0, imgui.Vec2{}) {
			dialog.info.ValueNames = map[int16]string{0: "Zero"}
			dialog.selectedEnum = 0
		}
		imgui.EndCombo()
	}
	if len(dialog.info.ValueNames) > 0 {
		gui.StepSliderInt("Index", &dialog.selectedEnum, math.MinInt16, math.MaxInt16)
		if dialog.varOverride {
			text := dialog.info.ValueNames[int16(dialog.selectedEnum)]
			if imgui.InputTextV("Value", &text, textFlags, nil) {
				if len(text) > 0 {
					dialog.info.ValueNames[int16(dialog.selectedEnum)] = text
				} else {
					delete(dialog.info.ValueNames, int16(dialog.selectedEnum))
					if len(dialog.info.ValueNames) == 0 {
						dialog.info.ValueNames[0] = "Zero"
					}
				}
			}
		} else {
			imgui.LabelText("Value", dialog.info.ValueNames[int16(dialog.selectedEnum)])
		}

		resetValue := dialog.info.ResetValueInt()
		enumEntries := toEnumEntries(dialog.info.ValueNames)
		linearInitIndex := linearIndex(enumEntries, resetValue)
		initUnifier := values.UnifierFor(linearInitIndex)
		values.RenderUnifiedCombo(dialog.info.Hardcoded || !dialog.varOverride, false, "Reset Value", initUnifier,
			func(u values.Unifier) int { return u.Unified().(int) },
			func(value int) string {
				if value < 0 {
					return fmt.Sprintf("%d: (unmapped)", resetValue)
				}
				entry := enumEntries[value]
				return fmt.Sprintf("%d: '%s'", entry.index, entry.value)
			},
			len(enumEntries),
			func(newValue int) {
				dialog.info.InitValue = &enumEntries[newValue].index
			})
	} else {
		minValue := math.MinInt16
		maxValue := math.MaxInt16

		if dialog.info.Limits != nil {
			minValue = int(dialog.info.Limits.Minimum)
			maxValue = int(dialog.info.Limits.Maximum)
		}
		if minValue > maxValue {
			maxValue = minValue
		}
		if dialog.varOverride {
			if gui.StepSliderInt("Minimum", &minValue, math.MinInt16, math.MaxInt16) {
				if minValue > maxValue {
					maxValue = minValue
				}
			}
			if gui.StepSliderInt("Maximum", &maxValue, math.MinInt16, math.MaxInt16) {
				if maxValue < minValue {
					minValue = maxValue
				}
			}
			if (minValue == math.MinInt16) && (maxValue == math.MaxInt16) {
				dialog.info.Limits = nil
			} else {
				if dialog.info.Limits == nil {
					dialog.info.Limits = &archive.GameVariableLimits{}
				}
				dialog.info.Limits.Minimum = int16(minValue)
				dialog.info.Limits.Maximum = int16(maxValue)
			}
		} else {
			imgui.LabelText("Minimum", fmt.Sprintf("%d", minValue))
			imgui.LabelText("Maximum", fmt.Sprintf("%d", maxValue))
		}

		resetValue := int(dialog.info.ResetValueInt())
		if dialog.info.Hardcoded || !dialog.varOverride {
			imgui.LabelText("Reset Value", fmt.Sprintf("%d", resetValue))
		} else if gui.StepSliderInt("Reset Value", &resetValue, minValue, maxValue) {
			resetValue16 := int16(resetValue)
			dialog.info.InitValue = &resetValue16
		}
	}

	imgui.PopItemWidth()
}

type enumEntry struct {
	index int16
	value string
}

func toEnumEntries(mapped map[int16]string) []enumEntry {
	minIndex := math.MaxInt16
	maxIndex := math.MinInt16
	for index := range mapped {
		if int(index) < minIndex {
			minIndex = int(index)
		}
		if int(index) > maxIndex {
			maxIndex = int(index)
		}
	}
	orderedEntries := make([]enumEntry, 0, len(mapped))
	for index, value := range mapped {
		orderedEntries = append(orderedEntries, enumEntry{
			index: index,
			value: value,
		})
	}
	sort.Slice(orderedEntries, func(a, b int) bool { return orderedEntries[a].index < orderedEntries[b].index })
	return orderedEntries
}

func linearIndex(entries []enumEntry, index int16) int {
	for linear, entry := range entries {
		if entry.index == index {
			return linear
		}
	}
	return -1
}
