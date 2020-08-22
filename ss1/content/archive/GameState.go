package archive

import (
	"fmt"
	"math"

	"github.com/inkyblackness/hacked/ss1/content/archive/level"
	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/text"
)

// GameStateSize specifies the byte count of a the GameState in an archive.
// Note that in the original archive.dat file, the resource only has 0x054D bytes.
// This does not matter much, as in the original engine the resource was ignored.
const GameStateSize = 0x0595

const (
	stateHackerNameSize = 20
	// InventoryWeaponSlots is the number of weapons available to the character.
	InventoryWeaponSlots   = 7
	weaponSlotsStartOffset = 0x048B
	weaponSlotSize         = 5

	// GrenadeTypeCount is the number of grenades (explosives) available to the character.
	GrenadeTypeCount        = 7
	grenadeCountStartOffset = 0x0350
	grenadeTimerStartOffset = 0x04FF
	grenadeTimerDefault     = 70
	// GrenadeTimerMaximum is the maximum time that can be set for a grenade in game.
	GrenadeTimerMaximum = 571

	// PatchTypeCount is the number of patches available to the character.
	PatchTypeCount        = 7
	patchCountStartOffset = 0x0349

	// AmmoTypeCount is the number of types of ammo available to the character.
	AmmoTypeCount                   = 15
	ammoFullClipCountStartOffset    = 0x032B
	ammoExtraRoundsCountStartOffset = 0x033A

	// GeneralInventorySlotCount is the number of general inventory slots available.
	GeneralInventorySlotCount       = 14
	generalInventorySlotStartOffset = 0x007A

	// HardwareTypeCount is the number of available hardware.
	HardwareTypeCount          = 15
	hardwareVersionStartOffset = 0x0309
	hardwareStatusStartOffset  = 0x04AE

	engineTicksPerSecond = 280
	secondsPerMinute     = 60
	engineTicksPerMinute = secondsPerMinute * engineTicksPerSecond
	cyberspaceMinTime    = 90 * engineTicksPerSecond
	cyberspaceMaxTime    = 30 * secondsPerMinute * engineTicksPerSecond

	messageStatusReceived = 0x80
	messageStatusRead     = 0x40

	// BooleanVarCount is the number of available boolean variables.
	BooleanVarCount       = 512
	booleanVarStartOffset = 0x00B6
	// IntegerVarCount is the number of available integer variables.
	IntegerVarCount       = 64
	integerVarStartOffset = 0x00F6
)

// GameState contains all the necessary information for the current progress of the game.
// It is a combination of world information, character information, and configuration values.
type GameState struct {
	*interpreters.Instance
}

// InventoryWeaponSlot describes one entry for weapons in the inventory.
type InventoryWeaponSlot struct {
	Index int
	State *GameState
}

func freeWeaponSlot() []byte {
	return []byte{0xFF, 0x00, 0x00, 0x00, 0x00}
}

func (slot InventoryWeaponSlot) rawData() []byte {
	if (slot.Index < 0) || (slot.Index >= InventoryWeaponSlots) {
		return freeWeaponSlot()
	}
	startIndex := weaponSlotsStartOffset + (weaponSlotSize * slot.Index)
	return slot.State.Raw()[startIndex : startIndex+weaponSlotSize]
}

// IsInUse returns true if the slot is currently holding a weapon.
func (slot InventoryWeaponSlot) IsInUse() bool {
	return slot.rawData()[0] != 0xFF
}

// SetFree clears the weapon slot.
func (slot InventoryWeaponSlot) SetFree() {
	rawData := slot.rawData()
	copy(rawData, freeWeaponSlot())
}

// SetInUse inserts the given weapon in this slot. Other fields are reset to zero.
func (slot InventoryWeaponSlot) SetInUse(subclass object.Subclass, objType object.Type) {
	rawData := slot.rawData()
	copy(rawData, freeWeaponSlot())
	rawData[0] = byte(subclass)
	rawData[1] = byte(objType)
}

// Triple returns the type information of the weapon in this slot.
func (slot InventoryWeaponSlot) Triple() object.Triple {
	rawData := slot.rawData()
	return object.TripleFrom(int(object.ClassGun), int(rawData[0]), int(rawData[1]))
}

