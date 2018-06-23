package about

type viewModel struct {
	windowOpen bool
}

func freshViewModel() viewModel {
	return viewModel{}
}
