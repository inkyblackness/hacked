package world

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial/rle"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// ModResetCallback is called when the mod was reset.
type ModResetCallback func()

// Mod is the central object for a game-mod.
//
// It is based on a "static" world and adds its own changes. The world data itself is not static, it is merely the
// unchangeable background for the mod. Changes to the mod are kept in a separate layer, which can be loaded and saved.
type Mod struct {
	worldManifest    *Manifest
	resourcesChanged resource.ModificationCallback
	resetCallback    ModResetCallback

	modPath        string
	lastChangeTime time.Time
	changedFiles   map[string]struct{}

	data ModData
}

// NewMod returns a new instance.
func NewMod(resourcesChanged resource.ModificationCallback, resetCallback ModResetCallback) *Mod {
	mod := &Mod{
		resourcesChanged: resourcesChanged,
		resetCallback:    resetCallback,
		changedFiles:     make(map[string]struct{}),
	}
	mod.worldManifest = NewManifest(mod.worldChanged)
	mod.data.FileChangeCallback = mod.markFileChanged

	return mod
}

// World returns the static background to the mod. Changes in the returned manifest may cause change callbacks
// being forwarded.
func (mod Mod) World() *Manifest {
	return mod.worldManifest
}

// Path returns the path where the mod is loaded from/saved to.
func (mod Mod) Path() string {
	return mod.modPath
}

// SetPath sets the path where the mod is loaded from/saved to.
func (mod *Mod) SetPath(p string) {
	mod.modPath = p
}

// ModifiedResources returns the current modification state.
func (mod Mod) ModifiedResources() []*LocalizedResources {
	return mod.data.LocalizedResources
}

// ModifiedFilenames returns the list of all filenames suspected of change.
func (mod Mod) ModifiedFilenames() []string {
	result := make([]string, 0, len(mod.changedFiles))
	for filename := range mod.changedFiles {
		result = append(result, filename)
	}
	return result
}

// AllAbsoluteFilenames returns the list of all filenames currently loaded in the mod.
func (mod Mod) AllAbsoluteFilenames() []string {
	var result []string
	for _, res := range mod.data.LocalizedResources {
		result = append(result, res.File.AbsolutePathFrom(mod.modPath))
	}
	if len(mod.data.ObjectProperties) > 0 {
		result = append(result, filepath.Join(mod.modPath, ObjectPropertiesFilename))
	}
	if len(mod.data.TextureProperties) > 0 {
		result = append(result, filepath.Join(mod.modPath, TexturePropertiesFilename))
	}
	sort.Strings(result)
	return result
}

// LastChangeTime returns the timestamp of the last change. Zero if not modified.
func (mod *Mod) LastChangeTime() time.Time {
	return mod.lastChangeTime
}

// ResetLastChangeTime clears the last change timestamp. It will be set again at the next modification.
func (mod *Mod) ResetLastChangeTime() {
	mod.lastChangeTime = time.Time{}
}

// MarkSave clears the list of modified filenames.
func (mod *Mod) MarkSave() {
	mod.changedFiles = make(map[string]struct{})
	mod.lastChangeTime = time.Time{}
}

// ModifiedResource retrieves the resource of given language and ID.
// There is no fallback lookup, it will return the exact resource stored under the provided identifier.
// Returns nil if the resource does not exist.
func (mod Mod) ModifiedResource(lang resource.Language, id resource.ID) resource.View {
	return mod.modifiedResource(lang, id)
}

func (mod Mod) modifiedResource(lang resource.Language, id resource.ID) *resource.Resource {
	for _, entry := range mod.data.LocalizedResources {
		if entry.Language == lang {
			res, err := entry.Store.Resource(id)
			if err == nil {
				return res
			}
		}
	}
	return nil
}

// CreateBlockPatch creates delta information for a block witch static data length.
// The returned patch structure contains details for both modifying the current state to be equal the new state,
// as well as the reversal delta. These deltas are calculated using the rle compression package.
// The returned boolean indicates whether the data differs. This can be used to detect whether the patch is necessary.
// An error is returned if the resource or the block do not exist, or if the length of newData does not match that of the block.
//
// If no error is returned, both the patch and the boolean provide valid information - even if the data is equal.
func (mod Mod) CreateBlockPatch(lang resource.Language, id resource.ID, index int, newData []byte) (BlockPatch, bool, error) {
	patch := BlockPatch{
		ID:          id,
		BlockIndex:  -1,
		BlockLength: 0,
	}
	res := mod.modifiedResource(lang, id)
	if res == nil {
		return patch, false, errors.New("resource unknown")
	}
	oldData, err := res.BlockRaw(index)
	if err != nil {
		return patch, false, err
	}
	if len(oldData) != len(newData) {
		return patch, false, fmt.Errorf("block length mismatch: current=%d, newData=%d", len(oldData), len(newData))
	}

	forwardData := bytes.NewBuffer(nil)
	err = rle.Compress(forwardData, newData, oldData)
	if err != nil {
		return patch, false, err
	}
	patch.ForwardData = forwardData.Bytes()

	reverseData := bytes.NewBuffer(nil)
	err = rle.Compress(reverseData, oldData, newData)
	if err != nil {
		return patch, false, err
	}
	patch.ReverseData = reverseData.Bytes()

	patch.BlockIndex = index
	patch.BlockLength = len(oldData)

	return patch, !bytes.Equal(oldData, newData), nil
}

// ModifiedBlock retrieves the specific block identified by given parameter.
// Returns empty slice if the block (or resource) is not modified.
func (mod Mod) ModifiedBlock(lang resource.Language, id resource.ID, index int) (data []byte) {
	res := mod.modifiedResource(lang, id)
	if res == nil {
		return
	}
	raw, _ := res.BlockRaw(index)
	return mod.blockCopy(raw)
}

