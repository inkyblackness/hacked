package citadel

import "github.com/inkyblackness/hacked/ss1/content/archive"

var integerVariables = archive.GameVariables{
	3: archive.GameVariableInfoFor("Engine State").At(2).LimitedBy(0, 100),

	12: archive.GameVariableInfoFor("Number of Groves").At(3).LimitedBy(0, 4),

	33: archive.GameVariableInfoFor("Destroyed CPUs Level 1").At(0).LimitedBy(0, 100),
	34: archive.GameVariableInfoFor("Destroyed CPUs Level 2").At(0).LimitedBy(0, 100),
	35: archive.GameVariableInfoFor("Destroyed CPUs Level 3").At(0).LimitedBy(0, 100),
	36: archive.GameVariableInfoFor("Destroyed CPUs Level 4").At(0).LimitedBy(0, 100),
	37: archive.GameVariableInfoFor("Destroyed CPUs Level 5").At(0).LimitedBy(0, 100),
	38: archive.GameVariableInfoFor("Destroyed CPUs Level 6").At(0).LimitedBy(0, 100),

	39: archive.GameVariableInfoFor("Killed Cyberguards, largce CSpace").At(0).LimitedBy(0, 100),
}

var booleanVariables = archive.GameVariables{}

// IntegerVariable returns a variable info for integer variables.
func IntegerVariable(index int) archive.GameVariableInfo {
	return integerVariables.Lookup(index, func() archive.GameVariableInfo {
		info := archive.EngineIntegerVariable(index)
		if index == 31 {
			info.Name = "Reactor Code 1"
		}
		if index == 32 {
			info.Name = "Reactor Code 2"
		}
		return info
	})
}

// BooleanVariable returns a variable info for boolean variables.
func BooleanVariable(index int) archive.GameVariableInfo {
	return booleanVariables.Lookup(index, func() archive.GameVariableInfo {
		return archive.EngineBooleanVariable(index)
	})
}
