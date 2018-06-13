package text

// Codepage wraps the methods for serializing strings. The encoded byte array is expected
// to be terminated by one 0x00 byte.
type Codepage interface {
	// Encode converts the provided string to a byte array. It appends one 0x00 byte at the end.
	Encode(value string) []byte
	// Decode converts the provided byte slice to a string. It ignores trailing 0x00 bytes.
	Decode(data []byte) string
}
