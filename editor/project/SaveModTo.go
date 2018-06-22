package project

import (
	"os"
	"path/filepath"

	"github.com/inkyblackness/hacked/editor/model"
	"github.com/inkyblackness/hacked/ss1/resource/lgres"
)

func saveModResourcesTo(localized model.LocalizedResources, modPath string) error {
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

	for filename, list := range resByFile {
		err := saveResourcesTo(list, filepath.Join(modPath, filename))
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
		_ = file.Close()
	}()
	err = lgres.Write(file, list)
	return err
}
