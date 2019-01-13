package resource

// Store is a resource provider that can be modified.
type Store struct {
	ids       []ID
	retriever map[uint16]func() (View, error)
}

func cloneIDs(source []ID) []ID {
	cloned := make([]ID, len(source))
	copy(cloned, source)
	return cloned
}

// NewProviderBackedStore returns a new store based on the specified provider.
// The created store will have a snapshot of the current provider; Changes in
// the list of available IDs of the given provider are not reflected.
func NewProviderBackedStore(provider Provider) *Store {
	store := &Store{
		ids:       cloneIDs(provider.IDs()),
		retriever: make(map[uint16]func() (View, error))}

	defaultToProvider := func(id ID) func() (View, error) {
		return func() (View, error) { return provider.Resource(id) }
	}
	for _, id := range store.ids {
		store.retriever[id.Value()] = defaultToProvider(id)
	}

	return store
}

// IDs returns a list of available IDs this store currently contains.
func (store *Store) IDs() []ID {
	return cloneIDs(store.ids)
}

// Resource returns a resource for the given identifier.
func (store *Store) Resource(id ID) (View, error) {
	retriever, existing := store.retriever[id.Value()]
	if !existing {
		return nil, ErrResourceDoesNotExist(id)
	}
	return retriever()
}

// Del removes the resource with given identifier from the store.
func (store *Store) Del(id ID) {
	key := id.Value()
	if _, existing := store.retriever[key]; existing {
		delete(store.retriever, key)
		newIDs := make([]ID, 0, len(store.ids)-1)
		for _, oldID := range store.ids {
			if oldID.Value() != key {
				newIDs = append(newIDs, oldID)
			}
		}
		store.ids = newIDs
	}
}

// Put (re-)assigns an identifier with data. If no resource with given ID exists,
// then it is created. Existing resources are overwritten with the provided data.
func (store *Store) Put(id ID, res View) {
	key := id.Value()
	if _, existing := store.retriever[key]; !existing {
		store.ids = append(store.ids, id)
	}
	store.retriever[key] = func() (View, error) { return res, nil }
}
