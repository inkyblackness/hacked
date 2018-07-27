package model

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial/rle"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// ModResetCallback is called when the mod was reset.
type ModResetCallback func()

// Mod is the central object for a game-mod.
//
// It is based on a "static" world and adds its own changes. The world data itself is not static, it is merely the
// unchangeable background for the mod. Changes to the mod are kept in a separate layer, which can be loaded and saved.
type Mod struct {
	worldManifest    *world.Manifest
	resourcesChanged resource.ModificationCallback
	resetCallback    ModResetCallback

	modPath            string
	lastChangeTime     time.Time
	changedFiles       map[string]struct{}
	localizedResources LocalizedResources
	objectProperties   object.PropertiesTable
}

// NewMod returns a new instance.
func NewMod(resourcesChanged resource.ModificationCallback, resetCallback ModResetCallback) *Mod {
	mod := &Mod{
		resourcesChanged:   resourcesChanged,
		resetCallback:      resetCallback,
		localizedResources: NewLocalizedResources(),
		changedFiles:       make(map[string]struct{}),
	}
	mod.worldManifest = world.NewManifest(mod.worldChanged)

	return mod
}

// World returns the static background to the mod. Changes in the returned manifest may cause change callbacks
// being forwarded.
func (mod Mod) World() *world.Manifest {
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
func (mod Mod) ModifiedResources() LocalizedResources {
	return mod.localizedResources
}

// ModifiedFilenames returns the list of all filenames suspected of change.
func (mod Mod) ModifiedFilenames() []string {
	var result []string
	for filename := range mod.changedFiles {
		result = append(result, filename)
	}
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
	return mod.localizedResources[lang][id]
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
	res := mod.localizedResources[lang][id]
	if res == nil {
		return patch, false, errors.New("resource unknown")
	}
	if (index < 0) || (index >= res.BlockCount()) {
		return patch, false, errors.New("block index wrong")
	}
	oldData := res.blocks[index]
	if len(oldData) != len(newData) {
		return patch, false, fmt.Errorf("block length mismatch: current=%d, newData=%d", len(oldData), len(newData))
	}

	forwardData := bytes.NewBuffer(nil)
	err := rle.Compress(forwardData, newData, oldData)
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
	res := mod.localizedResources[lang][id]
	if res == nil {
		return
	}
	return mod.blockCopy(res.blocks[index])
}

// ModifiedBlocks returns all blocks of the modified resource.
func (mod Mod) ModifiedBlocks(lang resource.Language, id resource.ID) [][]byte {
	res := mod.localizedResources[lang][id]
	if res == nil {
		return nil
	}
	data := make([][]byte, res.blockCount)
	for index := 0; index < res.blockCount; index++ {
		data[index] = mod.blockCopy(res.blocks[index])
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
	if res, resExists := mod.localizedResources[resource.LangAny][id]; resExists {
		list = list.With(res)
	}
	for _, worldLang := range resource.Languages() {
		if worldLang.Includes(lang) {
			if res, resExists := mod.localizedResources[lang][id]; resExists {
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
		As:   world.ResourceViewStrategy(),
	}
}

// Modify requests to change the mod. The provided function will be called to collect all changes.
// After the modifier completes, all the requests will be applied and any changes notified.
func (mod *Mod) Modify(modifier func(*ModTransaction)) {
	var trans ModTransaction
	trans.modifiedIDs = make(resource.IDMarkerMap)
	modifier(&trans)
	mod.modifyAndNotify(func() {
		for _, action := range trans.actions {
			action(mod)
		}
	}, trans.modifiedIDs.ToList())
}

// ObjectProperties returns the table of object properties.
func (mod *Mod) ObjectProperties() object.PropertiesTable {
	if len(mod.objectProperties) > 0 {
		return mod.objectProperties
	}
	return mod.worldManifest.ObjectProperties()
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

func (mod *Mod) ensureResource(lang resource.Language, id resource.ID) *MutableResource {
	res, resExists := mod.localizedResources[lang][id]
	if !resExists {
		res = mod.newResource(lang, id)
		mod.localizedResources[lang][id] = res
	}
	return res
}

func (mod *Mod) newResource(lang resource.Language, id resource.ID) *MutableResource {
	compound := true
	contentType := resource.ContentType(0xFF) // Default to something completely unknown.
	compressed := false
	filename := "unknown.res"

	if info, known := ids.Info(id); known {
		compound = info.Compound
		contentType = info.ContentType
		compressed = info.Compressed
		filename = info.ResFile.For(lang)
	}

	return &MutableResource{
		filename:  filename,
		saveOrder: math.MaxInt32,

		compound:    compound,
		contentType: contentType,
		compressed:  compressed,
		blocks:      make(map[int][]byte),
	}
}

func (mod *Mod) delResource(lang resource.Language, id resource.ID) {
	deleteEntry := func(specificLang resource.Language, id resource.ID) {
		if lang.Includes(specificLang) {
			res, existing := mod.localizedResources[specificLang][id]
			if existing {
				mod.markFileChanged(res.filename)
				delete(mod.localizedResources[specificLang], id)
			}
		}
	}
	for _, worldLang := range resource.Languages() {
		deleteEntry(worldLang, id)
	}
	deleteEntry(resource.LangAny, id)
}

// Reset changes the mod to a new set of resources.
func (mod *Mod) Reset(newResources LocalizedResources, objectProperties object.PropertiesTable) {
	modifiedIDs := make(resource.IDMarkerMap)
	collectIDs := func(res LocalizedResources) {
		for _, resMap := range res {
			for id := range resMap {
				modifiedIDs.Add(id)
			}
		}
	}
	collectIDs(mod.localizedResources)
	collectIDs(newResources)

	mod.localizedResources = newResources
	mod.objectProperties = objectProperties
	mod.changedFiles = make(map[string]struct{})
	mod.lastChangeTime = time.Time{}
	mod.resetCallback()
	mod.resourcesChanged(modifiedIDs.ToList(), nil)
}

func (mod *Mod) markFileChanged(filename string) {
	mod.changedFiles[filename] = struct{}{}
	mod.lastChangeTime = time.Now()
}
