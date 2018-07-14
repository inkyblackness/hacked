package actions

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/content/interpreters"
)

func forType(typeID int) func(*interpreters.Instance) bool {
	return func(inst *interpreters.Instance) bool {
		return inst.Get("Type") == uint32(typeID)
	}
}

var transportHackerDetails = interpreters.New().
	With("TargetX", 0, 4).As(interpreters.RangedValue(1, 63)).
	With("TargetY", 4, 4).As(interpreters.RangedValue(1, 63)).
	With("TargetZ", 8, 1).As(interpreters.RangedValue(0, 255)).
	With("PreserveHeight", 9, 1).As(interpreters.EnumValue(map[uint32]string{0: "No", 0x40: "Yes"})).
	With("CrossLevelTransportDestination", 12, 1).As(interpreters.RangedValue(0, 15)).
	With("CrossLevelTransportFlag", 13, 1).As(interpreters.EnumValue(map[uint32]string{
	0x00: "Cross-Level", 0x10: "Same-Level (0x10)", 0x20: "Same-Level (0x20)", 0x22: "Same-Level (0x22)"}))

var changeHealthDetails = interpreters.New().
	With("HealthDelta", 4, 1).
	With("HealthChangeFlag", 6, 2).As(interpreters.EnumValue(map[uint32]string{0: "Remove Delta", 1: "Add Delta"})).
	With("PowerDelta", 8, 1).
	With("PowerChangeFlag", 10, 2).As(interpreters.EnumValue(map[uint32]string{0: "Remove Delta", 1: "Add Delta"}))

var cloneMoveObjectDetails = interpreters.New().
	With("ObjectIndex", 0, 2).As(interpreters.ObjectID()).
	With("MoveFlag", 2, 2).As(interpreters.EnumValue(map[uint32]string{
	0x0000: "Clone Object",
	0x0001: "Move Object (0x0001)",
	0x0002: "Move Object (0x0002)",
	0x0FFF: "Move Object (0x0FFF)",
	0xAAAA: "Move Object (0xAAAA)",
	0xFFFF: "Move Object (0xFFFF)"})).
	With("TargetX", 4, 4).As(interpreters.RangedValue(1, 63)).
	With("TargetY", 8, 4).As(interpreters.RangedValue(1, 63)).
	With("TargetHeight", 12, 1).As(interpreters.SpecialValue("ObjectHeight")).
	With("KeepSourceHeight", 13, 1).As(interpreters.EnumValue(map[uint32]string{0x00: "Set height", 0x40: "Keep height"}))

var setGameVariableDetails = interpreters.New().
	With("VariableKey", 0, 4).As(interpreters.SpecialValue("VariableKey")).
	With("Value", 4, 2).
	With("Operation", 6, 2).As(interpreters.EnumValue(map[uint32]string{0: "Set", 1: "Add", 2: "Subtract", 3: "Multiply", 4: "Divide", 5: "Modulo"})).
	With("Message1", 8, 4).As(interpreters.RangedValue(0, 511)).
	With("Message2", 12, 4).As(interpreters.RangedValue(0, 511))

var showCutsceneDetails = interpreters.New().
	With("CutsceneIndex", 0, 4).As(interpreters.EnumValue(map[uint32]string{0: "Death", 1: "Intro", 2: "Ending"})).
	With("EndGameFlag", 4, 4).As(interpreters.EnumValue(map[uint32]string{0: "No (not working)", 1: "Yes"}))

func pointOneSecond(value int64) string {
	return fmt.Sprintf("%.1f sec", float64(value)*0.1)
}

var triggerOtherObjectsDetails = interpreters.New().
	With("Object1Index", 0, 2).As(interpreters.ObjectID()).
	With("Object1Delay", 2, 2).As(interpreters.FormattedRangedValue(0, 6000, pointOneSecond)).
	With("Object2Index", 4, 2).As(interpreters.ObjectID()).
	With("Object2Delay", 6, 2).As(interpreters.FormattedRangedValue(0, 6000, pointOneSecond)).
	With("Object3Index", 8, 2).As(interpreters.ObjectID()).
	With("Object3Delay", 10, 2).As(interpreters.FormattedRangedValue(0, 6000, pointOneSecond)).
	With("Object4Index", 12, 2).As(interpreters.ObjectID()).
	With("Object4Delay", 14, 2).As(interpreters.FormattedRangedValue(0, 6000, pointOneSecond))

