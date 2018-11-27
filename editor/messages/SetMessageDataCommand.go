package messages

import (
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
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

func (cmd setMessageDataCommand) Do(trans cmd.Transaction) error {
	return cmd.perform(trans, func(entry messageDataEntry) [][]byte { return entry.newData })
}

func (cmd setMessageDataCommand) Undo(trans cmd.Transaction) error {
	return cmd.perform(trans, func(entry messageDataEntry) [][]byte { return entry.oldData })
}

func (cmd setMessageDataCommand) perform(trans cmd.Transaction, dataResolver func(messageDataEntry) [][]byte) error {
	cmd.saveEntries(trans, cmd.key.ID.Plus(cmd.key.Index), cmd.textEntries, dataResolver)
	cmd.saveEntries(trans, cmd.key.ID.Plus(cmd.key.Index).Plus(300), cmd.audioEntries, dataResolver)

	cmd.model.restoreFocus = true
	cmd.model.currentKey = cmd.key
	cmd.model.showVerboseText = cmd.showVerboseText
	return nil
}

func (cmd setMessageDataCommand) saveEntries(trans cmd.Transaction, id resource.ID,
	entries map[resource.Language]messageDataEntry,
	dataResolver func(messageDataEntry) [][]byte) {
	for lang, entry := range entries {
		data := dataResolver(entry)
		if len(data) > 0 {
			trans.SetResourceBlocks(lang, id, data)
		} else {
			trans.DelResource(lang, id)
		}
	}
}
