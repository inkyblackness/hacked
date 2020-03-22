package edit

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/content/movie"
	"github.com/inkyblackness/hacked/ss1/content/text"
	"github.com/inkyblackness/hacked/ss1/edit/media"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// MovieService provides read/write functionality.
type MovieService struct {
	cp text.Codepage

	movieViewer media.MovieViewerService
	movieSetter media.MovieSetterService
}

// NewMovieService returns a new instance based on given accessor.
func NewMovieService(codepage text.Codepage,
	movieViewer media.MovieViewerService, movieSetter media.MovieSetterService) MovieService {
	return MovieService{
		cp: codepage,

		movieViewer: movieViewer,
		movieSetter: movieSetter,
	}
}

// RestoreFunc creates a snapshot of the current movie and returns a function to restore it.
func (service MovieService) RestoreFunc(key resource.Key) func(setter media.MovieBlockSetter) {
	oldContainer, _ := service.movieViewer.Container(key)
	isModified := service.movieViewer.Modified(key)

	return func(setter media.MovieBlockSetter) {
		if isModified {
			service.movieSetter.Set(setter, key, oldContainer)
		} else {
			service.movieSetter.Remove(setter, key)
		}
	}
}

// Remove erases the movie from the resources.
func (service MovieService) Remove(setter media.MovieBlockSetter, key resource.Key) {
	service.movieSetter.Remove(setter, key)
}

// Video returns the video component of identified movie.
func (service MovieService) Video(key resource.Key) []movie.Scene {
	return service.movieViewer.Video(key)
}

// MoveSceneEarlier moves the given scene one step earlier.
func (service MovieService) MoveSceneEarlier(setter media.MovieBlockSetter, key resource.Key, scene int) {
	service.swapScenes(setter, key, scene, scene-1)
}

// MoveSceneLater moves the given scene one step later.
func (service MovieService) MoveSceneLater(setter media.MovieBlockSetter, key resource.Key, scene int) {
	service.swapScenes(setter, key, scene, scene+1)
}

func (service MovieService) swapScenes(setter media.MovieBlockSetter, key resource.Key, sceneA, sceneB int) {
	baseContainer := service.getBaseContainer(key)
	if (sceneA < 0) || (sceneA >= (len(baseContainer.Video.Scenes))) {
		return
	}
	if (sceneB < 0) || (sceneB >= (len(baseContainer.Video.Scenes))) {
		return
	}
	scenes := make([]movie.HighResScene, len(baseContainer.Video.Scenes))
	copy(scenes, baseContainer.Video.Scenes)
	scenes[sceneA] = baseContainer.Video.Scenes[sceneB]
	scenes[sceneB] = baseContainer.Video.Scenes[sceneA]
	baseContainer.Video.Scenes = scenes
	service.movieSetter.Set(setter, key, baseContainer)
}

// AddScene adds the given scene at the end of the movie.
func (service MovieService) AddScene(setter media.MovieBlockSetter, key resource.Key, scene movie.HighResScene) {
	baseContainer := service.getBaseContainer(key)
	baseContainer.Video.Scenes = append(baseContainer.Video.Scenes, scene)
	service.movieSetter.Set(setter, key, baseContainer)
}

// RemoveScene cuts out the given scene from the movie.
func (service MovieService) RemoveScene(setter media.MovieBlockSetter, key resource.Key, scene int) {
	baseContainer := service.getBaseContainer(key)
	if (scene < 0) || (scene >= len(baseContainer.Video.Scenes)) {
		return
	}
	scenes := make([]movie.HighResScene, len(baseContainer.Video.Scenes)-1)
	copy(scenes[0:scene], baseContainer.Video.Scenes[0:scene])
	copy(scenes[scene:], baseContainer.Video.Scenes[scene+1:])
	baseContainer.Video.Scenes = scenes
	service.movieSetter.Set(setter, key, baseContainer)
}

// Audio returns the audio component of identified movie.
func (service MovieService) Audio(key resource.Key) audio.L8 {
	return service.movieViewer.Audio(key)
}

// SetAudio sets the audio component of identified movie.
func (service MovieService) SetAudio(setter media.MovieBlockSetter, key resource.Key, soundData audio.L8) {
	baseContainer := service.getBaseContainer(key)
	baseContainer.Audio.Sound = soundData
	service.movieSetter.Set(setter, key, baseContainer)
}

// Subtitles returns the subtitles associated with the given key.
func (service MovieService) Subtitles(key resource.Key, language resource.Language) movie.SubtitleList {
	return service.movieViewer.Subtitles(key, language)
}

// SetSubtitles sets the subtitles of identified movie in given language.
func (service MovieService) SetSubtitles(setter media.MovieBlockSetter, key resource.Key,
	language resource.Language, subtitles movie.SubtitleList) {
	baseContainer := service.getBaseContainer(key)
	baseContainer.Subtitles.PerLanguage[language] = subtitles
	service.movieSetter.Set(setter, key, baseContainer)
}

func (service MovieService) getBaseContainer(key resource.Key) movie.Container {
	container, err := service.movieViewer.Container(key)
	if err != nil {
		container = movie.Container{
			Audio: movie.Audio{Sound: audio.L8{SampleRate: 22050}},
			Video: movie.Video{
				Width:  movie.HighResDefaultWidth,
				Height: movie.HighResDefaultHeight,
			},
		}
	}
	return container
}
