package model

import (
	"github.com/inkyblackness/hacked/ss1/resource"
)

// IdentifiedResources is a map of mutable resources by their identifier.
type IdentifiedResources map[resource.ID]*MutableResource

// IDs returns a list of all IDs in the map.
func (res IdentifiedResources) IDs() []resource.ID {
	result := make([]resource.ID, 0, len(res))
	for id := range res {
		result = append(result, id)
	}
	// TODO sort by source sequence
	return result
}

// Resource returns the mutable resource as a read-only view.
func (res IdentifiedResources) Resource(id resource.ID) (*resource.Resource, error) {
	entry, existing := res[id]
	if !existing {
		return nil, resource.ErrResourceDoesNotExist(id)
	}
	return &resource.Resource{
		Compressed:    entry.Compressed(),
		ContentType:   entry.ContentType(),
		Compound:      entry.Compound(),
		BlockProvider: entry,
	}, nil
}
