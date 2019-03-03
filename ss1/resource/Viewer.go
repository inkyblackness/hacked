package resource

// Viewer provides resources from some storage.
type Viewer interface {
	// IDs returns a list of available IDs this viewer can provide.
	IDs() []ID

	// View returns a read-only view to a resource for the given identifier.
	View(id ID) (View, error)
}
