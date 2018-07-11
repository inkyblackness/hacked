package levels

import (
	"fmt"

	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/editor/event"
	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/editor/render"
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
	"github.com/inkyblackness/imgui-go"
)

// ControlView is the core view for level editing.
type ControlView struct {
	mod *model.Mod

	guiScale      float32
	commander     cmd.Commander
	eventListener event.Listener

	model controlViewModel
}

// NewControlView returns a new instance.
func NewControlView(mod *model.Mod, guiScale float32, commander cmd.Commander, eventListener event.Listener, eventRegistry event.Registry) *ControlView {
	view := &ControlView{
		mod:           mod,
		guiScale:      guiScale,
		commander:     commander,
		eventListener: eventListener,
		model:         freshControlViewModel(),
	}
	eventRegistry.RegisterHandler(view.onLevelSelectionSetEvent)
	return view
}

// WindowOpen returns the flag address, to be used with the main menu.
func (view *ControlView) WindowOpen() *bool {
	return &view.model.windowOpen
}

// SelectedLevel returns the currently selected level.
func (view *ControlView) SelectedLevel() int {
	return view.model.selectedLevel
}

// Render renders the view.
func (view *ControlView) Render(lvl *level.Level) {
	if view.model.restoreFocus {
		imgui.SetNextWindowFocus()
		view.model.restoreFocus = false
		view.model.windowOpen = true
	}
	if view.model.windowOpen {
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 400 * view.guiScale, Y: 300 * view.guiScale}, imgui.ConditionOnce)
		title := "Level Control"
		readOnly := !view.editingAllowed(lvl.ID())
		if readOnly {
			title += " (read-only)"
		}
		if imgui.BeginV(title+"###Level Control", view.WindowOpen(), imgui.WindowFlagsNoCollapse) {
			view.renderContent(lvl, readOnly)
		}
		imgui.End()
	}
}

var levelHeights = []string{
	"32 Tiles",
	"16 Tiles",
	"8 Tiles",
	"4 Tiles",
	"2 Tiles",
	"1 Tile",
	"1/2 Tile",
	"1/4 Tile",
}

func (view *ControlView) renderContent(lvl *level.Level, readOnly bool) {
	imgui.PushItemWidth(-200 * view.guiScale)
	gui.StepSliderInt("Active Level", &view.model.selectedLevel, 0, archive.MaxLevels-1)
	imgui.Separator()
	levelType := "Real World"
	if lvl.IsCyberspace() {
		levelType = "Cyberspace"
	}
	imgui.LabelText("Type", levelType)
	view.renderLevelHeight(lvl, readOnly)

	if !lvl.IsCyberspace() {
		view.renderTextureAtlas(lvl, readOnly)
		view.renderSurveillanceObjects(lvl, readOnly)
		view.renderHazards(lvl, readOnly)
		view.renderTextureAnimations(lvl, readOnly)
	}

	imgui.PopItemWidth()
}

func (view *ControlView) renderLevelHeight(lvl *level.Level, readOnly bool) {
	_, _, currentShift := lvl.Size()
	if readOnly {
		imgui.LabelText("Level Height", levelHeights[currentShift])
	} else {
		if imgui.BeginCombo("Level Height", levelHeights[currentShift]) {
			for shift, height := range levelHeights {
				if imgui.SelectableV(height, shift == int(currentShift), 0, imgui.Vec2{}) {
					view.requestSetZShift(lvl, shift)
				}
			}
			imgui.EndCombo()
		}
	}
}

