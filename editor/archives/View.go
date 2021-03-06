package archives

import (
	"fmt"
	"math"
	"strings"

	"github.com/inkyblackness/imgui-go/v3"

	"github.com/inkyblackness/hacked/editor/values"
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlids"
	"github.com/inkyblackness/hacked/ss1/content/archive/level/lvlobj"
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/numbers"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/citadel"
	"github.com/inkyblackness/hacked/ss1/world/ids"
	"github.com/inkyblackness/hacked/ui/gui"
)

const (
	hintUnknown = "???"
)

var variableContextList = map[edit.VariableBaseContextIdentifier]struct {
	title string
	info  string
}{
	edit.VariableContextCitadel: {
		title: "Citadel Mission",
		info: "Variables are listed based on use on Citadel.\n" +
			"Useful for Citadel-based savegames and mods, such as 'New Game+' variants.",
	},
	edit.VariableContextEngine: {
		title: "Engine",
		info: "Variables are listed based on the hardcoded engine behaviour.\n" +
			"Useful for dedicated missions.",
	},
}

// View provides edit controls for the archive.
type View struct {
	registry         cmd.Registry
	gameStateService *edit.GameStateService
	mod              *world.Mod
	textCache        *text.Cache
	cp               text.Codepage

	modalStateMachine gui.ModalStateMachine
	guiScale          float32
	commander         cmd.Commander

	model viewModel
}

// NewArchiveView returns a new instance.
func NewArchiveView(registry cmd.Registry,
	gameStateService *edit.GameStateService, mod *world.Mod,
	textCache *text.Cache, cp text.Codepage,
	modalStateMachine gui.ModalStateMachine, guiScale float32, commander cmd.Commander) *View {
	view := &View{
		registry:         registry,
		gameStateService: gameStateService,
		mod:              mod,
		textCache:        textCache,
		cp:               cp,

		modalStateMachine: modalStateMachine,
		guiScale:          guiScale,
		commander:         commander,

		model: freshViewModel(),
	}
	return view
}

// WindowOpen returns the flag address, to be used with the main menu.
func (view *View) WindowOpen() *bool {
	return &view.model.windowOpen
}

// Render renders the view.
func (view *View) Render() {
	if view.model.restoreFocus {
		imgui.SetNextWindowFocus()
		view.model.restoreFocus = false
		view.model.windowOpen = true
	}
	if view.model.windowOpen {
		imgui.SetNextWindowSizeV(imgui.Vec2{X: 650 * view.guiScale, Y: 400 * view.guiScale}, imgui.ConditionFirstUseEver)
		if imgui.BeginV("Archive", view.WindowOpen(), imgui.WindowFlagsNoCollapse) {
			view.renderContent()
		}
		imgui.End()
	}
}

func (view *View) renderContent() {
	imgui.BeginTabBar("archive-tab")

	if imgui.BeginTabItem("Levels") {
		view.renderLevelsContent()
		imgui.EndTabItem()
	}
	if imgui.BeginTabItem("Game State") {
		view.renderGameStateContent()
		imgui.EndTabItem()
	}

	imgui.EndTabBar()
}

func (view *View) renderLevelsContent() {
	imgui.BeginChildV("Levels", imgui.Vec2{X: -100 * view.guiScale, Y: 0}, true, 0)
	for id := 0; id < archive.MaxLevels; id++ {
		inMod := view.hasLevelInMod(id)
		info := fmt.Sprintf("%d", id)
		if !inMod {
			if view.hasGameStateInMod() && (id == view.effectiveGameState().CurrentLevel()) {
				// current level should really be in the archive.
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1.0, Y: 0.0, Z: 0.0, W: 1.0})
			} else {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1.0, Y: 1.0, Z: 1.0, W: 0.8})
			}
			info += " (not in mod, read-only)"
		}
		if imgui.SelectableV(info, id == view.model.selectedLevel, 0, imgui.Vec2{}) {
			view.model.selectedLevel = id
		}
		if !inMod {
			imgui.PopStyleColor()
		}
	}
	imgui.EndChild()
	imgui.SameLine()
	imgui.BeginGroup()
	if imgui.ButtonV("Clear", imgui.Vec2{X: -1, Y: 0}) {
		view.requestClearLevel(view.model.selectedLevel)
	}
	if imgui.ButtonV("Remove", imgui.Vec2{X: -1, Y: 0}) {
		view.requestRemoveLevel(view.model.selectedLevel)
	}
	imgui.EndGroup()
}

