package messages

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
)

type messageDataEntry struct {
	oldData [][]byte
	newData [][]byte
}

type setMessageDataCommand struct {
	model *viewModel

	key             resource.Key
	showVerboseText bool

	textEntries  map[resource.Language]messageDataEntry
	audioEntries map[resource.Language]messageDataEntry
}

func (cmd setMessageDataCommand) Do(modder world.Modder) error {
	return cmd.perform(modder, func(entry messageDataEntry) [][]byte { return entry.newData })
}

func (cmd setMessageDataCommand) Undo(modder world.Modder) error {
	return cmd.perform(modder, func(entry messageDataEntry) [][]byte { return entry.oldData })
}

func (cmd setMessageDataCommand) perform(modder world.Modder, dataResolver func(messageDataEntry) [][]byte) error {
	cmd.saveEntries(modder, cmd.key.ID.Plus(cmd.key.Index), cmd.textEntries, dataResolver)
	cmd.saveEntries(modder, cmd.key.ID.Plus(cmd.key.Index).Plus(300), cmd.audioEntries, dataResolver)

	cmd.model.restoreFocus = true
	cmd.model.currentKey = cmd.key
	cmd.model.showVerboseText = cmd.showVerboseText
	return nil
}

func (cmd setMessageDataCommand) saveEntries(modder world.Modder, id resource.ID,
	entries map[resource.Language]messageDataEntry,
	dataResolver func(messageDataEntry) [][]byte) {
	for lang, entry := range entries {
		data := dataResolver(entry)
		if len(data) > 0 {
			modder.SetResourceBlocks(lang, id, data)
		} else {
			modder.DelResource(lang, id)
		}
	}
}