var changeLightingDetails = interpreters.New().
	Refining("ObjectExtent", 0, 2, interpreters.New().With("Index", 0, 2).As(interpreters.ObjectID()), func(inst *interpreters.Instance) bool {
		var lightType = inst.Get("LightType")

		return (lightType == 0x00) || (lightType == 0x01)
	}).
	Refining("RadiusExtent", 0, 2, interpreters.New().With("Tiles", 0, 2).As(interpreters.RangedValue(0, 31)), func(inst *interpreters.Instance) bool {
		return inst.Get("LightType") == 0x03
	}).
	With("ReferenceObjectIndex", 2, 2).As(interpreters.ObjectID()).
	With("TransitionType", 4, 2).As(interpreters.EnumValue(map[uint32]string{0x0000: "immediate", 0x0001: "fade", 0x0100: "flicker"})).
	With("LightModification", 7, 1).As(interpreters.EnumValue(map[uint32]string{0x00: "light on", 0x10: "light off"})).
	With("LightType", 8, 1).As(interpreters.EnumValue(map[uint32]string{0x00: "rectangular", 0x03: "circular gradient"})).
	With("LightSurface", 10, 2).As(interpreters.EnumValue(map[uint32]string{0: "floor", 1: "ceiling", 2: "floor and ceiling"})).
	Refining("Rectangular", 12, 2, interpreters.New().
		With("Off light value", 0, 1).As(interpreters.RangedValue(0, 15)).
		With("On light value", 1, 1).As(interpreters.RangedValue(0, 15)), func(inst *interpreters.Instance) bool {
		return inst.Get("LightType") == 0x00
	}).
	Refining("Gradient", 12, 4, interpreters.New().
		With("Off light begin intensity", 0, 1).As(interpreters.RangedValue(0, 127)).
		With("Off light end intensity", 1, 1).As(interpreters.RangedValue(0, 127)).
		With("On light begin intensity", 2, 1).As(interpreters.RangedValue(0, 127)).
		With("On light end intensity", 3, 1).As(interpreters.RangedValue(0, 127)), func(inst *interpreters.Instance) bool {
		var lightType = inst.Get("LightType")

		return (lightType == 0x01) || (lightType == 0x03)
	})

var effectDetails = interpreters.New().
	With("SoundIndex", 0, 2).As(interpreters.RangedValue(0, 512)).
	With("SoundPlayCount", 2, 2).As(interpreters.RangedValue(0, 100)).
	With("VisualEffect", 4, 2).As(interpreters.EnumValue(map[uint32]string{
	0: "none",
	1: "power on",
	2: "quake",
	3: "escape pod",
	4: "red static",
	5: "interference"})).
	With("AdditionalVisualEffect", 8, 2).As(interpreters.EnumValue(map[uint32]string{
	0: "none",
	1: "white flash",
	2: "pink flash",
	3: "gray static (endless, don't use)",
	4: "vertical panning (endless, don't use)"}))

var changeTileHeightsDetails = interpreters.New().
	With("TileX", 0, 4).As(interpreters.RangedValue(1, 63)).
	With("TileY", 4, 4).As(interpreters.RangedValue(1, 63)).
	With("TargetFloorHeight", 8, 2).As(interpreters.SpecialValue("MoveTileHeight")).
	With("TargetCeilingHeight", 10, 2).As(interpreters.SpecialValue("MoveTileHeight")).
	With("Ignored000C", 12, 4).As(interpreters.SpecialValue("Ignored"))

var randomTimerDetails = interpreters.New().
	With("ObjectIndex", 0, 4).As(interpreters.ObjectID()).
	With("TimeInterval", 4, 4).As(interpreters.FormattedRangedValue(0, 6000,
	func(value int64) string {
		return fmt.Sprintf("%.1fs", float64(value)/10.0)
	})).
	With("ActivationValue", 8, 4).As(interpreters.EnumValue(map[uint32]string{0: "Off",
	0xFFFF: "On (0xFFFF)", 0x10000: "On (0x10000)", 0x11111: "On (0x11111)"})).
	With("Variance", 12, 2).As(interpreters.RangedValue(0, 512))

