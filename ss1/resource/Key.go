package resource

// Key identifies a particular resource out of a combination of a resource ID, a language, and a block index.
//
// Keys are used for retrieving particular content that is stored in resources.
// Not all resources contain only one atomic entity.
type Key struct {
	// ID identifies the resource.
	ID ID
	// Lang localizes the resource, should the actual content be localized.
	Lang Language
	// Index specifies the exact entry within a resource.
	Index int
}

// KeyOf is a factory function for creating a Key instance.
// It is done to avoid getting warnings when trying to do a similar construction via Key{id, lang, index}.
func KeyOf(id ID, lang Language, index int) Key {
	return Key{ID: id, Lang: lang, Index: index}
}
