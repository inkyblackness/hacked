package model

import (
	"math"

	"github.com/inkyblackness/hacked/ss1/resource"
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
	localizedResources LocalizedResources
}

// NewMod returns a new instance.
func NewMod(resourcesChanged resource.ModificationCallback, resetCallback ModResetCallback) *Mod {
	mod := &Mod{
		resourcesChanged:   resourcesChanged,
		resetCallback:      resetCallback,
		localizedResources: NewLocalizedResources(),
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

// ModifiedResource retrieves the resource of given language and ID.
// There is no fallback lookup, it will return the exact resource stored under the provided identifier.
// Returns nil if the resource does not exist.
func (mod Mod) ModifiedResource(lang resource.Language, id resource.ID) resource.View {
	return mod.localizedResources[lang][id]
}

// ModifiedBlock retrieves the specific block identified by given parameter.
// Returns empty slice if the block (or resource) is not modified.
func (mod Mod) ModifiedBlock(lang resource.Language, id resource.ID, index int) (data []byte) {
	res := mod.localizedResources[lang][id]
	if res == nil {
		return
	}
	return res.blocks[index]
}

// ModifiedBlocks returns all blocks of the modified resource.
func (mod Mod) ModifiedBlocks(lang resource.Language, id resource.ID) [][]byte {
	res := mod.localizedResources[lang][id]
	if res == nil {
		return nil
	}
	data := make([][]byte, res.blockCount)
	for index := 0; index < res.blockCount; index++ {
		data[index] = res.blocks[index]
	}
	return data
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
	for _, worldLang := range resource.Languages() {
		if lang.Includes(worldLang) {
			delete(mod.localizedResources[worldLang], id)
		}
	}
	if lang.Includes(resource.LangAny) {
		delete(mod.localizedResources[resource.LangAny], id)
	}
}

// Reset changes the mod to a new set of resources.
func (mod *Mod) Reset(newResources LocalizedResources) {
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
	mod.resetCallback()
	mod.resourcesChanged(modifiedIDs.ToList(), nil)
}
