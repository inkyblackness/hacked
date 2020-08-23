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

// DescribedAs returns an information with given text as description.
func (info GameVariableInfo) DescribedAs(text string) GameVariableInfo {
	info.Description = text
	return info
}

// GameVariables is a lookup map for information on game variables.
type GameVariables map[int]GameVariableInfo

// Lookup returns the information for given index. If the index is not known, nil is returned.
func (vars GameVariables) Lookup(index int) *GameVariableInfo {
	info, known := vars[index]
	if !known {
		return nil
	}
	return &info
}

const (
	securityValueDescription = "The current security value is re-calculated whenever a level is loaded up."
	randomCodeDescription    = "If both codes are equal at the start of a new game, they will be randomized."
	highscoreCodeDescription = "The highscore values are combined as one 32-bit integer, used for the MFD game MCOM."
)

var engineIntegerVariables = GameVariables{
	2: GameVariableInfoFor("Plastique explosion counter").At(0).LimitedBy(0, 100).
		DescribedAs("This value is incremented for each exploding plastique."),

	9: GameVariableInfoFor("Plot counter").At(0),

	13: GameVariableInfoFor("Difficulty: Mission").HardcodedConfig().LimitedBy(0, 3),
	14: GameVariableInfoFor("Difficulty: Cyber").HardcodedConfig().LimitedBy(0, 3),
	15: GameVariableInfoFor("Difficulty: Combat").HardcodedConfig().LimitedBy(0, 3),

	16: GameVariableInfoFor("Security Value: Level 0").HardcodedAt(0).LimitedBy(0, 1000).DescribedAs(securityValueDescription),
	17: GameVariableInfoFor("Security Value: Level 1").HardcodedAt(0).LimitedBy(0, 1000).DescribedAs(securityValueDescription),
	18: GameVariableInfoFor("Security Value: Level 2").HardcodedAt(0).LimitedBy(0, 1000).DescribedAs(securityValueDescription),
	19: GameVariableInfoFor("Security Value: Level 3").HardcodedAt(0).LimitedBy(0, 1000).DescribedAs(securityValueDescription),
	20: GameVariableInfoFor("Security Value: Level 4").HardcodedAt(0).LimitedBy(0, 1000).DescribedAs(securityValueDescription),
	21: GameVariableInfoFor("Security Value: Level 5").HardcodedAt(0).LimitedBy(0, 1000).DescribedAs(securityValueDescription),
	22: GameVariableInfoFor("Security Value: Level 6").HardcodedAt(0).LimitedBy(0, 1000).DescribedAs(securityValueDescription),
	23: GameVariableInfoFor("Security Value: Level 7").HardcodedAt(0).LimitedBy(0, 1000).DescribedAs(securityValueDescription),
	24: GameVariableInfoFor("Security Value: Level 8").HardcodedAt(0).LimitedBy(0, 1000).DescribedAs(securityValueDescription),
	25: GameVariableInfoFor("Security Value: Level 9").HardcodedAt(0).LimitedBy(0, 1000).DescribedAs(securityValueDescription),
	26: GameVariableInfoFor("Security Value: Level 10").HardcodedAt(0).LimitedBy(0, 1000).DescribedAs(securityValueDescription),
	27: GameVariableInfoFor("Security Value: Level 11").HardcodedAt(0).LimitedBy(0, 1000).DescribedAs(securityValueDescription),
	28: GameVariableInfoFor("Security Value: Level 12").HardcodedAt(0).LimitedBy(0, 1000).DescribedAs(securityValueDescription),
	29: GameVariableInfoFor("Security Value: Level 13").HardcodedAt(0).LimitedBy(0, 1000).DescribedAs(securityValueDescription),

	30: GameVariableInfoFor("Difficulty: Puzzle").HardcodedConfig().LimitedBy(0, 3),
	31: GameVariableInfoFor("Random Code 1").At(0).DescribedAs(randomCodeDescription),
	32: GameVariableInfoFor("Random Code 2").At(0).DescribedAs(randomCodeDescription),

	41: GameVariableInfoFor("Music Volume").HardcodedConfig().LimitedBy(0, 100),
	42: GameVariableInfoFor("Video Gamma").HardcodedConfig().LimitedBy(0, math.MaxInt16),
	43: GameVariableInfoFor("SFX Volume").HardcodedConfig().LimitedBy(0, 100),
	44: GameVariableInfoFor("Mouse Handedness").HardcodedConfig().Enumerated(
		map[int16]string{
			0: "Right Handed",
			1: "Left Handed",
		}),
	45: GameVariableInfoFor("Game: Highscore (pt1)").At(0).DescribedAs(highscoreCodeDescription),
	46: GameVariableInfoFor("Game: Highscore (pt2)").At(0).DescribedAs(highscoreCodeDescription),
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

var engineBooleanVariables = GameVariables{
	0: GameVariableInfoFor("Always False").At(0).
		DescribedAs("This variable should always stay at 'False'.\n" +
			"Default conditions assume boolean var 0 is zero to not block.\n" +
			"For example, default doors would all become locked if it were 'True'."),

	10: GameVariableInfoFor("Status HW: Delta Launch Enable").At(0),
	11: GameVariableInfoFor("Status HW: Alpha Launch Enable").At(0),
	12: GameVariableInfoFor("Status HW: Beta Launch Enable").At(0),
	15: GameVariableInfoFor("Status HW: Beta Launched").At(0),

	20: GameVariableInfoFor("Reactor on Destruct").At(0).
		DescribedAs("Note: Rumble will stop if boolean var 152 is 'True'."),

	145: GameVariableInfoFor("On-Line Help"),

	152: GameVariableInfoFor("Self-Destruct Rumble Stop").At(0).
		DescribedAs("Disables rumble caused by boolean var 20."),
	153: GameVariableInfoFor("Four or more plastique exploded").At(0).
		DescribedAs("Set to 'True' if integer variable 2 reached 4 (or higher)."),

	300: GameVariableInfoFor("New Message Flag").At(0),
}

// EngineIntegerVariable returns a variable info for integer variables.
// If the given index is not used by the engine, nil is returned.
func EngineIntegerVariable(index int) *GameVariableInfo {
	return engineIntegerVariables.Lookup(index)
}

// EngineBooleanVariable returns a variable info for boolean variables.
// If the given index is not used by the engine, nil is returned.
func EngineBooleanVariable(index int) *GameVariableInfo {
	return engineBooleanVariables.Lookup(index)
}

// IsRandomIntegerVariable returns true for the special variables that are randomized.
func IsRandomIntegerVariable(index int) bool {
	return (index == 31) || (index == 32)
}

// EngineVariables is a collector of engine-specific variable accessors.
type EngineVariables struct{}

var unusedVar = GameVariableInfoFor("(unused)").At(0)

// IntegerVariable returns a variable info for given index.
func (vars EngineVariables) IntegerVariable(index int) GameVariableInfo {
	varInfo := EngineIntegerVariable(index)
	if varInfo == nil {
		return unusedVar
	}
	return *varInfo
}

// BooleanVariable returns a variable info for given index.
func (vars EngineVariables) BooleanVariable(index int) GameVariableInfo {
	varInfo := EngineBooleanVariable(index)
	if varInfo == nil {
		return unusedVar
	}
	return *varInfo
}
