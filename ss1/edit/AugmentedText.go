package edit

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
	"github.com/inkyblackness/hacked/ss1/edit/media"
	"github.com/inkyblackness/hacked/ss1/resource"
)

// AugmentedTextBlockSetter modifies resource blocks.
type AugmentedTextBlockSetter interface {
	SetResourceBlock(lang resource.Language, id resource.ID, index int, data []byte)
	SetResourceBlocks(lang resource.Language, id resource.ID, data [][]byte)
	DelResource(lang resource.Language, id resource.ID)
}

func augmentedTextAudioKeyFor(key resource.Key) resource.Key {
	audioBase := KnownTexts().ByID(key.ID).audioBase
	if audioBase == 0 {
		return resource.KeyOf(0, resource.LangAny, 0)
	}

	soundKey := key
	soundKey.ID = audioBase.Plus(key.Index)
	soundKey.Index = 0
	return soundKey
}

// AugmentedTextService provides read/write functionality.
type AugmentedTextService struct {
	textViewer  media.TextViewerService
	textSetter  media.TextSetterService
	audioViewer media.AudioViewerService
	audioSetter media.AudioSetterService
}

// NewAugmentedTextService returns a new instance based on given accessor.
func NewAugmentedTextService(
	textViewer media.TextViewerService, textSetter media.TextSetterService,
	audioViewer media.AudioViewerService, audioSetter media.AudioSetterService) AugmentedTextService {
	return AugmentedTextService{
		textViewer:  textViewer,
		textSetter:  textSetter,
		audioViewer: audioViewer,
		audioSetter: audioSetter,
	}
}

// WithAudio returns true if the identified resource is one with an audio component.
func (service AugmentedTextService) WithAudio(key resource.Key) bool {
	return KnownTexts().ByID(key.ID).audioBase != 0
}

// Text returns the textual value of the identified text resource.
func (service AugmentedTextService) Text(key resource.Key) string {
	return service.textViewer.Text(key)
}

// SetText changes the textual value of a text resource.
func (service AugmentedTextService) SetText(setter AugmentedTextBlockSetter, key resource.Key, value string) {
	service.textSetter.Set(setter, key, value)
}

// RestoreTextFunc creates a snapshot of the current textual state and returns a function to restore it.
func (service AugmentedTextService) RestoreTextFunc(key resource.Key) func(setter AugmentedTextBlockSetter) {
	oldText := service.textViewer.Text(key)
	isModified := service.textViewer.Modified(key)

	return func(setter AugmentedTextBlockSetter) {
		if isModified {
			service.textSetter.Set(setter, key, oldText)
		} else {
			service.textSetter.Remove(setter, key)
		}
	}
}

// Sound returns the audio value of the identified text resource.
// In case the text resource has no audio, an empty sound will be returned.
func (service AugmentedTextService) Sound(key resource.Key) audio.L8 {
	if !service.WithAudio(key) {
		return audio.L8{}
	}
	return service.audioViewer.Audio(augmentedTextAudioKeyFor(key))
}

// SetSound changes the sound of a text resource.
// Should the text resource have no audio component, this call does nothing.
func (service AugmentedTextService) SetSound(setter AugmentedTextBlockSetter, key resource.Key, sound audio.L8) { // nolint: interfacer
	if service.WithAudio(key) {
		service.audioSetter.Set(setter, augmentedTextAudioKeyFor(key), sound)
	}
}

// RestoreSoundFunc creates a snapshot of the current sound and returns a function to restore it.
// In case the text resource has no audio, a stub method will be returned.
func (service AugmentedTextService) RestoreSoundFunc(key resource.Key) func(setter AugmentedTextBlockSetter) {
	if !service.WithAudio(key) {
		return func(setter AugmentedTextBlockSetter) {}
	}

	soundKey := augmentedTextAudioKeyFor(key)
	isSoundModified := service.audioViewer.Modified(soundKey)
	oldSound := service.audioViewer.Audio(soundKey)

	return func(setter AugmentedTextBlockSetter) {
		if isSoundModified {
			service.audioSetter.Set(setter, soundKey, oldSound)
		} else {
			service.audioSetter.Remove(setter, soundKey)
		}
	}
}

// Clear sets the text to an empty string and sets an empty sound if audio is associated.
func (service AugmentedTextService) Clear(setter AugmentedTextBlockSetter, key resource.Key) {
	service.textSetter.Clear(setter, key)
	if service.WithAudio(key) {
		service.audioSetter.Clear(setter, augmentedTextAudioKeyFor(key))
	}
}

// Remove erases the text and audio from the resources.
func (service AugmentedTextService) Remove(setter AugmentedTextBlockSetter, key resource.Key) {
	service.textSetter.Remove(setter, key)
	if service.WithAudio(key) {
		service.audioSetter.Remove(setter, augmentedTextAudioKeyFor(key))
	}
}

// RestoreFunc creates a snapshot of all associated media components and returns a function to restore it.
func (service AugmentedTextService) RestoreFunc(key resource.Key) func(setter AugmentedTextBlockSetter) {
	restoreText := service.RestoreTextFunc(key)
	restoreSound := service.RestoreSoundFunc(key)
	return func(setter AugmentedTextBlockSetter) {
		restoreText(setter)
		restoreSound(setter)
	}
}
