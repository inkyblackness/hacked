package world

import (
	"errors"

	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// Manifest contains all the data and information of concrete things in a world.
type Manifest struct {
	resourceChangeNotifier resource.ChangeNotifier
	entries                []*ManifestEntry
}

// NewManifest returns a new instance that notifies changes to the provided callback.
func NewManifest(modified resource.ModificationCallback) *Manifest {
	var manifest Manifest

	manifest.resourceChangeNotifier.Callback = modified
	manifest.resourceChangeNotifier.Localizer = &manifest

	return &manifest
}

var errIndexOutOfBounds = errors.New("index is out of bounds")
var errEntryIsNil = errors.New("entry is nil")

// EntryCount returns the number of entries currently in the manifest.
func (manifest Manifest) EntryCount() int {
	return len(manifest.entries)
}

// Entry returns the entry at given index.
func (manifest Manifest) Entry(at int) (*ManifestEntry, error) {
	if (at < 0) || (at >= len(manifest.entries)) {
		return nil, errIndexOutOfBounds
	}
	return manifest.entries[at], nil
}

// InsertEntry puts the provided entry at the specified index.
// Any entry at the identified index, and all those after that, are moved by one.
func (manifest *Manifest) InsertEntry(at int, entry *ManifestEntry) error {
	oldLen := len(manifest.entries)
	if (at < 0) || (at > oldLen) {
		return errIndexOutOfBounds
	}
	if entry == nil {
		return errEntryIsNil
	}
	manifest.changeAndNotify(func() {
		newEntries := make([]*ManifestEntry, oldLen+1)
		copy(newEntries[:at], manifest.entries[:at])
		copy(newEntries[at+1:oldLen+1], manifest.entries[at:oldLen])
		newEntries[at] = entry
		manifest.entries = newEntries
	}, manifest.listIDs(entry))
	return nil
}

// RemoveEntry removes the entry at given index.
func (manifest *Manifest) RemoveEntry(at int) error {
	oldLen := len(manifest.entries)
	if (at < 0) || (at >= oldLen) {
		return errIndexOutOfBounds
	}
	manifest.changeAndNotify(func() {
		copy(manifest.entries[at:oldLen-1], manifest.entries[at+1:oldLen])
		manifest.entries = manifest.entries[:oldLen-1]
	}, manifest.listIDs(manifest.entries[at]))
	return nil
}

// ReplaceEntry removes the current entry at the identified index and puts the provided one instead.
func (manifest *Manifest) ReplaceEntry(at int, entry *ManifestEntry) error {
	if (at < 0) || (at >= len(manifest.entries)) {
		return errIndexOutOfBounds
	}
	if entry == nil {
		return errEntryIsNil
	}
	manifest.changeAndNotify(func() {
		manifest.entries[at] = entry
	}, manifest.listIDs(manifest.entries[at]), manifest.listIDs(entry))

	return nil
}

// MoveEntry removes an entry and reinserts it at another index.
// Both indices are resolved before the move.
func (manifest *Manifest) MoveEntry(to, from int) error {
	curLen := len(manifest.entries)
	if (to < 0) || (to >= curLen) || (from < 0) || (from >= curLen) {
		return errIndexOutOfBounds
	}
	manifest.changeAndNotify(func() {
		entry := manifest.entries[from]
		if from >= to {
			copy(manifest.entries[to+1:from+1], manifest.entries[to:from])
		} else {
			copy(manifest.entries[from:to], manifest.entries[from+1:to+1])
		}
		manifest.entries[to] = entry
	}, manifest.listIDs(manifest.entries[from]))
	return nil
}

// Filter finds all resources in the world that match the given parameters.
func (manifest Manifest) Filter(lang resource.Language, id resource.ID) resource.List {
	var list resource.List
	for _, entry := range manifest.entries {
		list = list.Joined(entry.Resources.Filter(lang, id))
	}
	return list
}

// LocalizedResources produces a selector to retrieve resources for a specific language from the manifest.
// The returned selector has the strategy to merge the typical compound resource lists, such
// as the small textures, or string lookups. It is based on StandardResourceViewStrategy().
func (manifest *Manifest) LocalizedResources(lang resource.Language) resource.Selector {
	return resource.Selector{
		Lang: lang,
		From: manifest,
		As:   ResourceViewStrategy(),
	}
}

// ObjectProperties returns the table of object properties.
func (manifest *Manifest) ObjectProperties() object.PropertiesTable {
	var table object.PropertiesTable
	for _, entry := range manifest.entries {
		if len(entry.ObjectProperties) > 0 {
			table = entry.ObjectProperties
		}
	}
	return table
}

func (manifest *Manifest) listIDs(entry *ManifestEntry) (ids []resource.ID) {
	for _, res := range entry.Resources {
		singleIDs := res.Provider.IDs()
		ids = append(ids, singleIDs...)
	}
	return
}

func (manifest *Manifest) changeAndNotify(modifier func(), affectedIDs ...[]resource.ID) {
	manifest.resourceChangeNotifier.ModifyAndNotify(modifier, affectedIDs...)
}