func (view *View) renderGameStateContent() {
	data := view.gameStateData()
	modStateIsDefaulting := false

	imgui.PushItemWidth(-350 * view.guiScale)
	if data.isPresent() {
		modState := data.toInstance()
		if !modState.IsSavegame() {
			modStateIsDefaulting = modState.IsDefaulting()
			if modStateIsDefaulting {
				if imgui.Button("Override") {
					view.requestSetGameState(citadel.DefaultGameState().Raw())
				}
				if imgui.IsItemHovered() {
					imgui.SetTooltip("Values of new game archives are only considered by special engines.\n" +
						"Overriding them will only work in those.")
				}
			} else if imgui.Button("Remove") {
				view.requestSetGameState(archive.ZeroGameStateData())
			}
		}
	}

	readOnly := false
	var editState *archive.GameState
	if !data.isPresent() || modStateIsDefaulting {
		editState = citadel.DefaultGameState()
		readOnly = true
	} else {
		editState = data.imported().toInstance()
	}

	onChange := func() {
		view.requestSetGameState(editState.Raw())
	}
	if imgui.TreeNodeV("Game State: General", imgui.TreeNodeFlagsDefaultOpen|imgui.TreeNodeFlagsFramed) {
		imgui.LabelText("Hacker Name", "\""+editState.HackerName(view.cp)+"\"")
		view.createPropertyControls(readOnly, editState.Instance, func(key string, modifier func(uint32) uint32) {
			view.setInterpreterValueKeyed(editState.Instance, key, modifier)
			onChange()
		})
		imgui.TreePop()
	}
	if imgui.TreeNodeV("Inventory", imgui.TreeNodeFlagsFramed) {
		view.createInventoryControls(readOnly, editState, onChange)
		imgui.TreePop()
	}
	if imgui.TreeNodeV("Hard-/Software", imgui.TreeNodeFlagsFramed) {
		view.createWareControls(readOnly, editState, onChange)
		imgui.TreePop()
	}
	if imgui.TreeNodeV("Game State: Variables", imgui.TreeNodeFlagsFramed) {
		view.createVariableControls(readOnly, editState, onChange)
		imgui.TreePop()
	}
	if imgui.TreeNodeV("Game State: Messages", imgui.TreeNodeFlagsFramed) {
		view.createMessageControls(readOnly, editState, onChange)
		imgui.TreePop()
	}

	imgui.PopItemWidth()
}

type gameStateData struct{ raw []byte }

func (data gameStateData) isPresent() bool {
	return len(data.raw) > 0
}

func (data gameStateData) toInstance() *archive.GameState {
	return archive.NewGameState(data.raw)
}

func (data gameStateData) snapshot() gameStateData {
	copied := make([]byte, len(data.raw))
	copy(copied, data.raw)
	return gameStateData{raw: copied}
}

func (data gameStateData) imported() gameStateData {
	// Avoid truncating data when working with savegames.
	if len(data.raw) >= archive.GameStateSize {
		return data.snapshot()
	}
	copied := archive.ZeroGameStateData()
	copy(copied, data.raw)
	return gameStateData{raw: copied}
}

func (view *View) gameStateData() gameStateData {
	raw := view.mod.ModifiedBlock(resource.LangAny, ids.GameState, 0)
	return gameStateData{raw: raw}
}

func (view *View) effectiveGameState() *archive.GameState {
	state := view.gameStateData().toInstance()
	if state.IsDefaulting() {
		return citadel.DefaultGameState()
	}
	return state
}

func (view *View) hasGameStateInMod() bool {
	return len(view.mod.ModifiedBlocks(resource.LangAny, ids.GameState)) > 0
}

func (view *View) hasLevelInMod(id int) bool {
	return len(view.mod.ModifiedBlocks(resource.LangAny, ids.LevelResourcesStart.Plus(lvlids.PerLevel*id+lvlids.FirstUsed))) > 0
}

