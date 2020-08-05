package archive

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/content/interpreters"
	"github.com/inkyblackness/hacked/ss1/content/text"
)

// GameStateSize specifies the byte count of a serialized GameState.
const GameStateSize = 0x054D
const (
	stateHackerNameSize = 20
)

type GameState struct {
	*interpreters.Instance
}

var gameStateDesc = interpreters.New().
	With("Difficulty: Combat", 0x0015, 1).As(interpreters.RangedValue(0, 3)).
	With("Difficulty: Mission", 0x0016, 1).As(interpreters.RangedValue(0, 3)).
	With("Difficulty: Puzzle", 0x0017, 1).As(interpreters.RangedValue(0, 3)).
	With("Difficulty: Cyber", 0x0018, 1).As(interpreters.RangedValue(0, 3)).
	With("Current Level", 0x0039, 1).As(interpreters.RangedValue(0, 15)).
	With("Health", 0x009C, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int) string {
		return fmt.Sprintf("%.2f%%%%", float64(value)/255.0)
	})).
	With("Power", 0x00AC, 1).As(interpreters.FormattedRangedValue(0, 255,
	func(value int) string {
		return fmt.Sprintf("%.2f%%%%", float64(value)/255.0)
	})).
	With("Out-of-power flag", 0x00AF, 1).As(interpreters.EnumValue(
	map[uint32]string{
		0: "Powered",
		1: "Out-of-power",
	}))

func NewGameState(raw []byte) GameState {
	return GameState{Instance: gameStateDesc.For(raw)}
}

func (state *GameState) HackerName(cp text.Codepage) string {
	return cp.Decode(state.Raw()[:stateHackerNameSize])
}

func (state *GameState) SetHackerName(name string, cp text.Codepage) {
	raw := state.Raw()
	var zeroName [stateHackerNameSize]byte
	copy(raw[:stateHackerNameSize], zeroName[:])
	copy(raw[:stateHackerNameSize], cp.Encode(name))
	raw[stateHackerNameSize-1] = 0
}
