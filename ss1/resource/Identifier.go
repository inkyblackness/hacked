package resource

import "fmt"

// ID represents an integer key of resources.
type ID uint16

// Value returns the numerical value of the identifier.
func (id ID) Value() uint16 {
	return uint16(id)
}

// Plus adds the given offset and returns the resulting ID.
func (id ID) Plus(offset int) ID {
	return ID(int(id) + offset)
}

// String returns the textual representation.
func (id ID) String() string {
	return fmt.Sprintf("%04X", uint16(id))
}

// IDMarkerMap is used to collect IDs.
type IDMarkerMap struct {
	ids map[ID]struct{}
}

// Add adds the given ID to the map.
func (marker *IDMarkerMap) Add(id ID) {
	if marker.ids == nil {
		marker.ids = make(map[ID]struct{})
	}
	marker.ids[id] = struct{}{}
}

// ToList converts the map to a de-duplicated list.
func (marker IDMarkerMap) ToList() []ID {
	result := make([]ID, 0, len(marker.ids))
	for id := range marker.ids {
		result = append(result, id)
	}
	return result
}
