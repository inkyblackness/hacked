package resource

import (
	"bytes"
	"crypto/md5" // nolint: gas
	"encoding/binary"
	"io"
)

type resourceHash []byte
type resourceHashes map[ID]resourceHash
type resourceHashSnapshot map[Language]resourceHashes

// IDMarkerMap is used to collect IDs.
type IDMarkerMap map[ID]struct{}

// Add adds the given ID to the map.
func (marker IDMarkerMap) Add(id ID) {
	marker[id] = struct{}{}
}

// ToList converts the map to a de-duplicated list.
func (marker IDMarkerMap) ToList() []ID {
	result := make([]ID, 0, len(marker))
	for id := range marker {
		result = append(result, id)
	}
	return result
}

// ModificationCallback is a callback function to notify a change in resources.
type ModificationCallback func(modifiedIDs []ID, failedIDs []ID)

// Localizer provides selectors to resources for specific languages.
type Localizer interface {
	LocalizedResources(lang Language) Selector
}

// ChangeNotifier is a utility that assists in detecting changes in modifying resources.
// A callback is called for any resource identifier that refers to resources that are different after a modification.
//
// Use this utility in combination to resource lists where resources can be overwritten by other entries.
// Changes in order, or content, affects how a resource is resolved.
type ChangeNotifier struct {
	Localizer Localizer
	Callback  ModificationCallback
}

// ModifyAndNotify must be called with a range of affected IDs that will change during the call of the modifier.
// A hash snapshot is taken before and after the modifier, considering the affected IDs.
// Any change is then reported to the callback, listing all IDs that have different hashes.
//
// Hashing the resources considers all languages, as well as the meta-information of the resources.
func (notifier ChangeNotifier) ModifyAndNotify(modifier func(), affectedIDs ...[]ID) {
	var flattenedIDs []ID
	for _, idList := range affectedIDs {
		flattenedIDs = append(flattenedIDs, idList...)
	}
	oldHashes, _ := notifier.hashAll(flattenedIDs)
	modifier()
	newHashes, failedIDs := notifier.hashAll(flattenedIDs)
	modifiedResourceIDs := notifier.modifiedIDs(oldHashes, newHashes)
	notifier.Callback(modifiedResourceIDs, failedIDs)
}

func (notifier ChangeNotifier) hashAll(ids []ID) (hashes resourceHashSnapshot, failedIDs []ID) {
	failedMap := make(IDMarkerMap)
	hashes = make(resourceHashSnapshot)
	for _, lang := range Languages() {
		hashByResourceID := make(resourceHashes)
		hashes[lang] = hashByResourceID
		selector := notifier.Localizer.LocalizedResources(lang)

		for _, id := range ids {
			hash, hashErr := notifier.hashResource(&selector, id)
			if hashErr == nil {
				hashByResourceID[id] = hash
			} else {
				failedMap.Add(id)
			}
		}
	}
	failedIDs = failedMap.ToList()
	return
}

func (notifier ChangeNotifier) hashResource(selector *Selector, id ID) (resourceHash, error) {
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

func (notifier ChangeNotifier) modifiedIDs(oldHashes resourceHashSnapshot, newHashes resourceHashSnapshot) []ID {
	modifiedIDs := make(IDMarkerMap)
	for _, lang := range Languages() {
		oldHashesByResourceID := oldHashes[lang]
		newHashesByResourceID := newHashes[lang]

		for newID, newHash := range newHashesByResourceID {
			oldHash, oldExisting := oldHashesByResourceID[newID]
			if !oldExisting || !bytes.Equal(newHash, oldHash) {
				modifiedIDs.Add(newID)
			}
		}
		for oldID := range oldHashesByResourceID {
			if _, newExisting := newHashesByResourceID[oldID]; !newExisting {
				modifiedIDs.Add(oldID)
			}
		}
	}
	return modifiedIDs.ToList()
}
