package project

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/resource/lgres"
	"github.com/inkyblackness/hacked/ss1/serial"
	"github.com/inkyblackness/hacked/ss1/world"
)

func saveModResourcesTo(mod *world.Mod, modPath string) error {
	localized := mod.ModifiedResources()
	filenamesToSave := mod.ModifiedFilenames()

	shallBeSaved := func(filename string) bool {
		for _, toSave := range filenamesToSave {
			if toSave == filename {
				return true
			}
		}
		return false
	}

	for _, loc := range localized {
		if shallBeSaved(loc.Filename) {
			err := saveResourcesTo(loc.Store, filepath.Join(modPath, loc.Filename))
			if err != nil {
				return err
			}
		}
	}

	if shallBeSaved(world.TexturePropertiesFilename) {
		err := saveTexturePropertiesTo(mod.TextureProperties(), filepath.Join(modPath, world.TexturePropertiesFilename))
		if err != nil {
			return err
		}
	}
	if shallBeSaved(world.ObjectPropertiesFilename) {
		err := saveObjectPropertiesTo(mod.ObjectProperties(), filepath.Join(modPath, world.ObjectPropertiesFilename))
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
