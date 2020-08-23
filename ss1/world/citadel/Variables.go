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

var booleanVariables = archive.GameVariables{
	1:   archive.GameVariableInfoFor("(unknown)").At(1),
	2:   archive.GameVariableInfoFor("(unknown)").At(1),
	3:   archive.GameVariableInfoFor("(unknown)").At(1),
	16:  archive.GameVariableInfoFor("(unknown)").At(1),
	18:  archive.GameVariableInfoFor("(unknown)").At(1),
	21:  archive.GameVariableInfoFor("(unknown)").At(1),
	22:  archive.GameVariableInfoFor("(unknown)").At(1),
	23:  archive.GameVariableInfoFor("(unknown)").At(1),
	24:  archive.GameVariableInfoFor("(unknown)").At(1),
	25:  archive.GameVariableInfoFor("(unknown)").At(1),
	26:  archive.GameVariableInfoFor("(unknown)").At(1),
	32:  archive.GameVariableInfoFor("(unknown)").At(1),
	33:  archive.GameVariableInfoFor("(unknown)").At(1),
	36:  archive.GameVariableInfoFor("(unknown)").At(1),
	37:  archive.GameVariableInfoFor("(unknown)").At(1),
	75:  archive.GameVariableInfoFor("(unknown)").At(1),
	76:  archive.GameVariableInfoFor("(unknown)").At(1),
	77:  archive.GameVariableInfoFor("(unknown)").At(1),
	78:  archive.GameVariableInfoFor("(unknown)").At(1),
	79:  archive.GameVariableInfoFor("(unknown)").At(1),
	80:  archive.GameVariableInfoFor("(unknown)").At(1),
	81:  archive.GameVariableInfoFor("(unknown)").At(1),
	82:  archive.GameVariableInfoFor("(unknown)").At(1),
	83:  archive.GameVariableInfoFor("(unknown)").At(1),
	84:  archive.GameVariableInfoFor("(unknown)").At(1),
	85:  archive.GameVariableInfoFor("(unknown)").At(1),
	86:  archive.GameVariableInfoFor("(unknown)").At(1),
	87:  archive.GameVariableInfoFor("(unknown)").At(1),
	88:  archive.GameVariableInfoFor("(unknown)").At(1),
	89:  archive.GameVariableInfoFor("(unknown)").At(1),
	90:  archive.GameVariableInfoFor("(unknown)").At(1),
	91:  archive.GameVariableInfoFor("(unknown)").At(1),
	92:  archive.GameVariableInfoFor("(unknown)").At(1),
	93:  archive.GameVariableInfoFor("(unknown)").At(1),
	94:  archive.GameVariableInfoFor("(unknown)").At(1),
	95:  archive.GameVariableInfoFor("(unknown)").At(1),
	112: archive.GameVariableInfoFor("(unknown)").At(1),
	113: archive.GameVariableInfoFor("(unknown)").At(1),
	114: archive.GameVariableInfoFor("(unknown)").At(1),
	115: archive.GameVariableInfoFor("(unknown)").At(1),
	116: archive.GameVariableInfoFor("(unknown)").At(1),
	117: archive.GameVariableInfoFor("(unknown)").At(1),
	118: archive.GameVariableInfoFor("(unknown)").At(1),
	119: archive.GameVariableInfoFor("(unknown)").At(1),
	120: archive.GameVariableInfoFor("(unknown)").At(1),
	121: archive.GameVariableInfoFor("(unknown)").At(1),
	122: archive.GameVariableInfoFor("(unknown)").At(1),
	123: archive.GameVariableInfoFor("(unknown)").At(1),
	124: archive.GameVariableInfoFor("(unknown)").At(1),
	125: archive.GameVariableInfoFor("(unknown)").At(1),
	126: archive.GameVariableInfoFor("(unknown)").At(1),
	127: archive.GameVariableInfoFor("(unknown)").At(1),

	152: archive.GameVariableInfoFor("Bridge Separated").At(0),

	160: archive.GameVariableInfoFor("(unknown)").At(1),
	161: archive.GameVariableInfoFor("(unknown)").At(1),
	162: archive.GameVariableInfoFor("(unknown)").At(1),
	163: archive.GameVariableInfoFor("(unknown)").At(1),
	164: archive.GameVariableInfoFor("(unknown)").At(1),
	165: archive.GameVariableInfoFor("(unknown)").At(1),
	166: archive.GameVariableInfoFor("(unknown)").At(1),
	167: archive.GameVariableInfoFor("(unknown)").At(1),
	168: archive.GameVariableInfoFor("(unknown)").At(1),
	169: archive.GameVariableInfoFor("(unknown)").At(1),
	192: archive.GameVariableInfoFor("(unknown)").At(1),
	193: archive.GameVariableInfoFor("(unknown)").At(1),
	194: archive.GameVariableInfoFor("(unknown)").At(1),
	195: archive.GameVariableInfoFor("(unknown)").At(1),
	196: archive.GameVariableInfoFor("(unknown)").At(1),
	197: archive.GameVariableInfoFor("(unknown)").At(1),
	198: archive.GameVariableInfoFor("(unknown)").At(1),
	199: archive.GameVariableInfoFor("(unknown)").At(1),
	200: archive.GameVariableInfoFor("(unknown)").At(1),
	201: archive.GameVariableInfoFor("(unknown)").At(1),
	202: archive.GameVariableInfoFor("(unknown)").At(1),
	203: archive.GameVariableInfoFor("(unknown)").At(1),
	204: archive.GameVariableInfoFor("(unknown)").At(1),
	205: archive.GameVariableInfoFor("(unknown)").At(1),
	206: archive.GameVariableInfoFor("(unknown)").At(1),
	207: archive.GameVariableInfoFor("(unknown)").At(1),
	225: archive.GameVariableInfoFor("(unknown)").At(1),
	227: archive.GameVariableInfoFor("(unknown)").At(1),
	229: archive.GameVariableInfoFor("(unknown)").At(1),
	231: archive.GameVariableInfoFor("(unknown)").At(1),
	233: archive.GameVariableInfoFor("(unknown)").At(1),
	235: archive.GameVariableInfoFor("(unknown)").At(1),
	237: archive.GameVariableInfoFor("(unknown)").At(1),
	239: archive.GameVariableInfoFor("(unknown)").At(1),
	241: archive.GameVariableInfoFor("(unknown)").At(1),
	243: archive.GameVariableInfoFor("(unknown)").At(1),
	245: archive.GameVariableInfoFor("(unknown)").At(1),
	247: archive.GameVariableInfoFor("(unknown)").At(1),
	249: archive.GameVariableInfoFor("(unknown)").At(1),
	251: archive.GameVariableInfoFor("(unknown)").At(1),
	253: archive.GameVariableInfoFor("(unknown)").At(1),
	255: archive.GameVariableInfoFor("(unknown)").At(1),
	257: archive.GameVariableInfoFor("(unknown)").At(1),
	259: archive.GameVariableInfoFor("(unknown)").At(1),
	261: archive.GameVariableInfoFor("(unknown)").At(1),
	263: archive.GameVariableInfoFor("(unknown)").At(1),
	265: archive.GameVariableInfoFor("(unknown)").At(1),
	267: archive.GameVariableInfoFor("(unknown)").At(1),
	269: archive.GameVariableInfoFor("(unknown)").At(1),
	271: archive.GameVariableInfoFor("(unknown)").At(1),
	273: archive.GameVariableInfoFor("(unknown)").At(1),
	275: archive.GameVariableInfoFor("(unknown)").At(1),
	277: archive.GameVariableInfoFor("(unknown)").At(1),
	279: archive.GameVariableInfoFor("(unknown)").At(1),
	281: archive.GameVariableInfoFor("(unknown)").At(1),
	283: archive.GameVariableInfoFor("(unknown)").At(1),
	285: archive.GameVariableInfoFor("(unknown)").At(1),
	287: archive.GameVariableInfoFor("(unknown)").At(1),
	289: archive.GameVariableInfoFor("(unknown)").At(1),
	291: archive.GameVariableInfoFor("(unknown)").At(1),
	293: archive.GameVariableInfoFor("(unknown)").At(1),
	295: archive.GameVariableInfoFor("(unknown)").At(1),
	297: archive.GameVariableInfoFor("(unknown)").At(1),
	299: archive.GameVariableInfoFor("(unknown)").At(1),
}

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