// WeaponState returns the two state values. Their interpretation is subclass specific.
func (slot InventoryWeaponSlot) WeaponState() (byte, byte) {
	rawData := slot.rawData()
	return rawData[2], rawData[3]
}

// SetWeaponState updates the two state values.
func (slot InventoryWeaponSlot) SetWeaponState(a, b byte) {
	rawData := slot.rawData()
	rawData[2] = a
	rawData[3] = b
}

// InventoryGrenade describes properties of grenades by type
type InventoryGrenade struct {
	Index int
	State *GameState
}

func (grenade InventoryGrenade) isValid() bool {
	return (grenade.Index >= 0) && (grenade.Index < GrenadeTypeCount) && (grenade.State != nil)
}

// Count returns how many explosives the character is carrying.
func (grenade InventoryGrenade) Count() int {
	if !grenade.isValid() {
		return 0
	}
	return int(grenade.State.Raw()[grenadeCountStartOffset+grenade.Index])
}

// SetCount updates the amount of explosives the character is carrying. Invalid values are clamped.
func (grenade InventoryGrenade) SetCount(value int) {
	if !grenade.isValid() {
		return
	}
	count := byte(value)
	if value < 0 {
		count = 0
	} else if value > math.MaxUint8 {
		count = math.MaxUint8
	}

	grenade.State.Raw()[grenadeCountStartOffset+grenade.Index] = count
}

// GrenadeTimerSetting is the value for explosive timers.
type GrenadeTimerSetting uint16

// InSeconds returns the amount in unit of seconds.
func (setting GrenadeTimerSetting) InSeconds() float32 {
	return float32(setting) / 10.0
}

func (grenade InventoryGrenade) rawTimer() []byte {
	if !grenade.isValid() {
		return []byte{0x00, 0x00}
	}
	startIndex := grenadeTimerStartOffset + (2 * grenade.Index)
	return grenade.State.Raw()[startIndex : startIndex+2]
}

// TimerSetting returns the generic timer used for this type of grenade.
func (grenade InventoryGrenade) TimerSetting() GrenadeTimerSetting {
	raw := grenade.rawTimer()
	return GrenadeTimerSetting((uint16(raw[1]) << 8) | (uint16(raw[0]) << 0))
}

// SetTimerSetting updates the generic timer used for this type of grenade.
func (grenade InventoryGrenade) SetTimerSetting(value GrenadeTimerSetting) {
	raw := grenade.rawTimer()
	raw[0] = byte(value >> 0)
	raw[1] = byte(value >> 8)
}

// PatchState describes properties of patches by type
type PatchState struct {
	Index int
	State *GameState
}

func (patch PatchState) isValid() bool {
	return (patch.Index >= 0) && (patch.Index < PatchTypeCount) && (patch.State != nil)
}

// Count returns how many patches the character is carrying.
func (patch PatchState) Count() int {
	if !patch.isValid() {
		return 0
	}
	return int(patch.State.Raw()[patchCountStartOffset+patch.Index])
}

// SetCount updates the amount of patches the character is carrying. Invalid values are clamped.
func (patch PatchState) SetCount(value int) {
	if !patch.isValid() {
		return
	}
	count := byte(value)
	if value < 0 {
		count = 0
	} else if value > math.MaxUint8 {
		count = math.MaxUint8
	}

	patch.State.Raw()[patchCountStartOffset+patch.Index] = count
}

// InventoryAmmo describes properties of ammo by type
type InventoryAmmo struct {
	Index int
	State *GameState
}

func (ammo InventoryAmmo) isValid() bool {
	return (ammo.Index >= 0) && (ammo.Index < AmmoTypeCount) && (ammo.State != nil)
}

// FullClipCount returns how many full clips the character is carrying.
func (ammo InventoryAmmo) FullClipCount() int {
	if !ammo.isValid() {
		return 0
	}
	return int(ammo.State.Raw()[ammoFullClipCountStartOffset+ammo.Index])
}

