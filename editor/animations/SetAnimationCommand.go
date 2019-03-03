package animations

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
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

func (command setAnimationCommand) Do(modder world.Modder) error {
	return command.perform(modder, command.newAnimation, command.newFrames)
}

func (command setAnimationCommand) Undo(modder world.Modder) error {
	return command.perform(modder, command.oldAnimation, command.oldFrames)
}

func (command setAnimationCommand) perform(modder world.Modder, animData []byte, frames [][]byte) error {
	if len(frames) == 0 {
		modder.DelResource(command.animationKey.Lang, command.animationKey.ID.Plus(command.animationKey.Index))
		modder.DelResource(command.animationKey.Lang, command.framesID)
	} else {
		modder.SetResourceBlock(command.animationKey.Lang, command.animationKey.ID.Plus(command.animationKey.Index), 0, animData)
		modder.SetResourceBlocks(command.animationKey.Lang, command.framesID, frames)
	}

	command.model.restoreFocus = true
	command.model.currentKey = command.animationKey
	return nil
}