func (view *ControlView) renderTextureAtlas(lvl *level.Level, readOnly bool) {
	imgui.Separator()

	atlas := lvl.TextureAtlas()
	fcwSelected := ""
	if (view.model.selectedAtlasIndex >= 0) && (view.model.selectedAtlasIndex < len(atlas)) {
		fcwSelected = fmt.Sprintf("%d", view.model.selectedAtlasIndex)
	}
	if imgui.BeginComboV("Level Textures", fcwSelected, imgui.ComboFlagHeightLarge) {
		for i := 0; i < len(atlas); i++ {
			key := resource.KeyOf(ids.LargeTextures.Plus(int(atlas[i])), resource.LangAny, 0)
			textureID := render.TextureIDForBitmapTexture(key)
			if imgui.SelectableV(fmt.Sprintf("%2d", i), view.model.selectedAtlasIndex == i, 0, imgui.Vec2{X: 0, Y: 66 * view.guiScale}) {
				view.model.selectedAtlasIndex = i
			}
			imgui.SameLine()
			imgui.Image(textureID, imgui.Vec2{X: 64 * view.guiScale, Y: 64 * view.guiScale})
			imgui.SameLine()
			textureType := "(F/C/W)"
			if i >= level.FloorCeilingTextureLimit {
				textureType = "(walls only)"
			}
			imgui.Text(textureType)
		}
		imgui.EndCombo()
	}
	gameTextureIndex := -1
	if (view.model.selectedAtlasIndex >= 0) && (view.model.selectedAtlasIndex < len(atlas)) {
		gameTextureIndex = int(atlas[view.model.selectedAtlasIndex])
	}
	if !readOnly && imgui.BeginComboV("Game Textures", fmt.Sprintf("%d", gameTextureIndex), imgui.ComboFlagHeightLarge) {
		for i := 0; i < world.MaxWorldTextures; i++ {
			key := resource.KeyOf(ids.LargeTextures.Plus(i), resource.LangAny, 0)
			textureID := render.TextureIDForBitmapTexture(key)
			if imgui.SelectableV(fmt.Sprintf("%3d", i), gameTextureIndex == i, 0, imgui.Vec2{X: 0, Y: 66 * view.guiScale}) {
				view.requestSetLevelTexture(lvl, view.model.selectedAtlasIndex, i)
			}
			imgui.SameLine()
			imgui.Image(textureID, imgui.Vec2{X: 64 * view.guiScale, Y: 64 * view.guiScale})
		}
		imgui.EndCombo()
	}
}

func (view *ControlView) renderSurveillanceObjects(lvl *level.Level, readOnly bool) {
	imgui.Separator()

	if imgui.BeginCombo("Surveillance Object", fmt.Sprintf("Object %d", view.model.selectedSurveillanceObjectIndex)) {
		for i := 0; i < level.SurveillanceObjectCount; i++ {
			if imgui.SelectableV(fmt.Sprintf("Object %d", i), i == view.model.selectedSurveillanceObjectIndex, 0, imgui.Vec2{}) {
				view.model.selectedSurveillanceObjectIndex = i
			}
		}
		imgui.EndCombo()
	}
	sources := lvl.SurveillanceSources()
	surrogates := lvl.SurveillanceSurrogates()
	limit := lvl.ObjectLimit()

	view.renderSliderInt(readOnly, "Surveillance Source", int(sources[view.model.selectedSurveillanceObjectIndex]),
		func(int) string { return "%d" },
		0, int(limit),
		func(newValue int) {
			view.requestSetSurveillanceSource(lvl, view.model.selectedSurveillanceObjectIndex, level.ObjectID(newValue))
		})
	view.renderSliderInt(readOnly, "Surveillance Surrogate", int(surrogates[view.model.selectedSurveillanceObjectIndex]),
		func(int) string { return "%d" },
		0, int(limit),
		func(newValue int) {
			view.requestSetSurveillanceSurrogate(lvl, view.model.selectedSurveillanceObjectIndex, level.ObjectID(newValue))
		})
}

func (view *ControlView) renderHazards(lvl *level.Level, readOnly bool) {
	imgui.Separator()

	parameters := lvl.Parameters()
	currentCeiling := currentCeilingHazard(parameters)
	if readOnly {
		imgui.LabelText("Ceiling Hazard", currentCeiling.title)
	} else if imgui.BeginCombo("Ceiling Hazard", currentCeiling.title) {
		for _, info := range ceilingHazards {
			if imgui.SelectableV(info.title, info.title == currentCeiling.title, 0, imgui.Vec2{}) {
				view.requestSetCeilingHazard(lvl, info)
			}
		}
		imgui.EndCombo()
	}
	view.renderSliderInt(readOnly, "Ceiling Hazard Level", int(parameters.CeilingHazardLevel),
		currentCeiling.formatter, 0, 255,
		func(newValue int) {
			view.requestSetCeilingHazardLevel(lvl, byte(newValue))
		})

	currentFloor := currentFloorHazard(parameters)
	if readOnly {
		imgui.LabelText("Floor Hazard", currentFloor.title)
	} else if imgui.BeginCombo("Floor Hazard", currentFloor.title) {
		for _, info := range floorHazards {
			if imgui.SelectableV(info.title, info.title == currentFloor.title, 0, imgui.Vec2{}) {
				view.requestSetFloorHazard(lvl, info)
			}
		}
		imgui.EndCombo()
	}
	view.renderSliderInt(readOnly, "Floor Hazard Level", int(parameters.FloorHazardLevel),
		currentFloor.formatter, 0, 255,
		func(newValue int) {
			view.requestSetFloorHazardLevel(lvl, byte(newValue))
		})
}

