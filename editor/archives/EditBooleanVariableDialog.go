package archives

import (
	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ui/gui"
)

type editBooleanVariableDialog struct {
	machine gui.ModalStateMachine
	view    *View

	service     *edit.GameStateService
	index       int
	varOverride bool
	info        archive.GameVariableInfo

	opened bool
}

func (dialog *editBooleanVariableDialog) Render() {
	if !dialog.opened {
		dialog.opened = true
		dialog.varOverride = dialog.service.BooleanVariableOverride(dialog.index)
		dialog.info = dialog.service.BooleanVariable(dialog.index)
		imgui.OpenPopup("Edit boolean variable")
	}

	imgui.SetNextWindowSize(imgui.Vec2{X: 600, Y: 0}.Times(dialog.view.guiScale))
	if imgui.BeginPopupModalV("Edit boolean variable", nil,
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

func (dialog *editBooleanVariableDialog) applyChanges() {
	if dialog.varOverride {
		_ = dialog.service.SetBooleanVariable(dialog.index, dialog.info)
	} else {
		_ = dialog.service.DefaultBooleanVariable(dialog.index)
	}
}

func (dialog *editBooleanVariableDialog) close() {
	dialog.machine.SetState(nil)
	imgui.CloseCurrentPopup()
}

func (dialog *editBooleanVariableDialog) renderControls() {
	imgui.PushItemWidth(-150 * dialog.view.guiScale)
	values.RenderUnifiedCheckboxCombo(false, "Override", values.UnifierFor(dialog.varOverride), func(newValue bool) {
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
	defaultValues := map[int]string{0: "False", 1: "True"}
	valueLabels := map[bool]string{false: "Default", true: "Special"}
	hasSpecialValues := dialog.info.ValueNames != nil
	if dialog.varOverride {
		if imgui.BeginCombo("Values", valueLabels[hasSpecialValues]) {
			if imgui.SelectableV(valueLabels[false], !hasSpecialValues, 0, imgui.Vec2{}) {
				dialog.info.ValueNames = nil
			}
			if imgui.SelectableV(valueLabels[true], hasSpecialValues, 0, imgui.Vec2{}) {
				dialog.info.ValueNames = map[int16]string{0: defaultValues[0], 1: defaultValues[1]}
			}
			imgui.EndCombo()
		}
	} else {
		imgui.LabelText("Values", valueLabels[hasSpecialValues])
	}
	hasSpecialValues = dialog.info.ValueNames != nil
	if hasSpecialValues {
		for i := int16(0); i < 2; i++ {
			text := dialog.info.ValueNames[i]
			if imgui.InputTextV(defaultValues[int(i)], &text, textFlags, nil) {
				dialog.info.ValueNames[i] = text
			}
		}
	} else {
		imgui.LabelText(defaultValues[0], defaultValues[0])
		imgui.LabelText(defaultValues[1], defaultValues[1])
	}

	initUnifier := values.UnifierFor(dialog.info.ResetValueInt())
	values.RenderUnifiedCombo(dialog.info.Hardcoded || !dialog.varOverride, "Reset Value", initUnifier,
		func(u values.Unifier) int { return int(u.Unified().(int16)) },
		func(value int) string {
			name, found := dialog.info.ValueNames[int16(value)]
			if found {
				return name
			}
			return defaultValues[value]
		},
		2,
		func(newValue int) {
			value16 := int16(newValue)
			dialog.info.InitValue = &value16
		})

	imgui.PopItemWidth()
}
