package movie

import "github.com/inkyblackness/hacked/ss1/content/bitmap"

// MediaHandler is called from a MediaDispatcher on various media entries.
type MediaHandler interface {
	// OnAudio is called for an audio entry.
	OnAudio(timestamp Timestamp, samples []byte)
	// OnSubtitle is called for a subtitle entry.
	OnSubtitle(timestamp Timestamp, control SubtitleControl, text string)
	// OnVideo is called for a video entry.
	OnVideo(timestamp Timestamp, frame bitmap.Bitmap)
}
