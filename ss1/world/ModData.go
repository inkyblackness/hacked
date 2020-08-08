package world

import (
	"bytes"

	"github.com/inkyblackness/hacked/ss1/content/object"
	"github.com/inkyblackness/hacked/ss1/content/texture"
	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/serial/rle"
	"github.com/inkyblackness/hacked/ss1/world/ids"
)

// LocalizedResources associates a language with a resource store under a specific filename.
type LocalizedResources struct {
	File     FileLocation
	Language resource.Language
	Store    resource.Store
}

// ModData contains the core information about a mod.
type ModData struct {
	FileChangeCallback func(string)

	LocalizedResources []*LocalizedResources
	ObjectProperties   object.PropertiesTable
	TextureProperties  texture.PropertiesList
}

// SetResourceBlock changes the block data of a resource.
//
// If the block data is not empty, then:
// If the resource does not exist, it will be created with default meta information.
// If the block does not exist, the resource is extended to allow its addition.
//
// If the block data is empty (or nil), then the block is cleared.
func (data *ModData) SetResourceBlock(lang resource.Language, id resource.ID, index int, blockData []byte) {
	loc, res := data.ensureResource(lang, id)
	res.SetBlock(index, blockData)
	data.notifyFileChanged(loc.File.Name)
}

// PatchResourceBlock modifies an existing block.
// This modification assumes the block already exists and can take the given patch data.
// The patch data is expected to be produced by rle.Compress().
func (data *ModData) PatchResourceBlock(lang resource.Language, id resource.ID, index int, expectedLength int, patch []byte) {
	loc, res := data.ensureResource(lang, id)
	raw, err := res.BlockRaw(index)
	if (err == nil) && (len(raw) == expectedLength) {
		_ = rle.Decompress(bytes.NewReader(patch), raw)
		data.notifyFileChanged(loc.File.Name)
	}
}

// SetResourceBlocks sets the entire list of block data of a resource.
// This method is primarily meant for compound non-list resources (e.g. text pages).
func (data *ModData) SetResourceBlocks(lang resource.Language, id resource.ID, blocks [][]byte) {
	loc, res := data.ensureResource(lang, id)
	res.Set(blocks)
	data.notifyFileChanged(loc.File.Name)
}

// DelResource removes a resource from the mod in the given language.
func (data *ModData) DelResource(lang resource.Language, id resource.ID) {
	for _, loc := range data.LocalizedResources {
		if (loc.Language == lang) && loc.Store.Del(id) {
			data.notifyFileChanged(loc.File.Name)
		}
	}
}

// SetTextureProperties updates the properties of a specific texture.
func (data *ModData) SetTextureProperties(index int, properties texture.Properties) {
	if (index >= 0) && (index < len(data.TextureProperties)) {
		data.TextureProperties[index] = properties
		data.notifyFileChanged(TexturePropertiesFilename)
	}
}

// SetObjectProperties updates the properties of a specific object.
func (data *ModData) SetObjectProperties(triple object.Triple, properties object.Properties) {
	entry, err := data.ObjectProperties.ForObject(triple)
	if err != nil {
		return
	}
	*entry = properties.Clone()
	data.notifyFileChanged(ObjectPropertiesFilename)
}

func (data *ModData) ensureResource(lang resource.Language, id resource.ID) (*LocalizedResources, *resource.Resource) {
	for _, loc := range data.LocalizedResources {
		if loc.Language == lang {
			res, err := loc.Store.Resource(id)
			if err == nil {
				return loc, res
			}
		}
	}

	return data.newResource(lang, id)
}

func (data *ModData) newResource(lang resource.Language, id resource.ID) (*LocalizedResources, *resource.Resource) {
	compound := true
	contentType := resource.ContentType(0xFF) // Default to something completely unknown.
	compressed := false
	filename := "unknown.res"

	if info, known := ids.Info(id); known {
		compound = info.Compound
		contentType = info.ContentType
		compressed = info.Compressed
		filename = info.ResFile.For(lang)
	}

	loc := data.ensureStore(lang, filename)
	_ = loc.Store.Put(id, resource.Resource{
		Properties: resource.Properties{
			Compound:    compound,
			ContentType: contentType,
			Compressed:  compressed,
		},
	})
	res, _ := loc.Store.Resource(id)
	return loc, res
}

func (data *ModData) ensureStore(lang resource.Language, filename string) *LocalizedResources {
	for _, loc := range data.LocalizedResources {
		if loc.Language == lang && loc.File.Name == filename {
			return loc
		}
	}
	loc := &LocalizedResources{
		File:     FileLocation{DirPath: ".", Name: filename},
		Language: lang,
	}
	data.LocalizedResources = append(data.LocalizedResources, loc)
	return loc
}

func (data ModData) notifyFileChanged(filename string) {
	if data.FileChangeCallback != nil {
		data.FileChangeCallback(filename)
	}
}
