package project

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/inkyblackness/hacked/ss1/resource/lgres"
	"github.com/inkyblackness/hacked/ss1/serial"
	"github.com/inkyblackness/hacked/ss1/world"
)

func saveModResourcesTo(mod *model.Mod, modPath string) error {
	localized := mod.ModifiedResources()
	filenamesToSave := mod.ModifiedFilenames()
	resByFile := make(map[string]model.IdentifiedResources)
	for _, identifiedIn := range localized {
		for id, res := range identifiedIn {
			identifiedOut, exist := resByFile[res.Filename()]
			if !exist {
				identifiedOut = make(model.IdentifiedResources)
				resByFile[res.Filename()] = identifiedOut
			}
			identifiedOut[id] = res
		}
	}

	shallBeSaved := func(filename string) bool {
		for _, toSave := range filenamesToSave {
			if toSave == filename {
				return true
			}
		}
		return false
	}

	for filename, list := range resByFile {
		if shallBeSaved(filename) {
			err := saveResourcesTo(list, filepath.Join(modPath, filename))
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

func saveResourcesTo(list model.IdentifiedResources, absFilename string) error {
	file, err := os.Create(absFilename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close() // nolint: gas
	}()
	err = lgres.Write(file, list)
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
