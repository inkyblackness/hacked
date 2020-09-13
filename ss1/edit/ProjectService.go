package edit

import "github.com/inkyblackness/hacked/ss1/world"

// ProjectService handles the overall information about the active mod.
type ProjectService struct {
	mod *world.Mod
}

// NewProjectService returns a new instance of a service for given mod.
func NewProjectService(mod *world.Mod) *ProjectService {
	return &ProjectService{mod: mod}
}

// Mod returns the currently active mod in the project.
func (service ProjectService) Mod() *world.Mod {
	return service.mod
}
