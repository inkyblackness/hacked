package media

import (
	"bytes"

	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// MovieBlockSetter modifies storage of raw resource data.
type MovieBlockSetter interface {
	SetResourceBlocks(lang resource.Language, id resource.ID, data [][]byte)
	DelResource(lang resource.Language, id resource.ID)
}

// MovieSetterService can be used to set movie data.
type MovieSetterService struct {
	cp text.Codepage
}

// NewMovieSetterService returns a new instance.
func NewMovieSetterService(cp text.Codepage) MovieSetterService {
	return MovieSetterService{cp: cp}
}

// Remove deletes any movie resource for given key.
func (service MovieSetterService) Remove(setter MovieBlockSetter, key resource.Key) {
	setter.DelResource(key.Lang, key.ID)
}

// Set exports the given container.
func (service MovieSetterService) Set(setter MovieBlockSetter, key resource.Key, container movie.Container) {
	buf := bytes.NewBuffer(nil)
	_ = movie.Write(buf, container, service.cp)
	setter.SetResourceBlocks(key.Lang, key.ID, [][]byte{buf.Bytes()})
}