var cycleObjectsDetails = interpreters.New().
	With("ObjectIndex1", 0, 4).As(interpreters.ObjectID()).
	With("ObjectIndex2", 4, 4).As(interpreters.ObjectID()).
	With("ObjectIndex3", 8, 4).As(interpreters.ObjectID()).
	With("NextObject", 12, 4).As(interpreters.RangedValue(0, 2))

var deleteObjectsDetails = interpreters.New().
	With("ObjectIndex1", 0, 2).As(interpreters.ObjectID()).
	With("ObjectIndex2", 4, 2).As(interpreters.ObjectID()).
	With("ObjectIndex3", 8, 2).As(interpreters.ObjectID()).
	With("MessageIndex", 12, 2).As(interpreters.RangedValue(0, 511))

var receiveEmailDetails = interpreters.New().
	With("EmailIndex", 0, 2).As(interpreters.RangedValue(0, 1000)).
	With("DelaySec", 4, 2).As(interpreters.RangedValue(0, 600))

var changeEffectDetails = interpreters.New().
	With("DeltaValue", 0, 2).As(interpreters.RangedValue(0, 1000)).
	With("EffectChangeFlag", 2, 2).As(interpreters.EnumValue(map[uint32]string{0: "Add Delta", 1: "Remove Delta"})).
	With("EffectType", 4, 4).As(interpreters.EnumValue(map[uint32]string{4: "Radiation poisoning", 8: "Bio contamination"}))

var setObjectParameterDetails = interpreters.New().
	With("ObjectIndex", 0, 4).As(interpreters.ObjectID()).
	With("Value1", 4, 4).
	With("Value2", 8, 4).
	With("Value3", 12, 4)

var setScreenPictureDetails = interpreters.New().
	With("ScreenObjectIndex1", 0, 2).As(interpreters.ObjectID()).
	With("ScreenObjectIndex2", 2, 2).As(interpreters.ObjectID()).
	With("SingleSequenceSource", 4, 4).
	With("LoopSequenceSource", 8, 4)

var setCritterStateDetails = interpreters.New().
	With("ReferenceObjectIndex1", 4, 2).As(interpreters.ObjectID()).
	With("ReferenceObjectIndex2", 6, 2).As(interpreters.ObjectID()).
	With("NewState", 8, 1).As(interpreters.EnumValue(map[uint32]string{
	0: "docile",
	1: "cautious",
	2: "hostile",
	3: "cautious (?)",
	4: "attacking",
	5: "sleeping",
	6: "tranquilized",
	7: "confused"}))

var trapMessageDetails = interpreters.New().
	With("BackgroundImageIndex", 0, 4).As(interpreters.RangedValue(-2, 500)).
	With("MessageIndex", 4, 4).As(interpreters.RangedValue(0, 511)).
	With("TextColor", 8, 4).As(interpreters.RangedValue(0, 255)).
	With("MfdSuppressionFlag", 12, 4).As(interpreters.EnumValue(map[uint32]string{0: "Show in MFD", 1: "Show only in HUD"}))

var spawnObjectsDetails = interpreters.New().
	With("ObjectType", 0, 4).As(interpreters.SpecialValue("ObjectType")).
	With("ReferenceObject1Index", 4, 2).As(interpreters.ObjectID()).
	With("ReferenceObject2Index", 6, 2).As(interpreters.ObjectID()).
	With("NumberOfObjects", 8, 4).As(interpreters.RangedValue(0, 100)).
	With("Unknown000C", 12, 1).As(interpreters.SpecialValue("Unknown"))

var changeObjectTypeDetails = interpreters.New().
	With("ObjectIndex", 0, 4).As(interpreters.ObjectID()).
	With("NewType", 4, 2).As(interpreters.RangedValue(0, 16)).
	With("ResetMask", 6, 2).As(interpreters.RangedValue(0, 15))

// Change state block

var toggleRepulsorChange = interpreters.New().
	With("ObjectIndex", 0, 4).As(interpreters.ObjectID()).
	With("OffTextureIndex", 4, 1).As(interpreters.SpecialValue("LevelTexture")).
	With("OnTextureIndex", 5, 1).As(interpreters.SpecialValue("LevelTexture")).
	With("ToggleType", 8, 1).As(interpreters.EnumValue(map[uint32]string{0: "Toggle On/Off", 1: "Toggle On, Stay On", 2: "Toggle Off, Stay Off"}))

