package resource

// Filter filters for language and id to produce a list of matching resources.
type Filter interface {
	Filter(lang Language, id ID) List
}

// Selector provides a merged view of resources according to a language.
type Selector struct {
	// Lang specifies the language to filter by.
	Lang Language

	// From specifies from where the resources shall be taken.
	From Filter

	// As defines how the found resources should be viewed in case more than one matches.
	// By default, the last resource will be used.
	As ViewStrategy
}

// Select provides a collected view on one resource.
func (merger Selector) Select(id ID) (view View, err error) {
	list := merger.From.Filter(merger.Lang, id)
	if len(list) == 0 {
		return nil, ErrResourceDoesNotExist(id)
	}
	if (merger.As == nil) || !merger.As.IsCompoundList(id) {
		view = list[len(list)-1]
	} else {
		view = &listMerger{list}
	}

	return view, nil
}
