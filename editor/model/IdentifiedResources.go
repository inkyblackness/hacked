package model

import (
	"sort"

	"github.com/inkyblackness/hacked/ss1/resource"
)

// IdentifiedResources is a map of mutable resources by their identifier.
type IdentifiedResources map[resource.ID]*MutableResource

// Add puts the entries of the provided container into this container.
func (res IdentifiedResources) Add(other IdentifiedResources) {
	for id, entry := range other {
		res[id] = entry
	}
}

// IDs returns a list of all IDs in the map.
func (res IdentifiedResources) IDs() []resource.ID {
	result := make([]resource.ID, 0, len(res))
	for id := range res {
		result = append(result, id)
	}
	sort.Slice(result, func(a, b int) bool {
		idA := result[a]
		idB := result[b]
		entryA := res[idA]
		entryB := res[idB]
		if entryA.saveOrder == entryB.saveOrder {
			return idA < idB
		}
		return entryA.saveOrder < entryB.saveOrder
	})
	return result
}

// View returns the mutable resource as a read-only view.
func (res IdentifiedResources) View(id resource.ID) (resource.View, error) {
	entry, existing := res[id]
	if !existing {
		return nil, resource.ErrResourceDoesNotExist(id)
	}
	return entry, nil
}