func (view *View) requestClearLevel(id int) {
	if (id >= 0) && (id < archive.MaxLevels) {
		command := setArchiveDataCommand{
			model:         &view.model,
			selectedLevel: id,
			newData:       make(map[resource.ID][]byte),
			oldData:       make(map[resource.ID][]byte),
		}

		if !view.hasGameStateInMod() {
			command.newData[ids.ArchiveName] = text.DefaultCodepage().Encode("Starting Game | by InkyBlackness HackEd")
			command.newData[ids.GameState] = archive.ZeroGameStateData()
		}

		param := level.EmptyLevelParameters{
			Cyberspace: world.IsConsideredCyberspaceByDefault(id),
		}
		effectiveGameState := view.effectiveGameState()
		if id == effectiveGameState.CurrentLevel() {
			posX, posY := effectiveGameState.HackerMapPosition()
			param.MapModifier = func(x, y int, entry *level.TileMapEntry) {
				if (x == int(posX.Tile())) && (y == int(posY.Tile())) {
					entry.Type = level.TileTypeOpen
				}
			}
		}

		levelData := level.EmptyLevelData(param)
		levelIDBegin := ids.LevelResourcesStart.Plus(lvlids.PerLevel * id)
		for offset, newData := range &levelData {
			if (offset < lvlids.FirstUsed) || (offset >= lvlids.FirstUnused) {
				continue
			}
			resourceID := levelIDBegin.Plus(offset)
			oldData := view.mod.ModifiedBlock(resource.LangAny, resourceID, 0)
			if len(oldData) > 0 {
				command.oldData[resourceID] = oldData
			}
			if len(newData) > 0 {
				command.newData[resourceID] = newData
			}
		}

		view.commander.Queue(command)
	}
}

func (view *View) requestRemoveLevel(id int) {
	if (id >= 0) && (id < archive.MaxLevels) && view.hasLevelInMod(id) {
		command := setArchiveDataCommand{
			model:         &view.model,
			selectedLevel: id,
			newData:       make(map[resource.ID][]byte),
			oldData:       make(map[resource.ID][]byte),
		}

		levelIDBegin := ids.LevelResourcesStart.Plus(lvlids.PerLevel * id)
		for offset := lvlids.FirstUsed; offset < lvlids.PerLevel; offset++ {
			resourceID := levelIDBegin.Plus(offset)
			oldData := view.mod.ModifiedBlock(resource.LangAny, resourceID, 0)
			if len(oldData) > 0 {
				command.oldData[resourceID] = oldData
			}
		}

		view.commander.Queue(command)
	}
}

func (view *View) requestSetGameState(newData []byte) {
	command := setArchiveDataCommand{
		model:         &view.model,
		selectedLevel: view.model.selectedLevel,
		newData:       make(map[resource.ID][]byte),
		oldData:       make(map[resource.ID][]byte),
	}

	command.oldData[ids.GameState] = view.gameStateData().snapshot().raw
	command.newData[ids.GameState] = newData

	view.commander.Queue(command)
}

func (view *View) createPropertyControls(readOnly bool, rootInterpreter *interpreters.Instance,
	updater func(string, func(uint32) uint32)) {
	objTypeRenderer := values.ObjectTypeControlRenderer{
		Meta:      view.mod.ObjectProperties(),
		TextCache: view.textCache,
	}

	var processInterpreter func(string, *interpreters.Instance)
	processInterpreter = func(path string, interpreter *interpreters.Instance) {
		for _, key := range interpreter.Keys() {
			fullKey := path + key
			unifier := values.NewUnifier()
			unifier.Add(int32(interpreter.Get(key)))
			simplifier := values.StandardSimplifier(readOnly, fullKey, unifier,
				func(modifier func(uint32) uint32) {
					updater(fullKey, modifier)
				}, objTypeRenderer)

			interpreter.Describe(key, simplifier)
		}

		for _, key := range interpreter.ActiveRefinements() {
			fullKey := path + key
			if len(fullKey) > 0 {
				imgui.Separator()
				imgui.Text(fullKey + ":")
			}
			processInterpreter(fullKey+".", interpreter.Refined(key))
		}
	}
	processInterpreter("", rootInterpreter)
}

func (view *View) setInterpreterValueKeyed(instance *interpreters.Instance, key string, modifier func(uint32) uint32) {
	resolvedInterpreter := instance
	keys := strings.Split(key, ".")
	keyCount := len(keys)
	if len(keys) > 1 {
		for _, subKey := range keys[:keyCount-1] {
			resolvedInterpreter = resolvedInterpreter.Refined(subKey)
		}
	}
	resolvedInterpreter.Set(keys[keyCount-1], modifier(resolvedInterpreter.Get(keys[keyCount-1])))
}

