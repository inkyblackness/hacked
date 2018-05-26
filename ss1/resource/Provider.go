package resource

// Provider provides resources from some storage.
type Provider interface {
	// IDs returns a list of available IDs this provider can provide.
	IDs() []ID

	// Resource returns a resource for the given identifier.
	Resource(id ID) (*Resource, error)
}
