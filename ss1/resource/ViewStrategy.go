package resource

// ViewStrategy defines how selected resources shall be viewed.
type ViewStrategy interface {
	// IsCompoundList returns true for compound resources where each contained block is a separate entity.
	// Separate entities are those that can be replaced without affecting others.
	// Examples are the small game textures and the list of object names.
	IsCompoundList(id ID) bool
}