func (view *View) createInventoryControls(readOnly bool, gameState *archive.GameState, onChange func()) {
	if imgui.TreeNodeV("Weapons", imgui.TreeNodeFlagsFramed) {
		for i := 0; i < archive.InventoryWeaponSlots; i++ {
			if i > 0 {
				imgui.Separator()
			}
			imgui.PushID(fmt.Sprintf("weapon%d", i))
			view.createInventoryWeaponSlotControls(readOnly, gameState.InventoryWeaponSlot(i), onChange)
			imgui.PopID()
		}
		imgui.TreePop()
	}
	if imgui.TreeNodeV("Ammo", imgui.TreeNodeFlagsFramed) {
		for i := 0; i < archive.AmmoTypeCount; i++ {
			imgui.Separator()
			imgui.PushID(fmt.Sprintf("ammo%d", i))
			view.createInventoryAmmoControls(readOnly, gameState.InventoryAmmo(i), onChange)
			imgui.PopID()
		}
		imgui.TreePop()
	}
	if imgui.TreeNodeV("Grenades", imgui.TreeNodeFlagsFramed) {
		for i := 0; i < archive.GrenadeTypeCount; i++ {
			imgui.Separator()
			imgui.PushID(fmt.Sprintf("grenade%d", i))
			view.createInventoryGrenadeControls(readOnly, gameState.InventoryGrenade(i), onChange)
			imgui.PopID()
		}
		imgui.TreePop()
	}
	if imgui.TreeNodeV("Patches", imgui.TreeNodeFlagsFramed) {
		for i := 0; i < archive.GrenadeTypeCount; i++ {
			imgui.PushID(fmt.Sprintf("patch%d", i))
			view.createPatchStateControls(readOnly, gameState.PatchState(i), onChange)
			imgui.PopID()
		}
		imgui.TreePop()
	}
	if imgui.TreeNodeV("General", imgui.TreeNodeFlagsFramed) {
		for i := 0; i < archive.GeneralInventorySlotCount; i++ {
			imgui.PushID(fmt.Sprintf("general%d", i))
			view.createGeneralInventoryControls(readOnly, gameState.GeneralInventorySlot(i), onChange)
			imgui.PopID()
		}
		imgui.TreePop()
	}
}

func (view *View) createWareControls(readOnly bool, gameState *archive.GameState, onChange func()) {
	if imgui.TreeNodeV("Hardware", imgui.TreeNodeFlagsFramed) {
		for i := 0; i < archive.HardwareTypeCount; i++ {
			imgui.PushID(fmt.Sprintf("hw%d", i))
			view.createHardwareControls(readOnly, gameState.HardwareState(i), onChange)
			imgui.PopID()
		}
		imgui.TreePop()
	}

	if imgui.TreeNodeV("Software", imgui.TreeNodeFlagsFramed) {
		for i := 0; i < archive.SoftwareOffenseTypeCount; i++ {
			imgui.PushID(fmt.Sprintf("swoffense%d", i))
			view.createVersionedSoftwareControls(readOnly, gameState.VersionedOffenseSoftwareState(i), 0, onChange)
			imgui.PopID()
		}
		for i := 0; i < archive.SoftwareDefenseTypeCount; i++ {
			imgui.PushID(fmt.Sprintf("swdefense%d", i))
			view.createVersionedSoftwareControls(readOnly, gameState.VersionedDefenseSoftwareState(i), 1, onChange)
			imgui.PopID()
		}
		for i := 0; i < archive.SoftwareOneshotTypeCount; i++ {
			imgui.PushID(fmt.Sprintf("swoneshot%d", i))
			view.createCountedSoftwareControls(readOnly, gameState.OneshotSoftwareState(i), 2, onChange)
			imgui.PopID()
		}
		imgui.TreePop()
	}
}

func (view *View) createHardwareControls(readOnly bool, state archive.HardwareState, onChange func()) {
	values.RenderUnifiedSliderInt(readOnly,
		fmt.Sprintf("Version (%s)", view.indexedName(object.ClassHardware, state.Index)),
		values.UnifierFor(state.Version()),
		func(u values.Unifier) int {
			return u.Unified().(int)
		}, func(value int) string {
			return "%d"
		}, 0, 4,
		func(newValue int) {
			state.SetVersion(newValue)
			onChange()
		})

	values.RenderUnifiedCheckboxCombo(readOnly, "Is Active", values.UnifierFor(state.IsActive()), func(newActive bool) {
		state.SetActive(newActive)
		onChange()
	})
}

