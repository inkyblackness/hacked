package levels

type objectsViewModel struct {
	restoreFocus bool
	windowOpen   bool
}

func freshObjectsViewModel() objectsViewModel {
	return objectsViewModel{}
}
