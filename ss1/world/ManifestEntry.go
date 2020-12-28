package world

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// ManifestEntry describes a set of localized resources under a collective identifier.
type ManifestEntry struct {
	Origin    []string
	ID        string
	Resources resource.LocalizedResourcesList

	ObjectProperties  object.PropertiesTable
	TextureProperties texture.PropertiesList
}

// NewManifestEntryFrom attempts to create a manifest in memory from the given set of files.
// It uses LoadFiles() to load the files with given filenames into memory, allowing archives as well.
func NewManifestEntryFrom(names []string) (*ManifestEntry, error) {
	loaded := LoadFiles(true, names)

	if len(loaded.Resources) == 0 {
		return nil, errNoResourcesFound
	}

	entry := &ManifestEntry{
		Origin: names,
		ID:     names[0],
	}

	for location, viewer := range loaded.Resources {
		localized := resource.LocalizedResources{
			ID:       location.Name,
			Language: ids.LocalizeFilename(location.Name),
			Viewer:   viewer,
		}
		entry.Resources = append(entry.Resources, localized)
	}
	entry.ObjectProperties = loaded.ObjectProperties
	entry.TextureProperties = loaded.TextureProperties
	return entry, nil
}

// LocalizedResources produces a selector to retrieve resources for a specific language from this entry.
func (entry ManifestEntry) LocalizedResources(lang resource.Language) resource.Selector {
	return resource.Selector{
		From: entry.Resources,
		Lang: lang,
	}
}
