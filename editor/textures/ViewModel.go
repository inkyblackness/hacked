package textures

type viewModel struct {
	windowOpen   bool
	restoreFocus bool

	currentIndex int
}

func freshViewModel() viewModel {
	return viewModel{
		currentIndex: 0,
	}
}