func (view *View) createVersionedSoftwareControls(readOnly bool, sw archive.VersionedSoftwareState, subclass int, onChange func()) {
	values.RenderUnifiedSliderInt(readOnly,
		fmt.Sprintf("Version (%s)", view.tripleName(object.TripleFrom(int(object.ClassSoftware), subclass, sw.Index))),
		values.UnifierFor(sw.Version()),
		func(u values.Unifier) int {
			return u.Unified().(int)
		}, func(value int) string {
			return "%d"
		}, 0, 0xFF,
		func(newValue int) {
			sw.SetVersion(newValue)
			onChange()
		})
}

func (view *View) createCountedSoftwareControls(readOnly bool, sw archive.CountedSoftwareState, subclass int, onChange func()) {
	values.RenderUnifiedSliderInt(readOnly,
		fmt.Sprintf("Count (%s)", view.tripleName(object.TripleFrom(int(object.ClassSoftware), subclass, sw.Index))),
		values.UnifierFor(sw.Count()),
		func(u values.Unifier) int {
			return u.Unified().(int)
		}, func(value int) string {
			return "%d"
		}, 0, 0xFF,
		func(newValue int) {
			sw.SetCount(newValue)
			onChange()
		})
}

func (view *View) createGeneralInventoryControls(readOnly bool, slot archive.GeneralInventorySlot, onChange func()) {
	values.RenderUnifiedSliderInt(readOnly,
		fmt.Sprintf("ObjectIndex (Slot %d)", slot.Index+1),
		values.UnifierFor(slot.ObjectID()),
		func(u values.Unifier) int {
			return int(u.Unified().(level.ObjectID))
		}, func(value int) string {
			return "%d"
		}, 0, 871, // TODO: get ID limit of current level
		func(newValue int) {
			slot.SetObjectID(level.ObjectID(newValue))
			onChange()
		})
}

func (view *View) createPatchStateControls(readOnly bool, patch archive.PatchState, onChange func()) {
	values.RenderUnifiedSliderInt(readOnly,
		fmt.Sprintf("Patch Count (%s)", view.indexedName(object.ClassDrug, patch.Index)),
		values.UnifierFor(patch.Count()),
		func(u values.Unifier) int {
			return u.Unified().(int)
		}, func(value int) string {
			return "%d"
		}, 0, 0xFF,
		func(newValue int) {
			patch.SetCount(newValue)
			onChange()
		})
}

func (view *View) createInventoryGrenadeControls(readOnly bool, grenade archive.InventoryGrenade, onChange func()) {
	values.RenderUnifiedSliderInt(readOnly,
		fmt.Sprintf("Grenade Count (%s)", view.indexedName(object.ClassGrenade, grenade.Index)),
		values.UnifierFor(grenade.Count()),
		func(u values.Unifier) int {
			return u.Unified().(int)
		}, func(value int) string {
			return "%d"
		}, 0, 0xFF,
		func(newValue int) {
			grenade.SetCount(newValue)
			onChange()
		})

	values.RenderUnifiedSliderInt(readOnly, "Timer", values.UnifierFor(grenade.TimerSetting()),
		func(u values.Unifier) int {
			return int(u.Unified().(archive.GrenadeTimerSetting))
		}, func(value int) string {
			setting := archive.GrenadeTimerSetting(value)
			return fmt.Sprintf("%.01fs  -- raw: %%d", setting.InSeconds())
		}, 0, archive.GrenadeTimerMaximum,
		func(newValue int) {
			grenade.SetTimerSetting(archive.GrenadeTimerSetting(newValue))
			onChange()
		})
}

func (view *View) createInventoryAmmoControls(readOnly bool, ammo archive.InventoryAmmo, onChange func()) {
	values.RenderUnifiedSliderInt(readOnly,
		fmt.Sprintf("Clip Count (%s)", view.indexedName(object.ClassAmmo, ammo.Index)),
		values.UnifierFor(ammo.FullClipCount()),
		func(u values.Unifier) int {
			return u.Unified().(int)
		}, func(value int) string {
			return "%d"
		}, 0, 0xFF,
		func(newValue int) {
			ammo.SetFullClipCount(newValue)
			onChange()
		})
	values.RenderUnifiedSliderInt(readOnly,
		"Extra rounds",
		values.UnifierFor(ammo.ExtraRoundsCount()),
		func(u values.Unifier) int {
			return u.Unified().(int)
		}, func(value int) string {
			return "%d"
		}, 0, 0xFF,
		func(newValue int) {
			ammo.SetExtraRoundsCount(newValue)
			onChange()
		})
}

