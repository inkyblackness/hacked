package project

type viewModel struct {
	restoreFocus          bool
	windowOpen            bool
	selectedManifestEntry int
}

func freshViewModel() viewModel {
	return viewModel{
		windowOpen:            false,
		selectedManifestEntry: -1,
	}
}
