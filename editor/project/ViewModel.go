package project

type viewModel struct {
	selectedManifestEntry int
}

func freshViewModel() viewModel {
	return viewModel{
		selectedManifestEntry: -1,
	}
}