func (view *View) createInventoryWeaponSlotControls(readOnly bool, slot archive.InventoryWeaponSlot, onChange func()) {
	isInUse := slot.IsInUse()
	inUseName := fmt.Sprintf("Weapon %d In Use", slot.Index+1)

	values.RenderUnifiedCheckboxCombo(readOnly, inUseName, values.UnifierFor(isInUse), func(newUse bool) {
		if newUse {
			slot.SetInUse(0, 0)
		} else {
			slot.SetFree()
		}
		onChange()
	})
	if !isInUse {
		return
	}

	currentWeaponTriple := slot.Triple()
	if imgui.BeginCombo("Type", view.tripleName(currentWeaponTriple)) {
		allTypes := view.mod.ObjectProperties().TriplesInClass(currentWeaponTriple.Class)
		for _, triple := range allTypes {
			if imgui.SelectableV(view.tripleName(triple), triple == currentWeaponTriple, 0, imgui.Vec2{}) {
				slot.SetInUse(triple.Subclass, triple.Type)
				onChange()
			}
		}
		imgui.EndCombo()
	}
	isEnergyWeapon := currentWeaponTriple.Subclass >= 4
	if isEnergyWeapon {
		temperature, chargeAndOverload := slot.WeaponState()
		values.RenderUnifiedSliderInt(readOnly, "Temperature", values.UnifierFor(int(temperature)),
			func(u values.Unifier) int {
				return u.Unified().(int)
			}, func(value int) string {
				return fmt.Sprintf("%.02f%%%%  -- raw: %%d", (float32(value)/float32(lvlobj.EnergyWeaponMaxTemperature))*100)
			}, 0x00, lvlobj.EnergyWeaponMaxTemperature,
			func(newValue int) {
				slot.SetWeaponState(byte(newValue), chargeAndOverload)
				onChange()
			})

		overloadState := chargeAndOverload & 0x80
		charge := chargeAndOverload & 0x7F
		values.RenderUnifiedCheckboxCombo(readOnly, "Overload", values.UnifierFor(overloadState != 0), func(newOverload bool) {
			if newOverload {
				slot.SetWeaponState(temperature, 0x80|charge)
			} else {
				slot.SetWeaponState(temperature, charge)
			}
			onChange()
		})
		values.RenderUnifiedSliderInt(readOnly, "Charge", values.UnifierFor(int(charge)),
			func(u values.Unifier) int {
				return u.Unified().(int)
			}, func(value int) string {
				return fmt.Sprintf("%.02f%%%%  -- raw: %%d", (float32(value)/float32(lvlobj.EnergyWeaponMaxEnergy))*100)
			}, 0x00, 0x7F,
			func(newValue int) {
				slot.SetWeaponState(temperature, overloadState|byte(newValue))
				onChange()
			})
	} else {
		rounds, ammoType := slot.WeaponState()
		values.RenderUnifiedSliderInt(readOnly, "Rounds", values.UnifierFor(int(rounds)),
			func(u values.Unifier) int {
				return u.Unified().(int)
			}, func(value int) string {
				return fmt.Sprintf("%d", value)
			}, 0x00, 0xFF,
			func(newValue int) {
				slot.SetWeaponState(byte(newValue), ammoType)
				onChange()
			})
		values.RenderUnifiedCombo(readOnly, "Ammo Type", values.UnifierFor(int(ammoType)),
			func(u values.Unifier) int {
				return u.Unified().(int)
			},
			func(value int) string {
				switch value {
				case 0:
					return "Standard"
				case 1:
					return "Special"
				default:
					return fmt.Sprintf("Unknown %d", ammoType)
				}
			}, 2,
			func(newValue int) {
				slot.SetWeaponState(rounds, byte(newValue))
				onChange()
			})
	}
}

