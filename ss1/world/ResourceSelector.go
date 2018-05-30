package world

import "github.com/inkyblackness/hacked/ss1/resource"

// ResourceFilter filters for language and id to produce a list of matching resources.
type ResourceFilter interface {
	Filter(lang Language, id resource.ID) resource.List
}

// ResourceSelector provides a merged view of resources according to a language.
type ResourceSelector struct {
	// Lang specifies the language to filter by.
	Lang Language

	// From specifies from where the resources shall be taken.
	From ResourceFilter

	// As defines how the found resources should be viewed in case more than one matches.
	// By default, the last resource will be used.
	As ResourceViewStrategy
}

// Select provides a collected view on one resource.
func (merger ResourceSelector) Select(id resource.ID) (view ResourceView, err error) {
	list := merger.From.Filter(merger.Lang, id)
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