// ModifiedBlocks returns all blocks of the modified resource.
func (mod Mod) ModifiedBlocks(lang resource.Language, id resource.ID) [][]byte {
	res := mod.modifiedResource(lang, id)
	if res == nil {
		return nil
	}
	data := make([][]byte, res.BlockCount())
	for index := 0; index < res.BlockCount(); index++ {
		raw, _ := res.BlockRaw(index)
		data[index] = mod.blockCopy(raw)
	}
	return data
}

func (mod Mod) blockCopy(data []byte) []byte {
	result := make([]byte, len(data))
	copy(result, data)
	return result
}

// Filter returns a list of resources that match the given parameters.
func (mod Mod) Filter(lang resource.Language, id resource.ID) resource.List {
	list := mod.worldManifest.Filter(lang, id)
	if res := mod.modifiedResource(resource.LangAny, id); res != nil {
		list = list.With(res)
	}
	for _, worldLang := range resource.Languages() {
		if worldLang.Includes(lang) {
			if res := mod.modifiedResource(lang, id); res != nil {
				list = list.With(res)
			}
		}
	}
	return list
}

// LocalizedResources returns a resource selector for a specific language.
func (mod Mod) LocalizedResources(lang resource.Language) resource.Selector {
	return resource.Selector{
		Lang: lang,
		From: mod,
		As:   ResourceViewStrategy(),
	}
}

// Modify requests to change the mod. The provided function will be called to collect all changes.
// After the modifier completes, all the requests will be applied and any changes notified.
func (mod *Mod) Modify(modifier func(Modder)) {
	var trans ModTransaction
	modifier(&trans)
	mod.modifyAndNotify(func() {
		for _, action := range trans.actions {
			action(&mod.data)
		}
	}, trans.modifiedIDs.ToList())
}

// ObjectProperties returns the table of object properties.
func (mod *Mod) ObjectProperties() object.PropertiesTable {
	if mod.HasModifiableObjectProperties() {
		return mod.data.ObjectProperties
	}
	return mod.worldManifest.ObjectProperties()
}

// HasModifiableObjectProperties returns true if the mod has dedicated object properties.
func (mod *Mod) HasModifiableObjectProperties() bool {
	return len(mod.data.ObjectProperties) > 0
}

// TextureProperties returns the list of texture properties.
func (mod *Mod) TextureProperties() texture.PropertiesList {
	if mod.HasModifiableTextureProperties() {
		return mod.data.TextureProperties
	}
	return mod.worldManifest.TextureProperties()
}

// HasModifiableTextureProperties returns true if the mod has dedicated texture properties.
func (mod *Mod) HasModifiableTextureProperties() bool {
	return len(mod.data.TextureProperties) > 0
}

func (mod *Mod) modifyAndNotify(modifier func(), modifiedIDs []resource.ID) {
	notifier := resource.ChangeNotifier{
		Callback:  mod.resourcesChanged,
		Localizer: mod,
	}
	notifier.ModifyAndNotify(modifier, modifiedIDs)
}

func (mod Mod) worldChanged(modifiedIDs []resource.ID, failedIDs []resource.ID) {
	// It would be great to also check whether the mod hides any of these changes.
	// Sadly, this is not possible:
	// a) At the point of this callback, we can't do a check on the previous state anymore.
	// b) Even when changing the world only within a modification enclosure of our own notifier, we can't determine
	//    the list of changed IDs before actually changing them. (Specifying ALL IDs is not a good idea due to performance.)
	// As a result, simply forward this list. I don't even expect any big performance gain through such a filter.
	// This would only be relevant to "full conversion" mods AND a change in a big list in the world. Hardly the case.
	mod.resourcesChanged(modifiedIDs, failedIDs)
}

// Reset changes the mod to a new set of resources.
func (mod *Mod) Reset(newResources []*LocalizedResources, objectProperties object.PropertiesTable, textureProperties texture.PropertiesList) {
	var modifiedIDs resource.IDMarkerMap
	collectIDs := func(res []*LocalizedResources) {
		for _, loc := range res {
			for _, id := range loc.Store.IDs() {
				modifiedIDs.Add(id)
			}
		}
	}
	collectIDs(mod.data.LocalizedResources)
	collectIDs(newResources)

	mod.data.LocalizedResources = newResources
	mod.data.ObjectProperties = objectProperties
	mod.data.TextureProperties = textureProperties
	mod.changedFiles = make(map[string]struct{})
	mod.lastChangeTime = time.Time{}
	mod.resetCallback()
	mod.resourcesChanged(modifiedIDs.ToList(), nil)
}

func (mod *Mod) markFileChanged(filename string) {
	mod.changedFiles[filename] = struct{}{}
	mod.lastChangeTime = time.Now()
}

// FixListResources ensures all resources that contain resource lists to
// have maximum size. This is done to ensure compatibility with layered modding in the
// Source Port branch of engines.
// These engines will have "lower" mods bleed through only if the block is empty in a
// "higher" mod. This even counts for entries past the last modified one in the higher mod.
func (mod *Mod) FixListResources() {
	mod.Modify(func(modder Modder) {
		for _, localized := range mod.data.LocalizedResources {
			for _, id := range localized.Store.IDs() {
				res, _ := localized.Store.Resource(id)
				info, known := ids.Info(id)
				if known && info.List {
					baseCount := res.BlockCount()
					required := info.MaxCount - baseCount
					for i := 0; i < required; i++ {
						modder.SetResourceBlock(localized.Language, id, baseCount+i, nil)
					}
				}
			}
		}
	})
}