// SetFullClipCount updates the amount of full clips the character is carrying. Invalid values are clamped.
func (ammo InventoryAmmo) SetFullClipCount(value int) {
	if !ammo.isValid() {
		return
	}
	count := byte(value)
	if value < 0 {
		count = 0
	} else if value > math.MaxUint8 {
		count = math.MaxUint8
	}

	ammo.State.Raw()[ammoFullClipCountStartOffset+ammo.Index] = count
}

// ExtraRoundsCount returns how many extra rounds the character is carrying.
func (ammo InventoryAmmo) ExtraRoundsCount() int {
	if !ammo.isValid() {
		return 0
	}
	return int(ammo.State.Raw()[ammoExtraRoundsCountStartOffset+ammo.Index])
}

// SetExtraRoundsCount updates the amount of extra rounds the character is carrying. Invalid values are clamped.
func (ammo InventoryAmmo) SetExtraRoundsCount(value int) {
	if !ammo.isValid() {
		return
	}
	count := byte(value)
	if value < 0 {
		count = 0
	} else if value > math.MaxUint8 {
		count = math.MaxUint8
	}

	ammo.State.Raw()[ammoExtraRoundsCountStartOffset+ammo.Index] = count
}

// GeneralInventorySlot describes properties of a general inventory slot.
type GeneralInventorySlot struct {
	Index int
	State *GameState
}

func (slot GeneralInventorySlot) isValid() bool {
	return (slot.Index >= 0) && (slot.Index < GeneralInventorySlotCount) && (slot.State != nil)
}

func (slot GeneralInventorySlot) rawID() []byte {
	if !slot.isValid() {
		return []byte{0x00, 0x00}
	}
	startIndex := generalInventorySlotStartOffset + (2 * slot.Index)
	return slot.State.Raw()[startIndex : startIndex+2]
}

// ObjectID returns the ID of the object in this inventory slot.
func (slot GeneralInventorySlot) ObjectID() level.ObjectID {
	raw := slot.rawID()
	return level.ObjectID((uint16(raw[1]) << 8) | (uint16(raw[0]) << 0))
}

// SetObjectID updates the ID in this inventory slot.
func (slot GeneralInventorySlot) SetObjectId(id level.ObjectID) {
	raw := slot.rawID()
	raw[0] = byte(id >> 0)
	raw[1] = byte(id >> 8)
}

// HardwareState describes properties of hardware by type.
type HardwareState struct {
	Index int
	State *GameState
}

func (hardware HardwareState) isValid() bool {
	return (hardware.Index >= 0) && (hardware.Index < HardwareTypeCount) && (hardware.State != nil)
}

// Version returns which version of the hardware is installed. Zero means not installed.
func (hardware HardwareState) Version() int {
	if !hardware.isValid() {
		return 0
	}
	return int(hardware.State.Raw()[hardwareVersionStartOffset+hardware.Index])
}

// SetVersion updates the version the character has installed.
func (hardware HardwareState) SetVersion(value int) {
	if !hardware.isValid() {
		return
	}
	count := byte(value)
	if value < 0 {
		count = 0
	} else if value > math.MaxUint8 {
		count = math.MaxUint8
	}

	hardware.State.Raw()[hardwareVersionStartOffset+hardware.Index] = count
}

// IsActive returns whether the hardware is currently active.
func (hardware HardwareState) IsActive() bool {
	if !hardware.isValid() {
		return false
	}
	return hardware.State.Raw()[hardwareStatusStartOffset+hardware.Index] != 0
}

// SetActive updates the active state of a hardware.
func (hardware HardwareState) SetActive(on bool) {
	if !hardware.isValid() {
		return
	}
	var newValue byte
	if on {
		newValue = 1
	}

	hardware.State.Raw()[hardwareStatusStartOffset+hardware.Index] = newValue
}