func (view *View) createVariableControls(readOnly bool, gameState *archive.GameState, onChange func()) {
	currentContextIdentifier := view.gameStateService.VariableBaseContext()
	if imgui.BeginCombo("Base Context", variableContextList[currentContextIdentifier].title) {
		contextIdentifiers := edit.VariableContextIdentifiers()
		for _, identifier := range contextIdentifiers {
			contextInfo := variableContextList[identifier]
			if imgui.SelectableV(contextInfo.title, identifier == currentContextIdentifier, 0, imgui.Vec2{}) {
				view.gameStateService.SetVariableBaseContext(identifier)
			}
			if imgui.IsItemHovered() {
				imgui.SetTooltip(contextInfo.info)
			}
		}
		imgui.EndCombo()
	}
	if imgui.IsItemHovered() {
		imgui.SetTooltip("The base context defines to what variables default to\n" +
			"in case they are not overridden by the project.")
	}

	if imgui.Button("Remove All Overrides") {
		_ = view.gameStateService.DefaultAllVariables()
	}
	if imgui.IsItemHovered() {
		imgui.SetTooltip("Removes all overrides and their parameters.")
	}
	if !readOnly {
		imgui.SameLine()
		if imgui.Button("Reset All") {
			for i := 0; i < archive.BooleanVarCount; i++ {
				gameState.SetBooleanVar(i, view.gameStateService.BooleanVariable(i).ResetValueBool())
			}
			for i := 0; i < archive.IntegerVarCount; i++ {
				gameState.SetIntegerVar(i, view.gameStateService.IntegerVariable(i).ResetValueInt())
			}
			onChange()
		}
		if imgui.IsItemHovered() {
			imgui.SetTooltip("Resets all variables to their configured initial value.\n" +
				"To reset a single variable, hover over the specific variable and use the context menu.")
		}
	}

	isSavegame := gameState.IsSavegame()
	if imgui.TreeNodeV("Boolean Variables", imgui.TreeNodeFlagsFramed) {
		intConverter := func(u values.Unifier) int {
			bValue := u.Unified().(bool)
			if bValue {
				return 1
			}
			return 0
		}
		defaultNames := map[int]string{0: "False", 1: "True"}

		for i := 0; i < archive.BooleanVarCount; i++ {
			varIndex := i
			info := view.gameStateService.BooleanVariable(varIndex)

			varReadOnly := readOnly || (info.Hardcoded && !isSavegame)
			varLabel := fmt.Sprintf("Var%03d", varIndex)
			varName := fmt.Sprintf("Var%03d: %s", varIndex, info.Name)
			varUnifier := values.UnifierFor(gameState.BooleanVar(varIndex))

			values.RenderUnifiedCombo(varReadOnly, varName, varUnifier, intConverter,
				func(value int) string {
					name, found := info.ValueNames[int16(value)]
					if found {
						return name
					}
					return defaultNames[value]
				},
				2,
				func(newValue int) {
					gameState.SetBooleanVar(varIndex, newValue != 0)
					onChange()
				})

			if imgui.BeginPopupContextItemV(varLabel+"-Popup", 1) {
				if imgui.Selectable("Edit info...") {
					view.modalStateMachine.SetState(&editBooleanVariableDialog{
						machine: view.modalStateMachine,
						view:    view,
						service: view.gameStateService,
						index:   varIndex,
					})
				}
				if !varReadOnly {
					imgui.Separator()
					if imgui.Selectable("Reset") {
						gameState.SetBooleanVar(varIndex, info.ResetValueBool())
						onChange()
					}
				}
				imgui.EndPopup()
			} else if imgui.IsItemHovered() && (len(info.Description) > 0) {
				imgui.SetTooltip(info.Description)
			}
		}
		imgui.TreePop()
	}
	if imgui.TreeNodeV("Integer Variables", imgui.TreeNodeFlagsFramed) {
		intConverter := func(u values.Unifier) int {
			return int(u.Unified().(int16))
		}
		for i := 0; i < archive.IntegerVarCount; i++ {
			varIndex := i
			info := view.gameStateService.IntegerVariable(varIndex)

			varReadOnly := readOnly || (info.Hardcoded && !isSavegame)
			varLabel := fmt.Sprintf("Var%02d", varIndex)
			varName := fmt.Sprintf("Var%02d: %s", varIndex, info.Name)
			varUnifier := values.UnifierFor(gameState.IntegerVar(varIndex))
			changeHandler := func(newValue int) {
				gameState.SetIntegerVar(varIndex, int16(newValue))
				onChange()
			}

			switch {
			case archive.IsRandomIntegerVariable(varIndex):
				values.RenderUnifiedSliderInt(readOnly, varName, varUnifier,
					func(u values.Unifier) int { return int(numbers.FromBinaryCodedDecimal(uint16(intConverter(u)))) },
					func(value int) string { return "%03d" },
					0, 999,
					func(newValue int) {
						changeHandler(int(numbers.ToBinaryCodedDecimal(uint16(newValue))))
					})
			case len(info.ValueNames) > 0:
				var max int16
				for key := range info.ValueNames {
					if key > max {
						max = key
					}
				}
				values.RenderUnifiedCombo(varReadOnly, varName, varUnifier, intConverter,
					func(value int) string {
						name, found := info.ValueNames[int16(value)]
						if found {
							return name
						}
						return fmt.Sprintf("Unknown %d", value)
					},
					int(max+1),
					changeHandler)
			default:
				min := math.MinInt16
				max := math.MaxInt16
				if info.Limits != nil {
					min = int(info.Limits.Minimum)
					max = int(info.Limits.Maximum)
				}

				values.RenderUnifiedSliderInt(varReadOnly, varName, varUnifier, intConverter,
					func(value int) string {
						return "%d"
					},
					min, max,
					changeHandler)
			}
			if imgui.BeginPopupContextItemV(varLabel+"-Popup", 1) {
				if imgui.Selectable("Edit info...") {
					view.modalStateMachine.SetState(&editIntegerVariableDialog{
						machine: view.modalStateMachine,
						view:    view,
						service: view.gameStateService,
						index:   varIndex,
					})
				}
				if !varReadOnly {
					imgui.Separator()
					if imgui.Selectable("Reset") {
						gameState.SetIntegerVar(varIndex, info.ResetValueInt())
						onChange()
					}
				}
				imgui.EndPopup()
			} else if imgui.IsItemHovered() && (len(info.Description) > 0) {
				imgui.SetTooltip(info.Description)
			}
		}
		imgui.TreePop()
	}
}

