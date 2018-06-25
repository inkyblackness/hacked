package archives

type viewModel struct {
	windowOpen   bool
	restoreFocus bool

	selectedLevel int
}

func freshViewModel() viewModel {
	return viewModel{}
}
