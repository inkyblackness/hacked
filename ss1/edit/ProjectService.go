package edit

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"github.com/inkyblackness/hacked/ss1"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/resource/lgres"
	"github.com/inkyblackness/hacked/ss1/serial"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

const (
	errNoResourcesFound     ss1.StringError = "no resources found"
	errNoStorageLocationSet ss1.StringError = "no storage location set"
)

// ProjectSettings describe the properties of a project.
type ProjectSettings struct {
	ModFiles []string
	Manifest []ManifestEntrySettings
}

// ManifestEntrySettings describe the properties of one manifest entry in a project.
type ManifestEntrySettings struct {
	Origin []string
}

const autosaveTimeoutSec = 5

// SaveStatus describes the current change state.
type SaveStatus struct {
	mod *world.Mod
	// FilesModified is the count of how many files are affected by the current pending change.
	FilesModified int
	// SavePending is set if the mod has a storage location and changes are to be saved.
	SavePending bool
	// SaveIn is the duration after which an auto-save should happen.
	SaveIn time.Duration
}

// ConfirmPendingSave marks that the recent auto-save status has been acknowledged and should no longer be notified.
// A new change will re-start the auto-save timer.
func (status SaveStatus) ConfirmPendingSave() {
	if status.mod == nil {
		return
	}
	status.mod.ResetLastChangeTime()
}

// ProjectService handles the overall information about the active project.
type ProjectService struct {
	commander cmd.Registry

	mod     *world.Mod
	modPath string

	stateFilename string
}

// NewProjectService returns a new instance of a service for given mod.
func NewProjectService(commander cmd.Registry, mod *world.Mod) *ProjectService {
	return &ProjectService{
		commander: commander,
		mod:       mod,
	}
}

// SaveStatus returns the current pending save information.
func (service ProjectService) SaveStatus() SaveStatus {
	var status SaveStatus
	status.mod = service.mod
	status.FilesModified = len(service.mod.ModifiedFilenames())
	if status.FilesModified > 0 {
		lastChangeTime := service.mod.LastChangeTime()

		if service.ModHasStorageLocation() && !lastChangeTime.IsZero() {
			status.SavePending = true
			saveAt := lastChangeTime.Add(time.Duration(autosaveTimeoutSec) * time.Second)
			status.SaveIn = time.Until(saveAt)
			if status.SaveIn <= 0 {
				status.SaveIn = 0
			}
		}
	}
	return status
}

// CurrentStateFilename returns the name for the state of the project.
func (service ProjectService) CurrentStateFilename() string {
	return service.stateFilename
}

// SetCurrentStateFilename updates the current filename.
func (service *ProjectService) SetCurrentStateFilename(value string) {
	service.stateFilename = value
}

// CurrentSettings returns the snapshot of the project.
func (service ProjectService) CurrentSettings() ProjectSettings {
	manifest := service.mod.World()
	settings := ProjectSettings{Manifest: make([]ManifestEntrySettings, manifest.EntryCount())}
	for i := 0; i < manifest.EntryCount(); i++ {
		entry, _ := manifest.Entry(i)
		settings.Manifest[i].Origin = entry.Origin
	}

	settings.ModFiles = service.relativeToSettings(service.mod.AllAbsoluteFilenames(service.modPath)...)

	return settings
}

func (service *ProjectService) relativeToSettings(filenames ...string) []string {
	loc := world.FileLocationFrom(service.stateFilename)
	relatives := make([]string, 0, len(filenames))
	for _, filename := range filenames {
		relatives = append(relatives, world.FileLocationFrom(filename).NestedRelativeTo(loc.DirPath))
	}
	return relatives
}

func (service *ProjectService) absoluteFromSettings(filenames ...string) []string {
	loc := world.FileLocationFrom(service.stateFilename)
	absolutes := make([]string, 0, len(filenames))
	for _, filename := range filenames {
		absolutes = append(absolutes, world.FileLocationFrom(filename).AbsolutePathFrom(loc.DirPath))
	}
	return absolutes
}

