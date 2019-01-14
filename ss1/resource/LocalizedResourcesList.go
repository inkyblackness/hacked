package resource

// LocalizedResourcesList is a collection of localized resources.
// it exists to provide a typical implementation of a Selector.
type LocalizedResourcesList []LocalizedResources

// Filter returns all resources that match the given parameters.
func (list LocalizedResourcesList) Filter(lang Language, id ID) List {
	var result List
	for _, localized := range list {
		if localized.Language.Includes(lang) {
			if res, err := localized.Provider.View(id); err == nil {
				result = result.With(res)
			}
		}
	}
	return result
}
