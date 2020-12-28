package format

const (
	// HeaderString is the start of an LG Res file.
	HeaderString = "LG Res File v2\r\n"
	// CommentTerminator terminates the header string.
	CommentTerminator = byte(0x1A)
	// ResourceDirectoryFileOffsetPos is the byte offset for the pointer to the directory.
	ResourceDirectoryFileOffsetPos = 0x7C
	// BoundarySize is the byte interval at which entries must be placed.
	BoundarySize = 4

	// ResourceTypeFlagCompressed indicates a compressed resource.
	ResourceTypeFlagCompressed = byte(0x01)
	// ResourceTypeFlagCompound indicates a compound resource.
	ResourceTypeFlagCompound = byte(0x02)
)
