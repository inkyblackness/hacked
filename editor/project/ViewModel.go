package project

type viewModel struct {
	restoreFocus          bool
	selectedManifestEntry int
}

func freshViewModel() viewModel {
	return viewModel{
		selectedManifestEntry: -1,
	}
}