func (view *ControlView) renderTextureAnimations(lvl *level.Level, readOnly bool) {
	imgui.Separator()

	animations := lvl.TextureAnimations()
	selectedText := ""
	if (view.model.selectedTextureAnimationIndex >= 1) && (view.model.selectedTextureAnimationIndex < len(animations)) {
		selectedText = fmt.Sprintf("%d", view.model.selectedTextureAnimationIndex)
	}
	if imgui.BeginCombo("Texture Animation Group", selectedText) {
		for i := 1; i < len(animations); i++ {
			if imgui.SelectableV(fmt.Sprintf("%2d", i), view.model.selectedTextureAnimationIndex == i, 0, imgui.Vec2{}) {
				view.model.selectedTextureAnimationIndex = i
			}
		}
		imgui.EndCombo()
	}

	if (view.model.selectedTextureAnimationIndex >= 1) && (view.model.selectedTextureAnimationIndex < len(animations)) {
		animation := animations[view.model.selectedTextureAnimationIndex]
		view.renderSliderInt(readOnly, "Texture Animation Time", int(animation.FrameTime),
			func(int) string { return "%d msec" },
			0, 1000,
			func(newValue int) {
				view.requestSetTextureAnimationTime(lvl, view.model.selectedTextureAnimationIndex, uint16(newValue))
			})
		view.renderSliderInt(readOnly, "Texture Animation Frame Count", int(animation.FrameCount),
			func(int) string { return "%d" },
			0, 10,
			func(newValue int) {
				view.requestSetTextureAnimationFrameCount(lvl, view.model.selectedTextureAnimationIndex, byte(newValue))
			})
		if readOnly {
			imgui.LabelText("Texture Animation Loop Type", animation.LoopType.String())
		} else if imgui.BeginCombo("Texture Animation Loop Type", animation.LoopType.String()) {
			loopTypes := level.TextureAnimationLoopTypes()
			for _, loopType := range loopTypes {
				if imgui.SelectableV(loopType.String(), loopType == animation.LoopType, 0, imgui.Vec2{}) {
					view.requestSetTextureAnimationType(lvl, view.model.selectedTextureAnimationIndex, loopType)
				}
			}
			imgui.EndCombo()
		}
	}
}

func (view *ControlView) editingAllowed(id int) bool {
	gameStateData := view.mod.ModifiedBlocks(resource.LangAny, ids.GameState)
	isSavegame := (len(gameStateData) == 1) && (len(gameStateData[0]) == archive.GameStateSize) && (gameStateData[0][0x009C] > 0)
	moddedLevel := len(view.mod.ModifiedBlocks(resource.LangAny, ids.LevelResourcesStart.Plus(lvlids.PerLevel*id+lvlids.FirstUsed))) > 0

	return moddedLevel && !isSavegame
}

func (view *ControlView) renderSliderInt(readOnly bool, label string, selectedValue int,
	formatter func(int) string, min, max int, changeHandler func(int)) {

	selectedString := formatter(selectedValue)
	labelValue := fmt.Sprintf(selectedString, selectedValue)
	if readOnly {
		imgui.LabelText(label, labelValue)
	} else {
		if gui.StepSliderIntV(label, &selectedValue, min, max, selectedString) {
			changeHandler(selectedValue)
		}
	}
}

func (view *ControlView) requestSetZShift(lvl *level.Level, newValue int) {
	lvl.SetHeightShift(level.HeightShift(newValue))
	view.patchLevelResources(lvl, func() {})
}

