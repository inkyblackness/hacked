package project

import (
	"os"
	"path/filepath"

	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/ss1/resource/lgres"
)

func saveModResourcesTo(localized model.LocalizedResources, modPath string, filenamesToSave []string) error {
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
