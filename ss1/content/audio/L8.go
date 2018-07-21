package audio

// L8 is a raw, linear 8-bit sound snippet.
type L8 struct {
	SampleRate float32
	Samples    []byte
}
