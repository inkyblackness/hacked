package serial

// Codable is something than can be serialized with a Coder.
type Codable interface {
	// Code requests to serialize the Codable via the provided Coder.
	Code(coder Coder)
}
