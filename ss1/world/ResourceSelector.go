package world

import "github.com/inkyblackness/hacked/ss1/resource"

type resourceFilter interface {
	filter(lang Language, id resource.ID) resource.List
}

// ResourceSelector provides a merged view of from according to a language.
type ResourceSelector struct {
	lang Language
	from resourceFilter

	// As defines how the found resources should be viewed in case more than one matches.
	// By default, the last resource will be used.
	As ResourceViewStrategy
}

// Select provides a collected view on one resource.
func (merger ResourceSelector) Select(id resource.ID) (view ResourceView, err error) {
	list := merger.from.filter(merger.lang, id)
	if len(list) == 0 {
		return nil, resource.ErrResourceDoesNotExist(id)
	}
	if (merger.As == nil) || !merger.As.IsCompoundList(id) {
		view = &resourceViewer{list[len(list)-1]}
	} else {
		view = &resourceListMerger{list}
	}

	return view, nil
}
