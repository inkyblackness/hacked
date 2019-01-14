package world

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"

	"io/ioutil"
)

// IsSavegame returns true for resources that most likely identify a savegame.
// A savegame is one that has a state resource (0x0FA1) and hacker's health is more than zero.
func IsSavegame(provider resource.Provider) bool {
	res, err := provider.View(ids.GameState)
	if err != nil {
		return false
	}
	if res.BlockCount() < 1 {
		return false
	}
	dataReader, err := res.Block(0)
	if err != nil {
		return false
	}
	data, err := ioutil.ReadAll(dataReader)
	if err != nil {
		return false
	}
	healthOffset := 0x009C

	return (len(data) > healthOffset) && data[healthOffset] > 0
}
