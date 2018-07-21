package messages

import (
	"github.com/inkyblackness/hacked/editor/cmd"
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

	entries map[resource.Language]messageDataEntry
}

func (cmd setMessageDataCommand) Do(trans cmd.Transaction) error {
	return cmd.perform(trans, func(entry messageDataEntry) [][]byte { return entry.newData })
}

func (cmd setMessageDataCommand) Undo(trans cmd.Transaction) error {
	return cmd.perform(trans, func(entry messageDataEntry) [][]byte { return entry.oldData })
}

func (cmd setMessageDataCommand) perform(trans cmd.Transaction, dataResolver func(messageDataEntry) [][]byte) error {
	for lang, entry := range cmd.entries {
		data := dataResolver(entry)
		id := cmd.key.ID.Plus(cmd.key.Index)
		if len(data) > 0 {
			trans.SetResourceBlocks(lang, id, data)
		} else {
			trans.DelResource(lang, id)
		}
	}
	cmd.model.restoreFocus = true
	cmd.model.currentKey = cmd.key
	cmd.model.showVerboseText = cmd.showVerboseText
	return nil
}
