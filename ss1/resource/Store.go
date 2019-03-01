package resource

import "io/ioutil"

// Store holds a set of resources. This set can be modified.
type Store struct {
	ids       []ID
	resources map[ID]*Resource
}

// IDs returns a list of available IDs this store currently contains.
func (store Store) IDs() []ID {
	ids := make([]ID, len(store.resources))
	for id := range store.resources {
		ids[store.findIDIndex(id)] = id
	}
	return ids
}

// Resource returns reference to the contained resource for given identifier.
func (store Store) Resource(id ID) (*Resource, error) {
	res, existing := store.resources[id]
	if !existing {
		return nil, ErrResourceDoesNotExist(id)
	}
	return res, nil
}

// View returns a read-only view on the resource for given identifier.
func (store Store) View(id ID) (View, error) {
	return store.Resource(id)
}

// Del removes the resource with given identifier from the store.
func (store *Store) Del(id ID) {
	if _, existing := store.resources[id]; existing {
		delete(store.resources, id)
	}
}

// Put (re-)assigns an identifier with data. If no resource with given ID exists,
// then it is created. Existing resources are overwritten with the provided data.
func (store *Store) Put(id ID, view View) error {
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
	if store.resources == nil {
		store.resources = make(map[ID]*Resource)
	}
	store.resources[id] = res
	if store.findIDIndex(id) < 0 {
		store.ids = append(store.ids, id)
	}
	return nil
}

func (store Store) findIDIndex(id ID) int {
	for index := 0; index < len(store.ids); index++ {
		if store.ids[index] == id {
			return index
		}
	}
	return -1
}