var gameStateDesc = interpreters.New().
	With("Difficulty: Combat", 0x0015, 1).As(interpreters.RangedValue(0, 3)).
	With("Difficulty: Mission", 0x0016, 1).As(interpreters.RangedValue(0, 3)).
	With("Difficulty: Puzzle", 0x0017, 1).As(interpreters.RangedValue(0, 3)).
	With("Difficulty: Cyber", 0x0018, 1).As(interpreters.RangedValue(0, 3)).
	With("Game time", 0x001C, 4).As(interpreters.SpecialValue("Internal")).
	With("Current Level", 0x0039, 1).As(interpreters.RangedValue(0, 15)).
	With("Health", 0x009C, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int) string {
		return fmt.Sprintf("%.2f%%%%", float64(value*100)/255.0)
	})).
	With("Cyberspace integrity", 0x009D, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int) string {
		return fmt.Sprintf("%.2f%%%%", float64(value*100)/255.0)
	})).
	With("Health regeneration per minute", 0x009E, 2).As(interpreters.RangedValue(0, 1000)).
	With("Power", 0x00AC, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int) string {
		return fmt.Sprintf("%.2f%%%%", float64(value*100)/255.0)
	})).
	With("Power Usage", 0x00AD, 1).As(interpreters.RangedValue(0, 0xFF)).
	With("Out-of-power flag", 0x00AF, 1).As(interpreters.EnumValue(
	map[uint32]string{
		0: "Powered",
		1: "Out-of-power",
	})).
	With("Cyberspace base time", 0x00B2, 4).As(interpreters.FormattedRangedValue(cyberspaceMinTime, cyberspaceMaxTime,
	func(value int) string {
		minutes := value / engineTicksPerMinute
		return fmt.Sprintf("%2dm %2.2fs", minutes, float64(value-(minutes*engineTicksPerMinute))/(engineTicksPerSecond))
	})).
	With("Fatigue regeneration", 0x0181, 2).As(interpreters.RangedValue(0, 400)).
	With("Fatigue regeneration base", 0x0183, 2).As(interpreters.RangedValue(0, 400)).
	With("Fatigue regeneration max", 0x0185, 2).As(interpreters.RangedValue(0, 400)).
	With("Accuracy", 0x0187, 1).As(interpreters.RangedValue(0, 255)).
	With("Shield absorb rate", 0x0188, 1).As(interpreters.RangedValue(0, 255)).
	With("Hacker Position X", 0x053F, 4).As(interpreters.FormattedRangedValue(0x010000, 0x3FFFFF,
	func(value int) string {
		return fmt.Sprintf("%2.03f", float32(value)/0x010000)
	})).
	With("Hacker Position Y", 0x0543, 4).As(interpreters.FormattedRangedValue(0x010000, 0x3FFFFF,
	func(value int) string {
		return fmt.Sprintf("%2.03f", float32(value)/0x010000)
	})).
	With("Hacker Position Z", 0x0547, 4).As(interpreters.FormattedRangedValue(0x010000, 0x1FFFFF,
	func(value int) string {
		return fmt.Sprintf("%2.03f tile(s)", float32(value)/0x010000)
	})).
	With("Hacker Yaw", 0x054B, 4).As(interpreters.RotationValue(0, edmsFullCircle()))

func edmsFullCircle() int64 {
	full := math.Pi * 2 * float64(0x10000)
	return int64(full) - 1
}

// NewGameState() returns a GameState instance for given raw data.
func NewGameState(raw []byte) *GameState {
	return &GameState{Instance: gameStateDesc.For(raw)}
}

// IsSavegame returns true if the state describes one during a running game.
func (state *GameState) IsSavegame() bool {
	// Picking the right property is tricky. Most properties could be pre-initialized
	// by an archive as well.
	// For now, assume the "game time" will not be touched by initial archives.
	// A possible other approach would be to take the "version" field, which
	// is always set when a savegame is created, but is zero in the original. Yet,
	// it is not clear whether or not a (pedantic) engine would respect this field and
	// ignore data when it's zero.
	return state.Get("Game time") > 0
}

// IsDefaulting returns true if the state describes an archive from which the engine
// should initialize the state for a new game.
func (state *GameState) IsDefaulting() bool {
	return state.Get("Hacker Position X") == 0
}

// HackerName returns the name interpreted with given codepage.
func (state *GameState) HackerName(cp text.Codepage) string {
	return cp.Decode(state.Raw()[:stateHackerNameSize])
}

// SetHackerName stores the given name using given codepage, up to the internal available amount of bytes.
func (state *GameState) SetHackerName(name string, cp text.Codepage) {
	raw := state.Raw()
	var zeroName [stateHackerNameSize]byte
	copy(raw[:stateHackerNameSize], zeroName[:])
	copy(raw[:stateHackerNameSize], cp.Encode(name))
	raw[stateHackerNameSize-1] = 0
}