var showGameCodeDigitChange = interpreters.New().
	With("ScreenObjectIndex", 0, 4).As(interpreters.ObjectID()).
	With("DigitNumber", 4, 4).As(interpreters.RangedValue(1, 6))

var setParameterFromVariableChange = interpreters.New().
	With("ObjectIndex", 0, 4).As(interpreters.ObjectID()).
	With("ParameterNumber", 4, 4).As(interpreters.RangedValue(0, 16)).
	With("VariableIndex", 8, 4).As(interpreters.SpecialValue("VariableKey"))

var setButtonStateChange = interpreters.New().
	With("ObjectIndex", 0, 4).As(interpreters.ObjectID()).
	With("NewState", 4, 4).As(interpreters.EnumValue(map[uint32]string{0: "Off", 1: "On"}))

var doorControlChange = interpreters.New().
	With("ObjectIndex", 0, 4).As(interpreters.ObjectID()).
	With("ControlValue", 4, 4).As(interpreters.EnumValue(map[uint32]string{1: "open door", 2: "close door", 3: "toggle door", 4: "suppress auto-close"}))

var rotateObjectChange = interpreters.New().
	With("ObjectIndex", 0, 4).As(interpreters.ObjectID()).
	With("Amount", 4, 1).As(interpreters.RangedValue(0, 255)).
	With("RotationType", 5, 1).As(interpreters.EnumValue(map[uint32]string{0: "Endless", 1: "Back and forth"})).
	With("Direction", 6, 1).As(interpreters.EnumValue(map[uint32]string{0: "Forward", 1: "Backward"})).
	With("Axis", 7, 1).As(interpreters.EnumValue(map[uint32]string{0: "Z (Yaw)", 1: "X (Pitch)", 2: "Y (Roll)"})).
	With("ForwardLimit", 8, 1).As(interpreters.RangedValue(0, 255)).
	With("BackwardLimit", 9, 1).As(interpreters.RangedValue(0, 255))

var removeObjectsChange = interpreters.New().
	With("ObjectType", 0, 4).As(interpreters.SpecialValue("ObjectType")).
	With("Amount", 4, 1).As(interpreters.RangedValue(0, 255))

var setConditionChange = interpreters.New().
	With("ObjectIndex", 0, 4).As(interpreters.ObjectID()).
	With("Condition", 4, 4)

var makeItemRadioactiveChange = interpreters.New().
	With("ObjectIndex", 0, 4).As(interpreters.ObjectID()).
	With("WatchedObjectIndex", 4, 2).As(interpreters.ObjectID()).
	With("WatchedObjectTriggerState", 6, 2)

var orientedTriggerObjectChange = interpreters.New().
	With("HorizontalDirection", 0, 2).As(interpreters.RangedValue(0, 0xFFFF)).
	With("ObjectIndex", 4, 2).As(interpreters.ObjectID())

var closeDataMfdChange = interpreters.New().
	With("ObjectIndex", 0, 4).As(interpreters.ObjectID())

var changeObjectTypeGlobalChange = interpreters.New().
	With("ObjectType", 0, 4).As(interpreters.SpecialValue("ObjectType")).
	With("NewType", 4, 1).As(interpreters.RangedValue(0, 16))

