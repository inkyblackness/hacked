package world

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// Modder describes actions meant to modify resources.
type Modder interface {

	// SetResourceBlock changes the block data of a resource.
	//
	// If the block data is not empty, then:
	// If the resource does not exist, it will be created with default meta information.
	// If the block does not exist, the resource is extended to allow its addition.
	//
	// If the block data is empty (or nil), then the block is cleared.
	// If the resource is a compound list, then the underlying data will become visible again.
	SetResourceBlock(lang resource.Language, id resource.ID, index int, data []byte)

	// PatchResourceBlock modifies an existing block.
	// This modification assumes the block already exists and can take the given patch data.
	// The patch data is expected to be produced by rle.Compress(). (see also: Mod.CreateBlockPatch)
	PatchResourceBlock(lang resource.Language, id resource.ID, index int, expectedLength int, patch []byte)

	// SetResourceBlocks sets the entire list of block data of a resource.
	// This method is primarily meant for compound non-list resources (e.g. text pages).
	SetResourceBlocks(lang resource.Language, id resource.ID, data [][]byte)

	// DelResource removes a resource from the mod in the given language.
	//
	// After the deletion, all the underlying data of the world will become visible again.
	DelResource(lang resource.Language, id resource.ID)

	// SetTextureProperties updates the properties of a specific texture.
	SetTextureProperties(textureIndex int, properties texture.Properties)

	// SetObjectProperties updates the properties of a specific object.
	SetObjectProperties(triple object.Triple, properties object.Properties)
}