func (view *View) createMessageControls(readOnly bool, gameState *archive.GameState, onChange func()) {
	view.createMessageControlsFor(readOnly, "EMail",
		&view.model.emailIndex, archive.EMailCount, gameState.EMailState, onChange)
	view.createMessageControlsFor(readOnly, "Log",
		&view.model.logIndex, archive.LogCount, gameState.LogState, onChange)
	view.createMessageControlsFor(readOnly, "Fragment",
		&view.model.fragmentIndex, archive.FragmentCount, gameState.FragmentState, onChange)
}

func (view *View) createMessageControlsFor(readOnly bool, name string, index *int, limit int,
	retriever func(int) archive.MessageState, onChange func()) {
	imgui.PushID(name)
	gui.StepSliderInt(name+" Index", index, 0, limit-1)
	messageState := retriever(*index)
	values.RenderUnifiedCheckboxCombo(readOnly, "Received", values.UnifierFor(messageState.Received()),
		func(newState bool) {
			messageState.SetReceived(newState)
			onChange()
		})
	values.RenderUnifiedCheckboxCombo(readOnly, "Read", values.UnifierFor(messageState.Read()),
		func(newState bool) {
			messageState.SetRead(newState)
			onChange()
		})
	imgui.PopID()
}

func (view *View) tripleName(triple object.Triple) string {
	suffix := hintUnknown
	linearIndex := view.mod.ObjectProperties().TripleIndex(triple)
	if linearIndex >= 0 {
		key := resource.KeyOf(ids.ObjectLongNames, resource.LangDefault, linearIndex)
		objName, err := view.textCache.Text(key)
		if err == nil {
			suffix = singleLineString(objName)
		}
	}
	return triple.String() + ": " + suffix
}

func (view *View) indexedName(class object.Class, index int) string {
	startIndex := view.mod.ObjectProperties().TripleIndex(object.TripleFrom(int(class), 0, 0))
	if startIndex >= 0 {
		key := resource.KeyOf(ids.ObjectLongNames, resource.LangDefault, startIndex+index)
		objName, err := view.textCache.Text(key)
		if err == nil {
			return singleLineString(objName)
		}
	}
	return hintUnknown
}

func singleLineString(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "\r", ""), "\n", "")
}