// CurrentLevel returns the number of the level hacker is currently in.
func (state GameState) CurrentLevel() int {
	return int(state.Get("Current Level"))
}

// HackerPosition returns the rough X/Y location on the map.
func (state GameState) HackerMapPosition() (level.Coordinate, level.Coordinate) {
	x := state.Get("Hacker Position X")
	y := state.Get("Hacker Position Y")
	return level.CoordinateAt(byte(x>>16), byte(x>>8)), level.CoordinateAt(byte(y>>16), byte(y>>8))
}

// BooleanVar returns the state of the boolean variable at given index. Unsupported indices return 0.
func (state GameState) BooleanVar(index int) bool {
	if (index < 0) || (index >= BooleanVarCount) {
		return false
	}
	return state.Raw()[booleanVarStartOffset+(index/8)]&(0x01<<(uint(index%8))) != 0
}

// SetBooleanVar sets the state of the boolean variable at given index. Unsupported indices are ignored.
func (state *GameState) SetBooleanVar(index int, set bool) {
	if (index < 0) || (index >= BooleanVarCount) {
		return
	}
	byteIndex := booleanVarStartOffset + (index / 8)
	bitMask := byte(0x01 << uint(index%8))
	temp := state.Raw()[byteIndex]
	temp &= ^bitMask
	if set {
		temp |= bitMask
	}
	state.Raw()[byteIndex] = temp
}

// IntegerVar returns the value of the integer variable at given index. Unsupported indices return 0.
func (state *GameState) IntegerVar(index int) int16 {
	if (index < 0) || (index >= IntegerVarCount) {
		return 0
	}
	startOffset := integerVarStartOffset + (2 * index)
	val := uint16(state.Raw()[startOffset+1])<<8 | uint16(state.Raw()[startOffset+0])
	return int16(val)
}

// SetIntegerVar sets the value of the integer variable at given index. Unsupported indices are ignored.
func (state *GameState) SetIntegerVar(index int, value int16) {
	if (index < 0) || (index >= IntegerVarCount) {
		return
	}
	startOffset := integerVarStartOffset + (2 * index)
	state.Raw()[startOffset+0] = byte((value >> 0) & 0xFF)
	state.Raw()[startOffset+1] = byte((value >> 8) & 0xFF)
}

// InventoryWeaponSlot returns an accessor for the identified weapon slot.
// Index should be within the range of [0..InventoryWeaponSlots[.
func (state *GameState) InventoryWeaponSlot(index int) InventoryWeaponSlot {
	return InventoryWeaponSlot{
		Index: index,
		State: state,
	}
}

// InventoryGrenade returns an accessor for the identified grenade type index.
// Index should be within the range of [0..GrenadeTypeCount[
func (state *GameState) InventoryGrenade(index int) InventoryGrenade {
	return InventoryGrenade{
		Index: index,
		State: state,
	}
}

// PatchState returns an accessor for the identified patch type index.
// Index should be within the range of [0..PatchTypeCount[
func (state *GameState) PatchState(index int) PatchState {
	return PatchState{
		Index: index,
		State: state,
	}
}

// InventoryAmmo returns an accessor for the identified ammo type index.
// Index should be within the range of [0..AmmoTypeCount[
func (state *GameState) InventoryAmmo(index int) InventoryAmmo {
	return InventoryAmmo{
		Index: index,
		State: state,
	}
}

// GeneralInventorySlot returns an accessor for the identified general slot.
// Index should be within the range of [0..GeneralInventorySlotCount[.
func (state *GameState) GeneralInventorySlot(index int) GeneralInventorySlot {
	return GeneralInventorySlot{
		Index: index,
		State: state,
	}
}

// HardwareState returns an accessor for the identified hardware type index.
// Index should be within the range of [0..HardwareTypeCount[
func (state *GameState) HardwareState(index int) HardwareState {
	return HardwareState{
		Index: index,
		State: state,
	}
}

