package archives

type viewModel struct {
	windowOpen   bool
	restoreFocus bool

	selectedLevel int

	variableContextIndex int
}

func freshViewModel() viewModel {
	return viewModel{}
}
