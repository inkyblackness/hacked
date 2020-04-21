package ids

// SoundEffectInfo describes one sound effect in the game.
type SoundEffectInfo struct {
	// Name is the unique identifier for the effect source.
	Name string
	// Index refers to the audio index. Multiple effects may use the same audio. -1 indicates no audio mapped.
	AudioIndex int
}

// SoundEffectsForAudio returns the effect information that share the same audio.
func SoundEffectsForAudio(index int) []SoundEffectInfo {
	var result []SoundEffectInfo
	for _, effect := range soundEffects {
		if effect.AudioIndex == index {
			result = append(result, effect)
		}
	}
	return result
}

// The following table is based on the constants found in sfxlist.h (non-demo case).
var soundEffects = []SoundEffectInfo{
	// Doors
	{Name: "DoorMetal", AudioIndex: 3},
	{Name: "DoorNormal", AudioIndex: 5},
	{Name: "DoorIris", AudioIndex: 6},
	{Name: "DoorBulkhead", AudioIndex: 67},
	{Name: "DoorGrating", AudioIndex: 90},

	// Ambient
	{Name: "Bridge1", AudioIndex: -1},
	{Name: "Bridge2", AudioIndex: -1},
	{Name: "Bridge3", AudioIndex: -1},
	{Name: "Bridge4", AudioIndex: -1},
	{Name: "Maint1", AudioIndex: 9},
	{Name: "Maint2", AudioIndex: -1},
	{Name: "Maint3", AudioIndex: -1},
	{Name: "Maint4", AudioIndex: -1},
	{Name: "Grove1", AudioIndex: 43},

	// Critters
	{Name: "Death1", AudioIndex: 11},
	{Name: "Death2", AudioIndex: 49},
	{Name: "Death3", AudioIndex: 50},
	{Name: "Death4", AudioIndex: 51},
	{Name: "Death5", AudioIndex: 53},
	{Name: "Death6", AudioIndex: 54},
	{Name: "Death7", AudioIndex: 68},
	{Name: "Death8", AudioIndex: 69},
	{Name: "Death9", AudioIndex: 88},
	{Name: "Death10", AudioIndex: 93},
	{Name: "Death11", AudioIndex: 101},
	{Name: "Attack1", AudioIndex: 12},
	{Name: "Attack4", AudioIndex: 46},
	{Name: "Attack5", AudioIndex: 48},
	{Name: "Attack6", AudioIndex: 52},
	{Name: "Attack7", AudioIndex: 55},
	{Name: "Attack8", AudioIndex: 63},
	{Name: "Attack9", AudioIndex: 16},
	{Name: "Notice1", AudioIndex: 58},
	{Name: "Notice2", AudioIndex: 59},
	{Name: "Notice3", AudioIndex: 74},
	{Name: "Notice4", AudioIndex: 75},
	{Name: "Notice5", AudioIndex: 100},
	{Name: "Near1", AudioIndex: 73},
	{Name: "Near2", AudioIndex: 56},
	{Name: "Near3", AudioIndex: 47},
	{Name: "Near4", AudioIndex: 25},

	// Wacky Objects
	{Name: "Repulsor", AudioIndex: -1},
	{Name: "ForceBridge", AudioIndex: 72},
	{Name: "TerrainElevLoop", AudioIndex: 15},
	{Name: "SparkingCable", AudioIndex: 87},
	{Name: "SurgeryMachine", AudioIndex: 102},

	// Combat
	{Name: "GunMinipistol", AudioIndex: 39},
	{Name: "GunDartpistol", AudioIndex: 86},
	{Name: "GunMagnum", AudioIndex: 40},
	{Name: "GunAssault", AudioIndex: 17},
	{Name: "GunRiot", AudioIndex: 41},
	{Name: "GunFlechette", AudioIndex: 38},
	{Name: "GunSkorpion", AudioIndex: 65},
	{Name: "GunMagpulse", AudioIndex: 45},
	{Name: "GunRailgun", AudioIndex: 29},
	{Name: "GunPipeHitMeat", AudioIndex: 4},
	{Name: "GunPipeHitMetal", AudioIndex: 21},
	{Name: "GunPipeMiss", AudioIndex: 24},
	{Name: "GunLaserepeeHit", AudioIndex: 31},
	{Name: "GunLaserepeeMiss", AudioIndex: 34},
	{Name: "GunPhaser", AudioIndex: 18},
	{Name: "GunBlaster", AudioIndex: 94},
	{Name: "GunIonbeam", AudioIndex: 95},
	{Name: "GunStungun", AudioIndex: 19},
	{Name: "GunPlasma", AudioIndex: 97},

	{Name: "PlayerHurt", AudioIndex: 64},
	{Name: "Shield1", AudioIndex: 32},
	{Name: "Shield2", AudioIndex: 20},
	{Name: "ShieldUp", AudioIndex: 96},
	{Name: "ShieldDown", AudioIndex: 42},
	{Name: "MetalSpang", AudioIndex: 89},
	{Name: "Radiation", AudioIndex: 2},

	{Name: "Reload1", AudioIndex: 22},
	{Name: "Reload2", AudioIndex: 23},
	{Name: "GrenadeArm", AudioIndex: 7},
	{Name: "BatteryUse", AudioIndex: 28},

	{Name: "Explosion1", AudioIndex: 44},
	{Name: "Rumble", AudioIndex: 106},
	{Name: "Teleport", AudioIndex: 103},

	{Name: "MonitorExplode", AudioIndex: 57},
	{Name: "CameraExplode", AudioIndex: 8},
	{Name: "CpuExplode", AudioIndex: 10},
	{Name: "DestroyCrate", AudioIndex: 37},
	{Name: "DestroyBarrel", AudioIndex: 109},

	// Cspace
	{Name: "Pulser", AudioIndex: -1},
	{Name: "Drill", AudioIndex: -1},
	{Name: "Disc", AudioIndex: -1},
	{Name: "Datastorm", AudioIndex: -1},

	{Name: "Recall", AudioIndex: -1},
	{Name: "Turbo", AudioIndex: -1},
	{Name: "Fakeid", AudioIndex: -1},
	{Name: "Decoy", AudioIndex: -1},

	{Name: "EnterCspace", AudioIndex: 27},
	{Name: "OttoShodan", AudioIndex: 30},
	{Name: "CyberDamage", AudioIndex: -1},
	{Name: "Cyberheal", AudioIndex: -1},
	{Name: "Cybertoggle", AudioIndex: -1},
	{Name: "IceDefense", AudioIndex: -1},

	{Name: "CyberAttack1", AudioIndex: -1},
	{Name: "CyberAttack2", AudioIndex: -1},
	{Name: "CyberAttack3", AudioIndex: -1},

	// MFD and UI Wackiness
	{Name: "VideoDown", AudioIndex: 33},
	{Name: "MfdButton", AudioIndex: 35},
	{Name: "InventButton", AudioIndex: 79},
	{Name: "InventSelect", AudioIndex: 80},
	{Name: "InventAdd", AudioIndex: 81},
	{Name: "InventWare", AudioIndex: 82},
	{Name: "PatchUse", AudioIndex: 83},
	{Name: "ZoomBox", AudioIndex: 84},
	{Name: "MapZoom", AudioIndex: 85},
	{Name: "MfdKeypad", AudioIndex: 76},
	{Name: "MfdBuzz", AudioIndex: 77},
	{Name: "MfdSuccess", AudioIndex: 78},
	{Name: "Goggle", AudioIndex: 0},
	{Name: "Hudfrob", AudioIndex: 1},
	{Name: "Static", AudioIndex: 2},
	{Name: "Email", AudioIndex: 107},

	// SHODAN
	{Name: "ShodanBark", AudioIndex: 30},
	{Name: "ShodanWeak", AudioIndex: 30},
	{Name: "ShodanStrong", AudioIndex: 30},

	// Other
	{Name: "PanelSuccess", AudioIndex: 78},
	{Name: "PowerOut", AudioIndex: 26},
	{Name: "EnergyDrain", AudioIndex: 14},
	{Name: "EnergyRecharge", AudioIndex: 98},
	{Name: "Surge", AudioIndex: 60},
	{Name: "Vmail", AudioIndex: 92},
	{Name: "DropItem", AudioIndex: 99},

	// 3d World
	{Name: "Button", AudioIndex: 71},
	{Name: "MechButton", AudioIndex: 36},
	{Name: "Bigbutton", AudioIndex: 61},
	{Name: "NormalLever", AudioIndex: 62},
	{Name: "Biglever", AudioIndex: 62},
	{Name: "Klaxon", AudioIndex: 70},

	// Plot
	{Name: "GroveJett", AudioIndex: 66},

	// Extra
	{Name: "PanelConfirm", AudioIndex: 13},
	{Name: "Death12", AudioIndex: 104},
	{Name: "Death13", AudioIndex: 105},
	{Name: "Death14", AudioIndex: 108},
	{Name: "Attack10", AudioIndex: 91},
	{Name: "Attack11", AudioIndex: 110},
	{Name: "Near5", AudioIndex: 111},
	{Name: "Death15", AudioIndex: 112},
	{Name: "FunPackAttackingUs", AudioIndex: 113},
}
