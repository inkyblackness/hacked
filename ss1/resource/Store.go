package resource

import "io/ioutil"

// Store holds a set of resources. This set can be modified.
type Store struct {
	ids       []ID
	retriever map[uint16]func() (*Resource, error)
}

func cloneIDs(source []ID) []ID {
	cloned := make([]ID, len(source))
	copy(cloned, source)
	return cloned
}

// NewStore returns a new store instance.
func NewStore() *Store {
	store := &Store{
		ids:       nil,
		retriever: make(map[uint16]func() (*Resource, error))}

	return store
}

// IDs returns a list of available IDs this store currently contains.
func (store *Store) IDs() []ID {
	return cloneIDs(store.ids)
}

// Resource returns reference to the contained resource for given identifier.
func (store *Store) Resource(id ID) (*Resource, error) {
	retriever, existing := store.retriever[id.Value()]
	if !existing {
		return nil, ErrResourceDoesNotExist(id)
	}
	return retriever()
}

// View returns a read-only view on the resource for given identifier.
func (store *Store) View(id ID) (View, error) {
	res, err := store.Resource(id)
	return res, err
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
func (store *Store) Put(id ID, view View) error {
	key := id.Value()
	if _, existing := store.retriever[key]; !existing {
		store.ids = append(store.ids, id)
	}
	data := make([][]byte, view.BlockCount())
	for index := 0; index < len(data); index++ {
		reader, err := view.Block(index)
		if err != nil {
			return err
		}
		data[index], err = ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
	}

	res := &Resource{
		Properties: Properties{
			Compound:    view.Compound(),
			ContentType: view.ContentType(),
			Compressed:  view.Compressed(),
		},
		Blocks: BlocksFrom(data),
	}
	store.retriever[key] = func() (*Resource, error) { return res, nil }
	return nil
}
