package resource

// Properties describe the meta information about a resource.
type Properties struct {
	// Compound tells whether the resource should be serialized with a directory.
	// Compound resources can have zero, one, or more blocks.
	// Simple resources always have exactly one block.
	Compound bool

	// ContentType describes how the block data shall be interpreted.
	ContentType ContentType

	// Compressed tells whether the data shall be serialized in compressed form.
	Compressed bool
}

// Resource provides meta information as well as access to its contained blocks.
type Resource struct {
	Properties
	Blocks
}

// Compound tells whether the resource should be serialized with a directory.
// Compound resources can have zero, one, or more blocks.
// Simple resources always have exactly one block.
func (res Resource) Compound() bool {
	return res.Properties.Compound
}

// ContentType describes how the block data shall be interpreted.
func (res Resource) ContentType() ContentType {
	return res.Properties.ContentType
}

// Compressed tells whether the data shall be serialized in compressed form.
func (res Resource) Compressed() bool {
	return res.Properties.Compressed
}
