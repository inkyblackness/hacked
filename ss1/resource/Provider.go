package resource

// Provider provides resources from some storage.
type Provider interface {
	// IDs returns a list of available IDs this provider can provide.
	IDs() []ID

	// View returns a read-only view to a resource for the given identifier.
	View(id ID) (View, error)
}
