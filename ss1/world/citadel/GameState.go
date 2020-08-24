package citadel

import "github.com/inkyblackness/hacked/ss1/content/archive"

// DefaultGameState returns a new instance of a game state that the engine creates for a standard game.
func DefaultGameState() *archive.GameState {
	state := archive.DefaultGameState()

	state.Set("Current Level", 1)

	// location: in the neurosurgery chamber, looking West
	state.Set("Hacker Position X", (30<<16)+0x8000)
	state.Set("Hacker Position Y", (22<<16)+0x8000)
	state.Set("Hacker Position Z", 0x01BD00)
	state.Set("Hacker Yaw", 0x03243E)

	// set first message
	// TODO: use proper accessor
	data := state.Raw()
	data[0x0357+26] = archive.MessageStatusReceived // Rebecca Lansing's first message
	data[0x0519+9] = 0xFF                           // HUD active email -- set for similarity, has no effect.

	for index, info := range booleanVariables {
		if info.InitValue == nil {
			continue
		}
		state.SetBooleanVar(index, *info.InitValue != 0)
	}
	for index, info := range integerVariables {
		if info.InitValue == nil {
			continue
		}
		state.SetIntegerVar(index, *info.InitValue)
	}

	return state
}
