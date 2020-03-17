package movie

import (
	"github.com/inkyblackness/hacked/ss1/content/audio"
)

const audioEntrySize = 0x2000

// Audio represents the sound part of a movie.
type Audio struct {
	Sound audio.L8
}

// Duration returns the length of the audio stream.
func (a Audio) Duration() Timestamp {
	if a.Sound.SampleRate <= 0 {
		return Timestamp{}
	}
	return TimestampFromSeconds(a.Sound.Duration())
}

// Encode creates a list of buckets for writing a stream.
func (a Audio) Encode() []EntryBucket {
	buckets := make([]EntryBucket, 0, (len(a.Sound.Samples)/audioEntrySize)+1)
	addBucket := func(ts Timestamp, samples []byte) {
		buckets = append(buckets,
			EntryBucket{
				Priority:  EntryBucketPriorityAudio,
				Timestamp: ts,
				Entries: []Entry{
					{
						Timestamp: ts,
						Data: AudioEntryData{
							Samples: samples,
						},
					},
				},
			})
	}

	startOffset := 0
	for (startOffset + audioEntrySize) <= len(a.Sound.Samples) {
		ts := TimestampFromSeconds(float32(startOffset) / a.Sound.SampleRate)
		endOffset := startOffset + audioEntrySize
		addBucket(ts, a.Sound.Samples[startOffset:endOffset])
		startOffset = endOffset
	}
	if startOffset < len(a.Sound.Samples) {
		ts := TimestampFromSeconds(float32(startOffset) / a.Sound.SampleRate)
		addBucket(ts, a.Sound.Samples[startOffset:])
	}

	return buckets
}