func (view *ControlView) requestSetLevelTexture(lvl *level.Level, atlasIndex, worldTextureIndex int) {
	lvl.SetTextureAtlasEntry(atlasIndex, level.TextureIndex(worldTextureIndex))
	view.patchLevelResources(lvl, func() {
		view.model.selectedAtlasIndex = atlasIndex
	})
}

func (view *ControlView) requestSetSurveillanceSource(lvl *level.Level, objectIndex int, objectID level.ObjectID) {
	lvl.SetSurveillanceSource(objectIndex, objectID)
	view.patchLevelResources(lvl, func() {
		view.model.selectedSurveillanceObjectIndex = objectIndex
	})
}

func (view *ControlView) requestSetSurveillanceSurrogate(lvl *level.Level, objectIndex int, objectID level.ObjectID) {
	lvl.SetSurveillanceSurrogate(objectIndex, objectID)
	view.patchLevelResources(lvl, func() {
		view.model.selectedSurveillanceObjectIndex = objectIndex
	})
}

func (view *ControlView) requestSetCeilingHazard(lvl *level.Level, info ceilingHazardInfo) {
	parameters := lvl.Parameters()
	parameters.RadiationRegister = 0
	if info.radiationRegister {
		parameters.RadiationRegister = 2
	}
	view.patchLevelResources(lvl, func() {})
}

func (view *ControlView) requestSetCeilingHazardLevel(lvl *level.Level, value byte) {
	lvl.Parameters().CeilingHazardLevel = value
	view.patchLevelResources(lvl, func() {})
}

func (view *ControlView) requestSetFloorHazard(lvl *level.Level, info floorHazardInfo) {
	parameters := lvl.Parameters()
	parameters.BiohazardRegister = 0
	parameters.FloorHazardIsGravity = 0
	if info.isGravity {
		parameters.FloorHazardIsGravity = 1
	} else if info.biohazardRegister {
		parameters.BiohazardRegister = 2
	}
	view.patchLevelResources(lvl, func() {})
}

func (view *ControlView) requestSetFloorHazardLevel(lvl *level.Level, value byte) {
	lvl.Parameters().FloorHazardLevel = value
	view.patchLevelResources(lvl, func() {})
}

func (view *ControlView) requestSetTextureAnimationTime(lvl *level.Level, index int, value uint16) {
	lvl.TextureAnimations()[index].FrameTime = value
	view.patchLevelResources(lvl, func() {
		view.model.selectedTextureAnimationIndex = index
	})
}

func (view *ControlView) requestSetTextureAnimationFrameCount(lvl *level.Level, index int, value byte) {
	lvl.TextureAnimations()[index].FrameCount = value
	view.patchLevelResources(lvl, func() {
		view.model.selectedTextureAnimationIndex = index
	})
}

func (view *ControlView) requestSetTextureAnimationType(lvl *level.Level, index int, value level.TextureAnimationLoopType) {
	lvl.TextureAnimations()[index].LoopType = value
	view.patchLevelResources(lvl, func() {
		view.model.selectedTextureAnimationIndex = index
	})
}

func (view *ControlView) patchLevelResources(lvl *level.Level, extraRestoreState func()) {

	command := patchLevelDataCommand{
		restoreState: func() {
			view.model.restoreFocus = true
			view.setSelectedLevel(lvl.ID())
			extraRestoreState()
		},
	}

	newDataSet := lvl.EncodeState()
	for id, newData := range newDataSet {
		if len(newData) > 0 {
			resourceID := ids.LevelResourcesStart.Plus(lvlids.PerLevel*lvl.ID() + id)
			patch, changed, err := view.mod.CreateBlockPatch(resource.LangAny, resourceID, 0, newData)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				// TODO how to handle this? We're not expecting this, so crash and burn?
			} else if changed {
				command.patches = append(command.patches, patch)
			}
		}
	}

	view.commander.Queue(command)
}

func (view *ControlView) setSelectedLevel(id int) {
	view.eventListener.Event(LevelSelectionSetEvent{id: id})
}

func (view *ControlView) onLevelSelectionSetEvent(evt LevelSelectionSetEvent) {
	view.model.selectedLevel = evt.id
}
