package world

import "github.com/inkyblackness/hacked/ss1/resource"

// LocalizedResources associates a language with a resource provider under a specific identifier.
type LocalizedResources struct {
	// ID is the identifier of the provider. This could be a filename for instance.
	ID string
	// Language specifies for which language the provider has from.
	Language Language
	// Provider is the actual container of the from.
	Provider resource.Provider
}
