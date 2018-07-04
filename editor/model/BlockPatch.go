package model

import "github.com/inkyblackness/hacked/ss1/resource"

// BlockPatch describes a raw delta-change of a resource block.
// The change can be applied, and reverted.
type BlockPatch struct {
	// ID identifier of the resource
	ID resource.ID
	// BlockIndex identifies the block for which the patch applies.
	BlockIndex int
	// BlockLength describes the expected length of the block.
	BlockLength int

	// ForwardData contains the delta information to patch to a new version.
	ForwardData []byte
	// ReverseData contains the delta information to change a new version back to the original.
	ReverseData []byte
}
