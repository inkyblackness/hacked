package archive

import "math"

// GameVariableLimits describes limits of a game variable.
type GameVariableLimits struct {
	Minimum int16
	Maximum int16
}

// GameVariableInfo describes a variable in the game state of the archive.
type GameVariableInfo struct {
	// InitValue is nil if the initial value is system dependent.
	InitValue *int16
	// Name is the short identifier for the variable.
	Name string
	// Description is an optional text with further information on the variable.
	Description string
	// Limits may be provided to describe the range of possible values.
	Limits *GameVariableLimits
	// ValueNames may be set for enumerated values.
	ValueNames map[int16]string
	// Hardcoded variables are fixed for a new game.
	Hardcoded bool
}

// GameVariableInfoFor returns a new instance with given name.
func GameVariableInfoFor(name string) GameVariableInfo {
	return GameVariableInfo{Name: name}
}

// At returns an information based on the current one, with the given initial value.
func (info GameVariableInfo) At(value int16) GameVariableInfo {
	info.InitValue = &value
	return info
}

// HardcodedConfig returns an information that is marked as being a hardcoded initialized configuration value.
func (info GameVariableInfo) HardcodedConfig() GameVariableInfo {
	info.Hardcoded = true
	return info
}

// HardcodedAt returns an information that is marked as being a hardcoded variable with given initial value.
func (info GameVariableInfo) HardcodedAt(value int16) GameVariableInfo {
	info.InitValue = &value
	info.Hardcoded = true
	return info
}

// Enumerated returns an information with given value names.
func (info GameVariableInfo) Enumerated(names map[int16]string) GameVariableInfo {
	info.ValueNames = names
	return info
}

// Boolean returns an enumerated information with No/Yes as possible values.
func (info GameVariableInfo) Boolean() GameVariableInfo {
	return info.Enumerated(map[int16]string{
		0: "No",
		1: "Yes",
	})
}

// LimitedBy returns an information with given minimum and maximum values.
func (info GameVariableInfo) LimitedBy(min, max int16) GameVariableInfo {
	info.Limits = &GameVariableLimits{
		Minimum: min,
		Maximum: max,
	}
	return info
}

// GameVariables is a lookup map for information on game variables.
type GameVariables map[int]GameVariableInfo

// Lookup returns the information for given index. If the index is not known, the unknown function is used.
func (vars GameVariables) Lookup(index int, unknown func() GameVariableInfo) GameVariableInfo {
	info, known := vars[index]
	if known {
		return info
	}
	return unknown()
}

var engineIntegerVariables = GameVariables{
	2: GameVariableInfoFor("Destroyed Antenna Count").At(0).LimitedBy(0, 100),

	9: GameVariableInfoFor("Plot counter").At(0),

	13: GameVariableInfoFor("Difficulty: Mission").HardcodedConfig().LimitedBy(0, 3),
	14: GameVariableInfoFor("Difficulty: Cyber").HardcodedConfig().LimitedBy(0, 3),
	15: GameVariableInfoFor("Difficulty: Combat").HardcodedConfig().LimitedBy(0, 3),

	16: GameVariableInfoFor("Security Value: Level 0").HardcodedAt(0).LimitedBy(0, 1000),
	17: GameVariableInfoFor("Security Value: Level 1").HardcodedAt(0).LimitedBy(0, 1000),
	18: GameVariableInfoFor("Security Value: Level 2").HardcodedAt(0).LimitedBy(0, 1000),
	19: GameVariableInfoFor("Security Value: Level 3").HardcodedAt(0).LimitedBy(0, 1000),
	20: GameVariableInfoFor("Security Value: Level 4").HardcodedAt(0).LimitedBy(0, 1000),
	21: GameVariableInfoFor("Security Value: Level 5").HardcodedAt(0).LimitedBy(0, 1000),
	22: GameVariableInfoFor("Security Value: Level 6").HardcodedAt(0).LimitedBy(0, 1000),
	23: GameVariableInfoFor("Security Value: Level 7").HardcodedAt(0).LimitedBy(0, 1000),
	24: GameVariableInfoFor("Security Value: Level 8").HardcodedAt(0).LimitedBy(0, 1000),
	25: GameVariableInfoFor("Security Value: Level 9").HardcodedAt(0).LimitedBy(0, 1000),
	26: GameVariableInfoFor("Security Value: Level 10").HardcodedAt(0).LimitedBy(0, 1000),
	27: GameVariableInfoFor("Security Value: Level 11").HardcodedAt(0).LimitedBy(0, 1000),
	28: GameVariableInfoFor("Security Value: Level 12").HardcodedAt(0).LimitedBy(0, 1000),
	29: GameVariableInfoFor("Security Value: Level 13").HardcodedAt(0).LimitedBy(0, 1000),

	30: GameVariableInfoFor("Difficulty: Puzzle").HardcodedConfig().LimitedBy(0, 3),
	31: GameVariableInfoFor("Random Code 1").At(0),
	32: GameVariableInfoFor("Random Code 2").At(0),

	41: GameVariableInfoFor("Music Volume").HardcodedConfig().LimitedBy(0, 100),
	42: GameVariableInfoFor("Video Gamma").HardcodedConfig().LimitedBy(0, math.MaxInt16),
	43: GameVariableInfoFor("SFX Volume").HardcodedConfig().LimitedBy(0, 100),
	44: GameVariableInfoFor("Mouse Handedness").HardcodedConfig().Enumerated(
		map[int16]string{
			0: "Right Handed",
			1: "Left Handed",
		}),
	47: GameVariableInfoFor("Double-Click Speed").HardcodedConfig().LimitedBy(0, math.MaxInt16),
	48: GameVariableInfoFor("Language").HardcodedConfig().Enumerated(map[int16]string{
		0: "Default",
		1: "French",
		2: "German",
	}),
	49: GameVariableInfoFor("Audiolog Volume").HardcodedConfig().LimitedBy(0, 100),
	50: GameVariableInfoFor("Screen Mode").HardcodedConfig(),
	51: GameVariableInfoFor("Joystick Sensitivity").At(0x100).LimitedBy(0, 0x100),
	52: GameVariableInfoFor("Show Fullscreen Icons").HardcodedConfig().Boolean(),
	53: GameVariableInfoFor("Audio Messages").HardcodedConfig().Enumerated(map[int16]string{
		0: "Text Only",
		1: "Speech Only",
		2: "Text and Speech",
	}),
	54: GameVariableInfoFor("Show Fullscreen Vitals").HardcodedConfig().Boolean(),
	55: GameVariableInfoFor("Show Map Notes").HardcodedConfig().Boolean(),
	56: GameVariableInfoFor("Game: Wing Level").At(0).LimitedBy(0, 100),
	57: GameVariableInfoFor("HUD Color Bank").HardcodedConfig().LimitedBy(0, 2),
	58: GameVariableInfoFor("Audio Channels").HardcodedConfig().Enumerated(map[int16]string{
		0: "2",
		1: "4",
		2: "8",
	}),
}

var engineBooleanVariables = GameVariables{}

func unusedVariable() GameVariableInfo {
	return GameVariableInfoFor("(unused)").At(0)
}

// EngineIntegerVariable returns a variable info for integer variables.
func EngineIntegerVariable(index int) GameVariableInfo {
	return engineIntegerVariables.Lookup(index, unusedVariable)
}

// EngineBooleanVariable returns a variable info for boolean variables.
func EngineBooleanVariable(index int) GameVariableInfo {
	return engineBooleanVariables.Lookup(index, unusedVariable)
}
