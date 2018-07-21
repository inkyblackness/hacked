package movie

import "github.com/inkyblackness/hacked/ss1/content/bitmap"

// MediaHandler is called from a MediaDispatcher on various media entries.
type MediaHandler interface {
	// OnAudio is called for an audio entry.
	OnAudio(timestamp float32, samples []byte)
	// OnSubtitle is called for a subtitle entry.
	OnSubtitle(timestamp float32, control SubtitleControl, text string)
	// OnVideo is called for a video entry.
	OnVideo(timestamp float32, frame bitmap.Bitmap)
}
