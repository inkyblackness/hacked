package animations

import (
	"github.com/inkyblackness/hacked/editor/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

type setAnimationCommand struct {
	model *viewModel

	animationKey resource.Key
	framesID     resource.ID

	oldAnimation []byte
	newAnimation []byte

	oldFrames [][]byte
	newFrames [][]byte
}

func (command setAnimationCommand) Do(trans cmd.Transaction) error {
	return command.perform(trans, command.newAnimation, command.newFrames)
}

func (command setAnimationCommand) Undo(trans cmd.Transaction) error {
	return command.perform(trans, command.oldAnimation, command.oldFrames)
}

func (command setAnimationCommand) perform(trans cmd.Transaction, animData []byte, frames [][]byte) error {
	if len(frames) == 0 {
		trans.DelResource(command.animationKey.Lang, command.animationKey.ID.Plus(command.animationKey.Index))
		trans.DelResource(command.animationKey.Lang, command.framesID)
	} else {
		trans.SetResourceBlock(command.animationKey.Lang, command.animationKey.ID.Plus(command.animationKey.Index), 0, animData)
		trans.SetResourceBlocks(command.animationKey.Lang, command.framesID, frames)
	}

	command.model.restoreFocus = true
	command.model.currentKey = command.animationKey
	return nil
}
