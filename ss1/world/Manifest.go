package world

import (
	"errors"

	"bytes"
	"crypto/md5" // nolint: gas
	"encoding/binary"
	"github.com/inkyblackness/hacked/ss1/resource"
	"io"
)

type resourceHash []byte
type resourceHashes map[resource.ID]resourceHash
type resourceHashSnapshot map[Language]resourceHashes
type idMarkerMap map[resource.ID]bool

func (marker idMarkerMap) toList() []resource.ID {
	result := make([]resource.ID, 0, len(marker))
	for id := range marker {
		result = append(result, id)
	}
	return result
}

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

func (manifest *Manifest) changeAndNotify(modifier func(), affectedIDs ...[]resource.ID) {
	var flattenedIDs []resource.ID
	for _, idList := range affectedIDs {
		flattenedIDs = append(flattenedIDs, idList...)
	}
	oldHashes, _ := manifest.hashAll(flattenedIDs)
	modifier()
	newHashes, failedIDs := manifest.hashAll(flattenedIDs)
	modifiedResourceIDs := manifest.modifiedIDs(oldHashes, newHashes)
	manifest.Modified(modifiedResourceIDs, failedIDs)
}

func (manifest Manifest) filter(lang Language, id resource.ID) resource.List {
	var list resource.List
	for _, entry := range manifest.entries {
		list = list.Joined(entry.filter(lang, id))
	}
	return list
}

// LocalizedResources produces a selector to retrieve from for a specific language from the manifest.
// The returned selector has the strategy to merge the typical compound resource lists, such
// as the small textures, or string lookups. It is based on StandardResourceViewStrategy().
func (manifest *Manifest) LocalizedResources(lang Language) ResourceSelector {
	return ResourceSelector{
		lang: lang,
		from: manifest,
		As:   StandardResourceViewStrategy(),
	}
}

func (manifest *Manifest) listIDs(entry *ManifestEntry) (ids []resource.ID) {
	for _, res := range entry.Resources {
		singleIDs := res.Provider.IDs()
		ids = append(ids, singleIDs...)
	}
	return
}

func (manifest *Manifest) hashAll(ids []resource.ID) (hashes resourceHashSnapshot, failedIDs []resource.ID) {
	failedMap := make(idMarkerMap)
	hashes = make(resourceHashSnapshot)
	for _, lang := range Languages() {
		hashByResourceID := make(resourceHashes)
		hashes[lang] = hashByResourceID
		selector := manifest.LocalizedResources(lang)

		for _, id := range ids {
			hash, hashErr := manifest.hashResource(&selector, id)
			if hashErr == nil {
				hashByResourceID[id] = hash
			} else {
				failedMap[id] = true
			}
		}
	}
	failedIDs = failedMap.toList()
	return
}

func (manifest *Manifest) hashResource(selector *ResourceSelector, id resource.ID) (resourceHash, error) {
	view, viewErr := selector.Select(id)
	if viewErr != nil {
		return nil, viewErr
	}
	hasher := md5.New() // nolint: gas
	blocks := view.BlockCount()
	for index := 0; index < blocks; index++ {
		blockReader, blockErr := view.Block(index)
		if blockErr != nil {
			return nil, blockErr
		}
		written, hashErr := io.Copy(hasher, blockReader)
		if hashErr != nil {
			return nil, hashErr
		}
		binary.Write(hasher, binary.LittleEndian, &written) // nolint: errcheck
	}
	binary.Write(hasher, binary.LittleEndian, int64(blocks))      // nolint: errcheck
	binary.Write(hasher, binary.LittleEndian, view.Compound())    // nolint: errcheck
	binary.Write(hasher, binary.LittleEndian, view.ContentType()) // nolint: errcheck
	binary.Write(hasher, binary.LittleEndian, view.Compressed())  // nolint: errcheck

	return hasher.Sum(nil), nil
}

func (manifest *Manifest) modifiedIDs(oldHashes resourceHashSnapshot, newHashes resourceHashSnapshot) []resource.ID {
	modifiedIDs := make(idMarkerMap)
	for _, lang := range Languages() {
		oldHashesByResourceID := oldHashes[lang]
		newHashesByResourceID := newHashes[lang]

		for newID, newHash := range newHashesByResourceID {
			oldHash, oldExisting := oldHashesByResourceID[newID]
			if !oldExisting || !bytes.Equal(newHash, oldHash) {
				modifiedIDs[newID] = true
			}
		}
		for oldID := range oldHashesByResourceID {
			if _, newExisting := newHashesByResourceID[oldID]; !newExisting {
				modifiedIDs[oldID] = true
			}
		}
	}
	return modifiedIDs.toList()
}
