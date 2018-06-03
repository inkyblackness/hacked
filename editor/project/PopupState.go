package project

type popupState interface {
	Render()
	HandleFiles(names []string)
}

type idlePopupState struct{}

func (state idlePopupState) Render() {
}

func (state idlePopupState) HandleFiles(names []string) {

}
