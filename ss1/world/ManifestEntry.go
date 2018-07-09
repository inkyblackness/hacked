package world

import (
	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// ManifestEntry describes a set of localized resources under a collective identifier.
type ManifestEntry struct {
	ID        string
	Resources resource.LocalizedResourcesList

	ObjectProperties object.PropertiesTable
	// TODO: add texture properties
}

// LocalizedResources produces a selector to retrieve resources for a specific language from this entry.
func (entry ManifestEntry) LocalizedResources(lang resource.Language) resource.Selector {
	return resource.Selector{
		From: entry.Resources,
		Lang: lang,
	}
}
