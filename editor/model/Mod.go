package model

import (
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world"
)

// Mod is the central object for a game-mod.
// It is based on a "static" world and adds its own changes. The world data itself is not static, it is merely the
// unchangeable background for the mod. Changes to the mod are kept in a separate layer, which can be loaded and saved.
type Mod struct {
	worldManifest    *world.Manifest
	resourcesChanged world.ResourceModificationCallback
}

// NewMod returns a new instance.
func NewMod(resourcesChanged world.ResourceModificationCallback) *Mod {
	var mod *Mod
	mod = &Mod{
		resourcesChanged: resourcesChanged,
	}
	mod.worldManifest = world.NewManifest(mod.worldChanged)

	return mod
}

// World returns the static background to the mod. Changes in the returned manifest may cause change callbacks
// being forwarded.
func (mod Mod) World() *world.Manifest {
	return mod.worldManifest
}

// LocalizedResources returns a resource selector for a specific language.
func (mod Mod) LocalizedResources(lang world.Language) world.ResourceSelector {
	// TODO requires find, which requires public visibility of Finder.
	return mod.worldManifest.LocalizedResources(lang)
}

// Modify requests to change the mod. The provided function will be called to collect all changes.
// After the modifier completes, all the requests will be applied and any changes notified.
func (mod *Mod) Modify(modifier func(*ModTransaction)) {

}

func (mod Mod) worldChanged(modifiedIDs []resource.ID, failedIDs []resource.ID) {
	mod.resourcesChanged(modifiedIDs, failedIDs)
}
