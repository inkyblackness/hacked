package citadel

import (
	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/world"
)

// DefaultGameState returns a new instance of a game state that the engine creates for a standard game.
func DefaultGameState() *archive.GameState {
	state := archive.DefaultGameState()

	// location: medical level
	state.Set("Current Level", world.StartingLevel)

	// location: in the neurosurgery chamber, looking West
	state.Set("Hacker Position X", (world.StartingTileX<<16)+0x8000)
	state.Set("Hacker Position Y", (world.StartingTileY<<16)+0x8000)
	state.Set("Hacker Position Z", 0x01BD00)
	state.Set("Hacker Yaw", 0x03243E)

	// set first message
	state.EMailState(26).SetReceived(true) // Rebecca Lansing's first message
	// The engine also sets the active state of email, though this field index is not used.

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
