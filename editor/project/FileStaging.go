package project

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/resource/lgres"
	"github.com/inkyblackness/hacked/ss1/serial"
	"github.com/inkyblackness/hacked/ss1/world"
)

type fileStaging struct {
	failedFiles int
	savegames   map[string]resource.Provider
	resources   map[string]resource.Provider

	objectProperties object.PropertiesTable
}

func (staging *fileStaging) stage(name string, isOnlyStagedFile bool) {
	fileInfo, err := os.Stat(name)
	if err != nil {
		staging.failedFiles++
		return
	}
	file, err := os.Open(name)
	if err != nil {
		return
	}
	defer file.Close() // nolint: errcheck

	if fileInfo.IsDir() {
		if isOnlyStagedFile {
			subNames, _ := file.Readdirnames(0)
			for _, subName := range subNames {
				staging.stage(filepath.Join(name, subName), false)
			}
		}
	} else {
		fileData, err := ioutil.ReadAll(file)
		if err != nil {
			staging.failedFiles++
		}

		reader, err := lgres.ReaderFrom(bytes.NewReader(fileData))
		filename := filepath.Base(name)
		if (err == nil) && (isOnlyStagedFile || fileWhitelist.Matches(filename)) {
			if world.IsSavegame(reader) {
				staging.savegames[filename] = reader
			} else {
				staging.resources[filename] = reader
			}
		}
		if strings.ToLower(filename) == "objprop.dat" {
			decoder := serial.NewDecoder(bytes.NewReader(fileData))
			properties := object.StandardPropertiesTable()
			properties.Code(decoder)
			err = decoder.FirstError()
			if err == nil {
				staging.objectProperties = properties
			}
		}

		if err != nil {
			staging.failedFiles++
		}
	}
}
