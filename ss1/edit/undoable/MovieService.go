package undoable

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/edit"
	"github.com/inkyblackness/hacked/ss1/edit/undoable/cmd"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// MovieService provides read/write functionality with undo capability.
type MovieService struct {
	wrapped   edit.MovieService
	commander cmd.Commander
}

// NewMovieService returns a new instance of a service.
func NewMovieService(wrapped edit.MovieService, commander cmd.Commander) MovieService {
	return MovieService{
		wrapped:   wrapped,
		commander: commander,
	}
}

// Audio returns the audio component of identified movie.
func (service MovieService) Audio(key resource.Key) audio.L8 {
	return service.wrapped.Audio(key)
}
