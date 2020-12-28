package world

import (
	"archive/zip"
	"bytes"
	"io"
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
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

var fileAllowlist = ids.FilenameList{
	ids.Archive,
	ids.CybStrng,
	ids.CitALog,
	ids.CitBark,
	ids.CitMat,
	ids.DigiFX,
	ids.GamePal,
	ids.GameScr,
	ids.MfdArt,
	ids.Obj3D,
	ids.ObjArt,
	ids.ObjArt2,
	ids.ObjArt3,
	ids.SideArt,
	ids.Splash,
	ids.SvgaDeth,
	ids.SvgaEnd,
	ids.SvgaIntr,
	ids.Texture,
	ids.VidMail,
}

type fileLoader struct {
	resultMutex sync.Mutex

	allowZips bool

	result FileLoadResult
}

// FileLoadResult contains all the information of a LoadFiles attempt.
type FileLoadResult struct {
	FailedFiles int
	Savegames   map[FileLocation]resource.Viewer
	Resources   map[FileLocation]resource.Viewer

	ObjectProperties  object.PropertiesTable
	TextureProperties texture.PropertiesList
}

// LoadFiles attempts to load compatible files from the given set of filenames.
func LoadFiles(allowZips bool, names []string) FileLoadResult {
	loader := fileLoader{
		allowZips: allowZips,
		result: FileLoadResult{
			Resources: make(map[FileLocation]resource.Viewer),
			Savegames: make(map[FileLocation]resource.Viewer),
		},
	}
	loader.loadAll(names)
	return loader.result
}

func (loader *fileLoader) loadAll(names []string) {
	loader.loadList(names, len(names) == 1)
}

func (loader *fileLoader) loadList(names []string, isOnlyRequestedFile bool) {
	var wg sync.WaitGroup

	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			loader.load(name, isOnlyRequestedFile)
		}(name)
	}
	wg.Wait()
}

func (loader *fileLoader) load(name string, isOnlyRequestedFile bool) {
	fileInfo, err := os.Stat(name)
	if err != nil {
		loader.markFailedFile()
		return
	}
	file, err := os.Open(name)
	if err != nil {
		loader.markFailedFile()
		return
	}
	defer func() {
		_ = file.Close() // nolint: gas
	}()

	if fileInfo.IsDir() {
		if isOnlyRequestedFile {
			subNames, _ := file.Readdirnames(0)
			joinedSubNames := make([]string, len(subNames))
			for index, subName := range subNames {
				joinedSubNames[index] = filepath.Join(name, subName)
			}
			loader.loadList(joinedSubNames, false)
		}
	} else {
		tryDirect := true
		if loader.allowZips && isOnlyRequestedFile {
			tryDirect = false
			zipReader, err := zip.NewReader(file, fileInfo.Size())
			if err != nil {
				_, _ = file.Seek(0, io.SeekStart)
				tryDirect = true
			} else {
				loader.loadFileArchive(zipReader)
			}
		}

		if tryDirect {
			loader.loadFile(name, isOnlyRequestedFile, file)
		}
	}
}

func (loader *fileLoader) loadFileArchive(zipReader *zip.Reader) {
	var wg sync.WaitGroup

	for _, file := range zipReader.File {
		wg.Add(1)
		go func(entry *zip.File) {
			defer wg.Done()
			entryFile, err := entry.Open()
			if err != nil {
				return
			}
			defer func() {
				_ = entryFile.Close()
			}()

			loader.loadFile(entry.Name, false, entryFile)
		}(file)
	}
	wg.Wait()
}

func (loader *fileLoader) loadFile(name string, isOnlyStagedFile bool, file io.Reader) {
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		loader.markFailedFile()
	}

	reader, err := lgres.ReaderFrom(bytes.NewReader(fileData))
	filename := filepath.Base(name)
	if (err == nil) && (isOnlyStagedFile || fileAllowlist.Matches(filename)) {
		location := FileLocation{DirPath: filepath.Dir(name), Name: filename}
		loader.modify(func() {
			if stateView, stateErr := reader.View(ids.GameState); (stateErr == nil) && archive.IsSavegame(stateView) {
				loader.result.Savegames[location] = reader
			} else {
				loader.result.Resources[location] = reader
			}
		})
	}
	if strings.ToLower(filename) == ObjectPropertiesFilename {
		decoder := serial.NewDecoder(bytes.NewReader(fileData))
		properties := object.StandardPropertiesTable()
		properties.Code(decoder)
		err = decoder.FirstError()
		if err == nil {
			loader.modify(func() { loader.result.ObjectProperties = properties })
		}
	}
	if strings.ToLower(filename) == TexturePropertiesFilename && (len(fileData) > 4) {
		decoder := serial.NewDecoder(bytes.NewReader(fileData))
		entryCount := (len(fileData) - 4) / texture.PropertiesSize
		properties := make(texture.PropertiesList, entryCount)
		properties.Code(decoder)
		err = decoder.FirstError()
		if err == nil {
			loader.modify(func() { loader.result.TextureProperties = properties })
		}
	}

	if err != nil {
		loader.markFailedFile()
	}
}

func (loader *fileLoader) markFailedFile() {
	loader.modify(func() { loader.result.FailedFiles++ })
}

func (loader *fileLoader) modify(modifier func()) {
	loader.resultMutex.Lock()
	defer loader.resultMutex.Unlock()
	modifier()
}
