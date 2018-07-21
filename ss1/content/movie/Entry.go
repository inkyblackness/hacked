package movie

// Entry describes a block from a MOVI container.
type Entry interface {
	// Timestamp marks the beginning time of the entry, in seconds.
	Timestamp() float32
	// Type describes the content type of the data.
	Type() DataType
	// Data returns the raw bytes of the entry.
	Data() []byte
}
