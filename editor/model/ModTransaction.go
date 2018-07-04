package model

import (
	"bytes"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial/rle"
)

type modAction func(mod *Mod)

// ModTransaction is used to modify a mod. It allows modifications of related resources in one atomic action.
type ModTransaction struct {
	actions     []modAction
	modifiedIDs resource.IDMarkerMap
}

// SetResource changes the meta information about a resource.
// Should the resource exist in multiple languages, all are modified.
//
// This is a low-level function and should not be required on a regular basis.
// Mods typically extend on already existing resources, and/or the editor itself should have a list of templates for
// new resources.
func (trans *ModTransaction) SetResource(id resource.ID,
	compound bool, contentType resource.ContentType, compressed bool) {
	setResource := func(res *MutableResource) {
		res.compound = compound
		res.contentType = contentType
		res.compressed = compressed
	}
	trans.actions = append(trans.actions, func(mod *Mod) {
		for _, lang := range resource.Languages() {
			setResource(mod.ensureResource(lang, id))
		}
		setResource(mod.ensureResource(resource.LangAny, id))
	})
	trans.modifiedIDs.Add(id)
}

// SetResourceBlock changes the block data of a resource.
//
// If the block data is not empty, then:
// If the resource does not exist, it will be created with default meta information.
// If the block does not exist, the resource is extended to allow its addition.
//
// If the block data is empty (or nil), then the block is cleared.
// If the resource is a compound list, then the underlying data will become visible again.
func (trans *ModTransaction) SetResourceBlock(lang resource.Language, id resource.ID, index int, data []byte) {
	trans.actions = append(trans.actions, func(mod *Mod) {
		res := mod.ensureResource(lang, id)
		res.setBlock(index, data)
	})
	trans.modifiedIDs.Add(id)
}

// PatchResourceBlock modifies an existing block.
// This modification assumes the block already exists and can take the given patch data.
// The patch data is expected to be produced by rle.Compress().
func (trans *ModTransaction) PatchResourceBlock(lang resource.Language, id resource.ID, index int, expectedLength int, patch []byte) {
	trans.actions = append(trans.actions, func(mod *Mod) {
		res := mod.ensureResource(lang, id)
		if res.isBlockIndexValid(index) && (len(res.blocks[index]) == expectedLength) {
			rle.Decompress(bytes.NewReader(patch), res.blocks[index])
		}
	})
	trans.modifiedIDs.Add(id)
}

// SetResourceBlocks sets the entire list of block data of a resource.
// This method is primarily meant for compound non-list resources (e.g. text pages).
func (trans *ModTransaction) SetResourceBlocks(lang resource.Language, id resource.ID, data [][]byte) {
	trans.actions = append(trans.actions, func(mod *Mod) {
		res := mod.ensureResource(lang, id)
		res.setBlocks(data)
	})
	trans.modifiedIDs.Add(id)
}

// DelResource removes a resource from the mod in the given language.
//
// After the deletion, all the underlying data of the world will become visible again.
func (trans *ModTransaction) DelResource(lang resource.Language, id resource.ID) {
	trans.actions = append(trans.actions, func(mod *Mod) {
		mod.delResource(lang, id)
	})
	trans.modifiedIDs.Add(id)
}