// DefaultGameStateData returns the state block initialized as if the engine started a new default game.
func DefaultGameStateData() []byte {
	data := ZeroGameStateData()
	state := NewGameState(data)

	state.Set("Health", 212)
	state.Set("Cyberspace integrity", 255)
	state.Set("Health regeneration per minute", 0)
	state.Set("Power", 255)
	state.Set("Cyberspace base time", 30*60*engineTicksPerSecond)
	state.Set("Shield absorb rate", 0)
	state.Set("Fatigue regeneration", 0)
	state.Set("Fatigue regeneration base", 100)
	state.Set("Fatigue regeneration max", 400)
	state.Set("Accuracy", 100)

	for i := 0; i < InventoryWeaponSlots; i++ {
		state.InventoryWeaponSlot(i).SetFree()
	}
	for i := 0; i < GrenadeTypeCount; i++ {
		state.InventoryGrenade(i).SetTimerSetting(grenadeTimerDefault)
	}

	setInitialCitadelHackerState(state)
	setInitialCitadelVariables(state)

	return data
}

// ZeroGameStateData returns the state block reset for default engine behaviour.
func ZeroGameStateData() []byte {
	return make([]byte, GameStateSize)
}

func setInitialCitadelHackerState(state *GameState) {
	data := state.Raw()

	state.Set("Current Level", 1)

	// location: in the neurosurgery chamber, looking West
	state.Set("Hacker Position X", (30<<16)+0x8000)
	state.Set("Hacker Position Y", (22<<16)+0x8000)
	state.Set("Hacker Position Z", 0x01BD00)
	state.Set("Hacker Yaw", 0x03243E)

	// set first message
	data[0x0357+26] = messageStatusReceived // Rebecca Lansing's first message
	data[0x0519+9] = 0xFF                   // HUD active email -- set for similarity, has no effect.
}

func setInitialCitadelVariables(state *GameState) {
	// The following set is taken from player.c
	initialBooleanVariables := []int{
		0x001, 0x002, 0x003, 0x010, 0x012, 0x015, 0x016, 0x017, 0x018, 0x019, 0x01A,
		0x020, 0x021, 0x024, 0x025,
		0x04B, 0x04C, 0x04D, 0x04E, 0x04F,
		0x050, 0x051, 0x052, 0x053, 0x054, 0x055, 0x056, 0x057, 0x058, 0x059, 0x05A, 0x05B, 0x05C, 0x05D, 0x05E, 0x05F,
		0x070, 0x071, 0x072, 0x073, 0x074, 0x075, 0x076, 0x077, 0x078, 0x079, 0x07A, 0x07B, 0x07C, 0x07D, 0x07E, 0x07F,
		0x0A0, 0x0A1, 0x0A2, 0x0A3, 0x0A4, 0x0A5, 0x0A6, 0x0A7, 0x0A8, 0x0A9,
		0x0C0, 0x0C1, 0x0C2, 0x0C3, 0x0C4, 0x0C5, 0x0C6, 0x0C7, 0x0C8, 0x0C9, 0x0CA, 0x0CB, 0x0CC, 0x0CD, 0x0CE, 0x0CF,
		0x0E1, 0x0E3, 0x0E5, 0x0E7, 0x0E9, 0x0EB, 0x0ED, 0x0EF,
		0x0F1, 0x0F3, 0x0F5, 0x0F7, 0x0F9, 0x0FB, 0x0FD, 0x0FF,
		0x101, 0x103, 0x105, 0x107, 0x109, 0x10B, 0x10D, 0x10F,
		0x111, 0x113, 0x115, 0x117, 0x119, 0x11B, 0x11D, 0x11F,
		0x121, 0x123, 0x125, 0x127, 0x129, 0x12B,
	}
	initialIntegerVariables := map[int]int16{
		0x03: 2,     // engine state
		0x0C: 3,     // number of available groves
		0x33: 0x100, // joystick sensitivity
	}
	for i := 0; i < BooleanVarCount; i++ {
		state.SetBooleanVar(i, false)
	}
	for _, index := range initialBooleanVariables {
		state.SetBooleanVar(index, true)
	}
	for i := 0; i < IntegerVarCount; i++ {
		state.SetIntegerVar(i, initialIntegerVariables[i])
	}
}
