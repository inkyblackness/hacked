package project

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/inkyblackness/hacked/ss1/content/archive"
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/resource/lgres"
	"github.com/inkyblackness/hacked/ss1/serial"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

type fileStaging struct {
	resultMutex sync.Mutex

	failedFiles int
	savegames   map[world.FileLocation]resource.Viewer
	resources   map[world.FileLocation]resource.Viewer

	objectProperties  object.PropertiesTable
	textureProperties texture.PropertiesList
}

func newFileStaging() *fileStaging {
	return &fileStaging{
		resources: make(map[world.FileLocation]resource.Viewer),
		savegames: make(map[world.FileLocation]resource.Viewer),
	}
}

func (staging *fileStaging) stageAll(names []string) {
	staging.stageList(names, len(names) == 1)
}

func (staging *fileStaging) stageList(names []string, isOnlyStagedFile bool) {
	var wg sync.WaitGroup

	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			staging.stage(name, isOnlyStagedFile)
		}(name)
	}
	wg.Wait()
}

func (staging *fileStaging) stage(name string, isOnlyStagedFile bool) {
	fileInfo, err := os.Stat(name)
	if err != nil {
		staging.markFailedFile()
		return
	}
	file, err := os.Open(name)
	if err != nil {
		staging.markFailedFile()
		return
	}
	defer func() {
		_ = file.Close() // nolint: gas
	}()

	if fileInfo.IsDir() {
		if isOnlyStagedFile {
			subNames, _ := file.Readdirnames(0)
			joinedSubNames := make([]string, len(subNames))
			for index, subName := range subNames {
				joinedSubNames[index] = filepath.Join(name, subName)
			}
			staging.stageList(joinedSubNames, false)
		}
	} else {
		fileData, err := ioutil.ReadAll(file)
		if err != nil {
			staging.markFailedFile()
			return
		}

		reader, err := lgres.ReaderFrom(bytes.NewReader(fileData))
		filename := filepath.Base(name)
		if (err == nil) && (isOnlyStagedFile || fileAllowlist.Matches(filename)) {
			location := world.FileLocation{DirPath: filepath.Dir(name), Name: filename}
			staging.modify(func() {
				if stateView, stateErr := reader.View(ids.GameState); (stateErr == nil) && archive.IsSavegame(stateView) {
					staging.savegames[location] = reader
				} else {
					staging.resources[location] = reader
				}
			})
		}
		if strings.ToLower(filename) == world.ObjectPropertiesFilename {
			decoder := serial.NewDecoder(bytes.NewReader(fileData))
			properties := object.StandardPropertiesTable()
			properties.Code(decoder)
			err = decoder.FirstError()
			if err == nil {
				staging.modify(func() { staging.objectProperties = properties })
			}
		}
		if strings.ToLower(filename) == world.TexturePropertiesFilename && (len(fileData) > 4) {
			decoder := serial.NewDecoder(bytes.NewReader(fileData))
			entryCount := (len(fileData) - 4) / texture.PropertiesSize
			properties := make(texture.PropertiesList, entryCount)
			properties.Code(decoder)
			err = decoder.FirstError()
			if err == nil {
				staging.modify(func() { staging.textureProperties = properties })
			}
		}

		if err != nil {
			staging.markFailedFile()
		}
	}
}

func (staging *fileStaging) markFailedFile() {
	staging.modify(func() { staging.failedFiles++ })
}

func (staging *fileStaging) modify(modifier func()) {
	staging.resultMutex.Lock()
	defer staging.resultMutex.Unlock()
	modifier()
}
