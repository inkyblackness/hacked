package model

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
)

type identifiedResources map[resource.ID]*resource.Resource

// Mod is the central object for a game-mod.
//
// It is based on a "static" world and adds its own changes. The world data itself is not static, it is merely the
// unchangeable background for the mod. Changes to the mod are kept in a separate layer, which can be loaded and saved.
type Mod struct {
	worldManifest    *world.Manifest
	resourcesChanged world.ResourceModificationCallback

	localizedResources map[world.Language]identifiedResources
}

// NewMod returns a new instance.
func NewMod(resourcesChanged world.ResourceModificationCallback) *Mod {
	var mod *Mod
	mod = &Mod{
		resourcesChanged:   resourcesChanged,
		localizedResources: make(map[world.Language]identifiedResources),
	}
	mod.worldManifest = world.NewManifest(mod.worldChanged)
	for _, lang := range world.Languages() {
		mod.localizedResources[lang] = make(identifiedResources)
	}
	mod.localizedResources[world.LangAny] = make(identifiedResources)

	return mod
}

// World returns the static background to the mod. Changes in the returned manifest may cause change callbacks
// being forwarded.
func (mod Mod) World() *world.Manifest {
	return mod.worldManifest
}

// Filter returns a list of resources that match the given parameters.
func (mod Mod) Filter(lang world.Language, id resource.ID) resource.List {
	list := mod.worldManifest.Filter(lang, id)
	if res, resExists := mod.localizedResources[world.LangAny][id]; resExists {
		list = list.With(res)
	}
	for _, worldLang := range world.Languages() {
		if worldLang.Includes(lang) {
			if res, resExists := mod.localizedResources[lang][id]; resExists {
				list = list.With(res)
			}
		}
	}
	return list
}

// LocalizedResources returns a resource selector for a specific language.
func (mod Mod) LocalizedResources(lang world.Language) world.ResourceSelector {
	return world.ResourceSelector{
		Lang: lang,
		From: mod,
		As:   world.StandardResourceViewStrategy(),
	}
}

// Modify requests to change the mod. The provided function will be called to collect all changes.
// After the modifier completes, all the requests will be applied and any changes notified.
func (mod *Mod) Modify(modifier func(*ModTransaction)) {
	notifier := world.ResourceChangeNotifier{
		Callback:  mod.resourcesChanged,
		Localizer: mod,
	}
	var trans ModTransaction
	trans.modifiedIDs = make(idMarkerMap)
	modifier(&trans)
	notifier.ModifyAndNotify(func() {
		for _, action := range trans.actions {
			action(mod)
		}
	}, trans.modifiedIDs.toList())
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

func (mod *Mod) ensureResource(lang world.Language, id resource.ID) *resource.Resource {
	res, resExists := mod.localizedResources[lang][id]
	if !resExists {
		res = mod.newResource(lang, id)
		mod.localizedResources[lang][id] = res
	}
	return res
}

func (mod *Mod) newResource(lang world.Language, id resource.ID) *resource.Resource {
	// TODO: if not even existing, create based on defaults
	compound := false
	contentType := resource.Text
	compressed := false

	list := mod.worldManifest.Filter(lang, id)
	if len(list) > 0 {
		existing := list[0]
		compound = existing.Compound
		contentType = existing.ContentType
		compressed = existing.Compressed
	}

	return &resource.Resource{
		Compound:      compound,
		ContentType:   contentType,
		Compressed:    compressed,
		BlockProvider: resource.MemoryBlockProvider(nil),
	}
}

func (mod *Mod) delResource(lang world.Language, id resource.ID) {
	for _, worldLang := range world.Languages() {
		if lang.Includes(worldLang) {
			delete(mod.localizedResources[worldLang], id)
		}
	}
	if lang.Includes(world.LangAny) {
		delete(mod.localizedResources[world.LangAny], id)
	}
}
