package resource

import "io/ioutil"

// IsSavegame returns true for resources that most likely identify a savegame.
// A savegame is one that has a state resource (0x0FA1) and hacker's health is more than zero.
func IsSavegame(provider Provider) bool {
	stateID := ID(0x0FA1)
	res, err := provider.Resource(stateID)
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