var changeStateDetails = interpreters.New().
	With("Type", 0, 4).As(interpreters.EnumValue(map[uint32]string{
	1:  "Toggle Repulsor",
	2:  "Show Game Code Digit",
	3:  "Set Parameter from Variable",
	4:  "Set Button State",
	5:  "Door Control",
	6:  "Return to Menu",
	7:  "Rotate Objects",
	8:  "Remove Objects",
	9:  "SHODAN Pixelation",
	10: "Set Condition",
	11: "Show System Analyzer",
	12: "Make Item Radioactive",
	13: "Oriented Trigger Object",
	14: "Close Data MFD",
	15: "Earth Destruction by Laser",
	16: "Change Objects Type (Level)"})).
	Refining("ToggleRepulsor", 4, 12, toggleRepulsorChange, forType(1)).
	Refining("ShowGameCodeDigit", 4, 12, showGameCodeDigitChange, forType(2)).
	Refining("SetParameterFromVariable", 4, 12, setParameterFromVariableChange, forType(3)).
	Refining("SetButtonState", 4, 12, setButtonStateChange, forType(4)).
	Refining("DoorControl", 4, 12, doorControlChange, forType(5)).
	Refining("ReturnToMenu", 4, 12, interpreters.New(), forType(6)).
	Refining("RotateObject", 4, 12, rotateObjectChange, forType(7)).
	Refining("RemoveObjects", 4, 12, removeObjectsChange, forType(8)).
	Refining("ShodanPixelation", 4, 12, interpreters.New(), forType(9)).
	Refining("SetCondition", 4, 12, setConditionChange, forType(10)).
	Refining("ShowSystemAnalyzer", 4, 12, interpreters.New(), forType(11)).
	Refining("MakeItemRadioactive", 4, 12, makeItemRadioactiveChange, forType(12)).
	Refining("OrientedTriggerObject", 4, 12, orientedTriggerObjectChange, forType(13)).
	Refining("CloseDataMfd", 4, 12, closeDataMfdChange, forType(14)).
	Refining("EarthDestructionByLaser", 4, 12, interpreters.New(), forType(15)).
	Refining("ChangeObjectsType", 4, 12, changeObjectTypeGlobalChange, forType(16))

var unconditionalAction = interpreters.New().
	With("Type", 0, 1).As(interpreters.EnumValue(map[uint32]string{
	0:  "Nothing",
	1:  "Transport Hacker",
	2:  "Change Health",
	3:  "Clone/Move Object",
	4:  "Set Game Variable",
	5:  "Show Cutscene",
	6:  "Trigger Other Objects",
	7:  "Change Lighting",
	8:  "Effect",
	9:  "Change Tile Heights",
	10: "Unknown (10)",
	11: "Random Timer",
	12: "Cycle Objects",
	13: "Delete Objects",
	14: "Unknown (14)",
	15: "Receive E-Mail",
	16: "Change Effect",
	17: "Set Object Parameter",
	18: "Set Screen Picture",
	19: "Change State",
	20: "Unknown (20)",
	21: "Set Critter State",
	22: "Trap Message",
	23: "Spawn Objects",
	24: "Change Object Type"})).
	With("UsageQuota", 1, 1).
	Refining("TransportHacker", 6, 16, transportHackerDetails, forType(1)).
	Refining("ChangeHealth", 6, 16, changeHealthDetails, forType(2)).
	Refining("CloneMoveObject", 6, 16, cloneMoveObjectDetails, forType(3)).
	Refining("SetGameVariable", 6, 16, setGameVariableDetails, forType(4)).
	Refining("ShowCutscene", 6, 16, showCutsceneDetails, forType(5)).
	Refining("TriggerOtherObjects", 6, 16, triggerOtherObjectsDetails, forType(6)).
	Refining("ChangeLighting", 6, 16, changeLightingDetails, forType(7)).
	Refining("Effect", 6, 16, effectDetails, forType(8)).
	Refining("ChangeTileHeights", 6, 16, changeTileHeightsDetails, forType(9)).
	// 10 unknown
	Refining("RandomTimer", 6, 16, randomTimerDetails, forType(11)).
	Refining("CycleObjects", 6, 16, cycleObjectsDetails, forType(12)).
	Refining("DeleteObjects", 6, 16, deleteObjectsDetails, forType(13)).
	// 14 unknown
	Refining("ReceiveEmail", 6, 16, receiveEmailDetails, forType(15)).
	Refining("ChangeEffect", 6, 16, changeEffectDetails, forType(16)).
	Refining("SetObjectParameter", 6, 16, setObjectParameterDetails, forType(17)).
	Refining("SetScreenPicture", 6, 16, setScreenPictureDetails, forType(18)).
	Refining("ChangeState", 6, 16, changeStateDetails, forType(19)).
	// 20 unknown
	Refining("SetCritterState", 6, 16, setCritterStateDetails, forType(21)).
	Refining("TrapMessage", 6, 16, trapMessageDetails, forType(22)).
	Refining("SpawnObjects", 6, 16, spawnObjectsDetails, forType(23)).
	Refining("ChangeObjectType", 6, 16, changeObjectTypeDetails, forType(24))

// Unconditional returns the description of actions without a condition.
func Unconditional() *interpreters.Description {

	return unconditionalAction
}