// RestoreProject sets internal data based on the given settings.
func (service *ProjectService) RestoreProject(settings ProjectSettings, stateFilename string) {
	service.ResetProject()

	service.stateFilename = stateFilename

	manifest := service.mod.World()
	for _, entrySettings := range settings.Manifest {
		entry, err := world.NewManifestEntryFrom(entrySettings.Origin)
		if err != nil {
			continue
		}
		err = manifest.InsertEntry(manifest.EntryCount(), entry)
		if err != nil {
			continue
		}
	}

	_ = service.TryLoadModFrom(service.absoluteFromSettings(settings.ModFiles...))
}

// ResetProject clears the project and returns it to initial state.
func (service *ProjectService) ResetProject() {
	service.setActiveMod("", nil, nil, nil)
	service.mod.World().Reset()
	service.stateFilename = ""
}

// AddManifestEntry attempts to insert the given manifest entry at given index.
func (service *ProjectService) AddManifestEntry(at int, entry *world.ManifestEntry) error {
	return service.commander.Register(
		cmd.Named("AddManifestEntry"),
		cmd.Forward(func(modder world.Modder) error {
			return service.mod.World().InsertEntry(at, entry)
		}),
		cmd.Reverse(func(modder world.Modder) error {
			return service.mod.World().RemoveEntry(at)
		}),
	)
}

// RemoveManifestEntry attempts to remove the manifest entry at given index.
func (service *ProjectService) RemoveManifestEntry(at int) error {
	manifest := service.mod.World()
	entry, err := manifest.Entry(at)
	if err != nil {
		return err
	}
	return service.commander.Register(
		cmd.Named("RemoveManifestEntry"),
		cmd.Forward(func(modder world.Modder) error {
			return manifest.RemoveEntry(at)
		}),
		cmd.Reverse(func(modder world.Modder) error {
			return manifest.InsertEntry(at, entry)
		}),
	)
}

// MoveManifestEntry attempts to remove the manifest entry at given from index and re-insert it at given to index.
func (service *ProjectService) MoveManifestEntry(to, from int) error {
	return service.commander.Register(
		cmd.Named("MoveManifestEntry"),
		cmd.Forward(func(modder world.Modder) error {
			return service.mod.World().MoveEntry(to, from)
		}),
		cmd.Reverse(func(modder world.Modder) error {
			return service.mod.World().MoveEntry(from, to)
		}),
	)
}

// Mod returns the currently active mod in the project.
func (service ProjectService) Mod() *world.Mod {
	return service.mod
}

// NewMod resets the mod to a new state.
func (service *ProjectService) NewMod() {
	service.setActiveMod("", nil, nil, nil)
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
		return errNoResourcesFound
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
			if err != nil {
				// optimistic ignore
				continue
			}
			_ = loc.Store.Put(id, view)
		}
		locs = append(locs, loc)
	}

	service.setActiveMod(modPath, locs, loaded.ObjectProperties, loaded.TextureProperties)
	return nil
}

func (service *ProjectService) setActiveMod(modPath string, resources []*world.LocalizedResources,
	objectProperties object.PropertiesTable, textureProperties texture.PropertiesList) {
	service.setModPath(modPath)
	service.mod.Reset(resources, objectProperties, textureProperties)
	// fix list resources for any "old" mod.
	service.mod.FixListResources()
}

// ModifyModWith runs a function with the intent to alter the current mod.
func (service *ProjectService) ModifyModWith(modifier func(world.Modder) error) (err error) {
	service.mod.Modify(func(modder world.Modder) {
		err = modifier(modder)
	})
	return
}

// ModHasStorageLocation returns whether the mod has a place to be stored.
func (service ProjectService) ModHasStorageLocation() bool {
	return len(service.modPath) > 0
}

// ModPath returns the base path of the mod in the project.
func (service ProjectService) ModPath() string {
	return service.modPath
}

func (service *ProjectService) setModPath(value string) {
	service.modPath = value
}

// SaveMod will store the currently active mod in its current path.
func (service *ProjectService) SaveMod() error {
	if !service.ModHasStorageLocation() {
		return errNoStorageLocationSet
	}
	return service.SaveModUnder(service.modPath)
}

// SaveModUnder will store the currently active mod in the given path.
func (service *ProjectService) SaveModUnder(modPath string) error {
	service.mod.FixListResources()
	err := service.saveModResourcesTo(modPath)
	if err != nil {
		return err
	}
	service.setModPath(modPath)
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
