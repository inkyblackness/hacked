package edit

import (
	"fmt"

	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

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

// TryLoadModFrom attempts to set the active mod from given filenames.
func (service *ProjectService) TryLoadModFrom(names []string) error {
	loaded := world.LoadFiles(false, names)

	resourcesToTake := loaded.Resources
	isSavegame := false
	if (len(resourcesToTake) == 0) && (len(loaded.Savegames) == 1) {
		resourcesToTake = loaded.Savegames
		isSavegame = true
	}
	if len(resourcesToTake) == 0 {
		return fmt.Errorf("no resources found")
	}
	var locs []*world.LocalizedResources
	modPath := ""

	for location := range resourcesToTake {
		if (len(modPath) == 0) || (len(location.DirPath) < len(modPath)) {
			modPath = location.DirPath
		}
	}

	for location, viewer := range resourcesToTake {
		lang := ids.LocalizeFilename(location.Name)
		template := location.Name
		if isSavegame {
			template = string(ids.Archive)
		}
		loc := &world.LocalizedResources{
			File:     location,
			Template: template,
			Language: lang,
		}
		for _, id := range viewer.IDs() {
			view, err := viewer.View(id)
			if err == nil {
				_ = loc.Store.Put(id, view)
			}
			// TODO: handle error?
		}
		locs = append(locs, loc)
	}

	service.setActiveMod(modPath, locs, loaded.ObjectProperties, loaded.TextureProperties)
	return nil
}

func (service *ProjectService) setActiveMod(modPath string, resources []*world.LocalizedResources,
	objectProperties object.PropertiesTable, textureProperties texture.PropertiesList) {
	service.mod.SetPath(modPath)
	service.mod.Reset(resources, objectProperties, textureProperties)
	// fix list resources for any "old" mod.
	service.mod.FixListResources()
}
