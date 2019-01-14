package model

import (
	"io/ioutil"

	"github.com/inkyblackness/hacked/ss1/resource"
)

// MutableResource describes a resource open for modification.
type MutableResource struct {
	filename  string
	saveOrder int

	resource.Resource
}

// LocalizedResources is a map of identified resources by language.
type LocalizedResources map[resource.Language]IdentifiedResources

// NewLocalizedResources returns a new instance, prepared for all language keys.
func NewLocalizedResources() LocalizedResources {
	res := make(LocalizedResources)
	for _, lang := range resource.Languages() {
		res[lang] = make(IdentifiedResources)
	}
	res[resource.LangAny] = make(IdentifiedResources)
	return res
}

// MutableResourcesFromViewer returns MutableResource instances based on a viewer.
// This function retrieves all data from the viewer.
func MutableResourcesFromViewer(filename string, viewer resource.Provider) IdentifiedResources {
	ids := viewer.IDs()
	mutables := make(IdentifiedResources)
	for index, id := range ids {
		res, _ := viewer.View(id)
		mutable := &MutableResource{
			filename:  filename,
			saveOrder: index,

			Resource: resource.Resource{
				Properties: resource.Properties{
					Compound:    res.Compound(),
					ContentType: res.ContentType(),
					Compressed:  res.Compressed(),
				},
			},
		}
		blockCount := res.BlockCount()
		for blockIndex := 0; blockIndex < blockCount; blockIndex++ {
			mutable.SetBlock(blockIndex, readBlock(res, blockIndex))
		}
		mutables[id] = mutable
	}
	return mutables
}

func readBlock(provider resource.BlockProvider, index int) []byte {
	reader, err := provider.Block(index)
	if err != nil {
		return nil
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil
	}
	return data
}

// Filename returns the file this resource should be stored in.
func (res MutableResource) Filename() string {
	return res.filename
}
