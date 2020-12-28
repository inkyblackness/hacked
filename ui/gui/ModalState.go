package gui

// ModalState represents a modal dialog.
type ModalState interface {
	// Render renders the dialog.
	Render()
}

// FileHandler is accepting filenames.
type FileHandler interface {
	// HandleFiles is called for any dropped files.
	HandleFiles(names []string)
}
