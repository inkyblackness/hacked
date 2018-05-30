package world

import (
	"errors"

	"github.com/inkyblackness/hacked/ss1/resource"
)

// ManifestModificationCallback is a callback function to notify a change in the manifest.
type ManifestModificationCallback func(modifiedIDs []resource.ID, failedIDs []resource.ID)

// Manifest contains all the data and information of concrete things in a world.
type Manifest struct {
	// Modified is called whenever the final view on the manifest changes.
	Modified ManifestModificationCallback

	entries []*ManifestEntry
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
// Any entry at the identified index, and all those after that, are moved.
func (manifest *Manifest) InsertEntry(at int, entry *ManifestEntry) error {
	oldLen := len(manifest.entries)
	if (at < 0) || (at > oldLen) {
		return errIndexOutOfBounds
	}
	if entry == nil {
		return errEntryIsNil
	}
	newEntries := make([]*ManifestEntry, oldLen+1)
	copy(newEntries[:at], manifest.entries[:at])
	copy(newEntries[at+1:oldLen+1], manifest.entries[at:oldLen])
	newEntries[at] = entry
	manifest.entries = newEntries
	return nil
}

// RemoveEntry removes the entry at given index.
func (manifest *Manifest) RemoveEntry(at int) error {
	oldLen := len(manifest.entries)
	if (at < 0) || (at >= oldLen) {
		return errIndexOutOfBounds
	}
	copy(manifest.entries[at:oldLen-1], manifest.entries[at+1:oldLen])
	manifest.entries = manifest.entries[:oldLen-1]
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
	manifest.entries[at] = entry
	return nil
}

// MoveEntry removes an entry and reinserts it at another index.
// Both indices are resolved before the move.
func (manifest *Manifest) MoveEntry(to, from int) error {
	curLen := len(manifest.entries)
	if (to < 0) || (to >= curLen) || (from < 0) || (from >= curLen) {
		return errIndexOutOfBounds
	}
	entry := manifest.entries[from]
	if from >= to {
		copy(manifest.entries[to+1:from+1], manifest.entries[to:from])
	} else {
		copy(manifest.entries[from:to], manifest.entries[from+1:to+1])
	}
	manifest.entries[to] = entry
	return nil
}

// LocalizedResources produces a selector to retrieve from for a specific language from the manifest.
func (manifest *Manifest) LocalizedResources(lang Language) *ResourceSelector {
	return nil
}
