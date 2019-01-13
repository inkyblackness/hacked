package resource

type nullProvider struct{}

// NullProvider returns a Provider instance that is empty.
// It contains no IDs and will not provide any resource.
func NullProvider() Provider {
	return &nullProvider{}
}

func (*nullProvider) IDs() []ID {
	return nil
}

func (*nullProvider) Resource(id ID) (View, error) {
	return nil, ErrResourceDoesNotExist(id)
}
