package levels

type objectsViewModel struct {
	selectedObjects objectIDs

	restoreFocus bool
	windowOpen   bool
}

func freshObjectsViewModel() objectsViewModel {
	return objectsViewModel{}
}
