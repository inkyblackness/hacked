package archive

import (
	"io/ioutil"

	"github.com/inkyblackness/hacked/ss1/resource"
)

// IsSavegame returns true for resources that most likely identify a savegame.
// This function must be called with a resource from ID 0x0FA1.
func IsSavegame(res resource.View) bool {
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
	state := NewGameState(data)
	return state.IsSavegame()
}
