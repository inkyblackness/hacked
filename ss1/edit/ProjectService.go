package edit

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/resource/lgres"
	"github.com/inkyblackness/hacked/ss1/serial"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// ProjectService handles the overall information about the active mod.
type ProjectService struct {
	commander cmd.Registry
	mod       *world.Mod
}

// NewProjectService returns a new instance of a service for given mod.
func NewProjectService(commander cmd.Registry, mod *world.Mod) *ProjectService {
	return &ProjectService{
		commander: commander,
		mod:       mod,
	}
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

// ModHasStorageLocation returns whether the mod has a place to be stored.
func (service ProjectService) ModHasStorageLocation() bool {
	return len(service.mod.Path()) > 0
}

// SaveMod will store the currently active mod in its current path.
func (service *ProjectService) SaveMod() error {
	if !service.ModHasStorageLocation() {
		return fmt.Errorf("no storage location set")
	}
	return service.SaveModUnder(service.mod.Path())
}

// SaveModUnder will store the currently active mod in the given path.
func (service *ProjectService) SaveModUnder(modPath string) error {
	service.mod.FixListResources()
	err := service.saveModResourcesTo(modPath)
	if err != nil {
		return err
	}
	service.mod.SetPath(modPath)
	service.mod.MarkSave()
	return nil
}

func (service *ProjectService) saveModResourcesTo(modPath string) error {
	localized := service.mod.ModifiedResources()
	filenamesToSave := service.mod.ModifiedFilenames()

	shallBeSaved := func(filename string) bool {
		for _, toSave := range filenamesToSave {
			if toSave == filename {
				return true
			}
		}
		return false
	}

	for _, loc := range localized {
		if shallBeSaved(loc.File.Name) {
			err := saveResourcesTo(loc.Store, loc.File.AbsolutePathFrom(modPath))
			if err != nil {
				return err
			}
		}
	}

	if shallBeSaved(world.TexturePropertiesFilename) {
		err := saveTexturePropertiesTo(service.mod.TextureProperties(), filepath.Join(modPath, world.TexturePropertiesFilename))
		if err != nil {
			return err
		}
	}
	if shallBeSaved(world.ObjectPropertiesFilename) {
		err := saveObjectPropertiesTo(service.mod.ObjectProperties(), filepath.Join(modPath, world.ObjectPropertiesFilename))
		if err != nil {
			return err
		}
	}

	return nil
}

func saveResourcesTo(viewer resource.Viewer, absFilename string) error {
	file, err := os.Create(absFilename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close() // nolint: gas
	}()
	err = lgres.Write(file, viewer)
	return err
}

func saveTexturePropertiesTo(list texture.PropertiesList, absFilename string) error {
	return saveCodableTo(list, absFilename)
}

func saveObjectPropertiesTo(list object.PropertiesTable, absFilename string) error {
	return saveCodableTo(list, absFilename)
}

func saveCodableTo(codable serial.Codable, absFilename string) error {
	buffer := bytes.NewBuffer(nil)
	encoder := serial.NewEncoder(buffer)
	codable.Code(encoder)
	err := encoder.FirstError()
	if err != nil {
		return err
	}

	file, err := os.Create(absFilename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close() // nolint: gas
	}()
	_, err = file.Write(buffer.Bytes())
	return err
}
