package levels

import "github.com/inkyblackness/hacked/ss1/content/object"

type objectsViewModel struct {
	selectedObjects objectIDs

	newObjectTriple object.Triple

	restoreFocus bool
	windowOpen   bool
}

func freshObjectsViewModel() objectsViewModel {
	return objectsViewModel{}
}
