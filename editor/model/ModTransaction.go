package model

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
)

type modAction func(mod *Mod)
type idMarkerMap map[resource.ID]bool

func (marker idMarkerMap) toList() []resource.ID {
	result := make([]resource.ID, 0, len(marker))
	for id := range marker {
		result = append(result, id)
	}
	return result
}

// ModTransaction is used to modify a mod. It allows modifications of related resources in one atomic action.
type ModTransaction struct {
	actions     []modAction
	modifiedIDs idMarkerMap
}

// SetResource changes the meta information about a resource.
// Should the resource exist in multiple languages, all are modified.
//
// This is a low-level function and should not be required on a regular basis.
// Mods typically extend on already existing resources, and/or the editor itself should have a list of templates for
// new resources.
func (trans *ModTransaction) SetResource(id resource.ID,
	compound bool, contentType resource.ContentType, compressed bool) {
	setResource := func(res *resource.Resource) {
		res.Compound = compound
		res.ContentType = contentType
		res.Compressed = compressed
	}
	trans.actions = append(trans.actions, func(mod *Mod) {
		for _, lang := range world.Languages() {
			setResource(mod.ensureResource(lang, id))
		}
		setResource(mod.ensureResource(world.LangAny, id))
	})
	trans.modifiedIDs[id] = true
}

// SetResourceBlock changes the block data of a resource.
//
// If the block data is not empty, then:
// If the resource does not exist, it will be created with default meta information.
// If the block does not exist, the resource is extended to allow its addition.
//
// If the block data is empty (or nil), then the block is cleared.
// If the resource is a compound list, then the underlying data will become visible again.
func (trans *ModTransaction) SetResourceBlock(lang world.Language, id resource.ID, index int, data []byte) {
	trans.actions = append(trans.actions, func(mod *Mod) {
		res := mod.ensureResource(lang, id)
		res.SetBlock(index, data)
	})
	trans.modifiedIDs[id] = true
}

// DelResource removes a resource from the mod in the given language.
//
// After the deletion, all the underlying data of the world will become visible again.
func (trans *ModTransaction) DelResource(lang world.Language, id resource.ID) {
	trans.actions = append(trans.actions, func(mod *Mod) {
		mod.delResource(lang, id)
	})
	trans.modifiedIDs[id] = true
}
