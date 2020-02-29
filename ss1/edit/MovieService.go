package edit

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ss1/edit/media"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// MovieService provides read/write functionality.
type MovieService struct {
	movieViewer media.MovieViewerService
}

// NewMovieService returns a new instance based on given accessor.
func NewMovieService(
	movieViewer media.MovieViewerService) MovieService {
	return MovieService{
		movieViewer: movieViewer,
	}
}

// Audio returns the audio component of identified movie.
func (service MovieService) Audio(key resource.Key) audio.L8 {
	return service.movieViewer.Audio(key)
}

// Subtitles returns the subtitles associated with the given key.
func (service MovieService) Subtitles(key resource.Key, language resource.Language) movie.Subtitles {
	return service.movieViewer.Subtitles(key, language)
}
