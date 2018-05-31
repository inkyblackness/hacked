package resource

// ResourceViewStrategy defines how selected resources shall be viewed.
type ResourceViewStrategy interface {
	// IsCompoundList returns true for compound resources where each contained block is a separate entity.
	// Separate entities are those that can be replaced without affecting others.
	// Examples are the small game textures and the list of object names.
	IsCompoundList(id ID) bool
}
