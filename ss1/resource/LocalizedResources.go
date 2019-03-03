package resource

// LocalizedResources associates a language with a resource viewer under a specific identifier.
type LocalizedResources struct {
	// ID is the identifier of the viewer. This could be a filename for instance.
	ID string
	// Language specifies for which language the viewer has resources.
	Language Language
	// Viewer is the actual container of the resources.
	Viewer Viewer
}
