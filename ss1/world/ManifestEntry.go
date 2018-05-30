package world

import "github.com/inkyblackness/hacked/ss1/resource"

// ManifestEntry describes a set of localized resources under a collective identifier.
type ManifestEntry struct {
	ID        string
	Resources []LocalizedResources

	// TODO: add texture properties and object properties.
}

func (entry ManifestEntry) filter(lang Language, id resource.ID) resource.List {
	var list resource.List
	for _, localized := range entry.Resources {
		if (localized.Language == LangAny) || (localized.Language == lang) {
			if res, err := localized.Provider.Resource(id); err == nil {
				list = list.With(res)
			}
		}
	}
	return list
}

// LocalizedResources produces a selector to retrieve resources for a specific language from this entry.
func (entry ManifestEntry) LocalizedResources(lang Language) ResourceSelector {
	return ResourceSelector{
		from: &entry,
		lang: lang,
	}
}
