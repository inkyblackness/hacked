package archive

type viewModel struct {
	windowOpen   bool
	restoreFocus bool
}

func freshViewModel() viewModel {
	return viewModel{}
}
