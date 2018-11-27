package audio

// L8 is a raw, linear 8-bit sound snippet.
type L8 struct {
	SampleRate float32
	Samples    []byte
}

// Empty returns true if there are no samples to play.
func (sound L8) Empty() bool {
	return len(sound.Samples) == 0
}

// Duration returns the length of the sound in seconds.
// TODO return time.Duration type
func (sound L8) Duration() float32 {
	if sound.SampleRate <= 0 {
		return 0.0
	}
	return float32(len(sound.Samples)) / sound.SampleRate
}
